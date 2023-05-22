package helpers

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
)

func FetchUserByEmail(app configs.AppConfiger, email string) (models.User, error) {
	var user models.User
	db := app.GetDB()
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}
