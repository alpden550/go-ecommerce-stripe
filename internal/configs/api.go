package configs

import (
	"log"

	"github.com/alpden550/go-ecommerce-stripe/internal/models"
)

type ApiApplication struct {
	Config   Config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Version  string
	DB       models.DBModel
}

func (api *ApiApplication) GetDB() models.DBModel {
	return api.DB
}

func (api *ApiApplication) GetConfig() *Config {
	return &api.Config
}

func (api *ApiApplication) GetInfoLog() *log.Logger {
	return api.InfoLog
}

func (api *ApiApplication) GetErrorLog() *log.Logger {
	return api.ErrorLog
}
