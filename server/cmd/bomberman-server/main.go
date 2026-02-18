package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/N3moAhead/bomberman/server/internal/client"
	"github.com/N3moAhead/bomberman/server/internal/hub"
	"github.com/N3moAhead/bomberman/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8038", "http service address")

var l = logger.New("[Live-Server]")

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		l.Info("Checking origin: %s", r.Header.Get("Origin"))
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
			l.Error("WebSocket upgrade error: %v", err)
			return
		}
		l.Success("Client connected from: %s", conn.RemoteAddr())

		client := client.NewClient(hubInstance, conn, uuid.NewString())
		hubInstance.Register <- client
		client.StartPumps()
	})

	// Simple handler for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("[BOMBERMAN-SERVER] is running. Connect via WebSocket on /ws"))
		if err != nil {
			l.Errorln("Failed to write status message")
		}
	})

	l.Info("Bomberman-Server starting on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
