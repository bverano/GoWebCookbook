package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var routesVar = routes{
	route{
		"getEmployees",
		"GET",
		"/employees",
		getEmployees,
	},
	route{
		"getEmployee",
		"GET",
		"/employee/{id}",
		getEmployee,
	},
}

type employee struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type employees []employee

var employeesVar employees

func init() {
	employeesVar = employees{
		employee{
			ID: "1", FirstName: "Bruno", LastName: "Verano",
		},
		employee{
			ID: "2", FirstName: "Jacobo", LastName: "Verano",
		},
	}
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(employeesVar)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for _, employee := range employeesVar {
		if employee.ID == id {
			if err := json.NewEncoder(w).Encode(employee); err != nil {
				log.Print("error getting requested employee: ", err)
			}
		}
	}
}

func addRoutes(router *mux.Router) *mux.Router {
	for _, route := range routesVar {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	return router
}

func main() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	router := addRoutes(muxRouter)
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting the server: ", err)
		return
	}
}
