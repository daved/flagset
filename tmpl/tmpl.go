package tmpl

import (
	"bytes"
	"fmt"
	"text/template"
)

type Tmpl struct {
	Text string
	FMap template.FuncMap
	Data any
}

func New(text string, fMap template.FuncMap, data any) *Tmpl {
	return &Tmpl{
		Data: data,
		Text: text,
		FMap: fMap,
	}
}

func (t *Tmpl) String() string {
	tmpl := template.New("flagset").Funcs(t.FMap)

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(t.Text)
	if err != nil {
		fmt.Fprintf(buf, "%v\n", err)
		return buf.String()
	}

	if err := tmpl.Execute(buf, t.Data); err != nil {
		fmt.Fprintf(buf, "%v\n", err)
		return buf.String()
	}

	return buf.String()
}
