package configs

import (
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"log"
)

type AppConfiger interface {
	GetDB() models.DBModel
}

type BaseConfiger interface {
	GetConfig() *Config
	GetInfoLog() *log.Logger
	GetErrorLog() *log.Logger
}
