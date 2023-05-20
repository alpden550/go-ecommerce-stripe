package helpers

import (
	"errors"
	"go-ecommerce/internal/configs"
	"go-ecommerce/internal/models"
)

func SaveCustomer(a interface{}, firstName, lastName, email string) (int, error) {
	var id int
	var err error
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	switch app := a.(type) {
	case *configs.AppConfig:
		id, err = app.DB.InsertCustomer(customer)
	case *configs.ApiApplication:
		id, err = app.DB.InsertCustomer(customer)
	default:
		return 0, errors.New("invalid app pr api config")
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SaveTransaction(a interface{}, transaction models.Transaction) (int, error) {
	var id int
	var err error

	switch app := a.(type) {
	case *configs.AppConfig:
		id, err = app.DB.InsertTransaction(transaction)
	case *configs.ApiApplication:
		id, err = app.DB.InsertTransaction(transaction)
	default:
		return 0, errors.New("invalid app pr api config")
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SaveOrder(a interface{}, order models.Order) (int, error) {
	var id int
	var err error

	switch app := a.(type) {
	case *configs.AppConfig:
		id, err = app.DB.InsertOrder(order)
	case *configs.ApiApplication:
		id, err = app.DB.InsertOrder(order)
	default:
		return 0, errors.New("invalid app pr api config")
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
