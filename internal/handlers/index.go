package handlers

import (
	"fmt"
	"go-ecommerce/internal/configs"
	"go-ecommerce/internal/renders"
	"net/http"
)

var app *configs.AppConfig

func SetAppToHandlers(a *configs.AppConfig) {
	app = a
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "index", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
