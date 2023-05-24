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

	mux.Get("/api/widgets/{id}", handlers.GetWidgetByID)
	mux.Post("/api/payment-intent", handlers.GetPaymentIntent)
	mux.Post("/api/subscribe", handlers.Subscribe)

	return mux
}
