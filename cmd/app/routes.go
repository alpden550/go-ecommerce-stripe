package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", app.IndexPage)

	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/virtual-terminal-payment-succeeded", app.VirtualTerminalPaymentSucceed)
	mux.Get("/virtual-terminal-receipt", app.VirtualTerminalShowReceipt)

	mux.Get("/widgets/{id}", app.ChargeOnce)
	mux.Post("/payment-succeeded", app.PaymentSucceed)
	mux.Get("/receipt", app.ShowReceipt)

	mux.Get("/plans/bronze", app.BronzePlan)
	mux.Get("/receipt/bronze", app.BronzePlanShowReceipt)

	return mux
}
