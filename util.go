package main

import (
	"html/template"
	"net/http"
)

var templates = make(map[string]*template.Template)

// Compiles a template. The first string must be the "leaf" template of the tree, the order of the remaining templates
// used does not matter.
func compileTemplate(tmpls ...string) {
	tmpl := tmpls[0]
	for i := range tmpls {
		tmpls[i] = templateDir + tmpls[i] + templateExtension
	}
	templates[tmpl] = template.Must(template.ParseFiles(tmpls...))
}

// Takes a template name and renders the corresponding template. If there's an error in rendering, it handles it.
func renderTemplate(w http.ResponseWriter, tmpl string, ctx interface{}) {
	w.Header().Set("Content-Type", "text/html")
	err := templates[tmpl].ExecuteTemplate(w, "base", ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func contains(stringSlice []string, searchString string) bool {
	for _, str := range stringSlice {
		if str == searchString {
			return true
		}
	}
	return false
}
