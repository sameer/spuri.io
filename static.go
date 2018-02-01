package main

import "net/http"

var staticHandler = func() http.HandlerFunc {
	var staticHandlerFileServer = http.FileServer(http.Dir(staticDir))
	return func(w http.ResponseWriter, r *http.Request) {
		globalSetHeaders(w, r)
		staticHandlerFileServer.ServeHTTP(w, r)
	}
}()
