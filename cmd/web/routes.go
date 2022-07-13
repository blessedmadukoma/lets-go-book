package main

import "net/http"

func (log *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", log.home)
	mux.HandleFunc("/snippet/view", log.snippetView)
	mux.HandleFunc("/snippet/create", log.snippetCreate)
	
	return mux
}
