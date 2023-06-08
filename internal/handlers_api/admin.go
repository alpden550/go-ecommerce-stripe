package handlers_api

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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

func GetSale(writer http.ResponseWriter, request *http.Request) {
	orderId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	order, err := helpers.GetWidgetOrder(api, orderId)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, order)
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

func GetSubscription(writer http.ResponseWriter, request *http.Request) {
	orderId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	order, err := helpers.GetSubscriptionOrder(api, orderId)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, order)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}
}
