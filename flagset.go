// Package flagset provides simple flag and flag value handling using idiomatic
// techniques for advanced usage. Nomenclature and handling rules are based on
// POSIX standards. For example, all single hyphen prefixed arguments with
// multiple characters are exploded out as though they are their own flags (e.g.
// -abc = -a -b -c).
package flagset

import er "github.com/daved/flagset/fserrs"

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
// parsing occurs, all single hyphen prefixed arguments with multiple characters
// are exploded out as though they are their own flags (e.g. -abc = -a -b -c).
func (fs *FlagSet) Parse(args []string) error {
	fs.parsed = explodeShortArgs(args)

	ops, err := resolveFlags(fs.flags, fs.parsed)
	if err != nil {
		return er.NewError(er.NewParseError(err))
	}

	fs.ops = ops
	return nil
}

// Flag adds a flag option to the FlagSet.
// Valid values are:
//   - builtin: *string, *bool, error, *int, *int8, *int16, *int32, *int64,
//     *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64
//   - stdlib: *[time.Duration], [flag.Value]
//   - vtype: [vtype.TextMarshalUnmarshaler], [vtype.FlagCallback],
//     [vtype.FlagFunc], [vtype.FlagBoolFunc]
//
// Names can include multiple long and multiple short values. Each value should
// be separated by a pipe (|) character. If val has a usable non-zero value, it
// will be used as the default value for that flag option. Functions compatible
// with [vtype.FlagFunc] and [vtype.FlagBoolFunc] will be converted.
func (fs *FlagSet) Flag(val any, names, desc string) *Flag {
	flag := newFlag(val, names, desc)
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
