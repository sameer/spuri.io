package main

import (
	"net/http"
	"fmt"
)

const (
	cssHandlerPath    = "/style.css"
	cssFilePath       = "./static/style.css"
	blogHandlerPath   = "/blog/"
	c0dartPath        = "/c0dart/"
	templateExtension = ".html.tmpl"
	templateDir       = "./templates/"
)

func main() {
	fmt.Println("Launching...")
	compileTemplates()
	initGlobalContext()
	bindHandlers()
	fmt.Println("Ready!")
	http.ListenAndServe(":8000", nil)
}

func bindHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc(cssHandlerPath, cssHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc(blogHandlerPath, blogHandler)
	http.HandleFunc(c0dartPath, c0dartHandler)
}

func compileTemplates() {
	compileTemplate("error", "base")
	compileTemplate("index", "base")
	compileTemplate("blog_index", "base")
	compileTemplate("blog_page", "base")
	compileTemplate("c0dart", "base")
}
