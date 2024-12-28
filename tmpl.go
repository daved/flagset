package flagset

import (
	"strings"
	"text/template"
)

type TmplData struct {
	FlagSet *FlagSet
}

type TmplConfig struct {
	Text string
	FMap template.FuncMap
}

func NewDefaultTmplConfig() *TmplConfig {
	typeHintFn := func(t string) string {
		if t == "" {
			return ""
		}

		pre, post := "=", ""
		if t == "bool" {
			pre, post = "[=", "]"
		}

		return pre + strings.ToUpper(t) + post
	}

	defaultHintFn := func(d string) string {
		if d == "" {
			return ""
		}

		return "default: " + d
	}

	tmplFMap := template.FuncMap{
		"Join":        strings.Join,
		"TypeHint":    typeHintFn,
		"DefaultHint": defaultHintFn,
	}

	tmplText := strings.TrimSpace(`
{{- if .FlagSet.Flags -}}
Flags for {{.FlagSet.Name}}:
{{range $i, $flag := .FlagSet.Flags}}
  {{- if $flag.HideUsage}}{{continue}}{{end}}
  {{if .}}  {{end}}{{if $flag.Shorts}}-{{Join $flag.Shorts ", -"}}{{end}}
  {{- if and $flag.Shorts $flag.Longs}}, {{end}}
  {{- if $flag.Longs}}--{{Join $flag.Longs ", --"}}{{end}}
  {{- if $flag.TypeName}}  {{TypeHint $flag.TypeName}}{{end}}
  {{- if $flag.DefaultText}}    {{DefaultHint $flag.DefaultText}}{{end}}
        {{$flag.Description}}
{{end}}
{{else}}
{{- end}}
`)

	return &TmplConfig{
		Text: tmplText,
		FMap: tmplFMap,
	}
}
