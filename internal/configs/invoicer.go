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
