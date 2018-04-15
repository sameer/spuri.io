package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	devBindAddress            = "127.0.0.1:8000"
	devEnvironmentVariable    = "DEV"
	prodIpEnvironmentVariable = "PROD_IP"
	certEnvironmentVariable   = "CERT_FILE"
	keyEnvironmentVariable    = "KEY_FILE"
)

func init() {
	compileTemplates()
	goGracefulShutdownHandler()
}

func goGracefulShutdownHandler() {
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigchan
		log.Println("Gracefully shutting down (", sig, ")...")
		if server != nil {
			server.Shutdown(nil)
		}
	}()
}

var server *http.Server

func main() {
	log.Println("Launching...")

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

	log.Println("Initialized!")

	var bindAddress string
	if ip := os.Getenv(prodIpEnvironmentVariable); ip != "" {
		bindAddress = ip
		bindAddress += ":80"

	} else if os.Getenv(devEnvironmentVariable) != "" {
		bindAddress = devBindAddress
		log.Println("Environment is dev")
	} else {
		panic("Environment is unknown!")
	}
	bindAddressTLS := bindAddress + ":443"

	if cert, key := os.Getenv(certEnvironmentVariable), os.Getenv(keyEnvironmentVariable); cert != "" && key != "" {
		go func() {
			log.Println("Listening on", bindAddressTLS)
			http.ListenAndServeTLS(bindAddressTLS, cert, key, nil)
		}()
	}
	log.Println("Listening on", bindAddress)
	server = &http.Server{Addr: bindAddress, Handler: nil}
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
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
