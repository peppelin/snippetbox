package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// updated signature to return a http.Handler instead of *http.Servermux
func (app *application) routes() http.Handler {
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

	// Crteate the middleware chain
	standard := alice.New(app.recoveryPanic, app.logRequest, secureHeaders)

	//return the standard middleware with our mux using alice
	return standard.Then(mux)
}
