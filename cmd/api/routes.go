package main

import (
	handlers "github.com/alpden550/go-ecommerce-stripe/internal/handlers_api"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/api/auth/login", handlers.CreateAuthToken)
	mux.Post("/api/auth/is_authenticated", handlers.CheckAuthentication)
	mux.Post("/api/auth/forgot-password", handlers.SendPasswordResetEmail)
	mux.Post("/api/auth/reset-password", handlers.ResetPassword)

	mux.Get("/api/widgets/{id}", handlers.GetWidgetByID)
	mux.Post("/api/payment-intent", handlers.GetPaymentIntent)
	mux.Post("/api/subscribe", handlers.Subscribe)

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(MiddlewareAuth)

		mux.Post("/virtual-terminal-payment-succeeded", handlers.VirtualTerminalPaymentSucceeded)
		mux.Get("/all-sales", handlers.AllSales)
		mux.Get("/all-subscriptions", handlers.AllSubscriptions)
	})

	return mux
}
