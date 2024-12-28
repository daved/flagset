// Package flagset wraps the standard library flag package and focuses on the
// flag.FlagSet type. This is done to simplify usage, and to add handling for
// single hyphen (e.g -h) and double hyphen (e.g. --help) flags. Specifically,
// single hyphen prefixed values with multiple characters are exploded out as
// though they were their own flag (e.g. -abc = -a -b -c).
package flagset

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/daved/flagset/vtypes"
)

// TODO: add sensible errors

// FlagSet contains flag options and related information used for usage output.
type FlagSet struct {
	sfs    *flag.FlagSet
	flags  []*Flag
	parsed []string

	tmplCfg *TmplConfig

	HideTypeHints    bool
	HideDefaultHints bool
	Meta             map[string]any
}

// New constructs a FlagSet. In this package, it is conventional to name the
// flagset after the command that the options are being associated with.
func New(name string) *FlagSet {
	sfs := flag.NewFlagSet(name, flag.ContinueOnError)
	sfs.SetOutput(io.Discard)

	fs := &FlagSet{
		sfs:     sfs,
		tmplCfg: NewDefaultTmplConfig(),
		Meta:    map[string]any{},
	}

	return fs
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
	return fs.sfs.Arg(i)
}

// Operands returns the non-flag arguments.
func (fs *FlagSet) Operands() []string {
	return fs.sfs.Args()
}

// Name returns the name of the FlagSet set during construction.
func (fs *FlagSet) Name() string {
	return fs.sfs.Name()
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet are
// defined and before flags are accessed by the program. Before parsing occurs,
// all single hyphen flags with multiple characters are exploded out as though
// they were their own flag (e.g. -abc = -a -b -c).
func (fs *FlagSet) Parse(args []string) error {
	fs.parsed = explodeShortArgs(args)

	if err := fs.sfs.Parse(fs.parsed); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			return fmt.Errorf("flagset: parse: %w", mayWrapNotDefined(err))
		}

		if h, ok := findFirstHelp(args); ok {
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
	switch v := val.(type) {
	case vtypes.FlagFunc:
		val = vtypes.FlagCallback(v)
	case vtypes.FlagBoolFunc:
		val = vtypes.FlagCallback(v)
	case func(string) error:
		val = vtypes.FlagCallback(vtypes.FlagFunc(v))
	case func(bool) error:
		val = vtypes.FlagCallback(vtypes.FlagBoolFunc(v))
	}

	longs, shorts := longsAndShorts(names)

	for _, long := range longs {
		addFlagTo(fs.sfs, val, long, desc)
	}

	for _, short := range shorts {
		addFlagTo(fs.sfs, val, short, desc)
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
	case vtypes.TextMarshalUnmarshaler:
		fs.TextVar(v, flagName, v, desc)
	case flag.Value:
		fs.Var(v, flagName, desc)
	case vtypes.FlagCallback:
		if v.IsBool() {
			fs.BoolFunc(flagName, desc, v.OnFlag)
		} else {
			fs.Func(flagName, desc, v.OnFlag)
		}
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

// SetUsageTemplating allows callers to override the base template text, and
// provide a custom FuncMap. If a nil FuncMap is provided, no change will be
// made to the existing value.
func (fs *FlagSet) SetUsageTemplating(tmplCfg *TmplConfig) {
	fs.tmplCfg = tmplCfg
}

// Usage returns the parsed usage template. Each Flag type's Meta field is
// leveraged to convey detailed info/behavior. This method and related template
// can be used as an example for callers to wrap the FlagSet type and design
// their own usage output. For example, grouping, sorting, etc.
func (fs *FlagSet) Usage() string {
	return executeTmpl(fs.tmplCfg, &TmplData{FlagSet: fs})
}
