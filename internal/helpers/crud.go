package helpers

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
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

func SaveWidgetOrder(app configs.AppConfiger, order models.Order) (int, error) {
	db := app.GetDB()
	id, err := db.InsertOrder(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SaveSubscriptionOrder(app configs.AppConfiger, order models.Order) (int, error) {
	db := app.GetDB()
	id, err := db.InsertSubscriptionOrder(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func FetchDbWidgets(app configs.AppConfiger) ([]*models.Widget, error) {
	db := app.GetDB()
	widgets, err := db.GetAllWidgets()
	if err != nil {
		return nil, err
	}

	return widgets, nil
}

func SaveToken(app configs.AppConfiger, token *models.Token, user *models.User) error {
	db := app.GetDB()
	if err := db.InsertToken(token, user); err != nil {
		return err
	}

	return nil
}

func FetchAllWidgetOrder(app configs.AppConfiger) ([]*models.Order, error) {
	db := app.GetDB()
	orders, err := db.GetWidgetOrders()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func FetchAllSubscriptionsOrder(app configs.AppConfiger) ([]*models.Order, error) {
	db := app.GetDB()
	orders, err := db.GetSubscriptionsOrders()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetWidgetOrder(app configs.AppConfiger, id int) (*models.Order, error) {
	db := app.GetDB()
	order, err := db.GetWidgetOrderByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func GetSubscriptionOrder(app configs.AppConfiger, id int) (*models.Order, error) {
	db := app.GetDB()
	order, err := db.GetSubscriptionOrderByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}
