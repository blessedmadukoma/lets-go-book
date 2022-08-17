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

func (log *application) home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	// http.NotFound(w, r)
	// 	log.notFound(w)
	// 	return
	// }

	snippets, err := log.snippets.Latest()
	if err != nil {
		log.serverError(w, err)
		return
	}

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	data := log.newTemplateData(r)
	data.Snippets = snippets

	log.render(w, http.StatusOK, "home.html", data)
}

func (log *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		log.notFound(w)
		return
	}

	snippet, err := log.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			log.notFound(w)
		} else {
			log.serverError(w, err)
		}
		return
	}

	// flash := log.sessionManager.PopString(r.Context(), "flash")

	data := log.newTemplateData(r)
	data.Snippet = snippet
	// data.Flash = flash

	log.render(w, http.StatusOK, "view.html", data)

}

type snippetCreateForm struct {
	Title               string
	Content             string
	Expires             int
	validator.Validator `form:"-"`
}

func (log *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := log.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	log.render(w, http.StatusOK, "create.html", data)
}

func (log *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := log.decodePostForm(r, &form)
	if err != nil {
		log.clientError(w, http.StatusBadRequest)
		return
	}

	// fmt.Println("Postform map values:", r.PostForm)
	fmt.Printf("Title: %s\n", form.Title)

	form.CheckField(validator.NotBlank(form.Title), "title", "This title field cannot be blank!")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This title field cannot be more than 100 characters long!")
	form.CheckField(validator.NotBlank(form.Content), "content", "This content field cannot be blank!")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This expires field must equal 1, 7 or 365")

	if !form.Valid() {
		data := log.newTemplateData(r)
		data.Form = form
		log.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := log.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		log.serverError(w, err)
	}

	log.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
