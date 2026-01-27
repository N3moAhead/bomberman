const MessageType = {
    Welcome: "welcome",
    BackToLobby: "back_to_lobby",
    UpdateLobby: "update_lobby",
    PlayerStatusUpdate: "player_status_update",
    Error: "error",
    ClassicInput: "classic_input",
    ClassicState: "classic_state",
    GameStart: "game_start",
};

const PlayerMove = {
    DO_NOTHING: "nothing",
    MOVE_UP: "move_up",
    MOVE_RIGHT: "move_right",
    MOVE_DOWN: "move_down",
    MOVE_LEFT: "move_left",
    PLACE_BOMB: "place_bomb",
};

module.exports = {
    MessageType,
    PlayerMove,
};
