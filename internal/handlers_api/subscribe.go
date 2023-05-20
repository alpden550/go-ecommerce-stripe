package handlers_api

import (
	"encoding/json"
	"fmt"
	"go-ecommerce/internal/cards"
	"go-ecommerce/internal/helpers"
	"go-ecommerce/internal/models"
	"net/http"
	"strconv"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	okay := true
	transactionMessage := "Transaction successful"

	card := cards.Card{
		Secret: api.Config.Stripe.Secret,
		Key:    api.Config.Stripe.Key,
	}

	stripeCustomer, msg, err := card.CreateCustomer(payload.PaymentMethod, payload.Email)
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		okay = false
		transactionMessage = msg
	}

	if okay {
		_, err = card.SubscribeToPlan(stripeCustomer, payload.PlanID, payload.Email, payload.LastFour, "")
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
			okay = false
			transactionMessage = "Error subscribing a customer"
		}
	}

	if okay {
		productID, _ := strconv.Atoi(payload.ProductID)
		customerID, err := helpers.SaveCustomer(api, payload.FirstName, payload.LastName, payload.Email)
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
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
		transactionID, err := helpers.SaveTransaction(api, transaction)
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
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
		_, err = helpers.SaveOrder(api, order)
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}
	}

	response := jsonResponse{
		OK:      okay,
		Message: transactionMessage,
	}

	out, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
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
