package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
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

	(&handlerCoordinator{
		handlers: []handlerInterface{
			&staticHandler,
			&cssHandler,
			&blogHandler,
			&c0dartHandler,
			&aboutHandler,
			&studioStatisticsHandler,
			&indexHandler,
		},
	}).start(http.DefaultServeMux, time.Minute)

	fmt.Println("Initialized!")

	bindAddress := prodBindAddress
	if ip := os.Getenv(prodIpEnvironmentVariable); ip != "" {
		bindAddress = ip
	} else if os.Getenv(devEnvironmentVariable) != "" {
		bindAddress = devBindAddress
		fmt.Println("Environment is dev")
	} else {
		panic("Environment is unknown!")
	}
	fmt.Println("Listening on", bindAddress)
	if err := http.ListenAndServe(bindAddress, nil); err != nil {
		fmt.Printf("Error while launching %v\n", err)
	}
}

func compileTemplates() {
	templates := [6][2]string{
		{"error", "base"},
		{"index", "base"},
		{"blog_index", "base"},
		{"blog_page", "base"},
		{"c0dart_gallery", "base"},
		{"about", "base"},
	}
	for _, toCompile := range templates {
		compileTemplate(toCompile[0], toCompile[1])
	}
}
