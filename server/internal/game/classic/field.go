package classic

type Field [field_width * field_height]Tile

func NewField() *Field {
	f := Field{} // Will be initted with all air

	// Let's place some walls :)
	for x := range field_width {
		for y := range field_height {
			// left or right wall
			if x == 0 || x == field_width-1 {
				f.SetTile(x, y, WALL)
			}
			// top or bot wall
			if y == 0 || y == field_height-1 {
				f.SetTile(x, y, WALL)
			}

			// Labyrinth Walls
			if x%2 == 0 && y%2 == 0 {
				f.SetTile(x, y, WALL)
			}

			// TODO add box placement
		}
	}

	return &f
}

func (f *Field) GetTile(x, y int) Tile {
	return f[y*field_height+x]
}

func (f *Field) SetTile(x, y int, tile Tile) {
	f[y*field_height+x] = tile
}
