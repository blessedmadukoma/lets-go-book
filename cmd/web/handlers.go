package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
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

	data := log.newTemplateData(r)
	data.Snippet = snippet

	log.render(w, http.StatusOK, "view.html", data)

}

func (log *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet...\n"))
}

func (log *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		w.Header().Set("Content-Type", "application/json")
		// http.Error(w, `{"Message": "Method not allowed!"}`, http.StatusMethodNotAllowed)
		log.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := log.snippets.Insert(title, content, expires)
	if err != nil {
		log.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view%d", id), http.StatusSeeOther)
}
