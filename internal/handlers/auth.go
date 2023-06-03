package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"net/http"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "login", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func SubmitLogin(writer http.ResponseWriter, request *http.Request) {
	if err := app.Session.RenewToken(request.Context()); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	if err := request.ParseForm(); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	email, password := request.Form.Get("email"), request.Form.Get("password")
	id, err := app.DB.Authenticate(email, password)
	if err != nil {
		http.Redirect(writer, request, "/auth/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(request.Context(), "userId", id)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	app.Session.Destroy(request.Context())
	app.Session.RenewToken(request.Context())

	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func ForgotPassword(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer, request, "forgot-password", &renders.TemplateData{}, "nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
