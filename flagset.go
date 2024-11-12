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
	"strings"
	"time"
	"unicode/utf8"
)

// FlagSet contains flag options and related information used for usage output.
type FlagSet struct {
	fs      *flag.FlagSet
	opts    []Opt
	parsed  []string
	tmplTxt string

	MetaHideTypeHints    bool
	MetaHideDefaultHints bool
}

// New constructs a FlagSet. In this package, it is conventional to name the
// flagset after the command that the options are being associated with.
func New(name string) *FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	return &FlagSet{
		fs:      fs,
		tmplTxt: tmplText,
	}
}

// Opts returns all flag options that have been set.
func (fs *FlagSet) Opts() []Opt {
	return fs.opts
}

// Parsed returns the args provided when Parse was called with any single hyphen
// flags containing multiple characters exploded to their own entries. The
// returned value can be helpful for debugging.
func (fs *FlagSet) Parsed() []string {
	return fs.parsed
}

// Arg returns the i'th argument. Arg(0) is the first remaining argument after
// flags have been processed. Arg returns an empty string if the requested
// element does not exist. This value is determined after single hyphen flags
// with multiple characters have been exploded.
func (fs *FlagSet) Arg(i int) string {
	return fs.fs.Arg(i)
}

// Args returns the non-flag arguments.
func (fs *FlagSet) Args() []string {
	return fs.fs.Args()
}

// NArg is the number of arguments remaining after flags have been processed.
func (fs *FlagSet) NArg() int {
	return fs.fs.NArg()
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

// SetUsageTemplate allows callers to override the base template text.
func (fs *FlagSet) SetUsageTemplate(txt string) {
	fs.tmplTxt = txt
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

// Opt adds a flag option to the FlagSet.
// Valid values are: *string, *bool, *int, *int64, *uint, *uint64, *float64,
// *time.Duration, TextMarshalUnmarshaler, flag.Value, OptFunc
// Names can include multiple long and multiple short values. Each value should
// be separated by a pipe (|) character. If val has a usable non-zero value, it
// will be used as the default value for that flag option.
func (fs *FlagSet) Opt(val any, names, usage string) *Opt {
	if reflect.ValueOf(val).Kind() == reflect.Func {
		vto := reflect.TypeOf(val)
		errIface := reflect.TypeOf((*error)(nil)).Elem()
		if vto.In(0).Kind() == reflect.String && vto.Out(0).Implements(errIface) {
			val = OptFunc(val.(func(string) error))
		}
	}

	longs, shorts := longsAndShorts(names)

	for _, long := range longs {
		addOptTo(fs.fs, val, long, usage)
	}

	for _, short := range shorts {
		addOptTo(fs.fs, val, short, usage)
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

// TextMarshalUnmarshaler descibes types that are able to be marshaled to and
// unmarshaled from text.
type TextMarshalUnmarshaler interface {
	encoding.TextUnmarshaler
	encoding.TextMarshaler
}

// FlagValue is an alias for flag.Value and provided for visibility.
type FlagValue = flag.Value

// OptFunc describes functions that can be called when a flag option is
// succesfully parsed.
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
