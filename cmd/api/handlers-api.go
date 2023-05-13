package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-ecommerce/internal/cards"
	"go-ecommerce/internal/models"
	"net/http"
	"strconv"
)

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	ProductID     string `json:"product_id"`
	PlanID        string `json:"plan_id"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	LastFour      string `json:"last_four"`
	ExpireMonth   int    `json:"expire_month"`
	ExpireYear    int    `json:"expire_year"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true
	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(pi, "", "	")
		if err != nil {
			app.errorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
		}
		out, err := json.MarshalIndent(j, "", "	")
		if err != nil {
			app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
	}

}

func (app *application) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
	}

	out, err := json.MarshalIndent(widget, "", "	")
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}

func (app *application) Subscribe(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	okay := true
	transactionMessage := "Transaction successful"

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	stripeCustomer, msg, err := card.CreateCustomer(payload.PaymentMethod, payload.Email)
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		okay = false
		transactionMessage = msg
	}

	if okay {
		_, err = card.SubscribeToPlan(stripeCustomer, payload.PlanID, payload.Email, payload.LastFour, "")
		if err != nil {
			app.errorLog.Printf("%e", fmt.Errorf("%w", err))
			okay = false
			transactionMessage = "Error subscribing a customer"
		}
	}

	if okay {
		productID, _ := strconv.Atoi(payload.ProductID)
		customerID, err := app.SaveCustomer(payload.FirstName, payload.LastName, payload.Email)
		if err != nil {
			app.errorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}

		amount, _ := strconv.Atoi(payload.Amount)
		transaction := models.Transaction{
			Amount:              amount,
			Currency:            "usd",
			LastFour:            payload.LastFour,
			ExpireMonth:         payload.ExpireMonth,
			ExpireYear:          payload.ExpireYear,
			TransactionStatusID: 2,
		}
		transactionID, err := app.SaveTransaction(transaction)
		if err != nil {
			app.errorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}

		order := models.Order{
			WidgetID:      productID,
			TransactionID: transactionID,
			CustomerID:    customerID,
			StatusID:      1,
			Quantity:      1,
			Amount:        amount,
		}
		_, err = app.SaveOrder(order)
		if err != nil {
			app.errorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}
	}

	response := jsonResponse{
		OK:      okay,
		Message: transactionMessage,
	}

	out, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}
	if okay {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}
