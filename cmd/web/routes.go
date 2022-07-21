package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (log *application) routes() http.Handler {

	router := httprouter.New()

	// Handle not found error
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", log.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", log.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", log.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", log.snippetCreatePost)

	standard := alice.New(log.recoverPanic, log.logRequest, secureHeaders)

	// return log.recoverPanic(log.logRequest(secureHeaders(mux)))
	return standard.Then(router)
}
