package helpers

import (
	"errors"
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func PasswordMatcher(app configs.AppConfiger, hashed, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, errors.New("password doesn't match with user")
		default:
			return false, err
		}
	}

	return true, nil
}
