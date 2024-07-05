package flagset

import "strings"

var (
	MetaKeyType        = "Type"
	MetaKeyDefault     = "Default"
	MetaKeySkipUsage   = "SkipUsage"
	MetaKeyTypeHint    = "TypeHint"
	MetaKeyDefaultHint = "DefaultHint"

	defaultPrefix = "default: "
)

type metaOpts struct {
	HideTypeHint    bool
	HideDefaultHint bool
	Type            string
	Default         string
}

func makeMeta(opts metaOpts) map[string]any {
	m := map[string]any{
		MetaKeyType:    opts.Type,
		MetaKeyDefault: opts.Default,
	}

	if !opts.HideTypeHint {
		tHintPre, tHintPost := "=", ""
		if opts.Type == "bool" {
			tHintPre, tHintPost = "[=", "]"
		}
		m[MetaKeyTypeHint] = tHintPre + strings.ToUpper(opts.Type) + tHintPost
	}

	if !opts.HideDefaultHint {
		var dHint string
		if opts.Default != "" {
			dHint = defaultPrefix + opts.Default
		}
		m[MetaKeyDefaultHint] = dHint
	}

	return m
}
