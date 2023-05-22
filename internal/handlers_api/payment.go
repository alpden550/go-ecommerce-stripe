package handlers_api

import (
	"encoding/json"
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/cards"
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
