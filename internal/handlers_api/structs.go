package handlers_api

import "go-ecommerce/internal/configs"

var api *configs.ApiApplication

func SetApiAppToHandlers(a *configs.ApiApplication) {
	api = a
}

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	ProductID     string `json:"product_id"`
	PlanID        string `json:"plan_id"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	LastFour      string `json:"last_four"`
	ExpireMonth   int    `json:"expire_month"`
	ExpireYear    int    `json:"expire_year"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}
