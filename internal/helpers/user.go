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

func FetchUserByToken(app configs.AppConfiger, token string) (*models.User, error) {
	db := app.GetDB()
	user, err := db.GetUserForToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUserPassword(app configs.AppConfiger, user *models.User, hash string) error {
	db := app.GetDB()
	if err := db.SetUserPassword(user, hash); err != nil {
		return err
	}

	return nil
}

func FetchAllUsers(app configs.AppConfiger) ([]*models.User, error) {
	db := app.GetDB()
	users, err := db.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func FetchUserByID(app configs.AppConfiger, id int) (*models.User, error) {
	db := app.GetDB()
	user, err := db.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func EditUser(app configs.AppConfiger, user *models.User) error {
	db := app.GetDB()
	if err := db.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func SaveUser(app configs.AppConfiger, user *models.User, hash string) (int, error) {
	db := app.GetDB()
	id, err := db.AddUser(user, hash)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func RemoveUser(app configs.AppConfiger, id int) error {
	db := app.GetDB()
	err := db.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
