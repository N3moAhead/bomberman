package server

import (
	"log"
	"net/http"

	"github.com/N3moAhead/bomberman/server/internal/hub"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Checking origin: %s", r.Header.Get("Origin"))
		// TODO: Implement proper origin check for security
		return true
	},
}

func serveWs(hubInstance *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	log.Println("Client connected from:", conn.RemoteAddr())

	client := &hub.Client{
		Hub:          hubInstance,
		Conn:         conn,
		Send:         make(chan []byte, 256), // Use a buffered channel
		Id:           uuid.New().String(),
		Score:        0,
		SelectedGame: "",
	}

	client.Hub.Register <- client // Use the Register channel from the hub instance

	go client.WritePump()
	go client.ReadPump()
}
