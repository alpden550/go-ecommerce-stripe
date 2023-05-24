package handlers_api

import (
	"errors"
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"net/http"
	"strings"
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

func CheckAuthentication(writer http.ResponseWriter, request *http.Request) {
	user, err := authenticateToken(request)
	if err != nil {
		_ = helpers.InvalidCredentials(api, writer)
		return
	}
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("authenticated user %s", user.Email)
	err = helpers.WriteJSON(api, writer, http.StatusOK, payload)
}

func authenticateToken(request *http.Request) (*models.User, error) {
	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header found")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header found")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("token has wrong size")
	}

	user, err := helpers.FetchUserByToken(api, token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}
