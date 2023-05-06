package main

import (
	"github.com/go-chi/chi/v5"
	"go-ecommerce/internal/cards"
	"net/http"
	"strconv"
)

func (app *application) IndexPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(
		w, r, "index", &templateData{}, "nav",
	); err != nil {
		app.errorLog.Println(err)
	}
}

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
	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	pi, err := card.GetPaymentIntent(form.Get("payment_intent"))
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}
	pm, err := card.GetPaymentMethod(form.Get("payment_method"))
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	paymentData := map[string]interface{}{
		"cardholder":       form.Get("cardholder_name"),
		"email":            form.Get("cardholder_email"),
		"intent":           form.Get("payment_intent"),
		"method":           form.Get("payment_method"),
		"amount":           form.Get("payment_amount"),
		"currency":         form.Get("payment_currency"),
		"last_four":        pm.Card.Last4,
		"expire_month":     pm.Card.ExpMonth,
		"expire_year":      pm.Card.ExpYear,
		"latest_charge_id": pi.LatestCharge.ID,
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
