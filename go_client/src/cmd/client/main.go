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
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	panicOnError(err)

	conn, err := net.DialUDP("udp", nil, serverAddr)
	panicOnError(err)
	defer conn.Close()

	message := "Hey im the client nice to meet you ^^"
	_, err = conn.Write([]byte(message))
	panicOnError(err)

	fmt.Println("Message send:", message)

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	panicOnError(err)

	fmt.Println("Received answer from server:", string(buffer[:n]))
}
