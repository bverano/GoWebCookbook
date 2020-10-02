package main

import (
	"log"
	"net/http"

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

var postRquestHandler = http.HandlerFunc(
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
	router.Handle("/", getRequestHandler).Methods("GET")
	router.Handle("/post", postRquestHandler).Methods("POST")
	router.Handle("/hello/{name}", pathVariableHandler).Methods("GET", "PUT")

	log.Fatal(http.ListenAndServe(connHost+":"+connPort, router))
}
