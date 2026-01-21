package classic

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/N3moAhead/bomberman/server/internal/game"
	"github.com/N3moAhead/bomberman/server/internal/message"
	"github.com/N3moAhead/bomberman/server/pkg/types"
)

type Classic struct {
	gameFinisher game.GameFinisher
	stopChan     chan bool

	gameID    string
	players   map[string]*Player     // ClientID -> ClassicPlayer
	playerMap map[string]game.Player // ClientID -> game.Player
	playerMux sync.RWMutex

	field *Field

	isRunning  bool
	minPlayers int
	maxPLayers int

	ticker       *time.Ticker
	lastTickTime time.Time // for delta time
}

func NewClassic(finisher game.GameFinisher, id string) *Classic {
	return &Classic{
		gameFinisher: finisher,
		stopChan:     make(chan bool),

		gameID:    id,
		players:   make(map[string]*Player),
		playerMap: make(map[string]game.Player),

		field: NewField(),

		isRunning:  false,
		minPlayers: min_players,
		maxPLayers: max_players,
	}
}

func (c *Classic) GetID() string {
	return c.gameID
}

func (c *Classic) AddPlayer(player game.Player) error {
	c.playerMux.Lock()
	defer c.playerMux.Unlock()

	if len(c.players) >= c.maxPLayers {
		return fmt.Errorf(
			"[Game %s] Already full, can't add more players (%d/%d)!\n",
			c.gameID,
			len(c.players),
			c.maxPLayers,
		)
	}

	playerID := player.GetID()
	if _, exists := c.players[playerID]; exists {
		return fmt.Errorf("[Game %s] Player %s already exists.\n", c.gameID, playerID)
	}

	// TODO very the spawn position currently each player is getting
	// spawned in the same field
	spawnPos := types.NewVec2(1, 1)
	newPlayer := &Player{
		ID:       playerID,
		Pos:      spawnPos,
		Score:    0,
		Health:   initial_health,
		NextMove: DO_NOTHING,
	}
	c.players[playerID] = newPlayer
	c.playerMap[playerID] = player

	log.Printf("[Game %s] Player %s added.\n", c.gameID, playerID)
	return nil
}

func (c *Classic) RemovePlayer(player game.Player) {
	c.playerMux.Lock()
	defer c.playerMux.Unlock()

	playerID := player.GetID()
	if _, ok := c.players[playerID]; ok {
		delete(c.players, playerID)
		delete(c.playerMap, playerID)
		log.Printf("[Game %s] Player %s removed.\n", c.gameID, playerID)

		if len(c.players) < c.minPlayers && c.isRunning {
			log.Printf(
				"[Game %s] Not enough players remaining (%d/%d). Stopping game.\n",
				c.gameID,
				len(c.players),
				c.minPlayers,
			)
			go c.Stop()
		}
	}
}

func (c *Classic) Start() {
	c.playerMux.Lock()
	if len(c.players) < c.minPlayers {
		c.playerMux.Unlock()
		log.Printf("[Game %s] Cannot start, not enough players (%d/%d).", c.gameID, len(c.players), c.minPlayers)
		c.Stop()
	}

	c.isRunning = true
	c.lastTickTime = time.Now()
	c.ticker = time.NewTicker(TICK_RATE)

	c.playerMux.Unlock()

	log.Printf("[Game %s] Starting game loop.", c.gameID)
	defer func() {
		if c.ticker != nil {
			c.ticker.Stop()
		}
		log.Printf("[Game %s] Game loop stopped.", c.gameID)
	}()

	for {
		select {
		case <-c.ticker.C:
			// Update game state
			// Send game state to all players
			return
		case <-c.stopChan:
			// When receiving a stop signal we stop the goroutine
			return
		}
	}
}

func (c *Classic) Stop() {
	c.playerMux.Lock()
	if !c.isRunning {
		c.playerMux.Unlock()
		return
	}
	c.isRunning = false

	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}

	// closing the stop channel
	select {
	case <-c.stopChan:
	default:
		close(c.stopChan)
	}

	result := game.GameResult{
		Scores: make(map[string]int),
	}

	playersSnapshot := make([]game.Player, 0, len(c.playerMap))
	for _, p := range c.playerMap {
		playersSnapshot = append(playersSnapshot, p)
	}
	c.playerMux.Unlock()

	log.Printf("[Game %s] Stopping game.", c.gameID)

	// Inform the hub that the game is finished and retrieve all
	// players back to the lobby
	c.gameFinisher.GameFinished(c.gameID, result)
}

func (c *Classic) HandleMessage(player game.Player, msg message.Message) {
	playerID := player.GetID()

	switch msg.Type {
	case message.ClassicInput:
		var payload ClassicInputPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Printf("[Game %s] Error unmarshalling ClassInput from %s: %v\n", c.gameID, playerID, err)
			return
		}

		c.playerMux.Lock()
		pState, ok := c.players[playerID]
		if ok {
			pState.HandleInput(payload)
		} else {
			log.Printf(
				"[Game %s] Received input from player %s who is not in the internal state map.",
				c.gameID,
				playerID,
			)
		}
	default:
		log.Printf(
			"[Game %s] Received unhandled message type '%s' from player %s",
			c.gameID,
			msg.Type,
			playerID,
		)
	}
}

var _ game.Game = (*Classic)(nil)
