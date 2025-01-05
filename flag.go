package flagset

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/daved/flagset/vtype"
)

// Flag manages flag option data. Exported fields are able to be set for use in
// usage output templating (overriding existing behavior).
type Flag struct {
	longs  []string
	shorts []string
	desc   string

	TypeName    string
	DefaultText string
	HideUsage   bool
	Meta        map[string]any
}

func newFlag(val any, longs, shorts []string, desc string) *Flag {
	return &Flag{
		longs:       longs,
		shorts:      shorts,
		desc:        desc,
		TypeName:    typeName(val),
		DefaultText: defaultText(val),
		Meta:        map[string]any{},
	}
}

// Longs returns all long flag names.
func (f Flag) Longs() []string {
	return f.longs
}

// Shorts returns all short flag names.
func (f Flag) Shorts() []string {
	return f.shorts
}

// Description returns the description string.
func (f Flag) Description() string {
	return f.desc
}

func typeName(val any) string {
	var out string

	switch v := val.(type) {
	case vtype.FlagCallback:
		if v.IsBool() {
			out = "bool"
		} else {
			out = "value"
		}

	case vtype.TextMarshalUnmarshaler, flag.Value:
		out = "value"

	default:
		rv := reflect.ValueOf(val)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		out = rv.Type().Name()
	}

	return out
}

func defaultText(val any) string {
	var out string

	switch v := val.(type) {
	case vtype.TextMarshalUnmarshaler:
		t, err := v.MarshalText()
		if err != nil {
			return err.Error()
		}
		out = string(t)
	case vtype.FlagCallback:
		out = ""
	case fmt.Stringer:
		out = v.String()
	default:
		vo := reflect.ValueOf(val)
		if vo.Kind() == reflect.Ptr {
			vo = vo.Elem()
		}
		out = fmt.Sprint(vo)
	}

	return out
}
