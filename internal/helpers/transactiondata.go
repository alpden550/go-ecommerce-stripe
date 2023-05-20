package helpers

import (
	"go-ecommerce/internal/cards"
	"go-ecommerce/internal/configs"
	"net/http"
	"strconv"
)

type TransactionData struct {
	FirstName         string
	LastName          string
	Email             string
	PaymentIntentCode string
	PaymentMethodCode string
	BankReturnCode    string
	Amount            int
	Currency          string
	LastFour          string
	ExpireMonth       int
	ExpireYear        int
}

func GetTransactionData(app *configs.AppConfig, request *http.Request) (TransactionData, error) {
	var transactionData TransactionData
	err := request.ParseForm()
	if err != nil {
		app.ErrorLog.Printf("%e", err)
		return transactionData, err
	}

	form := request.Form
	amount, _ := strconv.Atoi(form.Get("payment_amount"))
	card := cards.Card{
		Secret: app.Config.Stripe.Secret,
		Key:    app.Config.Stripe.Key,
	}

	pi, err := card.GetPaymentIntent(form.Get("payment_intent"))
	if err != nil {
		app.ErrorLog.Printf("%e", err)
		return transactionData, err
	}
	pm, err := card.GetPaymentMethod(form.Get("payment_method"))
	if err != nil {
		app.ErrorLog.Printf("%e", err)
		return transactionData, err
	}

	transactionData = TransactionData{
		FirstName:         form.Get("first_name"),
		LastName:          form.Get("last_name"),
		Email:             form.Get("email"),
		PaymentIntentCode: pi.LatestCharge.ID,
		PaymentMethodCode: form.Get("payment_method"),
		BankReturnCode:    form.Get("payment_intent"),
		Amount:            amount,
		Currency:          form.Get("payment_currency"),
		LastFour:          pm.Card.Last4,
		ExpireMonth:       int(pm.Card.ExpMonth),
		ExpireYear:        int(pm.Card.ExpYear),
	}

	return transactionData, nil
}