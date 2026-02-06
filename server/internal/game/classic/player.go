package classic

import "github.com/N3moAhead/bomberman/server/pkg/types"

type Player struct {
	ID        string     `json:"id"`
	Pos       types.Vec2 `json:"pos"`
	Health    int        `json:"health"`
	Score     int        `json:"score"`
	AuthToken string
	NextMove  PlayerMove
}

func (p *Player) HandleInput(payload ClassicInputPayload) {
	p.NextMove = payload.Move
}
