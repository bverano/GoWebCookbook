package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var getRequestHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

var postRequestHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It's a post request"))
	})

var pathVariableHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		w.Write([]byte("Hi " + name))
	})

func main() {
	router := mux.NewRouter()
	router.Handle("/", handlers.LoggingHandler(os.Stdout, getRequestHandler)).Methods("GET")

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	/*
		os.O_WRONLY: Means that we open a file for writing only
		os.O_CREATE: If the file does not exists, it creates it
		os.O_APPEND: Append data to the file when writing
		0666: Read and write unix file-system permission
	*/

	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}

	router.Handle("/post", handlers.LoggingHandler(os.Stdout, postRequestHandler)).Methods("POST")
	router.Handle("/hello/{name}", handlers.CombinedLoggingHandler(logFile, pathVariableHandler)).Methods("GET")

	log.Fatal(http.ListenAndServe(connHost+":"+connPort, router))
}
