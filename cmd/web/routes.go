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

	mux.Route("/auth", func(mux chi.Router) {
		mux.Get("/login", handlers.Login)
		mux.Post("/login", handlers.SubmitLogin)
		mux.Get("/logout", handlers.Logout)
		mux.Get("/forgot-password", handlers.ForgotPassword)
		mux.Get("/reset-password", handlers.ShowResetPassword)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(MiddlewareAuth)
		mux.Get("/virtual-terminal", handlers.VirtualTerminal)
		mux.Get("/all-sales", handlers.AllSales)
		mux.Get("/all-sales/{id}", handlers.ShowSale)
		mux.Get("/all-subscriptions", handlers.AllSubscriptions)
		mux.Get("/all-subscriptions/{id}", handlers.ShowSubscription)
	})

	mux.Get("/widgets/{id}", handlers.WidgetChargeOnce)
	mux.Post("/widgets/payment-succeeded", handlers.WidgetPaymentSucceed)
	mux.Get("/widgets/receipt", handlers.WidgetShowReceipt)

	mux.Get("/plans/bronze", handlers.BronzePlan)
	mux.Get("/plans/bronze/receipt", handlers.BronzePlanShowReceipt)

	return mux
}
