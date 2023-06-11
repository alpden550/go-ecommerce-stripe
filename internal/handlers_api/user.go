package handlers_api

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"net/http"
)

func AllUsers(writer http.ResponseWriter, request *http.Request) {
	users, err := helpers.FetchAllUsers(api)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, users)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
	}
}
