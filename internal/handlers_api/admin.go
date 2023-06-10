package handlers_api

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var orderPayload struct {
	PageSize    int `json:"page_size"`
	CurrentPage int `json:"page"`
}

var orderResponse struct {
	CurrentPage int             `json:"page"`
	LastPage    int             `json:"last_page"`
	PageSize    int             `json:"page_size"`
	TotalOrders int             `json:"total_orders"`
	Orders      []*models.Order `json:"orders"`
}

func AllSales(writer http.ResponseWriter, request *http.Request) {
	if err := helpers.ReadJSON(api, writer, request, &orderPayload); err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	orders, lastPage, totalOrders, err := helpers.FetchAllWidgetOrder(api, orderPayload.PageSize, orderPayload.CurrentPage)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	orderResponse.CurrentPage = orderPayload.CurrentPage
	orderResponse.LastPage = lastPage
	orderResponse.PageSize = orderPayload.PageSize
	orderResponse.TotalOrders = totalOrders
	orderResponse.Orders = orders

	err = helpers.WriteJSON(api, writer, http.StatusOK, orderResponse)
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
	if err := helpers.ReadJSON(api, writer, request, &orderPayload); err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}
	orders, lastPage, totalOrders, err := helpers.FetchAllSubscriptionsOrder(api, orderPayload.PageSize, orderPayload.CurrentPage)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	orderResponse.CurrentPage = orderPayload.CurrentPage
	orderResponse.LastPage = lastPage
	orderResponse.PageSize = orderPayload.PageSize
	orderResponse.TotalOrders = totalOrders
	orderResponse.Orders = orders

	err = helpers.WriteJSON(api, writer, http.StatusOK, orderResponse)
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
