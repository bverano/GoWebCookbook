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

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login!!")
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Logout!!")
}

func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("errore while running the server: ", err)
		return
	}
}
