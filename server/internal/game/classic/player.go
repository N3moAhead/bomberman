package classic

import "github.com/N3moAhead/bomberman/server/pkg/types"

type Player struct {
	ID       string
	Pos      types.Vec2
	Health   int
	Score    int
	NextMove PlayerMove
}

func (p *Player) HandleInput(payload ClassicInputPayload) {

}
