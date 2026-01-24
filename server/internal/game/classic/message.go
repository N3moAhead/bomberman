package classic

import "github.com/N3moAhead/bomberman/server/pkg/types"

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
