package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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

func currentDb(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT current_database()")
	if err != nil {
		log.Print("error executing query :: ", err)
		return
	}
	var db string
	for rows.Next() {
		rows.Scan(&db)
	}
	fmt.Fprintf(w, "Current database is :: %s", db)
}

func main() {
	http.HandleFunc("/", currentDb)
	defer db.Close()
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}
