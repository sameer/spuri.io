package main

import (
	"net/http"
	"fmt"
)

const (
	staticHandlerPath = "/static/"
	cssHandlerPath    = "/style.css"
	indexHandlerPath  = "/"
	blogHandlerPath   = "/blog/"
	c0dartHandlerPath = "/c0dart/"

	staticDir         = "./static/"
	cssFilePath       = staticDir + "style.css"
	templateDir       = "./templates/"
	templateExtension = ".html.tmpl"
	c0dartDir         = staticDir + "c0dart/"
)

func main() {
	fmt.Println("Launching...")
	compileTemplates()
	initGlobalContext()
	bindHandlers()
	fmt.Println("Ready!")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("Error while launching %v", err)
	}
}

func bindHandlers() {
	http.Handle("/static/", http.StripPrefix(staticHandlerPath, http.FileServer(http.Dir(staticDir))))
	http.HandleFunc(cssHandlerPath, cssHandler)
	http.HandleFunc(indexHandlerPath, indexHandler)
	http.Handle(blogHandlerPath, http.StripPrefix(blogHandlerPath, http.HandlerFunc(blogHandler)))
	http.Handle(c0dartHandlerPath, http.StripPrefix(c0dartHandlerPath, http.HandlerFunc(c0dartHandler)))
}

func compileTemplates() {
	compileTemplate("error", "base")
	compileTemplate("index", "base")
	compileTemplate("blog_index", "base")
	compileTemplate("blog_page", "base")
	compileTemplate("c0dart_gallery", "base")
}
