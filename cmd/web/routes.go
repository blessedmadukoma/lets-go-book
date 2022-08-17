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

	dynamic := alice.New(log.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(log.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(log.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(log.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(log.snippetCreatePost))
	// router.HandlerFunc(http.MethodPost, "/snippet/create", log.snippetCreatePost)

	standard := alice.New(log.recoverPanic, log.logRequest, secureHeaders)

	// return log.recoverPanic(log.logRequest(secureHeaders(mux)))
	return standard.Then(router)
}
