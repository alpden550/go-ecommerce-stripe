package renders

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
)

var app *configs.AppConfig

var functions = template.FuncMap{
	"formatAmount": formatAmount,
}

//go:embed templates/*
var templateFS embed.FS

func SetAppToRender(a *configs.AppConfig) {
	app = a
}

func AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	widgets, _ := helpers.FetchDbWidgets(app)
	widgetsData := map[string]interface{}{
		"widgets": widgets,
	}

	td.API = app.Config.Api
	td.StripeSecret = app.Config.Stripe.Secret
	td.StripePublishKey = app.Config.Stripe.Key
	td.Widgets = widgetsData

	if app.Session.Exists(r.Context(), "userId") {
		td.IsAuthenticated = 1
		td.UserID = app.Session.GetInt(r.Context(), "userId")
	} else {
		td.IsAuthenticated = 0
	}

	return td
}

func RenderTemplate(
	writer http.ResponseWriter,
	request *http.Request,
	pagePath, pageName string,
	td *TemplateData,
	partials ...string,
) error {
	var t *template.Template
	var partialFiles []string
	var err error
	templateToRender := fmt.Sprintf("templates/%s", pagePath)

	if len(partials) > 0 {
		for _, x := range partials {
			partialFiles = append(partialFiles, fmt.Sprintf("templates/partials/%s.partial.gohtml", x))
		}
	}

	templateFiles := []string{
		"templates/base.layout.gohtml",
		templateToRender,
	}
	templateFiles = append(templateFiles, partialFiles...)

	_, templateInMap := app.TemplateCache[templateToRender]

	if templateInMap {
		t = app.TemplateCache[templateToRender]
	} else {
		t = template.Must(template.New(pageName).
			Funcs(functions).
			ParseFS(templateFS, templateFiles...),
		)
		app.TemplateCache[templateToRender] = t
	}

	if td == nil {
		td = &TemplateData{}
	}

	td = AddDefaultData(td, request)

	err = t.Execute(writer, td)
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return err
	}

	return nil
}
