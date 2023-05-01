package main

import (
	"net/http"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	stripeData := map[string]interface{}{
		"key": app.config.stripe.key,
	}
	if err := app.renderTemplate(w, r, "terminal", &templateData{Data: stripeData}); err != nil {
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

	if err := app.renderTemplate(w, r, "succeeded", &templateData{Data: paymentData}); err != nil {
		app.errorLog.Printf("%e", err)
		return
	}
}
