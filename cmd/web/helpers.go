package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form"
)

func (log *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.errorLog.Output(2, trace)

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

	_, err = buf.WriteTo(w)
	if err != nil {
		// error buf writing
		log.serverError(w, err)
		return
	}
	
}

func (log *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = log.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
