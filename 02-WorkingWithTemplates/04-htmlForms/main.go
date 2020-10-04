package main

import (
	"html/template"
	"log"
	"net/http"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

func login(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("templates/login-form.html")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		log.Printf("Error while rendering the template or writing it's output: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/", login)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting the server: ", err)
		return
	}
}
