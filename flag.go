package flagset

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

// Flag manages flag option data. The fields are mostly unexposed to emphasize
// that the fields should not be modified unexpectedly. It is appropriate to
// modify the Meta map to communicate info/behavior to the usage template.
type Flag struct {
	longs  []string
	shorts []string
	desc   string

	TypeHint    string
	DefaultHint string
	HideUsage   bool
	Meta        map[string]any
}

func newFlag(val any, longs, shorts []string, desc string) *Flag {
	return &Flag{
		longs:       longs,
		shorts:      shorts,
		desc:        desc,
		TypeHint:    typeName(val),
		DefaultHint: defaultText(val),
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

	switch val.(type) {
	case FlagFunc, FlagBoolFunc, TextMarshalUnmarshaler, flag.Value:
		out = "value"
	default:
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		out = v.Type().Name()
	}

	if out != "" {
		pre, post := "=", ""

		switch val.(type) {
		case *bool, FlagBoolFunc:
			pre, post = "[=", "]"
		}

		out = pre + strings.ToUpper(out) + post
	}

	return out
}

const defaultPrefix = "default: "

func defaultText(val any) string {
	var out string

	switch v := val.(type) {
	case TextMarshalUnmarshaler:
		t, err := v.MarshalText()
		if err != nil {
			return err.Error()
		}
		out = string(t)
	case FlagFunc, FlagBoolFunc:
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

	if out != "" {
		out = defaultPrefix + out
	}

	return out
}
