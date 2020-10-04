package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

type user struct {
	Username string
	Password string
}

func readForm(r *http.Request) *user {
	r.ParseForm()
	u := new(user)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(u, r.PostForm)
	if decodeErr != nil {
		log.Printf("error mapping parsed form data to struct: %v", decodeErr)
	}
	return u
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parsedTemplate, _ := template.ParseFiles("templates/login-form.html")
		err := parsedTemplate.Execute(w, nil)
		if err != nil {
			log.Printf("Error while rendering the template or writing it's output: %v", err)
			return
		}
	} else {
		u := readForm(r)
		fmt.Fprintf(w, "Hello "+u.Username+"!")
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
