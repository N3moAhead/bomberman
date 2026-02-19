package classic

import "github.com/N3moAhead/bombahead/server/pkg/types"

type PlayerMove string

const (
	NO_INPUT_DEFINED            = "undefined"
	DO_NOTHING       PlayerMove = "nothing" // Do nothing
	MOVE_UP          PlayerMove = "move_up"
	MOVE_RIGHT       PlayerMove = "move_right"
	MOVE_DOWN        PlayerMove = "move_down"
	MOVE_LEFT        PlayerMove = "move_left"
	PLACE_BOMB       PlayerMove = "place_bomb"
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

type PlayerHistoryEntry struct {
	PlayerState
	Move      PlayerMove `json:"move"`
	AuthToken string     `json:"authToken"`
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

// TickState represents the dynamic state of the game at a single tick for history purposes
type TickState struct {
	Players        []PlayerHistoryEntry `json:"players"`
	Bombs          []BombState          `json:"bombs"`
	Explosions     []types.Vec2         `json:"explosions"`
	DestroyedBoxes []types.Vec2         `json:"destroyed_boxes,omitempty"`
}

// GameHistory encapsulates the entire history of a game, with an initial field state
// and a sequence of state changes for each tick
type GameHistory struct {
	InitialField    FieldState  `json:"initial_field"`
	Ticks           []TickState `json:"ticks"`
	WinnerAuthToken string      `json:"winnerAuthToken"`
}
