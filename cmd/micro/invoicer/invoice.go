package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"

	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

var config configs.Config
var invoicer *configs.InvoiceApplication

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	flag.IntVar(&config.Port, "port", 5001, "Server port to listen on")
	flag.Parse()

	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Printf("%e", err)
	}
	config.SMTP.Host = os.Getenv("EMAIL_HOST")
	config.SMTP.Port = port
	config.SMTP.Username = os.Getenv("EMAIL_USERNAME")
	config.SMTP.Password = os.Getenv("EMAIL_PASSWORD")
	config.SMTP.EmailFrom = os.Getenv("EMAIL_FROM")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	invoicer = &configs.InvoiceApplication{
		Config:   config,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Version:  version,
	}

	err = helpers.CreateDir("./invoices")
	if err != nil {
		return
	}

	err = serve()
	if err != nil {
		log.Fatal(err)
	}
}

func serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", invoicer.Config.Port),
		Handler:           routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	invoicer.InfoLog.Println(fmt.Sprintf("Starting Invoice MicroService on port %d", invoicer.Config.Port))

	return srv.ListenAndServe()
}
