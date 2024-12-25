// Package flagset wraps the standard library flag package and focuses on the
// flag.FlagSet type. This is done to simplify usage, and to add handling for
// single hyphen (e.g -h) and double hyphen (e.g. --help) flags. Specifically,
// single hyphen prefixed values with multiple characters are exploded out as
// though they were their own flag (e.g. -abc = -a -b -c).
package flagset

import (
	"encoding"
	"errors"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"
)

// TODO: add sensible errors

// FlagSet contains flag options and related information used for usage output.
type FlagSet struct {
	fs     *flag.FlagSet
	flags  []*Flag
	parsed []string

	tmplFuncMap template.FuncMap
	tmplTxt     string

	HideTypeHints    bool
	HideDefaultHints bool
	Meta             map[string]any
}

// New constructs a FlagSet. In this package, it is conventional to name the
// flagset after the command that the options are being associated with.
func New(name string) *FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	return &FlagSet{
		fs: fs,
		tmplFuncMap: template.FuncMap{
			"Join": strings.Join,
		},
		tmplTxt: tmplText,
		Meta:    map[string]any{},
	}
}

// Flags returns all flag options that have been set.
func (fs *FlagSet) Flags() []*Flag {
	return fs.flags
}

// Parsed returns the args provided when Parse was called with any single hyphen
// flags containing multiple characters exploded to their own entries. The
// returned value can be helpful for debugging.
func (fs *FlagSet) Parsed() []string {
	return fs.parsed
}

// Operand returns the i'th operand. Operand(0) is the first remaining argument
// after flags have been processed. Operand returns an empty string if the
// requested element does not exist. This value is determined after single
// hyphen flags with multiple characters have been exploded.
func (fs *FlagSet) Operand(i int) string {
	return fs.fs.Arg(i)
}

// Operands returns the non-flag arguments.
func (fs *FlagSet) Operands() []string {
	return fs.fs.Args()
}

// NFlag returns the number of command-line flags that have been set. This
// value is determined after single hyphen flags with multiple characters have
// been exploded.
func (fs *FlagSet) NFlag() int {
	return fs.fs.NFlag()
}

// Name returns the name of the FlagSet set during construction.
func (fs *FlagSet) Name() string {
	return fs.fs.Name()
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet are
// defined and before flags are accessed by the program. Before parsing occurs,
// all single hyphen flags with multiple characters are exploded out as though
// they were their own flag (e.g. -abc = -a -b -c).
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
	}

	return nil
}

// Flag adds a flag option to the FlagSet.
// Valid values are: *string, *bool, *int, *int64, *uint, *uint64, *float64,
// *time.Duration, TextMarshalUnmarshaler, flag.Value, FlagFunc
// Names can include multiple long and multiple short values. Each value should
// be separated by a pipe (|) character. If val has a usable non-zero value, it
// will be used as the default value for that flag option.
func (fs *FlagSet) Flag(val any, names, desc string) *Flag {
	if reflect.ValueOf(val).Kind() == reflect.Func {
		vto := reflect.TypeOf(val)
		errIface := reflect.TypeOf((*error)(nil)).Elem()
		if vto.In(0).Kind() == reflect.String && vto.Out(0).Implements(errIface) {
			val = FlagFunc(val.(func(string) error))
		}
		if vto.In(0).Kind() == reflect.Bool && vto.Out(0).Implements(errIface) {
			val = FlagBoolFunc(val.(func(bool) error))
		}
	}

	longs, shorts := longsAndShorts(names)

	for _, long := range longs {
		addFlagTo(fs.fs, val, long, desc)
	}

	for _, short := range shorts {
		addFlagTo(fs.fs, val, short, desc)
	}

	flag := newFlag(val, longs, shorts, desc)
	fs.flags = append(fs.flags, flag)

	return flag
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

// TextMarshalUnmarshaler descibes types that are able to be marshaled to and
// unmarshaled from text.
type TextMarshalUnmarshaler interface {
	encoding.TextUnmarshaler
	encoding.TextMarshaler
}

// FlagFunc describes functions that can be called when a flag option is
// succesfully parsed. Currently, this cannot pass errors values back to callers
// as the stdlib flag pkg eats them.
type FlagFunc func(string) error

// FlagBoolFunc describes functions that can be called when a bool flag option
// is succesfully parsed. Currently, this cannot pass errors values back to
// callers as the stdlib flag pkg eats them.
type FlagBoolFunc func(bool) error

func addFlagTo(fs *flag.FlagSet, val any, flagName, desc string) {
	switch v := val.(type) {
	case *string:
		fs.StringVar(v, flagName, *v, desc)
	case *bool:
		fs.BoolVar(v, flagName, *v, desc)
	case *int:
		fs.IntVar(v, flagName, *v, desc)
	case *int64:
		fs.Int64Var(v, flagName, *v, desc)
	case *uint:
		fs.UintVar(v, flagName, *v, desc)
	case *uint64:
		fs.Uint64Var(v, flagName, *v, desc)
	case *float64:
		fs.Float64Var(v, flagName, *v, desc)
	case *time.Duration:
		fs.DurationVar(v, flagName, *v, desc)
	case TextMarshalUnmarshaler:
		fs.TextVar(v, flagName, v, desc)
	case flag.Value:
		fs.Var(v, flagName, desc)
	case FlagFunc:
		fs.Func(flagName, desc, v)
	case FlagBoolFunc:
		fn := func(s string) error {
			if s == "" {
				return v(true)
			}
			b, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			return v(b)
		}
		fs.BoolFunc(flagName, desc, fn)
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
