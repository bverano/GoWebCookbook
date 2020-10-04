package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/gorilla/schema"
)

const (
	connHost             = "localhost"
	connPort             = "8080"
	usernameErrorMessage = "Please enter a valid Username."
	passwordErrorMessage = "Please enter a valid Password."
	genericErrorMessage  = "Validation Error."
)

type user struct {
	Username string `valid:"alpha,required"`
	Password string `valid:"alpha,required"`
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

func validateUser(w http.ResponseWriter, r *http.Request, u *user) (bool, string) {
	valid, validationError := govalidator.ValidateStruct(u)
	if !valid {
		usernameError := govalidator.ErrorByField(validationError, "Username")
		passwordError := govalidator.ErrorByField(validationError, "Password")
		if usernameError != "" {
			log.Printf("username validation error: %v", usernameError)
			return valid, usernameErrorMessage
		}
		if passwordError != "" {
			log.Printf("password validation error: %v", passwordError)
			return valid, passwordErrorMessage
		}
	}
	return valid, genericErrorMessage
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
		valid, validationErrorMessage := validateUser(w, r, u)
		if !valid {
			fmt.Fprintf(w, validationErrorMessage)
			return
		}
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
