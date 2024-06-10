package flagset

type Opt struct {
	names  string
	longs  []string
	shorts []string
	typ    string
	defalt string
	usage  string
	meta   map[string]any
}

func makeOpt(fs *FlagSet, ns string, ls, ss []string, t, d, u string) Opt {
	m := metaFab{fs.HideTypeHint, fs.HideDefaultHint}.make(t, d)

	return Opt{
		names:  ns,
		longs:  ls,
		shorts: ss,
		typ:    t,
		defalt: d,
		usage:  u,
		meta:   m,
	}
}

func (o Opt) Names() string {
	return o.names
}

func (o Opt) Longs() []string {
	return o.longs
}

func (o Opt) Shorts() []string {
	return o.shorts
}

func (o Opt) Type() string {
	return o.typ
}

func (o Opt) Default() string {
	return o.defalt
}

func (o Opt) Usage() string {
	return o.usage
}

func (o Opt) Meta() map[string]any {
	return o.meta
}
