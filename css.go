package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const cssDefault = "* { font-family: sans-serif; }"

var cssHandler = handlerWithFinalState{
	handlerGenericAttributes: handlerGenericAttributes{cssHandlerPath, false},
	handler: func(w http.ResponseWriter, r *http.Request, s state) {
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Content-Length", strconv.Itoa(len(s.([]byte))))
		w.Write(s.([]byte))
	},
	initializer: func() state {

		// Keep the file in-memory because it's only several KB & lowers load time.
		cssFile, err := ioutil.ReadFile(cssFilePath)
		if err != nil {
			cssFile = []byte(cssDefault)
			log.Println(err)
		}
		return cssFile
	},
}
