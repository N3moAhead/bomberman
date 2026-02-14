package message

// Vec2 represents a 2-dimensional vector
type Vec2 struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Tile represents the type of a tile on the game field
type Tile string

// PlayerMove represents a move a player can make
type PlayerMove string

// PlayerState represents the state of a player at a certain point
type PlayerState struct {
	ID     string `json:"id"`
	Pos    Vec2   `json:"pos"`
	Health int    `json:"health"`
	Score  int    `json:"score"`
}

// PlayerHistoryEntry represents the state of a player and their move for a single tick
type PlayerHistoryEntry struct {
	PlayerState
	Move      PlayerMove `json:"move"`
	AuthToken string     `json:"authToken"`
}

// FieldState represents the static state of the game field
type FieldState struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Field  []Tile `json:"field"`
}

// BombState represents the state of a bomb
type BombState struct {
	Pos  Vec2 `json:"pos"`
	Fuse int  `json:"fuse"`
}

// TickState represents the dynamic state of the game at a single tick for history purposes
type TickState struct {
	Players        []PlayerHistoryEntry `json:"players"`
	Bombs          []BombState          `json:"bombs"`
	Explosions     []Vec2               `json:"explosions"`
	DestroyedBoxes []Vec2               `json:"destroyed_boxes,omitempty"`
}

// GameHistory encapsulates the entire history of a game
type GameHistory struct {
	InitialField    FieldState  `json:"initial_field"`
	Ticks           []TickState `json:"ticks"`
	WinnerAuthToken string      `json:"winnerAuthToken"`
}
