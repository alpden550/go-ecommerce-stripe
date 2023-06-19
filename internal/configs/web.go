package configs

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
)

type AppConfig struct {
	Config        Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
	DB            models.DBModel
	Session       *scs.SessionManager
}

func (app *AppConfig) GetDB() models.DBModel {
	return app.DB
}
