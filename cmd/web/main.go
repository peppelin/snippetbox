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
	// We initialize a new ServerMux and assign the handlers
	mux := http.NewServeMux()

	// Serving static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	// Initialize http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  mux,
	}
	// Use listen and serve to start the new server.
	app.infoLog.Printf("Starting server in port %s", *addr)
	// Calling our nbew http server
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
