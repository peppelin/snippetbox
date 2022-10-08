package main

import (
	"html/template"
	"path/filepath"

	"github.com/peppelin/snippetbox/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize our cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extracting the name of each page template
		name := filepath.Base(page)

		// parse the base template
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// parse all the partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// parse the templates
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Adding the template set to the cache
		cache[name] = ts
	}
	return cache, nil
}
