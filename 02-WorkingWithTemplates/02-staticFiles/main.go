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
	p := person{ID: "1", Name: "Jacobo"}
	parsedTemplate, _ := template.ParseFiles("templates/first-template.html")
	err := parsedTemplate.Execute(w, p)
	if err != nil {
		log.Printf("Error while executing the template or writing it's output: %v", err)
		return
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", renderTemplate)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
