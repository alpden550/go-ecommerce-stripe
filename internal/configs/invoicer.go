package configs

import (
	"log"
)

type InvoiceApplication struct {
	Config   Config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Version  string
}

func (i *InvoiceApplication) GetConfig() *Config {
	return &i.Config
}

func (i *InvoiceApplication) GetInfoLog() *log.Logger {
	return i.InfoLog
}

func (i *InvoiceApplication) GetErrorLog() *log.Logger {
	return i.ErrorLog
}
