package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.HandleFunc("/", renderTemplate).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error while starting the server: ", err)
		return
	}
}
