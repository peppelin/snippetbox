package main

import (
	"log"
	"net/http"
)

func main() {
	// We initialize a new ServerMux and assign the handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use listen and serve to start the new server.

	log.Print("Starting server in port :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
