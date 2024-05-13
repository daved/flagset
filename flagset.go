package flagset

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"
)

type Opt struct {
	Longs  []string
	Shorts []string
	Help   string
	Type   string
	Init   string
}

type FlagSet struct {
	fs   *flag.FlagSet
	opts map[string]Opt // placing arg as key
}

func (fs *FlagSet) Collected() map[string]Opt {
	return fs.opts
}

func New(name string) *FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	return &FlagSet{
		fs:   fs,
		opts: make(map[string]Opt),
	}
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
	return fs.fs.Parse(arguments)
}

func (fs *FlagSet) Visit(fn func(*flag.Flag)) {
	fs.fs.Visit(fn)
}

func (fs *FlagSet) VisitAll(fn func(*flag.Flag)) {
	fs.fs.VisitAll(fn)
}

func (fs *FlagSet) Opt(val any, flagNames, help, placing string) {
	longs, shorts := longsAndShorts(flagNames)
	v := reflect.ValueOf(val).Elem()

	fs.opts[flagNames] = Opt{
		Longs:  longs,
		Shorts: shorts,
		Help:   help,
		Type:   v.Type().Name(),
		Init:   fmt.Sprintf("%v", v),
	}

	for _, long := range longs {
		addOptTo(fs.fs, val, long, help)
	}
	for _, short := range shorts {
		addOptTo(fs.fs, val, short, help)
	}
}

func addOptTo(fs *flag.FlagSet, val any, flagName, help string) {
	switch v := val.(type) {
	case *string:
		fs.StringVar(v, flagName, *v, help)
	case *bool:
		fs.BoolVar(v, flagName, *v, help)
	}
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

func (fs *FlagSet) Help() string {
	buf := &bytes.Buffer{}
	fs.fs.SetOutput(buf)
	defer fs.fs.SetOutput(io.Discard)
	fs.fs.Usage()
	return buf.String()
}
