package main

import "net/http"

var staticHandler = handlerWithoutState{
	handlerGenericAttributes{staticHandlerPath, true},
	func() http.HandlerFunc {
		var staticHandlerFileServer = http.FileServer(http.Dir(staticDir))
		return func(w http.ResponseWriter, r *http.Request) {
			staticHandlerFileServer.ServeHTTP(w, r)
		}
	}(),
}
