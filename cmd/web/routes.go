package main

import (
	"github.com/go-chi/chi/v5"
	"go-ecommerce/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", handlers.IndexPage)

	mux.Get("/virtual-terminal", handlers.VirtualTerminal)
	mux.Post("/virtual-terminal/payment-succeeded", handlers.VirtualTerminalPaymentSucceed)
	mux.Get("/virtual-terminal/receipt", handlers.VirtualTerminalShowReceipt)

	mux.Get("/widgets/{id}", handlers.WidgetChargeOnce)
	mux.Post("/widgets/payment-succeeded", handlers.WidgetPaymentSucceed)
	mux.Get("/widgets/receipt", handlers.WidgetShowReceipt)

	mux.Get("/plans/bronze", handlers.BronzePlan)
	mux.Get("/plans/bronze/receipt", handlers.BronzePlanShowReceipt)

	return mux
}
