package main

import (
	"flag"

	"github.com/N3moAhead/bomberman/server/internal/server"
)

var addr = flag.String("addr", ":8038", "http service address")

func main() {
	flag.Parse()
	server.Run(addr)
}
