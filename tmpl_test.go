package learning

import (
	"os"
	"testing"
	"text/template"
)

var tmplStr = `{{if .}}Hello {{.}}{{else}}Hello there{{end}}
`

func TestTmpl(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(tmplStr))
	err := tmpl.Execute(os.Stdout, "Oyen")
	if err != nil {
		t.Error(err)
	}
}

var tmplRangeStr = `{{range $idx, $name := .}}{{.}}
{{end}}`

func TestTmplRange(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(tmplRangeStr))
	err := tmpl.Execute(os.Stdout, []string{"Evanbill", "Antonio", "Kore"})
	if err != nil {
		t.Error(err)
	}
}
