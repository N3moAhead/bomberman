package bomber

import (
	"encoding/json"

	"github.com/N3moAhead/bomberman/client_go/pkg/types"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MessageType string

const (
	Welcome            MessageType = "welcome"
	BackToLobby        MessageType = "back_to_lobby"
	UpdateLobby        MessageType = "update_lobby"
	PlayerStatusUpdate MessageType = "player_status_update"
	Error              MessageType = "error"
	ClassicInput       MessageType = "classic_input"
	ClassicState       MessageType = "classic_state"
	GameStart          MessageType = "game_start"
)

type GameInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WelcomeMessage struct {
	ClientID     string     `json:"clientId"`
	CurrentGames []GameInfo `json:"currentGames"`
}

type PlayerInfo struct {
	InGame  bool `json:"inGame"`
	IsReady bool `json:"isReady"`
	Score   int  `json:"score"`
}

type LobbyUpdateMessage struct {
	Players map[string]PlayerInfo `json:"players"`
}

type PlayerStatusUpdatePayload struct {
	IsReady bool `json:"isReady"`
}

type GameStartPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GameID      string `json:"gameId"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type PlayerMove string

const (
	DO_NOTHING PlayerMove = "nothing" // Do nothing
	MOVE_UP    PlayerMove = "move_up"
	MOVE_RIGHT PlayerMove = "move_right"
	MOVE_DOWN  PlayerMove = "move_down"
	MOVE_LEFT  PlayerMove = "move_left"
	PLACE_BOMB PlayerMove = "place_bomb"
)

type ClassicInputPayload struct {
	Move PlayerMove `json:"move"`
}

type PlayerState struct {
	ID     string     `json:"id"`
	Pos    types.Vec2 `json:"pos"`
	Health int        `json:"health"`
	Score  int        `json:"score"`
}

type FieldState struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Field  []Tile `json:"field"`
}

type BombState struct {
	Pos  types.Vec2 `json:"pos"`
	Fuse int        `json:"fuse"`
}

type ClassicStatePayload struct {
	Players    []PlayerState `json:"players"`
	Field      FieldState    `json:"field"`
	Bombs      []BombState   `json:"bombs"`
	Explosions []types.Vec2  `json:"explosions"`
}
