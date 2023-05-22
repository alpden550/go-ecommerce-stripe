package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"net/http"
)

func BronzePlan(writer http.ResponseWriter, request *http.Request) {
	sbcr, err := app.DB.GetSubscriptionByName("Bronze Plan")
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}
	data := map[string]interface{}{
		"subscription": sbcr,
	}
	if err := renders.RenderTemplate(
		writer, request, "bronze-plan", &renders.TemplateData{Data: data}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func BronzePlanShowReceipt(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "bronze-plan-receipt", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
