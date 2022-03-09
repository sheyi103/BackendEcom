package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

const version = "1.0.0"
const cssVersion = "1"

//configuration type
type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}

	stripe struct {
		secret string
		key    string
	}
}

//receiver for the various part of my application

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

//create a function that is a pointer to application

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	app.infoLog.Println("starting HTTP server in %s  mode on port %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	var cfg config

	//persing varible from command line
	flag.IntVar(&cfg.port, "port", 4000, "server port to listern on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")

	flag.Parse()

	//reading from environment variable
	cfg.stripe.key = os.Getenv("STRIP_KEY")
	cfg.stripe.secret = os.Getenv("STRIP_SECRET")

	//logging Function
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Map for template Cache
	tc := make(map[string]*template.Template)

	//application variable
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}

	err := app.serve()

	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

}
