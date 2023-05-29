package main

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", handlers.IndexPage)

	mux.Get("/auth/login", handlers.Login)
	mux.Post("/auth/login", handlers.SubmitLogin)
	mux.Get("/auth/logout", handlers.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(MiddlewareAuth)

		mux.Get("/virtual-terminal", handlers.VirtualTerminal)
	})

	mux.Get("/widgets/{id}", handlers.WidgetChargeOnce)
	mux.Post("/widgets/payment-succeeded", handlers.WidgetPaymentSucceed)
	mux.Get("/widgets/receipt", handlers.WidgetShowReceipt)

	mux.Get("/plans/bronze", handlers.BronzePlan)
	mux.Get("/plans/bronze/receipt", handlers.BronzePlanShowReceipt)

	return mux
}
