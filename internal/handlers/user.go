package handlers

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"net/http"
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
