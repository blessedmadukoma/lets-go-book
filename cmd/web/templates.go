package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"snippetbox/internal/models"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Parse current year in all the pages: this can be equivalent to parsing user profile picture or details after log in on all the pages
func (log *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// Adding human readable date instead of the UTC 0000 old vibe
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// store the humanDate function in a global variable
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get all files matching the *.html extension
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// iterate over the pages filepath one-by-one
	for _, page := range pages {
		name := filepath.Base(page)

		// register the global variable functions before calling the ParseFiles method
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// parse all .html file in partials folder into the template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Parse the page template to this template set ts
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map
		cache[name] = ts
	}

	return cache, nil
}
