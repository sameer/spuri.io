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

	prodBindAddress           = "[0:0:0:0:0:0:0:0]:80"
	devBindAddress            = "127.0.0.1:8000"
	devEnvironmentVariable    = "DEV"
	prodIpEnvironmentVariable = "PROD_IP"
)

func main() {
	fmt.Println("Launching...")
	compileTemplates()
	staticCtx.Store(staticContext{}.init())
	bindHandlers()
	fmt.Println("Ready!")

	bindAddress := prodBindAddress
	if ip := os.Getenv(prodIpEnvironmentVariable); ip != "" {
		bindAddress = ip
	} else if os.Getenv(devEnvironmentVariable) != "" {
		bindAddress = devBindAddress
	}
	err := http.ListenAndServe(bindAddress, nil)
	if err != nil {
		fmt.Printf("Error while launching %v\n", err)
	}
}

func bindHandler(path string, handler http.Handler) {
	http.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		globalSetHeaders(w, r)
		handler.ServeHTTP(w, r)
	}))
}

func bindHandlers() {
	bindHandler(staticHandlerPath, http.StripPrefix(staticHandlerPath, staticHandler))
	bindHandler(cssHandlerPath, cssHandler)
	bindHandler(blogHandlerPath, http.StripPrefix(blogHandlerPath, blogHandler))
	bindHandler(c0dartHandlerPath, http.StripPrefix(c0dartHandlerPath, c0dartHandler))
	bindHandler(aboutHandlerPath, aboutHandler)
	bindHandler(studioStatisticsHandlerPath, studioStatisticsHandler)
	bindHandler(indexHandlerPath, indexHandler)
}

func compileTemplates() {
	templates := [][]string{
		{"error", "base"},
		{"index", "base"},
		{"blog_index", "base"},
		{"blog_page", "base"},
		{"c0dart_gallery", "base"},
		{"about", "base"},
	}
	for _, toCompile := range templates {
		compileTemplate(toCompile...)
	}
}
