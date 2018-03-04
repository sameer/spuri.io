package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var cssHandler = func() http.HandlerFunc {
	const defaultCSS = "* { font-family: sans-serif; }"
	// Keep the file in-memory because it's only several KB & lowers load time.
	cssFile, err := ioutil.ReadFile(cssFilePath)
	if err != nil {
		cssFile = []byte(defaultCSS)
		fmt.Println(err)
	}

	cssFileLength := strconv.Itoa(len(cssFile))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Content-Length", cssFileLength)
		w.Write(cssFile)
	}
}()
