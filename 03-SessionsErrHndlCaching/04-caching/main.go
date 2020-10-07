package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var newCache *cache.Cache

func init() {
	newCache = cache.New(5*time.Minute, 10*time.Minute)
	newCache.Set("foo", "bar", cache.DefaultExpiration)
}

func getFormCache(w http.ResponseWriter, r *http.Request) {
	foo, found := newCache.Get("foo")
	if found {
		log.Print("Key found in Cache with value as :: ", foo.(string))
		fmt.Fprintf(w, "Hello "+foo.(string))
	} else {
		log.Print("Key not found in Cache :: ", "foo")
		fmt.Fprintf(w, "Key not found in Cache")
	}
}

func main() {
	http.HandleFunc("/", getFormCache)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting the server: ", err)
		return
	}
}
