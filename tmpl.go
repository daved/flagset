package flagset

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type tmplData struct {
	Name  string
	Flags []*Flag
}

var tmplText = strings.TrimSpace(`
{{- if .Flags -}}
Flags for {{.Name}}:
{{range $i, $flag := .Flags}}
  {{- if $flag.HideUsage}}{{continue}}{{end}}
  {{if .}}  {{end}}{{if $flag.Shorts}}-{{Join $flag.Shorts ", -"}}{{end}}
  {{- if and $flag.Shorts $flag.Longs}}, {{end}}
  {{- if $flag.Longs}}--{{Join $flag.Longs ", --"}}{{end}}
  {{- if $flag.TypeHint}}  {{$flag.TypeHint}}{{end}}
  {{- if $flag.DefaultHint}}    {{$flag.DefaultHint}}{{end}}
        {{$flag.Usage}}
{{end}}
{{else}}
{{- end}}
`)

// SetUsageTemplate allows callers to override the base template text.
func (fs *FlagSet) SetUsageTemplate(txt string) {
	fs.tmplTxt = txt
}

func (fs *FlagSet) UpdateUsageFuncMap(m template.FuncMap) {
	for k, v := range m {
		fs.tmplFuncMap[k] = v
	}
}

// Usage returns the parsed usage template. Each Flag type's Meta field is
// leveraged to convey detailed info/behavior. This method and related template
// can be used as an example for callers to wrap the FlagSet type and design
// their own usage output. For example, grouping, sorting, etc.
func (fs *FlagSet) Usage() string {
	data := &tmplData{
		Name:  fs.Name(),
		Flags: fs.Flags(),
	}

	tmpl := template.New("flagset").Funcs(fs.tmplFuncMap)

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(fs.tmplTxt)
	if err != nil {
		fmt.Fprintf(buf, "flagset: template error: %v\n", err)
		return buf.String()
	}

	if err := tmpl.Execute(buf, data); err != nil {
		fmt.Fprintf(buf, "flagset: template error: %v\n", err)
		return buf.String()
	}

	return buf.String()
}
