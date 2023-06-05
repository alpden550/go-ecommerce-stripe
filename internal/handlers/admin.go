package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"net/http"
)

func AllSales(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "all-sales", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func AllSubscriptions(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "all-subscriptions", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
