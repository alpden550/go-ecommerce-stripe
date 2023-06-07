package handlers_api

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"net/http"
)

func AllSales(writer http.ResponseWriter, request *http.Request) {
	sales, err := helpers.FetchAllWidgetOrder(api)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, sales)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}
}

func AllSubscriptions(writer http.ResponseWriter, request *http.Request) {
	sales, err := helpers.FetchAllSubscriptionsOrder(api)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, sales)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}
}
