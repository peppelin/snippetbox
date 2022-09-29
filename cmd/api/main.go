package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// Check if the path is different from the root folder and retunr
	// a not found f that's the case
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("hello snippetbox!"))
}
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here you can view some snippets"))
}
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
	w.Write([]byte("Here you can create some snippets"))
}
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
