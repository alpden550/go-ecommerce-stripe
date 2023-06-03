package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates
var emailTemplateFS embed.FS

func renderHTML(tmpl string, data interface{}) (string, error) {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)
	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
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

func renderPlainText(tmpl string, data interface{}) (string, error) {
	templateToRender := fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
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
