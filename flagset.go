// Package flagset provides simple, POSIX-friendly flag parsing. For example,
// multi-character single-hyphen flags like "-abc" are exploded as "-a -b -c".
package flagset

import (
	"slices"
	"unicode/utf8"

	er "github.com/daved/flagset/fserrs"
)

// FlagSet contains flag options and usage-related values. Exported fields are
// used for easy post-construction configuration.
type FlagSet struct {
	// Fields used for templating:
	Tmpl             *Tmpl // set to NewUsageTmpl by default
	HideTypeHints    bool
	HideDefaultHints bool
	Meta             map[string]any

	name   string
	flags  []*Flag
	parsed []string
	ops    []string
}

// New constructs a FlagSet. Package convention is to name the flagset after the
// command that the options are associated with.
func New(name string) *FlagSet {
	fs := &FlagSet{
		name: name,
		Meta: map[string]any{},
	}

	fs.Tmpl = NewUsageTmpl(fs)

	return fs
}

// Flags returns all flag options that have been set.
func (fs *FlagSet) Flags() []*Flag {
	return fs.flags
}

// Lookup returns the first encountered flag that matches the provided name.
func (fs *FlagSet) Lookup(name string) *Flag {
	return lookupFlag(fs.flags, name)
}

// Parsed returns the args provided to Parse with single-hyphen flags containing
// multiple characters exploded out.
func (fs *FlagSet) Parsed() []string {
	return fs.parsed
}

// Operand returns the i'th operand. For example, Operand(0) is the first
// argument after flags are parsed. Non-existant indexes return an empty string.
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

// Parse processes flags and flag values from the argument list, which must not
// include the initial command name. Parse must be called after all flags in the
// FlagSet are defined and before flag value access. Before parsing occurs,
// multi-character single-hyphen flags like "-abc" are exploded as "-a -b -c".
func (fs *FlagSet) Parse(args []string) error {
	fs.parsed = explodeShortArgs(args)

	ops, err := resolveFlags(fs.flags, fs.parsed)
	if err != nil {
		return er.NewError(er.NewParseError(err))
	}

	fs.ops = ops
	return nil
}

// Flag adds a flag option to the FlagSet. See [vtype.Hydrate] for details about
// which value types are supported. Names can include multiple long and multiple
// short values, each separated by a pipe (|) character. If val has a usable
// non-zero value, it will be used as the flag's default value. Functions
// compatible with [vtype] typed functions will be auto-converted.
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

// Usage returns usage text. The default template construction function
// ([NewUsageTmpl]) can be used as a reference for custom templates which should
// be used to set the "Tmpl" field on FlagSet.
func (fs *FlagSet) Usage() string {
	return fs.Tmpl.String()
}

func lookupFlag(flags []*Flag, name string) *Flag {
	for _, flag := range flags {
		ss := flag.shorts
		if utf8.RuneCountInString(name) > 1 {
			ss = flag.longs
		}
		if slices.Contains(ss, name) {
			return flag
		}
	}
	return nil
}
