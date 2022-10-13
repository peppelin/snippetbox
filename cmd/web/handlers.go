package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/peppelin/snippetbox/internal/models"
	"github.com/peppelin/snippetbox/internal/validator"
)

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
	Validator   validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the path is different from the root folder and retunr
	// a not found f that's the case
	//
	// New router takes care of thes piece of code below
	//
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
	}
	//using the newTemplateData
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Reading the params from the new router
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	//using the newTemplateData
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
	}

	//Checking for empty title
	form.Validator.CheckField(form.Validator.NotBlank(form.Title),"title", "title can't be empty")
	form.Validator.CheckField(form.Validator.MaxChars(form.Title, 100),"title", "title can't be empty")
	
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "title can't be longer than 100 chars"
	}

	// Checking for empty content
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "content can't be empty"
	}

	// Checking for invalid expiration date
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "expiration time should be 1, 7 or 365"

	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)

		return
	}
	// Pass the data to the SnippetModel.Insert() method, receiving the // ID of the new record back.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
