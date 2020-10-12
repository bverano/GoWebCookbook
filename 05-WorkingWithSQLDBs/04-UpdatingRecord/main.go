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

func updateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("going to update record in database for id :: ", id)
		stmt, err := db.Prepare("UPDATE employees SET name = $1 WHERE id = $2")
		if err != nil {
			log.Print("error ocurred while preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0], id)
		if err != nil {
			log.Print("error ocurred executing query :: ", err)
			return
		}
		rowsAffected, _ := result.RowsAffected()
		fmt.Fprintf(w, "Number of rows updated in Database are :: %d", rowsAffected)
	} else {
		fmt.Fprintf(w, "Error ocurred while updating record in database for id :: %s", id)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employee/update/{id}", updateRecord).Methods("PUT")
	defer db.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting the server :: ", err)
		return
	}
}
