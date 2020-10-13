package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	connHost = "localhost"
	connPort = "8080"
	mongoURL = "127.0.0.1"
)

var session *mgo.Session
var connectionError error

type employee struct {
	ID   int    `json:"uid"`
	Name string `json:"name"`
}

func readDocuments(w http.ResponseWriter, r *http.Request) {
	log.Print("reading documents from database")
	var employees []employee
	collection := session.DB("gocookbook").C("employees")
	err := collection.Find(bson.M{}).All(&employees)
	if err != nil {
		log.Print("error while reading documents from database :: ", err)
		return
	}
	json.NewEncoder(w).Encode(employees)
}

func init() {
	session, connectionError = mgo.Dial(mongoURL)
	if connectionError != nil {
		log.Fatal("error connecting to database :: ", connectionError)
	}
	session.SetMode(mgo.Monotonic, true)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employees", readDocuments).Methods("GET")
	defer session.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting the server", err)
		return
	}
}
