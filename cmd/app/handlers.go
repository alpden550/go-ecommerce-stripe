package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-ecommerce/internal/models"
	"net/http"
	"strconv"
)

func (app *application) IndexPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(
		w, r, "index", &templateData{}, "nav",
	); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(
		w, r, "terminal", &templateData{}, "stripe-js", "nav",
	); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ShowReceipt(w http.ResponseWriter, r *http.Request) {
	payment := app.Session.Get(r.Context(), "receipt").(TransactionData)
	paymentData := map[string]interface{}{
		"paymentData": payment,
	}
	app.Session.Remove(r.Context(), "receipt")
	if err := app.renderTemplate(w, r, "receipt", &templateData{Data: paymentData}, "nav"); err != nil {
		app.errorLog.Printf("%e", err)
		return
	}
}

func (app *application) VirtualTerminalShowReceipt(w http.ResponseWriter, r *http.Request) {
	payment := app.Session.Get(r.Context(), "virtual-terminal-receipt").(TransactionData)
	paymentData := map[string]interface{}{
		"paymentData": payment,
	}
	app.Session.Remove(r.Context(), "receipt")
	if err := app.renderTemplate(w, r, "receipt", &templateData{Data: paymentData}, "nav"); err != nil {
		app.errorLog.Printf("%e", err)
		return
	}
}

func (app *application) PaymentSucceed(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	form := r.Form
	widgetID, _ := strconv.Atoi(form.Get("product_id"))
	transactionData, err := app.GetTransactionData(r)
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	customerID, err := app.SaveCustomer(transactionData.FirstName, transactionData.LastName, transactionData.Email)
	if err != nil {
		app.errorLog.Printf("%e", err)
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

	transactionID, err := app.SaveTransaction(transaction)
	if err != nil {
		app.errorLog.Printf("%e", err)
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
	_, err = app.SaveOrder(order)
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	app.Session.Put(r.Context(), "receipt", transactionData)
	http.Redirect(w, r, "/receipt", http.StatusSeeOther)
}

func (app *application) VirtualTerminalPaymentSucceed(w http.ResponseWriter, r *http.Request) {
	transactionData, err := app.GetTransactionData(r)
	if err != nil {
		app.errorLog.Printf("%e", err)
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

	_, err = app.SaveTransaction(transaction)
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	app.Session.Put(r.Context(), "virtual-terminal-receipt", transactionData)
	http.Redirect(w, r, "/virtual-terminal-receipt", http.StatusSeeOther)
}

func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *application) SaveTransaction(transaction models.Transaction) (int, error) {
	id, err := app.DB.InsertTransaction(transaction)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *application) SaveOrder(order models.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}
	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Printf("%e", err)
		return
	}

	data := map[string]interface{}{
		"widget": widget,
	}

	if err := app.renderTemplate(
		w, r, "buy-once", &templateData{Data: data}, "stripe-js", "nav",
	); err != nil {
		app.errorLog.Printf("%e", err)
	}
}

func (app *application) BronzePlan(w http.ResponseWriter, r *http.Request) {
	sbcr, err := app.DB.GetSubscriptionByName("Bronze Plan")
	if err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}
	data := map[string]interface{}{
		"subscription": sbcr,
	}
	if err := app.renderTemplate(w, r, "bronze-plan", &templateData{Data: data}, "nav"); err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func (app *application) BronzePlanShowReceipt(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "bronze-plan-receipt", &templateData{}, "nav"); err != nil {
		app.errorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
