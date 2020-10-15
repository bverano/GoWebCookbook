package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	Message string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan message)
var upgrader = websocket.Upgrader{}

func broadcastMessageToClients() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error ocurred while writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleClients(w http.ResponseWriter, r *http.Request) {
	go broadcastMessageToClients()
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket :: ", err)
	}
	defer websocket.Close()
	clients[websocket] = true
	for {
		var msg message
		err := websocket.ReadJSON(&msg)
		if err != nil {
			log.Printf("error ocurred while reading message : %v", err)
			delete(clients, websocket)
			break
		}
		broadcast <- msg
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/echo", handleClients)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}
