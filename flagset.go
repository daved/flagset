package flagset

import (
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"
)

type Opt struct {
	Names   string
	Longs   []string
	Shorts  []string
	Type    string
	Default string
	Usage   string
	Meta    map[string]any
}

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
	return fs.fs.Parse(fs.parsed)
}

func (fs *FlagSet) Visit(fn func(*flag.Flag)) {
	fs.fs.Visit(fn)
}

func (fs *FlagSet) VisitAll(fn func(*flag.Flag)) {
	fs.fs.VisitAll(fn)
}

func (fs *FlagSet) Opt(val any, names, usage string, metas ...map[string]any) {
	longs, shorts := longsAndShorts(names)
	v := reflect.ValueOf(val).Elem()

	t := v.Type().Name()
	def := fmt.Sprintf("%v", v)
	m := conMeta{fs.HideTypeHint, fs.HideDefaultHint}.make(t, def)

	for _, meta := range metas {
		for k, v := range meta {
			m[k] = v
		}
	}

	opt := Opt{
		Names:   names,
		Longs:   longs,
		Shorts:  shorts,
		Type:    t,
		Default: def,
		Usage:   usage,
		Meta:    m,
	}
	fs.opts = append(fs.opts, opt)

	for _, long := range longs {
		addOptTo(fs.fs, val, long, usage)
	}

	for _, short := range shorts {
		addOptTo(fs.fs, val, short, usage)
	}
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

func addOptTo(fs *flag.FlagSet, val any, flagName, usage string) {
	switch v := val.(type) {
	case *string:
		fs.StringVar(v, flagName, *v, usage)
	case *bool:
		fs.BoolVar(v, flagName, *v, usage)
	}
}
