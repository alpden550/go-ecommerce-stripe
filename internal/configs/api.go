package configs

import (
	"go-ecommerce/internal/models"
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
