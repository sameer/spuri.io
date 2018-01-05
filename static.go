package main

import "net/http"

var staticHandlerFileServer = http.FileServer(http.Dir(staticDir))
func staticHandler(w http.ResponseWriter, r *http.Request) {
	globalSetHeaders(w, r)
	staticHandlerFileServer.ServeHTTP(w, r)
}