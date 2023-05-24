package main

import (
	handlers "github.com/alpden550/go-ecommerce-stripe/internal/handlers_api"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"net/http"
)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := handlers.AuthenticateToken(request)
		if err != nil {
			helpers.InvalidCredentials(api, writer)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
