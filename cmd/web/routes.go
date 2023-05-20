package main

import (
	"github.com/go-chi/chi/v5"
	handlers2 "go-ecommerce/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", handlers2.IndexPage)

	mux.Get("/virtual-terminal", handlers2.VirtualTerminal)
	mux.Post("/virtual-terminal/payment-succeeded", handlers2.VirtualTerminalPaymentSucceed)
	mux.Get("/virtual-terminal/receipt", handlers2.VirtualTerminalShowReceipt)

	mux.Get("/widgets/{id}", handlers2.WidgetChargeOnce)
	mux.Post("/widgets/payment-succeeded", handlers2.WidgetPaymentSucceed)
	mux.Get("/widgets/receipt", handlers2.WidgetShowReceipt)

	mux.Get("/plans/bronze", handlers2.BronzePlan)
	mux.Get("/plans/bronze/receipt", handlers2.BronzePlanShowReceipt)

	return mux
}
