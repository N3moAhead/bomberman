package classic

import "time"

const (
	// --- Field ---
	field_width    = 11
	field_height   = 11
	box_spawn_rate = 0.75

	// --- Game l---
	min_players = 2
	max_players = 4
	TICK_RATE   = 200 * time.Millisecond

	// --- Player ---
	initial_health = 3
)
