package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// initialize a buffer to check fi there's an error before rendering the frontend to the user
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		// error buf writing
		app.serverError(w, err)
		return
	}

}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// fmt.Println(r.FormValue("title"))
	// fmt.Println("Post form: ", r.PostForm)
	// fmt.Println("Just form: ", r.Form)

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

// isAuthenticated returns boolean value if a user is logged in or not
func (app *application) isAuthenticated(r *http.Request) bool  {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}