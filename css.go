package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

// Keep the file in-memory because it's only several KB & lowers load time.

var cssFile []byte = nil

func cssHandler(w http.ResponseWriter, r *http.Request) {
	if cssFile == nil {
		content, err := ioutil.ReadFile(cssFilePath)
		if err != nil {
			cssFile = []byte("* { font-family: sans-serif; }")
			fmt.Printf("%v", err)
		} else {
			cssFile = content
		}
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(cssFile)
}
