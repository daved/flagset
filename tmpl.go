package flagset

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Tmpl holds template configuration details.
type Tmpl struct {
	Text string
	FMap template.FuncMap
	Data any
}

// Execute parses the template text and funcmap, then executes it using the set
// data.
func (t *Tmpl) Execute() (string, error) {
	tmpl := template.New("clic").Funcs(t.FMap)

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(t.Text)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(buf, t.Data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// String calls the Execute method returning either the validly executed
// template output or error message text.
func (t *Tmpl) String() string {
	s, err := t.Execute()
	if err != nil {
		s = fmt.Sprintf("%v\n", err)
	}
	return s
}

// NewUsageTmpl returns the default template configuration. This can be used as
// an example of how to setup custom usage output templating.
func NewUsageTmpl(fs *FlagSet) *Tmpl {
	type tmplData struct {
		FlagSet *FlagSet
	}

	data := &tmplData{
		FlagSet: fs,
	}

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

	fMap := template.FuncMap{
		"Join":        strings.Join,
		"TypeHint":    typeHintFn,
		"DefaultHint": defaultHintFn,
	}

	text := strings.TrimSpace(`
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

	return &Tmpl{text, fMap, data}
}
