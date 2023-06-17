package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*
var emailTemplateFS embed.FS

func renderTemplate(tmpl string, data interface{}) (string, error) {
	templateToRender := fmt.Sprintf("templates/%s", tmpl)
	t, err := template.New("email").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		return "", err
	}
	formattedMessage := tpl.String()
	return formattedMessage, nil
}
