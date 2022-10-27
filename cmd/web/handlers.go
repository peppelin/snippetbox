package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/peppelin/snippetbox/internal/models"
	"github.com/peppelin/snippetbox/internal/validator"
)

// Adding the struct tags so form.Decoder can get them from the form.
type snippetCreateForm struct {
	Title     string              `form:"title"`
	Content   string              `form:"content"`
	Expires   int                 `form:"expires"`
	Validator validator.Validator `form:"-"`
}

// add userSignupForm
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// adding userLoginForm
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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
	var form snippetCreateForm
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Checking for empty title
	form.Validator.CheckField(form.Validator.NotBlank(form.Title), "title", "title can't be empty")
	form.Validator.CheckField(form.Validator.MaxChars(form.Title, 100), "title", "title can't be longer than 100 chars")
	form.Validator.CheckField(form.Validator.NotBlank(form.Content), "content", "content can't be empty")
	form.Validator.CheckField(form.Validator.PermitedInt(form.Expires, 1, 7, 365), "expires", "expiration time should be 1, 7 or 365")

	if !form.Validator.Valid() {
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

	//we'll use the put() method to store the message in the session for succesfully created snippet
	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully!")

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate a new user")
}
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Checking all validations
	form.Validator.CheckField(form.Validator.NotBlank(form.Email), "email", "email can't be blank")
	form.Validator.CheckField(form.Validator.NotBlank(form.Name), "name", "name can't be blank")
	form.Validator.CheckField(form.Validator.NotBlank(form.Password), "password", "password can't be blank")
	form.Validator.CheckField(form.Validator.MinChars(form.Password, 8), "password", "password must be at least 8 chars long")
	form.Validator.CheckField(form.Validator.Matches(form.Email, validator.EmailRX), "email", "email must be a valid email address")

	// if any checks before failed, we render the errors
	if !form.Validator.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Trying to create the record in the database
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Validator.AddFieldError("email", "email address already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// if all is successful, flash the ok message
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please sign in.")

	// Redirect to the login age
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "logout a user")
}
