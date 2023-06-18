package handlers_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alpden550/go-ecommerce-stripe/internal/cards"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/stripe/stripe-go/v74"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	okay := true
	var subscription *stripe.Subscription
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
		subscription, err = card.SubscribeToPlan(stripeCustomer, payload.PlanID, payload.Email, payload.LastFour, "")
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
			SubscriptionCode:    subscription.ID,
			PaymentMethodCode:   payload.PaymentMethod,
		}
		transactionID, err := helpers.SaveTransaction(api, transaction)
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}

		order := models.Order{
			SubscriptionID: productID,
			TransactionID:  transactionID,
			CustomerID:     customerID,
			StatusID:       1,
			Quantity:       1,
			Amount:         amount,
		}
		orderID, err := helpers.SaveSubscriptionOrder(api, order)
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}

		invoice := models.Invoice{
			ID:        orderID,
			Quantity:  order.Quantity,
			Amount:    order.Amount,
			Product:   "Subscription",
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			CreatedAt: time.Now(),
		}

		err = invoice.SendInvoice()
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
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

func CancelSubscription(writer http.ResponseWriter, request *http.Request) {
	var subscription struct {
		ID               int    `json:"id"`
		SubscriptionCode string `json:"sc"`
		Currency         string `json:"currency"`
	}

	if err := helpers.ReadJSON(writer, request, &subscription); err != nil {
		helpers.BadRequest(writer, request, err)
		return
	}

	card := cards.Card{
		Secret:   api.Config.Stripe.Secret,
		Key:      api.Config.Stripe.Key,
		Currency: subscription.Currency,
	}

	err := card.CancelSubscription(subscription.SubscriptionCode)
	if err != nil {
		helpers.BadRequest(writer, request, err)
		return
	}

	err = helpers.UpdateOrderStatus(api, subscription.ID, 3)
	if err != nil {
		helpers.BadRequest(writer, request, errors.New("the subscription was cancelled, but the the database was not updated"))
		return
	}

	response := jsonResponse{
		OK:      true,
		Message: "Subscription was cancelled",
	}

	err = helpers.WriteJSON(writer, http.StatusOK, response)
	if err != nil {
		return
	}

}
