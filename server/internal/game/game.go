package game

import "github.com/N3moAhead/bomberman/server/internal/message"

// The player struct defines the functions that a game
// awaits from a connected player
type Player interface {
	GetID() string
	SendMessage(msgType message.MessageType, payload any) error
}

// After a game is finished a game result should be returned
// To help us update all the scores
type GameResult struct {
	Winner string
	Scores map[string]int // Map from PlayerID to game scores
}

type GameFinisher interface {
	GameFinished(gameID string, result GameResult)
}

type Game interface {
	Start()                                           // Starts the game
	AddPlayer(player Player) error                    // Adds a new player to the game
	RemovePlayer(player Player)                       // Removes a playser from the game
	HandleMessage(player Player, msg message.Message) // Handles incoming user input
	Stop()                                            // Stops the game
	GetID() string                                    // Returns the game id
}
