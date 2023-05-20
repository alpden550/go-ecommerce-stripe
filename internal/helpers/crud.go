package helpers

import (
	"go-ecommerce/internal/configs"
	"go-ecommerce/internal/models"
)

func SaveCustomer(app configs.AppConfiger, firstName, lastName, email string) (int, error) {
	db := app.GetDB()
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := db.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SaveTransaction(app configs.AppConfiger, transaction models.Transaction) (int, error) {
	db := app.GetDB()
	id, err := db.InsertTransaction(transaction)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SaveOrder(app configs.AppConfiger, order models.Order) (int, error) {
	db := app.GetDB()
	id, err := db.InsertOrder(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}
