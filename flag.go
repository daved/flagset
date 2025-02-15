package flagset

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/daved/flagset/vtype"
)

// Flag manages flag option data. The exported fields are for easy
// post-construction configuration.
type Flag struct {
	// Fields used for templating:
	TypeName    string // is derived from the val type when possible
	DefaultText string // is derived from the val value when possible
	HideUsage   bool
	Meta        map[string]any

	val    any
	longs  []string
	shorts []string
	desc   string
}

func newFlag(val any, names string, desc string) *Flag {
	switch v := val.(type) {
	case func(string) error:
		val = vtype.FlagFunc(v)
	case func(bool) error:
		val = vtype.FlagBoolFunc(v)
	default:
		vo := reflect.ValueOf(val)
		if vo.Kind() == reflect.Ptr {
			vo = vo.Elem()
		}

		switch vo.Kind() {
		case reflect.Slice:
			vts := vtype.MakeSlice(val)
			val = &vts
		}
	}

	longs, shorts := longsAndShorts(names)

	return &Flag{
		val:         val,
		longs:       longs,
		shorts:      shorts,
		desc:        desc,
		TypeName:    typeName(val),
		DefaultText: defaultText(val),
		Meta:        map[string]any{},
	}
}

// Longs returns all long flag names.
func (f *Flag) Longs() []string {
	return f.longs
}

// Shorts returns all short flag names.
func (f *Flag) Shorts() []string {
	return f.shorts
}

// Description returns the description string.
func (f *Flag) Description() string {
	return f.desc
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

func typeName(val any) string {
	switch v := val.(type) {
	case interface{ UsageTypeName() string }:
		return v.UsageTypeName()

	case vtype.FlagCallback:
		if v.IsBool() {
			return "bool"
		}
		return "value"

	case vtype.TextMarshalUnmarshaler, flag.Value:
		return "value"

	case error:
		return ""

	default:
		rv := reflect.ValueOf(val)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		return rv.Type().Name()
	}
}

func defaultText(val any) string {
	switch v := val.(type) {
	case vtype.TextMarshalUnmarshaler:
		t, err := v.MarshalText()
		if err != nil {
			return err.Error()
		}
		return string(t)

	case vtype.FlagCallback, error:
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
