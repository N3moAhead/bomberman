package classic

import (
	"github.com/N3moAhead/bomberman/server/pkg/types"
)

func (c *Classic) update() {
	// Move all players
	c.applyPlayerInput()
	// Bombs Explode
	c.resetExplosions()
	c.updateBombs()
}

func (c *Classic) applyPlayerInput() {
	for _, player := range c.players {
		switch player.NextMove {
		case MOVE_UP:
			newPos := player.Pos.Add(types.Vec2{X: 0, Y: -1})
			// TODO Check if the tile contains a bomb the player can't walk over bombs
			if !c.field.isTileBlocked(newPos.X, newPos.Y) {
				player.Pos = newPos
			}
		case MOVE_RIGHT:
			newPos := player.Pos.Add(types.Vec2{X: 1, Y: 0})
			if !c.field.isTileBlocked(newPos.X, newPos.Y) {
				player.Pos = newPos
			}
		case MOVE_DOWN:
			newPos := player.Pos.Add(types.Vec2{X: 0, Y: 1})
			if !c.field.isTileBlocked(newPos.X, newPos.Y) {
				player.Pos = newPos
			}
		case MOVE_LEFT:
			newPos := player.Pos.Add(types.Vec2{X: -1, Y: 0})
			if !c.field.isTileBlocked(newPos.X, newPos.Y) {
				player.Pos = newPos
			}
		case PLACE_BOMB:
			if !c.containsBomb(player.Pos) {
				newBomb := NewBomb(player.Pos)
				c.bombs[newBomb.Pos.String()] = newBomb
			}
		default:
			// Is the do nothing move
		}
	}
}

func (c *Classic) updateBombs() {
	for _, bomb := range c.bombs {
		bomb.Fuse -= 1
		if bomb.Fuse < 1 {
			delete(c.bombs, bomb.Pos.String())
			c.explodeBomb(bomb.Pos, bomb_explosion_radius)
		}
	}
}

func (c *Classic) resetExplosions() {
	for k := range c.explosions {
		delete(c.explosions, k)
	}

}

func (c *Classic) explodeBomb(pos types.Vec2, distance int) {
	// Up
	c.createExplodePath(pos, types.NewVec2(0, -1), distance)
	// Right
	c.createExplodePath(pos, types.NewVec2(1, 0), distance)
	// Down
	c.createExplodePath(pos, types.NewVec2(0, 1), distance)
	// Left
	c.createExplodePath(pos, types.NewVec2(-1, 0), distance)
}

func (c *Classic) createExplodePath(pos types.Vec2, dir types.Vec2, distance int) {
	// I won't check if the explosion is out of bounds because a bomb can't be placed
	// out of bounds so it could not occur...
	if distance == 0 {
		return
	}
	tile := c.field.getTile(pos.X, pos.Y)
	if tile == WALL {
		return
	}
	if tile == BOX {
		c.field.setTile(pos.X, pos.Y, AIR)
		c.addExplosion(pos)
		return
	}
	if tile == AIR {
		// there could be a bomb
		if c.containsBomb(pos) {
			// Explosion explodes bomb
			c.explodeBomb(pos, bomb_explosion_radius)
		}
		c.addExplosion(pos)
		c.createExplodePath(pos.Add(dir), dir, distance-1)
	}
}

func (c *Classic) addExplosion(pos types.Vec2) {
	if _, ok := c.explosions[pos.String()]; !ok {
		c.explosions[pos.String()] = pos
	}
}

func (c *Classic) isGameOver() bool {
	alivePlayers := 0
	for _, player := range c.players {
		if player.Health > 0 {
			alivePlayers += 1
		}
	}

	// If there are one or less players left the
	// game is over
	return alivePlayers <= 1
}

func (c *Classic) getGameState() ClassicStatePayload {
	// Get Players
	pStates := []PlayerState{}
	for _, player := range c.players {
		pState := PlayerState{
			ID:     player.ID,
			Pos:    player.Pos,
			Health: player.Health,
			Score:  player.Score,
		}
		pStates = append(pStates, pState)
	}

	// Get Field
	field := []Tile{}
	for x := range field_width {
		for y := range field_height {
			field = append(field, c.field.getTile(x, y))
		}
	}
	fieldState := FieldState{
		Width:  field_width,
		Height: field_height,
		Field:  field,
	}
	// Get Bombs
	bombs := []BombState{}
	for _, bomb := range c.bombs {
		bombs = append(bombs, BombState{Pos: bomb.Pos, Fuse: bomb.Fuse})
	}
	// Get Explosions
	explosions := []types.Vec2{}
	for _, ePos := range c.explosions {
		explosions = append(explosions, ePos)
	}

	return ClassicStatePayload{
		Players:    pStates,
		Field:      fieldState,
		Bombs:      bombs,
		Explosions: explosions,
	}
}
