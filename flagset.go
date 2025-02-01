// Package flagset wraps the standard library [flag] package and focuses on the
// [flag.FlagSet] type. This is done to simplify usage, and to add handling for
// single hyphen (e.g -h) and double hyphen (e.g. --help) flags. Specifically,
// single hyphen prefixed values with multiple characters are exploded out as
// though each character was its own flag (e.g. -abc = -a -b -c).
package flagset

import (
	"strings"
	"unicode/utf8"

	"github.com/daved/flagset/vtype"
)

// FlagSet contains flag options and related information used for usage output.
// The exported fields are for easy post-construction configuration.
type FlagSet struct {
	// Fields used for templating:
	HideTypeHints    bool
	HideDefaultHints bool
	Meta             map[string]any

	name   string
	flags  []*Flag
	parsed []string
	ops    []string

	tmplCfg *TmplConfig
}

// New constructs a FlagSet. In this package, it is conventional to name the
// flagset after the command that the options are being associated with.
func New(name string) *FlagSet {
	fs := &FlagSet{
		name:    name,
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
	if i >= len(fs.ops) {
		return ""
	}
	return fs.ops[i]
}

// Operands returns the non-flag arguments.
func (fs *FlagSet) Operands() []string {
	return fs.ops
}

// Name returns the name of the FlagSet set during construction.
func (fs *FlagSet) Name() string {
	return fs.name
}

// Parse parses flag definitions from the argument list, which must not	include
// the initial command name. Parse must be called after all flags in the FlagSet
// are defined and before flag values are accessed by the program. Before
// parsing occurs, all single hyphen flags with multiple characters are exploded
// out as though they are their own flag (e.g. -abc = -a -b -c).
func (fs *FlagSet) Parse(args []string) error {
	fs.parsed = explodeShortArgs(args)

	ops, err := resolve(fs.flags, fs.parsed)
	if err != nil {
		return err
	}

	fs.ops = ops
	return nil
}

// Flag adds a flag option to the FlagSet.
// Valid values are:
//   - builtin: *string, *bool, *int, *int64, *uint, *uint64, *float64
//   - stdlib: *[time.Duration], [flag.Value]
//   - vtype: [vtype.TextMarshalUnmarshaler], [vtype.FlagCallback],
//     [vtype.FlagFunc], [vtype.FlagBoolFunc]
//
// Names can include multiple long and multiple short values. Each value should
// be separated by a pipe (|) character. If val has a usable non-zero value, it
// will be used as the default value for that flag option.
func (fs *FlagSet) Flag(val any, names, desc string) *Flag {
	switch v := val.(type) {
	case vtype.FlagFunc:
		val = vtype.FlagCallback(v)
	case vtype.FlagBoolFunc:
		val = vtype.FlagCallback(v)
	case func(string) error:
		val = vtype.FlagCallback(vtype.FlagFunc(v))
	case func(bool) error:
		val = vtype.FlagCallback(vtype.FlagBoolFunc(v))
	}

	longs, shorts := longsAndShorts(names)

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

// SetUsageTemplating is used to override the base template text, and provide a
// custom FuncMap. If a nil FuncMap is provided, no change will be made to the
// existing value.
func (fs *FlagSet) SetUsageTemplating(tmplCfg *TmplConfig) {
	fs.tmplCfg = tmplCfg
}

// Usage returns the executed usage template. Each Flag type's Meta field can
// be leveraged to convey detailed info/behavior in a custom template.
func (fs *FlagSet) Usage() string {
	return executeTmpl(fs.tmplCfg, &TmplData{FlagSet: fs})
}
