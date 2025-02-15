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
	typeHintFn := func(f *Flag) string {
		if f.TypeName == "" {
			return ""
		}

		_, isBool := boolValRaw(f.val)

		pre, post := "=", ""
		if len(f.Longs()) > 0 && isBool {
			pre, post = "[=", "]"
		}

		return pre + strings.ToUpper(f.TypeName) + post
	}

	defaultHintFn := func(f *Flag) string {
		if f.DefaultText == "" {
			return ""
		}

		return "default: " + f.DefaultText
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
  {{- if $flag.TypeName}}  {{TypeHint $flag}}{{end}}
  {{- if $flag.DefaultText}}    {{DefaultHint $flag}}{{end}}
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
