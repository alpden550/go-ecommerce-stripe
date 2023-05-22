package handlers_api

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"net/http"
	"time"
)

func CreateAuthToken(writer http.ResponseWriter, request *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := helpers.ReadJSON(api, writer, request, &userInput); err != nil {
		_ = helpers.BadRequest(api, writer, request, err)
		return
	}

	user, err := helpers.FetchUserByEmail(api, userInput.Email)
	if err != nil {
		_ = helpers.InvalidCredentials(api, writer)
		return
	}
	_, err = helpers.PasswordMatcher(api, user.Password, userInput.Password)
	if err != nil {
		_ = helpers.InvalidCredentials(api, writer)
		return
	}

	token, err := models.GenerateNewToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		_ = helpers.BadRequest(api, writer, request, err)
		return
	}
	err = helpers.SaveToken(api, token, &user)
	if err != nil {
		_ = helpers.BadRequest(api, writer, request, err)
		return
	}

	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Token   *models.Token `json:"token"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("token for %s created", user.Email)
	payload.Token = token
	err = helpers.WriteJSON(api, writer, http.StatusOK, payload)
	if err != nil {
		return
	}
}
