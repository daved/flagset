package flagset

import (
	"strings"
	"unicode/utf8"

	"github.com/daved/vtype"
)

// Flag manages flag option data. Exported fields are for easy post-construction
// configuration.
type Flag struct {
	// Templating
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
	val = vtype.ConvertCompatible(val)
	longs, shorts := longsAndShorts(names)

	return &Flag{
		val:         val,
		longs:       longs,
		shorts:      shorts,
		desc:        desc,
		TypeName:    vtype.ValueTypeName(val),
		DefaultText: vtype.DefaultValueText(val),
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
