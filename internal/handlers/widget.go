package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func WidgetChargeOnce(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	widgetID, err := strconv.Atoi(id)
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}
	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	data := map[string]interface{}{
		"widget": widget,
	}

	if err := renders.RenderTemplate(
		writer, request, "buy-once", &renders.TemplateData{Data: data}, "stripe-js", "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func WidgetPaymentSucceed(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	form := request.Form
	widgetID, _ := strconv.Atoi(form.Get("product_id"))
	transactionData, err := helpers.GetTransactionData(app, request)
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	customerID, err := helpers.SaveCustomer(app, transactionData.FirstName, transactionData.LastName, transactionData.Email)
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
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

	transactionID, err := helpers.SaveTransaction(app, transaction)
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	order := models.Order{
		WidgetID:      widgetID,
		TransactionID: transactionID,
		CustomerID:    customerID,
		StatusID:      1,
		Quantity:      1,
		Amount:        transactionData.Amount,
	}
	_, err = helpers.SaveOrder(app, order)
	if err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	app.Session.Put(request.Context(), "receipt", transactionData)
	http.Redirect(writer, request, "/widgets/receipt", http.StatusSeeOther)
}

func WidgetShowReceipt(writer http.ResponseWriter, request *http.Request) {
	payment, ok := app.Session.Get(request.Context(), "receipt").(helpers.TransactionData)
	if !ok {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
	paymentData := map[string]interface{}{
		"paymentData": payment,
	}
	app.Session.Remove(request.Context(), "receipt")
	if err := renders.RenderTemplate(
		writer, request, "receipt", &renders.TemplateData{Data: paymentData}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}
}
