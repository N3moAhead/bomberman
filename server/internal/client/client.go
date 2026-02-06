package client

import (
	"encoding/json"
	"log"
	"time"

	"github.com/N3moAhead/bomberman/server/internal/game"
	"github.com/N3moAhead/bomberman/server/internal/hub"
	"github.com/N3moAhead/bomberman/server/internal/message"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

// Client is a WebSocket client that communicates with a Hub
type Client struct {
	Hub       hub.HubConnection
	Conn      *websocket.Conn
	Send      chan []byte
	ID        string
	Score     int
	isReady   bool
	gameID    string
	authToken string // Is just important for async bot games and the one shot hub
}

// Assure Client implements the interface from the hub package
var _ hub.Client = (*Client)(nil)

// NewClient creates a new client instance
func NewClient(hub hub.HubConnection, conn *websocket.Conn, id string) *Client {
	return &Client{
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		ID:      id,
		isReady: false,
	}
}

// StartPumps starts the client's read and write pumps in separate goroutines
func (c *Client) StartPumps() {
	go c.WritePump()
	go c.ReadPump()
}

// GetID returns the client's unique identifier
func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) SetAuthToken(authToken string) {
	c.authToken = authToken
}

func (c *Client) GetAuthToken() string {
	return c.authToken
}

// IsReady indicates if the client is ready to start a game
func (c *Client) IsReady() bool {
	return c.isReady
}

// SetReady sets the client's ready status
func (c *Client) SetReady(ready bool) {
	c.isReady = ready
}

// GetScore returns the client's current score
func (c *Client) GetScore() int {
	return c.Score
}

// IncrementScore adds a value to the client's score
func (c *Client) IncrementScore(delta int) {
	c.Score += delta
}

// SetGameID sets the ID of the game the client is currently in
func (c *Client) SetGameID(id string) {
	c.gameID = id
}

// Close closes the client's send channel. The connection itself is closed
// by the read/write pumps when they exit
func (c *Client) Close() {
	close(c.Send)
}

// SendMessage formats and sends a structured message to the client
func (c *Client) SendMessage(msgType message.MessageType, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload for client %s: %v", c.ID, err)
		return err
	}
	message := message.Message{
		Type:    msgType,
		Payload: json.RawMessage(payloadBytes),
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message for client %s: %v", c.ID, err)
		return err
	}

	select {
	case c.Send <- messageBytes:
	default:
		log.Printf("Client %s send buffer full. Dropping message.", c.ID)
	}
	return nil
}

var _ game.Player = (*Client)(nil)

// ReadPump transfers messages from the WebSocket to the Hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
		log.Printf("Client %s disconnected (readPump closed)", c.ID)
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message for client %s: %v", c.ID, err)
			}
			break
		}

		var msg message.Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("error unmarshalling message from client %s: %v", c.ID, err)
			continue
		}

		c.Hub.HandleIncomingMessage(c, msg)
	}
}

// WritePump transfers messages from the Hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		log.Printf("Client %s writePump closed", c.ID)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Printf("Client %s send channel closed by hub", c.ID)
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error writing message to client %s: %v", c.ID, err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error sending ping to client %s: %v", c.ID, err)
				return
			}
		}
	}
}
