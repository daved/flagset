package flagset

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type tmplData struct {
	Name string
	Opts []Opt
}

var tmplText = strings.TrimSpace(`
{{- if .Opts -}}
Flags for {{.Name}}:
{{range $i, $opt := .Opts}}
  {{- if $opt.Meta.SkipUsage}}{{continue}}{{end}}
  {{if .}}  {{end}}{{if $opt.Shorts}}-{{Join $opt.Shorts ", -"}}{{end}}
  {{- if and $opt.Shorts $opt.Longs}}, {{end}}
  {{- if $opt.Longs}}--{{Join $opt.Longs ", --"}}{{end}}
  {{- if $opt.Meta.TypeHint}}  {{end}}{{$opt.Meta.TypeHint}}
  {{- if $opt.Meta.DefaultHint}}    {{$opt.Meta.DefaultHint}}{{end}}
	{{$opt.Usage}}
{{end}}
{{else}}
{{- end}}
`)

func (fs *FlagSet) Usage() string {
	data := &tmplData{
		Name: fs.Name(),
		Opts: fs.Opts(),
	}

	tmpl := template.New("flagset").Funcs(
		template.FuncMap{
			"Join": strings.Join,
		},
	)

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(tmplText)
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
