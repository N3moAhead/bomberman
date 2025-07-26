package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const FIVE_SECONDS time.Duration = 5 * time.Second

func main() {
	url := "ws://host.docker.internal:8038/ws"
	log.Printf("Connecting to %s", url)

	c, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Print("Handshake failed with error:")
		if resp != nil {
			log.Printf("HTTP Response Status: %s", resp.Status)
		}
		log.Fatalf("Dial error: %v", err)
	}
	// The client just stays for 5 seconds and closes the connection afterwards
	time.Sleep(FIVE_SECONDS)
	c.Close()
}
