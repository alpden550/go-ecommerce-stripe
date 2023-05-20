package helpers

import (
	"go-ecommerce/internal/configs"
	"go-ecommerce/internal/models"
)

func SaveCustomer(app *configs.AppConfig, firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SaveTransaction(app *configs.AppConfig, transaction models.Transaction) (int, error) {
	id, err := app.DB.InsertTransaction(transaction)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SaveOrder(app *configs.AppConfig, order models.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}
