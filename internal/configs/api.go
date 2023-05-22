package configs

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"log"
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
