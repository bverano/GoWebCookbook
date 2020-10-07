package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/securecookie"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func getUserName(r *http.Request) (userName string) {
	cookie, err := r.Cookie("session")
	if err == nil {
		cookieValue := make(map[string]string)
		err = cookieHandler.Decode("session", cookie.Value, &cookieValue)
		if err == nil {
			userName = cookieValue["username"]
		}
	}
	return userName
}

func setSession(userName string, w http.ResponseWriter) {
	value := map[string]string{"username": userName}
	encoded, err := cookieHandler.Encode("session", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	target := "/"
	if username != "" && password != "" {
		setSession(username, w)
		target = "/home"
	}
	http.Redirect(w, r, target, 302)
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("templates/login.html")
	parsedTemplate.Execute(w, nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	if userName != "" {
		data := map[string]interface{}{"userName": userName}
		parsedTemplate, _ := template.ParseFiles("templates/home.html")
		parsedTemplate.Execute(w, data)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", loginPage)
	router.HandleFunc("/home", homePage)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	http.Handle("/", router)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting the server: ", err)
		return
	}
}
