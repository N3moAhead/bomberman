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
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	panicOnError(err)

	conn, err := net.ListenUDP("udp", addr)
	panicOnError(err)

	fmt.Println("Bomberman-Server is running on Port 8080...")

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Printf("Received raw message from %s: %s\n", clientAddr, string(buffer[:n]))

		var player Player
		err = json.Unmarshal(buffer[:n], &player)
		if err != nil {
			fmt.Printf("Failed to decode JSON from %s: %v\n", clientAddr, err)
			_, _ = conn.WriteToUDP([]byte("Invalid JSON format!"), clientAddr)
			continue
		}

		fmt.Printf("Player received from %s: %+v\n", clientAddr, player)

		response := fmt.Sprintf("Hello, %s! Welcome to Bomberman!", player.Name)
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Println("Error while answering:", err)
		}
	}
}
