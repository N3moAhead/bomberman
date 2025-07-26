package message

import (
	"encoding/json"
)

// Message represents a generic message that is sent over WebSocket.
// The 'Type' helps the server or client understand how to interpret the 'Payload'.
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MessageType string

const (
	// Message types for the WebSocket communication
	Welcome     MessageType = "welcome"       // Sent when a client connects
	BackToLobby MessageType = "back_to_lobby" // Send when a player returns from a game back to the lobby
	UpdateLobby MessageType = "update_lobby"  // Sent to update the lobby state
)

type GameInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WelcomeMessage contains the ID of the new client and the list of available games
type WelcomeMessage struct {
	ClientID     string     `json:"clientId"`
	CurrentGames []GameInfo `json:"currentGames"`
}

type PlayerInfo struct {
	InGame bool `json:"inGame"`
}

// LobbyUpdateMessage contains the current state of the lobby
type LobbyUpdateMessage struct {
	Players map[string]PlayerInfo `json:"players"` // Maps client id to client infos
}

// ErrorMessage is sent in case of errors
type ErrorMessage struct {
	Message string `json:"message"`
}
