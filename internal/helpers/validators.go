package helpers

import (
	"fmt"
	"net/http"

	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
)

func FailedValidation(api configs.BaseConfiger, writer http.ResponseWriter, errors map[string]string) {
	errorLog := api.GetErrorLog()
	var payload struct {
		OK      bool              `json:"ok"`
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors"`
	}
	payload.OK = false
	payload.Message = "failed validation"
	payload.Errors = errors

	err := WriteJSON(writer, http.StatusUnprocessableEntity, payload)
	if err != nil {
		errorLog.Printf("%w", fmt.Errorf("%e", err))
		return
	}
}
