package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/peppelin/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// adding the snippetmodel to make it avaliable to the handlers
	snippets *models.SnippetModel
	//adding the cache template
	templateCache map[string]*template.Template
}

func main() {
	// Adding command line arguments
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=True", "MySQL data source name")

	// We need to parse the flag arguments, if not, it will get the default value
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Stablishing connection to MySQL
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize the cache

	templateCache, err := newTemplateCache()

	if err != nil {
		errorLog.Fatal(err)
	}

	// initialze our application
	app := &application{
		// create a new logger for ERROR logs
		errorLog: errorLog,
		// create a new logger for INFO logs
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	// Initialize http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}
	// Use listen and serve to start the new server.
	app.infoLog.Printf("Starting server in port %s", *addr)
	// Calling our nbew http server
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
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
