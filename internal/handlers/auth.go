package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"net/http"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "login", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
