package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/N3moAhead/bomberman/server/internal/hub"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8038", "http service address")

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Checking origin: %s", r.Header.Get("Origin"))
		// TODO: Implement proper origin check for security
		return true
	},
}

func main() {
	flag.Parse()

	hubInstance := hub.NewHub()
	go hubInstance.Run()

	// Register the WebSocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Pass the single hub instance to the handler
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		log.Println("Client connected from:", conn.RemoteAddr())

		client := &hub.Client{
			Hub:     hubInstance,
			Conn:    conn,
			Send:    make(chan []byte, 1024), // Use a buffered channel
			ID:      uuid.New().String(),
			Score:   0,
			IsReady: false,
		}

		client.Hub.Register <- client // Use the Register channel from the hub instance

		go client.WritePump()
		go client.ReadPump()
	})

	// Simple handler for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[BOMBERMAN-SERVER] is running. Connect via WebSocket on /ws"))
	})

	log.Printf("Bomberman-Server starting on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
