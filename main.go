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
	codeArtHandlerPath          = "/codeArt/"
	aboutHandlerPath            = "/about"
	studioStatisticsHandlerPath = "/studio-statistics.png"
	indexHandlerPath            = "/"

	staticDir         = "./static/"
	cssFilePath       = staticDir + "style.css"
	templateDir       = "./templates/"
	templateExtension = ".html.tmpl"
	codeArtDir        = staticDir + "codeArt/"
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
		if httpServer != nil {
			httpServer.Shutdown(nil)
		}
		if httpsServer != nil {
			httpsServer.Shutdown(nil)
		}
	}()
}

var httpServer *http.Server
var httpsServer *http.Server

func main() {
	log.Println("Launching...")

	(&handlerCoordinator{
		handlers: []handlerInterface{
			&staticHandler,
			&cssHandler,
			&blogHandler,
			&codeArtHandler,
			&aboutHandler,
			&studioStatisticsHandler,
			&indexHandler,
		},
	}).start(http.DefaultServeMux, time.Minute)

	log.Println("Initialized!")

	var bindAddress string
	var bindAddressTLS string
	if ip := os.Getenv(prodIpEnvironmentVariable); ip != "" {
		bindAddress = ip + ":80"
		bindAddressTLS = ip + ":443"
	} else if os.Getenv(devEnvironmentVariable) != "" {
		bindAddress = devBindAddress
		log.Println("Environment is dev")
	} else {
		panic("Environment is unknown!")
	}
	if cert, key := os.Getenv(certEnvironmentVariable), os.Getenv(keyEnvironmentVariable); cert != "" && key != "" {
		if bindAddressTLS == "" {
			log.Println("No TLS bind address provided, not starting")
		} else {
			httpsServer = &http.Server{Addr: bindAddressTLS, Handler: nil}
			go func() {
				log.Println("Listening on", bindAddressTLS)
				err := http.ListenAndServeTLS(bindAddressTLS, cert, key, nil)
				if err != nil && err != http.ErrServerClosed {
					log.Fatalln(err)
				}
			}()
		}
	}
	log.Println("Listening on", bindAddress)
	httpServer = &http.Server{Addr: bindAddress, Handler: nil}
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}

func compileTemplates() {
	templates := [6][2]string{
		{"error", "base"},
		{"index", "base"},
		{"blog_index", "base"},
		{"blog_page", "base"},
		{"codeArt_gallery", "base"},
		{"about", "base"},
	}
	for _, toCompile := range templates {
		compileTemplate(toCompile[0], toCompile[1])
	}
}
