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

func FetchDbWidgets(app configs.AppConfiger) ([]models.Widget, error) {
	db := app.GetDB()
	widgets, err := db.GetAllWidgets()
	if err != nil {
		return nil, err
	}

	return widgets, nil
}

func FetchUserByEmail(app configs.AppConfiger, email string) (models.User, error) {
	var user models.User
	db := app.GetDB()
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}
