package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

type nameNotFoundError struct {
	Code int
	Err  error
}

// With this function, the nameNotFoundError type implements the error interface
func (e nameNotFoundError) Error() string {
	return e.Err.Error()
}

//This type is getting the form of a handlerFunc
type wrapperHandler func(http.ResponseWriter, *http.Request) error

/*When implementing ServeHTTP as a method for the wrapperHandler type, we are implementing the Handler Interface, which may be use for responding a HTTP Request:
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request)*/
func (wrapper wrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := wrapper(w, r)
	if err != nil {
		switch e := err.(type) {
		case nameNotFoundError:
			log.Printf("HTTP %s - %d", e.Err, e.Code)
			http.Error(w, e.Err.Error(), e.Code)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func getName(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	name := vars["name"]
	if strings.EqualFold(name, "foo") {
		fmt.Fprintf(w, "Hello, "+name)
		return nil
	}
	return nameNotFoundError{404, errors.New("Name not Found")}
}

func main() {
	router := mux.NewRouter()
	router.Handle("/employee/get/{name}", wrapperHandler(getName)).Methods("GET")
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting the server: ", err)
		return
	}
}
