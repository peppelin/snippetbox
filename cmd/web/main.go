package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Adding command line arguments
	addr := flag.String("addr", ":4000", "HTTP network address")

	// We need to parse the flag arguments, if not, it will get the default value
	flag.Parse()

	// We initialize a new ServerMux and assign the handlers
	mux := http.NewServeMux()

	// Serving static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use listen and serve to start the new server.
	log.Printf("Starting server in port %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
