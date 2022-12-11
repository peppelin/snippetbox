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

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { app.notFound(w) })

	// Serving static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Unprotected application using dynamic chain
	// use noSurf middleware on all dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// Update routes with the new router
	//and change the type from router.handlerFunc to router.Handler
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	// User related routes
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	// Protected application using protectedThen func
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Crteate the middleware chain
	standard := alice.New(app.recoveryPanic, app.logRequest, secureHeaders)

	//return the standard middleware with our mux using alice
	return standard.Then(router)
}
