package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/tobi/mogrify-go"
)

type C0dartContext struct {
	*GlobalContext
	Images []C0dartImage
}

type C0dartImage struct {
	Name string
	Href string
}

var igs []mogrify.Image = nil

func c0dartHandler(w http.ResponseWriter, r *http.Request) {
	var c0dartImages []C0dartImage= nil
	if images, err := ioutil.ReadDir("static/c0dart"); err == nil {
		c0dartImages = make([]C0dartImage, len(images))
		for i, img := range images {
			c0dartImages[i] = C0dartImage{Name: img.Name(), Href: "/static/c0dart/" + img.Name()}
		}
	} else {
		fmt.Printf("Error reading c0dart directory: %v\n", err)
	}
	ctx := C0dartContext{ GlobalContext: globalContext, Images: c0dartImages}
	renderTemplate(w, "c0dart", ctx)
}
