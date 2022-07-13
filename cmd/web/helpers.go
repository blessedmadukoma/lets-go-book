package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (log *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (log *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (log *application) notFound(w http.ResponseWriter) {
	log.clientError(w, http.StatusNotFound)
}

func (log *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := log.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		log.serverError(w, err)
		return
	}

	// initialize a buffer to check fi there's an error before rendering the frontend to the user
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		log.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	// err := ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	log.serverError(w, err)
	// 	return
	// }
	buf.WriteTo(w)
}
