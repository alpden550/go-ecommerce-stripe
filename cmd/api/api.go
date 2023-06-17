package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/driver"
	handlers "github.com/alpden550/go-ecommerce-stripe/internal/handlers_api"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "1.0.0"

var config configs.Config
var api *configs.ApiApplication

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	flag.IntVar(&config.Port, "port", 4001, "Server port to listen on")
	flag.StringVar(&config.Env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&config.FrontEnd, "api", "http://0.0.0.0:4000", "URL to frontend")
	flag.Parse()

	conn, err := prepare()
	if err != nil {
		log.Fatalf("%#e", err)
	}
	defer conn.Close()

	err = serve()
	if err != nil {
		log.Fatal(err)
	}
}

func prepare() (*sql.DB, error) {
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Printf("%e", err)
	}

	config.Stripe.Key = os.Getenv("STRIPE_KEY")
	config.Stripe.Secret = os.Getenv("STRIPE_SECRET")
	config.DB.Dsn = os.Getenv("DSN")
	config.SMTP.Host = os.Getenv("EMAIL_HOST")
	config.SMTP.Port = port
	config.SMTP.Username = os.Getenv("EMAIL_USERNAME")
	config.SMTP.Password = os.Getenv("EMAIL_PASSWORD")
	config.SMTP.EmailFrom = os.Getenv("EMAIL_FROM")
	config.SecretKey = os.Getenv("SECRET_KEY")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(config.DB.Dsn)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	api = &configs.ApiApplication{
		Config:   config,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Version:  version,
		DB:       models.DBModel{DB: conn},
	}

	handlers.SetApiAppToHandlers(api)

	return conn, nil
}

func serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", api.Config.Port),
		Handler:           routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	api.InfoLog.Println(fmt.Sprintf("Starting Back end server in %s mode on port %d", api.Config.Env, api.Config.Port))

	return srv.ListenAndServe()
}
