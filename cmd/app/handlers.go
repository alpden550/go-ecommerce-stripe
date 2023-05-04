package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(
		w, r, "terminal", &templateData{}, "stripe-js", "nav",
	); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PaymentSucceed(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	form := r.Form

	paymentData := map[string]interface{}{
		"cardholder": form.Get("cardholder_name"),
		"email":      form.Get("cardholder_email"),
		"intent":     form.Get("payment_intent"),
		"method":     form.Get("payment_method"),
		"amount":     form.Get("payment_amount"),
		"currency":   form.Get("payment_currency"),
	}

	if err := app.renderTemplate(w, r, "succeeded", &templateData{Data: paymentData}, "nav"); err != nil {
		app.errorLog.Printf("%e", err)
		return
	}
}

func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}
	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	data := map[string]interface{}{
		"widget": widget,
	}

	if err := app.renderTemplate(
		w, r, "buy-once", &templateData{Data: data}, "stripe-js", "nav",
	); err != nil {
		app.errorLog.Printf("%e", err)
	}
}
