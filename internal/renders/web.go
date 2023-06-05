package renders

import (
	"embed"
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"html/template"
	"net/http"
)

var app *configs.AppConfig

var functions = template.FuncMap{
	"formatAmount": formatAmount,
}

//go:embed templates
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
	} else {
		td.IsAuthenticated = 0
	}

	return td
}

func RenderTemplate(writer http.ResponseWriter, request *http.Request, page string, td *TemplateData, partials ...string) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	_, templateInMap := app.TemplateCache[templateToRender]

	if app.Config.Env == "production" && templateInMap {
		t = app.TemplateCache[templateToRender]
	} else {
		t, err = parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
			return err
		}
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

func parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	// build partials
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/partials/%s.partial.gohtml", x)
		}
	}

	partialFiles := []string{
		"templates/base.layout.gohtml",
		templateToRender,
	}
	partialFiles = append(partialFiles, partials...)

	t, err = template.
		New(fmt.Sprintf("%s.page.gohtml", page)).
		Funcs(functions).
		ParseFS(templateFS, partialFiles...)

	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return nil, err
	}

	app.TemplateCache[templateToRender] = t
	return t, nil
}
