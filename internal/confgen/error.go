package confgen

import (
	"bytes"
	"html/template"
)

const tmplName = "test"

type IllegalFieldType struct {
	FieldType string
	Line      int
}

const ErrorTmpl_IllegalFieldType = "{{.FieldType}} is illegal at {{.Line}}"

func (e IllegalFieldType) Error() string {
	tmpl, err := template.New(tmplName).Parse(ErrorTmpl_IllegalFieldType)
	if err != nil {
		panic(err)
	}
	var errmsg bytes.Buffer
	if err := tmpl.Execute(&errmsg, e); err != nil {
		panic(err)
	}
	return errmsg.String()
}
