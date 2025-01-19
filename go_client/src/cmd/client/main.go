package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type Player struct {
	Name string `json:"name"`
}

func main() {
	var self Player = Player{Name: "Go Client"}

	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	panicOnError(err)

	conn, err := net.DialUDP("udp", nil, serverAddr)
	panicOnError(err)
	defer conn.Close()

	jsonData, err := json.Marshal(self)
	panicOnError(err)

	_, err = conn.Write(jsonData)
	panicOnError(err)

	fmt.Println("JSON data sent:", string(jsonData))

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	panicOnError(err)

	fmt.Println("Received answer from server:", string(buffer[:n]))
}
