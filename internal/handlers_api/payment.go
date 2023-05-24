package handlers_api

import (
	"encoding/json"
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/cards"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"net/http"
	"strconv"
)

func GetPaymentIntent(writer http.ResponseWriter, request *http.Request) {
	var payload stripePayload
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	card := cards.Card{
		Secret:   api.Config.Stripe.Secret,
		Key:      api.Config.Stripe.Key,
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
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(out)
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
		}
		out, err := json.MarshalIndent(j, "", "	")
		if err != nil {
			api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(out)
	}

}

func VirtualTerminalPaymentSucceeded(writer http.ResponseWriter, request *http.Request) {
	var paymentData struct {
		PaymentAmount   int    `json:"amount"`
		PaymentCurrency string `json:"currency"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		PaymentIntent   string `json:"payment_intent"`
		PaymentMethod   string `json:"payment_method"`
		BankReturnCode  string `json:"bank_return_code"`
		ExpiredMonth    int    `json:"expired_month"`
		ExpiredYear     int    `json:"expired_year"`
		LastFour        string `json:"last_four"`
	}

	if err := helpers.ReadJSON(api, writer, request, &paymentData); err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	card := cards.Card{
		Secret: api.Config.Stripe.Secret,
		Key:    api.Config.Stripe.Key,
	}
	pi, err := card.GetPaymentIntent(paymentData.PaymentIntent)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	pm, err := card.GetPaymentMethod(paymentData.PaymentMethod)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	paymentData.LastFour = pm.Card.Last4
	paymentData.ExpiredMonth = int(pm.Card.ExpMonth)
	paymentData.ExpiredYear = int(pm.Card.ExpYear)

	transaction := models.Transaction{
		Amount:              paymentData.PaymentAmount,
		Currency:            paymentData.PaymentCurrency,
		LastFour:            paymentData.LastFour,
		ExpireMonth:         paymentData.ExpiredMonth,
		ExpireYear:          paymentData.ExpiredYear,
		BankReturnCode:      pi.LatestCharge.ID,
		TransactionStatusID: 2,
		PaymentMethodCode:   paymentData.PaymentMethod,
		PaymentIntentCode:   paymentData.PaymentIntent,
	}

	_, err = helpers.SaveTransaction(api, transaction)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	helpers.WriteJSON(api, writer, http.StatusOK, paymentData)
}
