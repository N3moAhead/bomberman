package hub

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/N3moAhead/bomberman/server/internal/game"
	"github.com/N3moAhead/bomberman/server/internal/game/classic"
	"github.com/N3moAhead/bomberman/server/internal/message"
	"github.com/google/uuid"
)

// Client defines the interface for a client connecting to the hub.
// This allows the hub to manage clients without depending on a concrete implementation.
type Client interface {
	game.Player // Embeds GetID() and SendMessage()
	IsReady() bool
	SetReady(bool)
	GetScore() int
	IncrementScore(delta int)
	SetGameID(id string)
	Close()
	StartPumps()
	SetAuthToken(authToken string)
}

type hubMessage struct {
	client  Client
	message message.Message
}

// GameDefinition can later be used for diffrent versions
// or variants for bomberman games...
type GameDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Hub struct {
	clients        map[Client]bool
	incoming       chan hubMessage
	Register       chan Client
	unregister     chan Client
	activeGames    map[string]game.Game
	availableGames []message.GameInfo
	clientToGame   map[Client]string
	gameMutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		incoming:   make(chan hubMessage, 2048),
		Register:   make(chan Client),
		unregister: make(chan Client),
		availableGames: []message.GameInfo{
			{Name: "Classic", Description: "The classic and simple bomberman game!"},
		},
		clients:      make(map[Client]bool),
		activeGames:  make(map[string]game.Game),
		clientToGame: make(map[Client]string),
	}
}

// Unregister allows a client to request to be unregistered from the hub.
func (h *Hub) UnregisterClient(c Client) {
	h.unregister <- c
}

// HandleIncomingMessage is called by clients to pass a message to the hub for processing.
func (h *Hub) HandleIncomingMessage(c Client, msg message.Message) {
	h.incoming <- hubMessage{client: c, message: msg}
}

func (h *Hub) Run() {
	log.Info("Hub is running...")
	for {
		select {
		case client := <-h.Register:
			h.gameMutex.Lock()
			h.clients[client] = true
			h.gameMutex.Unlock()
			log.Info("Client %s registered. Total clients: %d", client.GetID(), len(h.clients))
			welcomePayload := message.WelcomeMessage{
				ClientID:     client.GetID(),
				CurrentGames: h.availableGames,
			}
			client.SendMessage(message.Welcome, welcomePayload)
			h.broadcastLobbyUpdate()

		case client := <-h.unregister:
			h.gameMutex.Lock()
			if _, ok := h.clients[client]; ok {
				gameID, inGame := h.clientToGame[client]
				if inGame {
					if activeGame, gameExists := h.activeGames[gameID]; gameExists {
						activeGame.RemovePlayer(client)
						log.Info("Removed client %s from game %s", client.GetID(), activeGame.GetID())
					}
					delete(h.clientToGame, client)
				}
				delete(h.clients, client)
				client.Close()
				log.Warn("Client %s unregistered. Total clients: %d", client.GetID(), len(h.clients))
			}
			h.gameMutex.Unlock()
			h.broadcastLobbyUpdate()
			h.checkAndPotentiallyStartGame()

		case hubMsg := <-h.incoming:
			h.gameMutex.RLock()
			gameID, inGame := h.clientToGame[hubMsg.client]
			h.gameMutex.RUnlock()

			if inGame {
				h.gameMutex.RLock()
				currentGame, gameExists := h.activeGames[gameID]
				h.gameMutex.RUnlock()

				if gameExists {
					currentGame.HandleMessage(hubMsg.client, hubMsg.message)
				} else {
					log.Error("Client %s mapped to game %s, but game does not exist.", hubMsg.client.GetID(), gameID)
					h.gameMutex.Lock()
					delete(h.clientToGame, hubMsg.client)
					h.gameMutex.Unlock()
				}
			} else {
				h.handleLobbyMessage(hubMsg.client, hubMsg.message)
			}
		}
	}
}

func (h *Hub) handleLobbyMessage(client Client, msg message.Message) {
	switch msg.Type {
	case message.PlayerStatusUpdate:
		var payload message.PlayerStatusUpdatePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Error("Error unmarshalling select_game payload from %s: %v", client.GetID(), err)
			client.SendMessage(message.Error, message.ErrorMessage{Message: "Invalid PlayerStatusUpdatePayload payload"})
			return
		}

		client.SetReady(payload.IsReady)
		h.broadcastLobbyUpdate()
		h.checkAndPotentiallyStartGame()
	default:
		log.Warn("Received unhandled lobby message type '%s' from client %s", msg.Type, client.GetID())
	}
}

func (h *Hub) selectAndStartGame() {
	h.gameMutex.Lock()

	gameInfo := h.availableGames[0]
	gameID := uuid.New().String()
	newGame := classic.NewClassic(h, gameID)
	h.activeGames[gameID] = newGame

	clientsInLobby := []Client{}
	for client := range h.clients {
		if _, inGame := h.clientToGame[client]; !inGame {
			clientsInLobby = append(clientsInLobby, client)
		}
	}

	clientsReady := []Client{}
	for _, client := range clientsInLobby {
		if client.IsReady() {
			clientsReady = append(clientsReady, client)
		}
	}

	if len(clientsReady) < 2 {
		log.Warn("Not enough players are ready and available to start a new game")
		h.gameMutex.Unlock()
		h.broadcastLobbyUpdate()
		return
	}

	if len(clientsReady) < len(clientsInLobby) {
		log.Info("Some players are in the lobby but still not ready we are going to wait for them")
		h.gameMutex.Unlock()
		h.broadcastLobbyUpdate()
		return
	}

	for _, client := range clientsInLobby {
		h.clientToGame[client] = gameID
		err := newGame.AddPlayer(client)
		if err != nil {
			log.Error("Error while trying to add player to game %s; %v", client.GetID(), err)
			delete(h.clientToGame, client)
		} else {
			client.SetGameID(gameID)
			client.SetReady(false)
			startPayload := message.GameStartPayload{Name: gameInfo.Name, Description: gameInfo.Description, GameID: gameID}
			client.SendMessage(message.GameStart, startPayload)
			log.Success("Added player %s to game %s", client.GetID(), h.availableGames[0].Name)
		}
	}

	go newGame.Start()
	log.Success("Started game %s (%s) in a new goroutine", gameInfo.Name, gameID)

	h.gameMutex.Unlock()

	h.broadcastLobbyUpdate()
}

func (h *Hub) GameFinished(gameID string, result game.GameResult) {
	h.gameMutex.Lock()
	defer h.gameMutex.Unlock()

	log.Info("Game %s finished. Processing results.", gameID)

	if _, exists := h.activeGames[gameID]; exists {
		delete(h.activeGames, gameID)
	} else {
		log.Warn("GameFinished called for non-existent or already finished game %s", gameID)
		return
	}

	clientsToRemove := []Client{}
	for client, gid := range h.clientToGame {
		if gid == gameID {
			clientsToRemove = append(clientsToRemove, client)
		}
	}

	for _, client := range clientsToRemove {
		delete(h.clientToGame, client)
		client.SetGameID("")
		client.SendMessage(message.BackToLobby, nil)
		log.Info("Client %s removed from finished game %s, returned to lobby.", client.GetID(), gameID)
	}

	if len(result.Scores) > 0 {
		h.updateScoresInternal(result.Scores)
	}

	if result.Winner != "" {
		log.Info("The winner of the game is %s", result.Winner)
	}

	// Using a goroutine to avoid blocking and potential deadlocks
	go func() {
		h.broadcastLobbyUpdate()
		time.AfterFunc(500*time.Millisecond, h.checkAndPotentiallyStartGame)
	}()
}

func (h *Hub) broadcastLobbyUpdate() {
	playerInfos := make(map[string]message.PlayerInfo)
	h.gameMutex.RLock()
	for client := range h.clients {
		_, inGame := h.clientToGame[client]
		playerInfos[client.GetID()] = message.PlayerInfo{
			InGame:  inGame,
			IsReady: client.IsReady(),
			Score:   client.GetScore(),
		}
	}
	h.gameMutex.RUnlock()
	payload := message.LobbyUpdateMessage{Players: playerInfos}

	h.broadcastMessageInternal(message.UpdateLobby, payload)
}

func (h *Hub) broadcastMessageInternal(msgType message.MessageType, payload any) {
	h.gameMutex.RLock()
	log.Info("Broadcasting message type '%s' to %d clients", msgType, len(h.clients))
	clientList := make([]Client, 0, len(h.clients))
	for client := range h.clients {
		clientList = append(clientList, client)
	}
	h.gameMutex.RUnlock()

	for _, client := range clientList {
		err := client.SendMessage(msgType, payload)
		if err != nil {
			log.Error("Error broadcasting message type %s to client %s: %v", msgType, client.GetID(), err)
		}
	}
}

func (h *Hub) checkAndPotentiallyStartGame() {
	h.gameMutex.RLock()
	lobbyClientsCount := 0
	for c := range h.clients {
		if _, inGame := h.clientToGame[c]; !inGame {
			lobbyClientsCount++
		}
	}
	canStart := lobbyClientsCount > 1
	h.gameMutex.RUnlock()

	if canStart {
		log.Info("Enough players %d/2 in lobby. Trying to start a game...", lobbyClientsCount)
		h.selectAndStartGame()
	} else {
		log.Warn("Not Enough players %d/2 in lobby. Waiting for more players to join...", lobbyClientsCount)
	}
}

func (h *Hub) updateScoresInternal(scores map[string]int) {
	for clientID, delta := range scores {
		var targetClient Client
		for c := range h.clients {
			if c.GetID() == clientID {
				targetClient = c
				break
			}
		}
		if targetClient != nil {
			targetClient.IncrementScore(delta)
			log.Info("Score updated for %s: new score %d", targetClient.GetID(), targetClient.GetScore())
		} else {
			log.Warn("Could not find client %s to update score", clientID)
		}
	}
}

type GameFinisher interface {
	GameFinished(gameID string, result game.GameResult)
}

// Checking if the hub implements the game finished interface correctly
var _ GameFinisher = (*Hub)(nil)
var _ HubConnection = (*Hub)(nil)
