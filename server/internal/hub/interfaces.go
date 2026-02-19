package hub

import "github.com/N3moAhead/bombahead/server/internal/message"

// HubConnection is the interface that a Hub must implement to be used by a Client
// It defines the methods a client can use to communicate back to the hub it's connected to
type HubConnection interface {
	UnregisterClient(c Client)
	HandleIncomingMessage(c Client, msg message.Message)
}
