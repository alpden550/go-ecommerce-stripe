package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/encryption"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"github.com/alpden550/go-ecommerce-stripe/internal/urlsigner"
	"net/http"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer,
		request,
		"auth/login.page.gohtml",
		"login.page.gohtml",
		&renders.TemplateData{},
		"nav",
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
		writer,
		request,
		"auth/forgot-password.page.gohtml",
		"forgot-password.page.gohtml",
		&renders.TemplateData{},
		"nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func ShowResetPassword(writer http.ResponseWriter, request *http.Request) {
	url := request.RequestURI
	testUrl := fmt.Sprintf("%s%s", app.Config.FrontEnd, url)
	email := request.URL.Query().Get("email")

	signer := urlsigner.Signer{Secret: []byte(app.Config.SecretKey)}
	verified := signer.VerifyToken(testUrl)
	if !verified {
		app.ErrorLog.Println("invalid hash for email ", email)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}

	if expired := signer.Expired(testUrl, 60); expired {
		app.ErrorLog.Println("link expired for email ", email)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}

	encryptor := encryption.Encryption{Key: []byte(app.Config.SecretKey)}
	encrEmail, err := encryptor.Encrypt(email)
	if err != nil {
		app.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}

	data := map[string]interface{}{
		"email": encrEmail,
	}
	if err := renders.RenderTemplate(
		writer,
		request,
		"auth/reset-password.page.gohtml",
		"reset-password.page.gohtml",
		&renders.TemplateData{Data: data},
		"nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
