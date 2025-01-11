package main

import (
	"fmt"
	"net"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
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

		fmt.Printf("Empfangen von %s: %s\n", clientAddr, string(buffer[:n]))

		_, err = conn.WriteToUDP([]byte("Message received!"), clientAddr)
		if err != nil {
			fmt.Println("Error while answering:", err)
		}
	}
}
