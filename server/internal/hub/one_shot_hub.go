package hub

import (
	"encoding/json"
	"sync"

	"github.com/N3moAhead/bomberman/server/internal/game"
	"github.com/N3moAhead/bomberman/server/internal/game/classic"
	"github.com/N3moAhead/bomberman/server/internal/message"
	"github.com/N3moAhead/bomberman/server/pkg/logger"
	"github.com/google/uuid"
)

var log = logger.New("[HUB]")

// OneShotHub is a hub that waits for two players,
// runs one game, and then shuts down
type OneShotHub struct {
	clients    map[Client]bool
	Register   chan Client
	unregister chan Client
	incoming   chan hubMessage
	game       game.Game
	gameMutex  sync.Mutex
	shutdown   chan struct{}
	Done       chan struct{}
}

// NewOneShotHub creates a new OneShotHub
func NewOneShotHub() *OneShotHub {
	return &OneShotHub{
		clients:    make(map[Client]bool),
		Register:   make(chan Client),
		unregister: make(chan Client),
		incoming:   make(chan hubMessage),
		shutdown:   make(chan struct{}),
		Done:       make(chan struct{}),
	}
}

// Run starts the hubs main loop
func (h *OneShotHub) Run() {
	defer close(h.Done)
	log.Info("Is running, waiting for 2 players...")
	gameStarted := false

	for {
		select {
		case client := <-h.Register:
			if len(h.clients) < 2 {
				h.clients[client] = true
				log.Info("Client %s registered. Total clients: %d/2", client.GetID(), len(h.clients))
				welcomePayload := message.WelcomeMessage{ClientID: client.GetID()}
				client.SendMessage(message.Welcome, welcomePayload)
			}
			if len(h.clients) == 2 && !gameStarted {
				log.Info("Two players connected, starting game...")
				h.startGame()
				gameStarted = true
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
				log.Warn("Client %s unregistered.", client.GetID())
				// If a client disconnects, the game should end, which will trigger shutdown
				// The game logic itself handles player disconnection
			}
		case hubMsg := <-h.incoming:
			h.gameMutex.Lock()
			if hubMsg.message.Type == message.PlayerStatusUpdate {
				var payload message.PlayerStatusUpdatePayload
				if err := json.Unmarshal(hubMsg.message.Payload, &payload); err != nil {
					log.Error("Error while unmarshalling PlayerStatusUpdate %v\n", err)
					return
				}
				hubMsg.client.SetAuthToken(payload.AuthToken)
				log.Success("Set auth Token %s for client %s", payload.AuthToken, hubMsg.client.GetID())
			} else {
				if h.game != nil {
					h.game.HandleMessage(hubMsg.client, hubMsg.message)
				} else {
					log.Warn("Message from %s received before game start, ignoring.", hubMsg.client.GetID())
				}
			}
			h.gameMutex.Unlock()
		case <-h.shutdown:
			log.Info("Game finished, OneShotHub is shutting down.")
			h.gameMutex.Lock()
			for client := range h.clients {
				client.Close()
			}
			h.clients = make(map[Client]bool) // Clear clients
			h.gameMutex.Unlock()
			return // Exit Run loop
		}
	}
}

func (h *OneShotHub) startGame() {
	h.gameMutex.Lock()
	defer h.gameMutex.Unlock()

	gameID := uuid.New().String()
	// The OneShotHub implements GameFinisher, so we pass 'h'
	newGame := classic.NewClassic(h, gameID)
	h.game = newGame

	for client := range h.clients {
		err := h.game.AddPlayer(client)
		if err != nil {
			log.Error("Error adding player %s to game: %v", client.GetID(), err)
		} else {
			log.Info("Added player %s to game %s", client.GetID(), gameID)
			startPayload := message.GameStartPayload{Name: "Classic (One-Shot)", Description: "The classic bomberman game, one-shot style!", GameID: gameID}
			client.SendMessage(message.GameStart, startPayload)
		}
	}
	go h.game.Start()
	log.Success("Started game %s in a new goroutine", gameID)
}

// --- Interface Implementations ---

// UnregisterClient implements the HubConnection interface.
func (h *OneShotHub) UnregisterClient(c Client) {
	h.unregister <- c
}

// HandleIncomingMessage implements the HubConnection interface.
func (h *OneShotHub) HandleIncomingMessage(c Client, msg message.Message) {
	h.incoming <- hubMessage{client: c, message: msg}
}

// GameFinished implements the GameFinisher interface.
func (h *OneShotHub) GameFinished(gameID string, result game.GameResult) {
	log.Success("Game %s finished in OneShotHub. Result: %+v. Signalling shutdown.", gameID, result)
	close(h.shutdown)
}

// Statically assert that OneShotHub implements the necessary interfaces
var _ HubConnection = (*OneShotHub)(nil)
var _ GameFinisher = (*OneShotHub)(nil)
