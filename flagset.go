package flagset

import (
	"encoding"
	"errors"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

type FlagSet struct {
	fs     *flag.FlagSet
	opts   []Opt
	parsed []string

	HideTypeHint    bool
	HideDefaultHint bool
}

func New(name string) *FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	return &FlagSet{
		fs: fs,
	}
}

func (fs *FlagSet) Opts() []Opt {
	return fs.opts
}

func (fs *FlagSet) Parsed() []string {
	return fs.parsed
}

func (fs *FlagSet) Arg(i int) string {
	return fs.fs.Arg(i)
}

func (fs *FlagSet) Args() []string {
	return fs.fs.Args()
}

func (fs *FlagSet) Lookup(name string) *flag.Flag {
	return fs.fs.Lookup(name)
}

func (fs *FlagSet) NArg() int {
	return fs.fs.NArg()
}

func (fs *FlagSet) NFlag() int {
	return fs.fs.NFlag()
}

func (fs *FlagSet) Name() string {
	return fs.fs.Name()
}

func (fs *FlagSet) Parse(arguments []string) error {
	fs.parsed = explodeShortArgs(arguments)

	if err := fs.fs.Parse(fs.parsed); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			return fmt.Errorf("flagset: parse: %w", mayWrapNotDefined(err))
		}

		if h, ok := findFirstHelp(arguments); ok {
			err = fmt.Errorf("flagset: parse: flag provided but not defined: %s", h)
			return mayWrapNotDefined(err)
		}

		return nil
	}

	return nil
}

func (fs *FlagSet) Visit(fn func(*flag.Flag)) {
	fs.fs.Visit(fn)
}

func (fs *FlagSet) VisitAll(fn func(*flag.Flag)) {
	fs.fs.VisitAll(fn)
}

func (fs *FlagSet) Opt(val any, names, usage string) *Opt {
	longs, shorts := longsAndShorts(names)

	for _, long := range longs {
		addOptTo(fs.fs, val, long, usage)
	}

	for _, short := range shorts {
		addOptTo(fs.fs, val, short, usage)
	}

	if reflect.ValueOf(val).Kind() == reflect.Func {
		vto := reflect.TypeOf(val)
		errIface := reflect.TypeOf((*error)(nil)).Elem()
		if vto.In(0).Kind() == reflect.String && vto.Out(0).Implements(errIface) {
			val = OptFunc(val.(func(string) error))
		}
	}

	typName := typeName(val)
	defTxt := defaultText(val)

	opt := makeOpt(fs, names, longs, shorts, typName, defTxt, usage)
	fs.opts = append(fs.opts, opt)

	return &opt
}

func explodeShortArgs(args []string) []string {
	var exed []string

	for _, arg := range args {
		if len(arg) > 1 && arg[0] == '-' && arg[1] != '-' {
			for _, a := range arg[1:] {
				exed = append(exed, "-"+string(a))
			}
			continue
		}

		exed = append(exed, arg)
	}

	return exed
}

func findFirstHelp(args []string) (string, bool) {
	for _, arg := range args {
		if arg == "-h" || arg == "--h" || arg == "--help" {
			return arg, true
		}
	}
	return "", false
}

func longsAndShorts(flags string) (longs, shorts []string) {
	fs := strings.Split(flags, "|")
	for _, f := range fs {
		if utf8.RuneCountInString(f) == 1 {
			shorts = append(shorts, f)
			continue
		}
		longs = append(longs, f)
	}
	return longs, shorts
}

type TextMarshalUnmarshaler interface {
	encoding.TextUnmarshaler
	encoding.TextMarshaler
}

type FlagValue = flag.Value

type OptFunc func(string) error

func addOptTo(fs *flag.FlagSet, val any, flagName, usage string) {
	switch v := val.(type) {
	case *string:
		fs.StringVar(v, flagName, *v, usage)
	case *bool:
		fs.BoolVar(v, flagName, *v, usage)
	case *int:
		fs.IntVar(v, flagName, *v, usage)
	case *int64:
		fs.Int64Var(v, flagName, *v, usage)
	case *uint:
		fs.UintVar(v, flagName, *v, usage)
	case *uint64:
		fs.Uint64Var(v, flagName, *v, usage)
	case *float64:
		fs.Float64Var(v, flagName, *v, usage)
	case *time.Duration:
		fs.DurationVar(v, flagName, *v, usage)
	case TextMarshalUnmarshaler:
		fs.TextVar(v, flagName, v, usage)
	case flag.Value:
		fs.Var(v, flagName, usage)
	case OptFunc:
		fs.Func(flagName, usage, v)
	}
}

func typeName(val any) string {
	switch val.(type) {
	case TextMarshalUnmarshaler, OptFunc:
		return ""
	case flag.Value:
		return "value"
	default:
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		return v.Type().Name()
	}
}

func defaultText(val any) string {
	switch v := val.(type) {
	case TextMarshalUnmarshaler, OptFunc:
		return ""
	case fmt.Stringer:
		return v.String()
	default:
		vo := reflect.ValueOf(val)
		if vo.Kind() == reflect.Ptr {
			vo = vo.Elem()
		}
		return fmt.Sprint(vo)
	}
}

type transparentError struct {
	err error
	msg string
}

func (e *transparentError) Error() string {
	return e.msg
}

func (e *transparentError) Unwrap() error {
	return e.err
}

func mayWrapNotDefined(err error) error {
	if !strings.Contains(err.Error(), "but not defined:") {
		return err
	}

	token := "defined: -"
	msg := err.Error()
	_, flag, ok := strings.Cut(msg, token)
	if ok && len(flag) > 1 && flag[0] != '-' {
		msg = strings.ReplaceAll(msg, token, token+"-")
	}

	return &transparentError{err, msg}
}
