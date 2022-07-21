package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (log *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", log.home)
	mux.HandleFunc("/snippet/view", log.snippetView)
	mux.HandleFunc("/snippet/create", log.snippetCreate)

	standard := alice.New(log.recoverPanic, log.logRequest, secureHeaders)

	// return log.recoverPanic(log.logRequest(secureHeaders(mux)))
	return standard.Then(mux)
}
