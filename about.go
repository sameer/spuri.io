package main

import "net/http"

var aboutHandler = handlerWithoutState{
	handlerGenericAttributes{aboutHandlerPath, false},
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "about", staticCtx.Load())
	}),
}
