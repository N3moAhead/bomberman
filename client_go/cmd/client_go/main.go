package main

import (
	"net/url"

	"github.com/N3moAhead/bomberman/client_go/pkg/bomber"
)

type Bot struct{}

// The code for your own bot
func (b *Bot) CalcNextMove(state bomber.ClassicStatePayload) bomber.PlayerMove {
	// Currently a pretty lazy player :(
	return bomber.PLACE_BOMB
}

var _ bomber.BomberBot = (*Bot)(nil)

func main() {
	newBot := &Bot{}
	b := bomber.NewBomber(newBot)
	b.Start(url.URL{Scheme: "ws", Host: "localhost:8038", Path: "/ws"})
}
