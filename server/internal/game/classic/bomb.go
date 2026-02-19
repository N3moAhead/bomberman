package classic

import "github.com/N3moAhead/bombahead/server/pkg/types"

type Bomb struct {
	Pos  types.Vec2 `json:"pos"`
	Fuse int        `json:"fuse"` // Countdown each tick explodes at 0
}

func NewBomb(pos types.Vec2) *Bomb {
	return &Bomb{
		Pos:  pos,
		Fuse: fuse_ticks,
	}
}

func (c *Classic) containsBomb(pos types.Vec2) bool {
	_, ok := c.bombs[pos.String()]
	return ok
}
