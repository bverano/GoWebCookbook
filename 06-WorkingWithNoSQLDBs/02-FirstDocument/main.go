package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
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

func init() {
	session, connectionError = mgo.Dial(mongoURL)
	if connectionError != nil {
		log.Fatal("error connecting to database :: ", connectionError)
	}
	session.SetMode(mgo.Monotonic, true)
}

func createDocument(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, nameOk := vals["name"]
	id, idOk := vals["id"]
	if nameOk && idOk {
		employeeID, err := strconv.Atoi(id[0])
		if err != nil {
			log.Print("error converting string id into int", err)
			return
		}
		log.Print("going to insert document in database for name :: ", name[0])
		collection := session.DB("gocookbook").C("employees")
		err = collection.Insert(&employee{employeeID, name[0]})
		if err != nil {
			log.Print("error while inserting document in database :: ", err)
			return
		}
		fmt.Fprintf(w, "last created document id is :: %s", id[0])
	} else {
		fmt.Fprintf(w, "error ocurred while creating the document in database for name :: %s", name[0])
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employee/create", createDocument).Methods("POST")
	defer session.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting the server :: ", err)
		return
	}
}
