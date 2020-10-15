package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/websocket"
)

func TestWebSocketServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleClients))
	defer server.Close()
	u := "ws" + strings.TrimPrefix(server.URL, "http")
	socket, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer socket.Close()
	m := message{Message: "hello"}
	if err := socket.WriteJSON(&m); err != nil {
		t.Fatalf("%v", err)
	}
	var msg message
	err = socket.ReadJSON(&msg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	assert.Equal(t, "hello", msg.Message, "they should be equal")
}
