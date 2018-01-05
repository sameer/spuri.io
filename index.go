package main

import "net/http"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	globalSetHeaders(w, r)
	if r.URL.Path != "/" {
		// TODO: render not found page
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "index", globalContext)
}