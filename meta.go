package flagset

import "strings"

type conMeta struct {
	HideTypeHint    bool
	HideDefaultHint bool
}

func (con conMeta) make(typ, defalt string) map[string]any {
	m := map[string]any{
		"Type":    typ,
		"Default": defalt,
	}

	if !con.HideTypeHint {
		tHintPre, tHintPost := "=", ""
		if typ == "bool" {
			tHintPre, tHintPost = "[=", "]"
		}
		m["TypeHint"] = tHintPre + strings.ToUpper(typ) + tHintPost
	}

	if !con.HideDefaultHint {
		var dHint string
		if defalt != "" {
			dHint = "default: " + defalt
		}
		m["DefaultHint"] = dHint
	}

	return m
}
