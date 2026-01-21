package classic

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
