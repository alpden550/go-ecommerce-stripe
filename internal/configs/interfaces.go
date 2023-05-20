package configs

import "go-ecommerce/internal/models"

type AppConfiger interface {
	GetDB() models.DBModel
}
