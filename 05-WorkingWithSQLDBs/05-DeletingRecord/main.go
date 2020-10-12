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

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("going to delete record in database for name :: ", name[0])
		stmt, err := db.Prepare("DELETE FROM employees WHERE name = $1")
		if err != nil {
			log.Print("error preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0])
		if err != nil {
			log.Print("error executing query :: ", err)
			return
		}
		rowsAffected, _ := result.RowsAffected()
		fmt.Fprintf(w, "Number of rows affected in database are :: %d", rowsAffected)
	} else {
		fmt.Fprintf(w, "Error ocurred while deleting record in database for name %s", name[0])
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employee/delete", deleteRecord).Methods("DELETE")
	defer db.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error running server :: ", err)
		return
	}
}
