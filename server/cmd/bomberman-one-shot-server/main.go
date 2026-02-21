package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/N3moAhead/bombahead/server/internal/client"
	"github.com/N3moAhead/bombahead/server/internal/hub"
	"github.com/N3moAhead/bombahead/server/pkg/logger"
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
	historyFilePath := os.Getenv("BOMBERMAN_MATCH_HISTORY_PATH")
	if historyFilePath == "" {
		log.Warn("History file path env missing")
	}

	oneShotHub := hub.NewOneShotHub(historyFilePath)
	go oneShotHub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorln(err)
			return
		}

		select {
		case <-oneShotHub.Done:
			log.Warn("Refusing new connection, server is shutting down.")
			conn.Close()
			return
		default:
		}

		client := client.NewClient(oneShotHub, conn, uuid.New().String())
		oneShotHub.Register <- client
		client.StartPumps()
	})

	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	go func() {
		log.Info("One-shot server starting on %s\n", *addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	<-oneShotHub.Done
	log.Info("Hub has shut down, initiating server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Failed:", err)
	}

	log.Success("Server has shut down gracefully.")
}
