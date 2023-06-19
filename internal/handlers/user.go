package handlers

import (
	"fmt"
	"net/http"

	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
)

func ShowAllUsers(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer,
		request,
		"user/all-users.page.gohtml",
		"all-users.page.gohtml",
		&renders.TemplateData{},
		"nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}

func ShowUser(writer http.ResponseWriter, request *http.Request) {
	if err := renders.RenderTemplate(
		writer,
		request,
		"user/user.page.gohtml",
		"user.page.gohtml",
		&renders.TemplateData{},
		"nav",
	); err != nil {
		app.ErrorLog.Printf("%e", fmt.Errorf("%w", err))
	}
}
