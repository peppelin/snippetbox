package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Adding command line arguments
	addr := flag.String("addr", ":4000", "HTTP network address")

	// We need to parse the flag arguments, if not, it will get the default value
	flag.Parse()

	// initialze our application
	app := &application{
		// create a new logger for ERROR logs
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		// create a new logger for INFO logs
		infoLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
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
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
