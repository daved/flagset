package flagset

import (
	"flag"
	"fmt"
	"reflect"
)

// Flag manages flag option data. The fields are mostly unexposed to emphasize
// that the fields should not be modified unexpectedly. It is appropriate to
// modify the Meta map to communicate info/behavior to the usage template.
type Flag struct {
	names  string
	longs  []string
	shorts []string
	typ    string
	defalt string
	usage  string
	Meta   map[string]any
}

func makeFlag(fs *FlagSet, ns string, ls, ss []string, t, d, u string) Flag {
	m := makeMeta(metaOpts{
		HideTypeHint:    fs.MetaHideTypeHints,
		HideDefaultHint: fs.MetaHideDefaultHints,
		Type:            t,
		Default:         d,
	})

	return Flag{
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
func (f Flag) Names() string {
	return f.names
}

// Longs returns all long flag names.
func (f Flag) Longs() []string {
	return f.longs
}

// Shorts returns all short flag names.
func (f Flag) Shorts() []string {
	return f.shorts
}

// Type returns the flag value type name.
func (f Flag) Type() string {
	return f.typ
}

// Default returns the flag default value.
func (f Flag) Default() string {
	return f.defalt
}

// Usage returns the usage string.
func (f Flag) Usage() string {
	return f.usage
}

func typeName(val any) string {
	switch val.(type) {
	case FlagFunc, FlagBoolFunc, TextMarshalUnmarshaler, flag.Value:
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
	case FlagFunc, FlagBoolFunc:
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
