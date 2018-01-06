package main

import "net/http"

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	globalSetHeaders(w, r)
	renderTemplate(w, "about", globalContext)
}