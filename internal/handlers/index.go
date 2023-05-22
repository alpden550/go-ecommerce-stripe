package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
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
