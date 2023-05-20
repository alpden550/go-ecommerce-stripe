package configs

import (
	"github.com/alexedwards/scs/v2"
	"go-ecommerce/internal/models"
	"html/template"
	"log"
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
