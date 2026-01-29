package main

import (
	"flag"
	"net/http"

	"github.com/N3moAhead/bomberman/server/internal/client"
	"github.com/N3moAhead/bomberman/server/internal/hub"
	"github.com/N3moAhead/bomberman/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8038", "http service address")
var log = logger.New("[OS-Server]")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	oneShotHub := hub.NewOneShotHub()
	go oneShotHub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorln(err)
			return
		}

		client := client.NewClient(oneShotHub, conn, uuid.New().String())
		oneShotHub.Register <- client
		client.StartPumps()
	})

	log.Info("One-shot server starting on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
