package classic

type Field [field_width * field_height]Tile

func NewField() *Field {
	f := Field{} // Will be initted with all air

	// Let's place some walls :)
	for x := range field_width {
		for y := range field_height {
			// left or right wall
			if x == 0 || x == field_width-1 {
				f.setTile(x, y, WALL)
			}
			// top or bot wall
			if y == 0 || y == field_height-1 {
				f.setTile(x, y, WALL)
			}

			// Labyrinth Walls
			if x%2 == 0 && y%2 == 0 {
				f.setTile(x, y, WALL)
			}

			// TODO add box placement
		}
	}

	return &f
}

func (f *Field) getTile(x, y int) Tile {
	return f[y*field_height+x]
}

func (f *Field) setTile(x, y int, tile Tile) {
	f[y*field_height+x] = tile
}

func (f *Field) isWall(x, y int) bool {
	tile := f.getTile(x, y)
	return tile == WALL
}

func (f *Field) isBox(x, y int) bool {
	tile := f.getTile(x, y)
	return tile == BOX
}

func (f *Field) isTileBlocked(x, y int) bool {
	tile := f.getTile(x, y)
	if tile == WALL || tile == BOX {
		return true
	}
	return false
}
