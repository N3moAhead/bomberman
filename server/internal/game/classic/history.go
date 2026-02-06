package classic

import "github.com/N3moAhead/bomberman/server/pkg/types"

// History manages the recording of a game's progression
type History struct {
	InitialField FieldState
	Ticks        []TickState
}

// NewHistory creates a new game history recorder, capturing the initial state of the field
func NewHistory(initialField FieldState) *History {
	return &History{
		InitialField: initialField,
		Ticks:        make([]TickState, 0),
	}
}

// RecordTick captures the dynamic state of the game for the current tick
func (h *History) RecordTick(
	players map[string]*Player,
	bombs map[string]*Bomb,
	explosions map[string]types.Vec2,
	destroyedBoxes []types.Vec2,
) {
	playerHistory := make([]PlayerHistoryEntry, 0, len(players))
	for _, p := range players {
		playerHistory = append(playerHistory, PlayerHistoryEntry{
			PlayerState: PlayerState{
				ID:     p.ID,
				Pos:    p.Pos,
				Health: p.Health,
				Score:  p.Score,
			},
			Move:      p.NextMove,
			AuthToken: p.AuthToken,
		})
	}

	bombStates := make([]BombState, 0, len(bombs))
	for _, b := range bombs {
		bombStates = append(bombStates, BombState{Pos: b.Pos, Fuse: b.Fuse})
	}

	explosionVecs := make([]types.Vec2, 0, len(explosions))
	for _, e := range explosions {
		explosionVecs = append(explosionVecs, e)
	}

	tick := TickState{
		Players:        playerHistory,
		Bombs:          bombStates,
		Explosions:     explosionVecs,
		DestroyedBoxes: destroyedBoxes,
	}

	h.Ticks = append(h.Ticks, tick)
}

// ToGameHistory converts the internal history representation to the serializable format
func (h *History) ToGameHistory(winnerAuthToken string) GameHistory {
	return GameHistory{
		InitialField:    h.InitialField,
		Ticks:           h.Ticks,
		WinnerAuthToken: winnerAuthToken,
	}
}
