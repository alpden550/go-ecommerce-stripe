package handlers

import (
	"fmt"
	"net/http"

	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
)

var app *configs.AppConfig

func SetAppToHandlers(a *configs.AppConfig) {
	app = a
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer,
		request,
		"index.page.gohtml",
		"index.page.gohtml",
		&renders.TemplateData{},
		"nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
