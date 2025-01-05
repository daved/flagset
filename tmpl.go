package flagset

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// TmplData is the structure used for usage output templating. Custom template
// string values should be based on this type.
type TmplData struct {
	FlagSet *FlagSet
}

// TmplConfig tracks the template string and function map used for usage output
// templating.
type TmplConfig struct {
	Text string
	FMap template.FuncMap
}

// NewDefaultTmplConfig returns the default TmplConfig value. This can be used
// as an example of how to setup custom usage output templating.
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
{{end}}{{else}}{{- end}}
`)

	return &TmplConfig{
		Text: tmplText,
		FMap: tmplFMap,
	}
}

func executeTmpl(tc *TmplConfig, data any) string {
	tmpl := template.New("flagset").Funcs(tc.FMap)

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(tc.Text)
	if err != nil {
		fmt.Fprintf(buf, "%v\n", err)
		return buf.String()
	}

	if err := tmpl.Execute(buf, data); err != nil {
		fmt.Fprintf(buf, "%v\n", err)
		return buf.String()
	}

	return buf.String()
}
