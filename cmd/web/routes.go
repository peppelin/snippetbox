package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// updated signature to return a http.Handler instead of *http.Servermux
func (app *application) routes() http.Handler {

	// Initializethe new router from julienschmidt/httprouter
	router := httprouter.New()

	// Serving static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/statis/*filepath", http.StripPrefix("/static", fileServer))

	// Update routes with the new router
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Crteate the middleware chain
	standard := alice.New(app.recoveryPanic, app.logRequest, secureHeaders)

	//return the standard middleware with our mux using alice
	return standard.Then(router)
}
