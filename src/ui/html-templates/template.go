package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func htmlTemplate(t string, vars interface{}) string {
	tmpl := template.New("template")
	tmpl, err := tmpl.Parse(t)
	if err != nil {
		return fmt.Sprintf("Error %v for tmpl.Parse(t)", err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, vars)
	if err != nil {
		return fmt.Sprintf("Error %v for tmpl.Execute()", err)
	}

	return b.String()
}

