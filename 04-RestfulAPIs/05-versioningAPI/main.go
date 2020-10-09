package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

type route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type routes []route

var routesVar = routes{
	route{
		Name:    "getEmployees",
		Method:  "GET",
		Pattern: "/employees",
		Handler: getEmployees,
	},
}

type employe struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type employees []employe

var employeesVar employees
var employeesVarV1 employees
var employeesVarV2 employees

func init() {
	employeesVar = employees{
		employe{ID: "1", FirstName: "Bruno", LastName: "Verano"},
	}
	employeesVarV1 = employees{
		employe{ID: "1", FirstName: "Bruno", LastName: "Verano"},
		employe{ID: "2", FirstName: "Jacobo", LastName: "Verano"},
	}
	employeesVarV2 = employees{
		employe{ID: "1", FirstName: "Jacobo", LastName: "Verano"},
		employe{ID: "2", FirstName: "Luis", LastName: "Catañeda"},
	}
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/v1") {
		json.NewEncoder(w).Encode(employeesVarV1)
	} else if strings.HasPrefix(r.URL.Path, "/v2") {
		json.NewEncoder(w).Encode(employeesVarV2)
	} else {
		json.NewEncoder(w).Encode(employeesVar)
	}
}

func addRoutes(router *mux.Router) *mux.Router {
	for _, route := range routesVar {
		router.Name(route.Name).Methods(route.Method).Path(route.Pattern).Handler(route.Handler)
	}
	return router
}

func main() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	router := addRoutes(muxRouter)
	//v1
	addRoutes(muxRouter.PathPrefix("/v1").Subrouter())
	//v2
	addRoutes(muxRouter.PathPrefix("/v2").Subrouter())
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error running the server :: ", err)
		return
	}
}
