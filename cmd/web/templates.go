package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox/internal/models"
	"snippetbox/ui"
	"time"
)

type templateData struct {
	Snippets        []*models.Snippet
	Snippet         *models.Snippet
	Form            any
	Flash           string
	CurrentYear     int
	IsAuthenticated bool
	CSRFToken       string
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
	// pages, err := filepath.Glob("./ui/html/pages/*.html")
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// iterate over the pages filepath one-by-one
	for _, page := range pages {
		name := filepath.Base(page)

		// slice containing the filepath patterns for the templates we want to parse
		patterns := []string {
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map
		cache[name] = ts
	}

	return cache, nil
}
