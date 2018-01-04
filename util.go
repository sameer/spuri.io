package main

import (
	"net/http"
	"html/template"
)

var templates = make(map[string]*template.Template)

func compileTemplate(tmpls ...string) {
	tmpl := tmpls[0]
	for i, _ := range tmpls {
		tmpls[i] = templateDir + tmpls[i] + templateExtension
	}
	templates[tmpl] = template.Must(template.ParseFiles(tmpls...))
}

func renderTemplate(w http.ResponseWriter, tmpl string, ctx interface{}) {
	w.Header().Set("Content-Type", "text/html")
	err := templates[tmpl].ExecuteTemplate(w, "base", ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
