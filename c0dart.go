package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

type C0dartContext struct {
	*GlobalContext
	Images []string
}

func c0dartHandler(w http.ResponseWriter, r *http.Request) {
	var names []string = nil
	if images, err := ioutil.ReadDir("static/c0dart"); err == nil {
		names = make([]string, len(images))
		for i, img := range images {
			names[i] = img.Name()
		}
	} else {
		fmt.Printf("Error reading c0dart directory: %v\n", err)
	}
	ctx := C0dartContext{ GlobalContext: globalContext, Images: names}
	renderTemplate(w, "c0dart", ctx)
}
