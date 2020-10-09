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
		Name:        "getEmployees",
		Method:      "GET",
		Pattern:     "/employees",
		HandlerFunc: getEmployees,
	},
	route{
		Name:        "getEmployee",
		Method:      "GET",
		Pattern:     "/employee/{id}",
		HandlerFunc: getEmployee,
	},
	route{
		Name:        "addEmployee",
		Method:      "POST",
		Pattern:     "/employee/add",
		HandlerFunc: addEmployee,
	},
	route{
		Name:        "updateEmployee",
		Method:      "PUT",
		Pattern:     "/employee/update",
		HandlerFunc: updateEmployee,
	},
	route{
		Name:        "deleteEmployee",
		Method:      "DELETE",
		Pattern:     "/employee/delete",
		HandlerFunc: deleteEmployee,
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
		employee{ID: "1", FirstName: "Bruno", LastName: "Verano"},
		employee{ID: "2", FirstName: "Alvaro", LastName: "Verano"},
		employee{ID: "3", FirstName: "Jacobo", LastName: "Verano"},
	}
}

func getIndex(id string) int {
	for i := 0; i < len(employeesVar); i++ {
		if employeesVar[i].ID == id {
			return i
		}
	}
	return -1
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(employeesVar)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for _, emp := range employeesVar {
		if emp.ID == id {
			if err := json.NewEncoder(w).Encode(emp); err != nil {
				log.Print("error getting requested employee :: ", err)
				return
			}
		}
	}
}

func addEmployee(w http.ResponseWriter, r *http.Request) {
	emp := employee{}
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Print("error decoding employee data :: ", err)
		return
	}
	log.Printf("adding employee %s - %s %s", emp.ID, emp.FirstName, emp.LastName)
	employeesVar = append(employeesVar, employee{ID: emp.ID, FirstName: emp.FirstName, LastName: emp.LastName})
	json.NewEncoder(w).Encode(employeesVar)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	emp := employee{}
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Print("error decoding employee data :: ", err)
		return
	}
	var isUpsert = true
	for idx, e := range employeesVar {
		if e.ID == emp.ID {
			isUpsert = false
			log.Printf("updating employee %s - %s %s", emp.ID, emp.FirstName, emp.LastName)
			employeesVar[idx].FirstName = emp.FirstName
			employeesVar[idx].LastName = emp.LastName
			break
		}
	}
	if isUpsert {
		log.Printf("updaserting employee %s - %s %s", emp.ID, emp.FirstName, emp.LastName)
		employeesVar = append(employeesVar, employee{ID: emp.ID, FirstName: emp.FirstName, LastName: emp.LastName})
	}
	json.NewEncoder(w).Encode(employeesVar)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	emp := employee{}
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Print("error decoding employee data :: ", err)
		return
	}
	log.Printf("deleting employee %s - %s %s", emp.ID, emp.FirstName, emp.LastName)
	index := getIndex(emp.ID)
	if index == -1 {
		log.Print("employee not found")
		return
	}
	employeesVar = append(employeesVar[:index], employeesVar[index+1:]...)
	json.NewEncoder(w).Encode(employeesVar)
}

func addRoutes(router *mux.Router) *mux.Router {
	for _, r := range routesVar {
		router.Name(r.Name).Methods(r.Method).Path(r.Pattern).Handler(r.HandlerFunc)
	}
	return router
}

func main() {
	muxRouter := mux.NewRouter()
	router := addRoutes(muxRouter)
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error running the server :: ", err)
		return
	}
}
