package main

import "net/http"

var indexHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// TODO: render not found page
		http.NotFound(w, r)
	} else {
		renderTemplate(w, "index", globalCtx.Load())
	}
}
