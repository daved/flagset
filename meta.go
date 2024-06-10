package flagset

import "strings"

type metaFab struct {
	HideTypeHint    bool
	HideDefaultHint bool
}

func (f metaFab) make(typ, defalt string) map[string]any {
	m := map[string]any{
		"Type":    typ,
		"Default": defalt,
	}

	if !f.HideTypeHint {
		tHintPre, tHintPost := "=", ""
		if typ == "bool" {
			tHintPre, tHintPost = "[=", "]"
		}
		m["TypeHint"] = tHintPre + strings.ToUpper(typ) + tHintPost
	}

	if !f.HideDefaultHint {
		var dHint string
		if defalt != "" {
			dHint = "default: " + defalt
		}
		m["DefaultHint"] = dHint
	}

	return m
}
