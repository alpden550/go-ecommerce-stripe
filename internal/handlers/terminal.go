package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"net/http"
)

func VirtualTerminal(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "terminal", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func VirtualTerminalPaymentSucceed(writer http.ResponseWriter, request *http.Request) {
	transactionData, err := helpers.GetTransactionData(app, request)
	if err != nil {
		app.ErrorLog.Printf("%e", err)
		return
	}

	transaction := models.Transaction{
		Amount:              transactionData.Amount,
		Currency:            transactionData.Currency,
		LastFour:            transactionData.LastFour,
		ExpireMonth:         transactionData.ExpireMonth,
		ExpireYear:          transactionData.ExpireYear,
		BankReturnCode:      transactionData.BankReturnCode,
		PaymentMethodCode:   transactionData.PaymentMethodCode,
		PaymentIntentCode:   transactionData.PaymentIntentCode,
		TransactionStatusID: 2,
	}

	_, err = helpers.SaveTransaction(app, transaction)
	if err != nil {
		app.ErrorLog.Printf("%e", err)
		return
	}

	app.Session.Put(request.Context(), "virtual-terminal-receipt", transactionData)
	http.Redirect(writer, request, "/virtual-terminal/receipt", http.StatusSeeOther)
}

func VirtualTerminalShowReceipt(writer http.ResponseWriter, request *http.Request) {
	payment, ok := app.Session.Get(request.Context(), "virtual-terminal-receipt").(helpers.TransactionData)
	if !ok {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
	paymentData := map[string]interface{}{
		"paymentData": payment,
	}
	app.Session.Remove(request.Context(), "virtual-terminal-receipt")
	if err := renders.RenderTemplate(
		writer, request, "receipt", &renders.TemplateData{Data: paymentData}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", err)
		return
	}
}
