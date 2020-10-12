package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	connHost     = "localhost"
	connPort     = "8080"
	dbDriverName = "postgres"
	dbHost       = "localhost"
	dbPort       = 5432
	dbUser       = "postgres"
	dbPassword   = "12345"
	dbName       = "gocookbook"
)

var db *sql.DB
var connectionError error

func init() {
	db, connectionError = sql.Open(dbDriverName, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName))
	if connectionError != nil {
		log.Fatal("error connecting to database :: ", connectionError)
	}
}

type employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func readRecords(w http.ResponseWriter, r *http.Request) {
	log.Print("reading records from database")
	rows, err := db.Query("SELECT * FROM employees")
	if err != nil {
		log.Print("error ocurred getting employees :: ", err)
		return
	}
	employees := []employee{}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Print("error ocurred parsing employees :: ", err)
			return
		}
		emp := employee{ID: id, Name: name}
		employees = append(employees, emp)
	}
	json.NewEncoder(w).Encode(employees)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employees", readRecords).Methods("GET")
	defer db.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting the server :: ", err)
		return
	}
}
