package handlers_api

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func AllUsers(writer http.ResponseWriter, request *http.Request) {
	users, err := helpers.FetchAllUsers(api)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, users)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
	}
}

func OneUser(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	user, err := helpers.FetchUserByID(api, userID)
	if err != nil {
		helpers.WriteJSON(api, writer, http.StatusNotFound, err)
		return
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, user)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
	}
}

func EditUser(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	var user models.User
	err = helpers.ReadJSON(api, writer, request, &user)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	response := jsonResponse{
		OK: true,
	}

	if userID > 0 {
		err = helpers.EditUser(api, &user)
		if err != nil {
			helpers.BadRequest(api, writer, request, err)
			return
		}

		if user.Password != "" {
			newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			if err != nil {
				helpers.BadRequest(api, writer, request, err)
				return
			}

			err = helpers.UpdateUserPassword(api, &user, string(newHash))
			if err != nil {
				helpers.BadRequest(api, writer, request, err)
				return
			}
		}
		response.Message = fmt.Sprintf("Updated user %d", userID)
	} else {
		newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			helpers.BadRequest(api, writer, request, err)
			return
		}
		newUserId, err := helpers.SaveUser(api, &user, string(newHash))
		if err != nil {
			helpers.BadRequest(api, writer, request, err)
			return
		}
		response.Message = fmt.Sprintf("Added user %d", newUserId)
	}

	err = helpers.WriteJSON(api, writer, http.StatusOK, response)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
	}

}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	err = helpers.RemoveUser(api, userID)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}
	err = helpers.DeleteTokenByUserID(api, userID)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
		return
	}

	response := jsonResponse{OK: true}
	err = helpers.WriteJSON(api, writer, http.StatusOK, response)
	if err != nil {
		helpers.BadRequest(api, writer, request, err)
	}
}
