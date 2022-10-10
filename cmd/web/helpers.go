package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError writes the error message and the stack trace to the errorLog
// and sends a generic internalServerError to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack)
	// we need to get the report from the original source, not from where the error is generated
	app.errorLog.Output(2, trace)

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

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	// Retrieve the right templateset based on the name.
	// if it doesn't exist, return a serverError()
	// Page with errors found in ui/html/pages/error.tmpl
	// rename it to view.tmpl

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer to store the rendered pages and check for errors
	buf := new(bytes.Buffer)

	// render the page into the buffer and check for errors
	err := ts.ExecuteTemplate(w, "base", data)

	if err != nil {
		app.serverError(w, err)
		return
	}

	// write the header data
	w.WriteHeader(status)

	// Write the buffer content into the httwriter
	buf.WriteTo(w)

}
