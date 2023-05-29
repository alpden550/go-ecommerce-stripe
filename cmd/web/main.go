package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	"github.com/alpden550/go-ecommerce-stripe/internal/driver"
	"github.com/alpden550/go-ecommerce-stripe/internal/handlers"
	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/alpden550/go-ecommerce-stripe/internal/renders"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
const cssVersion = "1"

var config configs.Config
var app *configs.AppConfig
var session *scs.SessionManager

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gob.Register(helpers.TransactionData{})

	flag.IntVar(&config.Port, "port", 4000, "Server port to listen on")
	flag.StringVar(&config.Env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&config.Api, "api", "http://localhost:4001", "URL to api")
	flag.Parse()

	conn, err := prepare()
	if err != nil {
		log.Fatalf("%#e", err)
	}
	defer conn.Close()

	err = serveApp(app)
	if err != nil {
		app.ErrorLog.Printf("%e", err)
		return
	}
}

func prepare() (*sql.DB, error) {
	config.Stripe.Key = os.Getenv("STRIPE_KEY")
	config.Stripe.Secret = os.Getenv("STRIPE_SECRET")
	config.DB.Dsn = os.Getenv("DSN")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(config.DB.Dsn)
	if err != nil {
		return nil, err
	}

	session = scs.New()
	session.Store = postgresstore.New(conn)

	tc := make(map[string]*template.Template)
	app = &configs.AppConfig{
		Config:        config,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		TemplateCache: tc,
		Version:       version,
		DB:            models.DBModel{DB: conn},
		Session:       session,
	}

	renders.SetAppToRender(app)
	handlers.SetAppToHandlers(app)

	return conn, nil
}

func serveApp(app *configs.AppConfig) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.Config.Port),
		Handler:           routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.InfoLog.Println(fmt.Sprintf("Starting HTTP server in %s mode on port %d", app.Config.Env, app.Config.Port))

	return srv.ListenAndServe()
}
