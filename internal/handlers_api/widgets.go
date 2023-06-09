package handlers_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/cards"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func GetWidgetByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := api.DB.GetWidget(widgetID)
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}

	out, err := json.MarshalIndent(widget, "", "	")
	if err != nil {
		api.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(out)
}

func RefundWidget(writer http.ResponseWriter, request *http.Request) {
	var chargeToRefund struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"`
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
	}

	if err := helpers.ReadJSON(api, writer, request, &chargeToRefund); err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	card := cards.Card{
		Secret:   api.Config.Stripe.Secret,
		Key:      api.Config.Stripe.Key,
		Currency: chargeToRefund.Currency,
	}

	if err := card.Refund(chargeToRefund.PaymentIntent, chargeToRefund.Amount); err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err := helpers.UpdateOrderStatus(api, chargeToRefund.ID, 2)
	if err != nil {
		helpers.BadRequest(api, writer, request, errors.New("the charge was refunded, but the the database was not updated"))
		return
	}

	response := jsonResponse{
		OK:      true,
		Message: "Charge refunded",
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, response)
	if err != nil {
		return
	}

}
