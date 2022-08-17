package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":8000", "HTTP Network Address")

	dsn := flag.String("dsn", "snippetuser:snippet@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Logging
	// create a file and log into it
	// logFile, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Fatalf("Error opening file: %s\n", err)
	// }

	// defer logFile.Close()

	// infoLog := log.New(logFile, "INFO:\t", log.Ldate|log.Ltime)
	// errorLog := log.New(logFile, "ERROR:\t", log.Ldate|log.Ltime|log.Llongfile)

	// Display log in terminal
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	router := app.routes()

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  router,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
