package main

import (
	"net/http"
)

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !app.Session.Exists(request.Context(), "userId") {
			http.Redirect(writer, request, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
