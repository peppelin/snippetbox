package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError writes the error message and the stack trace to the errorLog
// and sends a generic internalServerError to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack)
	app.errorLog.Print(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends the error and the message to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// implementation for the notFound. this is a wrapper around the clientError.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
