package configs

import "github.com/alpden550/go-ecommerce-stripe/internal/models"

type AppConfiger interface {
	GetDB() models.DBModel
}
