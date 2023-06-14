package handlers_api

import (
	"errors"
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/encryption"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/mailer"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/alpden550/go-ecommerce-stripe/internal/urlsigner"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func CreateAuthToken(writer http.ResponseWriter, request *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := helpers.ReadJSON(writer, request, &userInput); err != nil {
		_ = helpers.BadRequest(writer, request, err)
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
		_ = helpers.BadRequest(writer, request, err)
		return
	}
	err = helpers.SaveToken(api, token, &user)
	if err != nil {
		_ = helpers.BadRequest(writer, request, err)
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
	err = helpers.WriteJSON(writer, http.StatusOK, payload)
	if err != nil {
		return
	}
}

func CheckAuthentication(writer http.ResponseWriter, request *http.Request) {
	user, err := AuthenticateToken(request)
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
	err = helpers.WriteJSON(writer, http.StatusOK, payload)
}

func AuthenticateToken(request *http.Request) (*models.User, error) {
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

func SendPasswordResetEmail(writer http.ResponseWriter, request *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	if err := helpers.ReadJSON(writer, request, &payload); err != nil {
		_ = helpers.BadRequest(writer, request, err)
		return
	}

	if _, err := helpers.FetchUserByEmail(api, payload.Email); err != nil {
		response := jsonResponse{
			OK:      false,
			Message: "no matching email found",
		}
		_ = helpers.WriteJSON(writer, http.StatusAccepted, response)
		return
	}

	link := fmt.Sprintf("%s/auth/reset-password?email=%s", api.Config.FrontEnd, payload.Email)
	signer := urlsigner.Signer{Secret: []byte(api.Config.SecretKey)}
	signedLink := signer.GenerateTokenFromString(link)

	var data struct {
		Link string
	}
	data.Link = signedLink

	if err := mailer.SendEmail(
		api,
		api.Config.SMTP.EmailFrom,
		payload.Email,
		"Password Reset Requested",
		"password-reset.plain.tmpl",
		"password-reset.html.tmpl",
		data,
	); err != nil {
		api.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		_ = helpers.BadRequest(writer, request, err)
	}

	response := jsonResponse{
		OK:      true,
		Message: "sent",
	}
	_ = helpers.WriteJSON(writer, http.StatusOK, response)
}

func ResetPassword(writer http.ResponseWriter, request *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := helpers.ReadJSON(writer, request, &payload); err != nil {
		_ = helpers.BadRequest(writer, request, err)
		return
	}

	encryptor := encryption.Encryption{Key: []byte(api.Config.SecretKey)}
	realEmail, err := encryptor.Decrypt(payload.Email)
	if err != nil {
		_ = helpers.BadRequest(writer, request, err)
		return
	}

	user, err := helpers.FetchUserByEmail(api, realEmail)
	if err != nil {
		_ = helpers.BadRequest(writer, request, err)
		return
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		_ = helpers.BadRequest(writer, request, err)
		return
	}

	if err = helpers.UpdateUserPassword(api, &user, string(newHashedPassword)); err != nil {
		_ = helpers.BadRequest(writer, request, err)
		return
	}

	response := jsonResponse{
		OK:      true,
		Message: "Password was changed",
	}
	_ = helpers.WriteJSON(writer, http.StatusCreated, response)
}
