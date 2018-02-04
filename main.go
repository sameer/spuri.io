package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	staticHandlerPath           = "/static/"
	cssHandlerPath              = "/style.css"
	blogHandlerPath             = "/blog/"
	c0dartHandlerPath           = "/c0dart/"
	aboutHandlerPath            = "/about"
	studioStatisticsHandlerPath = "/studio-statistics.png"
	indexHandlerPath            = "/"

	staticDir         = "./static/"
	cssFilePath       = staticDir + "style.css"
	templateDir       = "./templates/"
	templateExtension = ".html.tmpl"
	c0dartDir         = staticDir + "c0dart/"
	blogDir           = "./blog/"

	prodBindAddress        = "[0:0:0:0:0:0:0:0]:80"
	devBindAddress         = "127.0.0.1:8000"
	devEnvironmentVariable = "DEV=1"
)

func main() {
	fmt.Println("Launching...")
	compileTemplates()
	globalCtx.refresh()
	bindHandlers()
	fmt.Println("Ready!")

	bindAddress := prodBindAddress
	if contains(os.Environ(), devEnvironmentVariable) {
		bindAddress = devBindAddress
	}
	err := http.ListenAndServe(bindAddress, nil)
	if err != nil {
		fmt.Printf("Error while launching %v\n", err)
	}
}

func bindHandlers() {
	http.Handle("/static/", http.StripPrefix(staticHandlerPath, staticHandler))
	http.HandleFunc(cssHandlerPath, cssHandler)
	http.Handle(blogHandlerPath, http.StripPrefix(blogHandlerPath, blogHandler))
	http.Handle(c0dartHandlerPath, http.StripPrefix(c0dartHandlerPath, c0dartHandler))
	http.HandleFunc(aboutHandlerPath, aboutHandler)
	http.HandleFunc(studioStatisticsHandlerPath, studioStatisticsHandler)
	http.HandleFunc(indexHandlerPath, indexHandler)
}

func compileTemplates() {
	compileTemplate("error", "base")
	compileTemplate("index", "base")
	compileTemplate("blog_index", "base")
	compileTemplate("blog_page", "base")
	compileTemplate("c0dart_gallery", "base")
	compileTemplate("about", "base")
}
