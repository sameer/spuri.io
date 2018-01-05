package main

import "net/http"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// TODO: render not found page
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "index", globalContext)
}