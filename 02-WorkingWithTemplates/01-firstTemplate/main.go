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

type person struct {
	ID   string
	Name string
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	person := person{ID: "1", Name: "Bruno"}
	parsedTemplate, _ := template.ParseFiles("templates/first-template.html")
	err := parsedTemplate.Execute(w, person)
	if err != nil {
		log.Printf("Error ocurred while executing the template or writing it's output: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/", renderTemplate)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
