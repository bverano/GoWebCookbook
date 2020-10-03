package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func main() {
	http.HandleFunc("/", helloWorld)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
