package flagset

import (
	"flag"
	"fmt"
	"reflect"
)

// Opt manages flag option data. The fields are mostly unexposed to emphasize
// that the fields should not be modified unexpectedly. It is appropriate to
// modify the Meta map to communicate info/behavior to the usage template.
type Opt struct {
	names  string
	longs  []string
	shorts []string
	typ    string
	defalt string
	usage  string
	Meta   map[string]any
}

func makeOpt(fs *FlagSet, ns string, ls, ss []string, t, d, u string) Opt {
	m := makeMeta(metaOpts{
		HideTypeHint:    fs.HideTypeHint,
		HideDefaultHint: fs.HideDefaultHint,
		Type:            t,
		Default:         d,
	})

	return Opt{
		names:  ns,
		longs:  ls,
		shorts: ss,
		typ:    t,
		defalt: d,
		usage:  u,
		Meta:   m,
	}
}

// Names returns a string of the defined flag names.
func (o Opt) Names() string {
	return o.names
}

// Longs returns all long flag names.
func (o Opt) Longs() []string {
	return o.longs
}

// Shorts returns all short flag names.
func (o Opt) Shorts() []string {
	return o.shorts
}

// Type returns the flag value type name.
func (o Opt) Type() string {
	return o.typ
}

// Default returns the flag default value.
func (o Opt) Default() string {
	return o.defalt
}

// Usage returns the usage string.
func (o Opt) Usage() string {
	return o.usage
}

func typeName(val any) string {
	switch val.(type) {
	case OptFunc, TextMarshalUnmarshaler, flag.Value:
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
	case TextMarshalUnmarshaler:
		t, err := v.MarshalText()
		if err != nil {
			t = []byte(err.Error())
		}
		return string(t)
	case OptFunc:
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
