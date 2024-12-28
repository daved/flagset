package flagset

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/daved/flagset/vtypes"
)

// Flag manages flag option data. The fields are mostly unexposed to emphasize
// that the fields should not be modified unexpectedly. It is appropriate to
// modify the Meta map to communicate info/behavior to the usage template.
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
	case vtypes.FlagCallback:
		if v.IsBool() {
			out = "bool"
		} else {
			out = "value"
		}

	case vtypes.TextMarshalUnmarshaler, flag.Value:
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
	case vtypes.TextMarshalUnmarshaler:
		t, err := v.MarshalText()
		if err != nil {
			return err.Error()
		}
		out = string(t)
	case vtypes.FlagCallback:
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
