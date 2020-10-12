package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	connHost   = "localhost"
	connPort   = "8080"
	driverName = "postgres"
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "12345"
	dbName     = "gocookbook"
)

var db *sql.DB
var connectionError error

func init() {
	db, connectionError = sql.Open(driverName, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName))
	if connectionError != nil {
		log.Fatal("error connecting to database :: ", connectionError)
	}
}

func createRecord(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("going to insert record in database for name: ", name[0])
		stmt, err := db.Prepare("INSERT INTO employees(name) VALUES ($1)")
		if err != nil {
			log.Print("error preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0])
		if err != nil {
			log.Print("error executing query: ", err)
			return
		}
		id, err := result.LastInsertId()
		fmt.Fprintf(w, "Last inserted record id is: %v", id)
	} else {
		fmt.Fprintf(w, "Error ocurred while created record in database for name :: %s", name[0])
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employees/create", createRecord).Methods("POST")
	defer db.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server ::", err)
		return
	}
}
