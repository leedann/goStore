package main

import (
	"log"
	"net/http"
	"os"

	"github.com/info344-s17/challenges-leedann/apiserver/handlers"
)

const defaultPort = "443"

const (
	apiRoot    = "/v1/"
	apiSummary = apiRoot + "summary"
)

//main is the main entry point for this program
func main() {
	//read and use the following environment variables
	//when initializing and starting your web server
	// PORT - port number to listen on for HTTP requests (if not set, use defaultPort)
	// HOST - host address to respond to (if not set, leave empty, which means any host)

	PORT := os.Getenv("PORT")
	if len(PORT) == 0 {
		PORT = defaultPort
	}
	HOST := os.Getenv("HOST")
	if len(HOST) == 0 {
		HOST = ""
	}

	certPath := os.Getenv("TLSCERT")
	keyPath := os.Getenv("TLSKEY")

	//add your handlers.SummaryHandler function as a handler
	//for the apiSummary route
	//HINT: https://golang.org/pkg/net/http/#HandleFunc
	http.HandleFunc(apiSummary, handlers.SummaryHandler)

	//start your web server and use log.Fatal() to log
	//any errors that occur if the server can't start
	//HINT: https://golang.org/pkg/net/http/#ListenAndServe
	addr := HOST + ":" + PORT
	log.Fatal(http.ListenAndServeTLS(addr, certPath, keyPath, nil))

}
