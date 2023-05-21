package helpers

import (
	"encoding/json"
	"errors"
	"go-ecommerce/internal/configs"
	"io"
	"net/http"
)

type JsonPayload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func ReadJSON(app configs.AppConfiger, writer http.ResponseWriter, request *http.Request, data interface{}) error {
	maxBytes := 1048576

	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}

	err := dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func WriteJSON(app configs.AppConfiger, writer http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			writer.Header()[k] = v
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(out)
	return nil
}

func BadRequest(app configs.AppConfiger, writer http.ResponseWriter, request *http.Request, err error) error {
	var payload JsonPayload

	payload.Error = true
	payload.Message = err.Error()
	out, errJson := json.MarshalIndent(payload, "", "\t")
	if errJson != nil {
		return err
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(out)
	return nil
}

func InvalidCredentials(app configs.AppConfiger, writer http.ResponseWriter) error {
	var payload JsonPayload
	payload.Error = true
	payload.Message = "Invalid authentication credentials"

	if err := WriteJSON(app, writer, http.StatusUnauthorized, payload); err != nil {
		return err
	}
	return nil
}
