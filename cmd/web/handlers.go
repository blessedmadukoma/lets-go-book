package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	// http.NotFound(w, r)
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
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

	// flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Snippet = snippet
	// data.Flash = flash

	app.render(w, http.StatusOK, "view.html", data)

}

type snippetCreateForm struct {
	Title               string
	Content             string
	Expires             int
	validator.Validator `form:"-"`
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	for i, v := range r.PostForm {
		for _, val := range v {
			switch i {
			case "title":
				form.Title = val
			case "content":
				form.Content = val
			case "expires":
				form.Expires, _ = strconv.Atoi(val)
			default:
				app.errorLog.Fatalln("Could not store the form values!")
			}
		}
	}
	// fmt.Printf("Form: %v\n", form)

	form.CheckField(validator.NotBlank(form.Title), "title", "This title field cannot be blank!")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This title field cannot be more than 100 characters long!")
	form.CheckField(validator.NotBlank(form.Content), "content", "This content field cannot be blank!")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This expires field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Created the HTML structure for signup, created signup struct, updated sign up handler to display the HTML sign up form

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
