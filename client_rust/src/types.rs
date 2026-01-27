use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Serialize, Deserialize, Debug)]
pub struct Message {
    #[serde(rename = "type")]
    pub type_of: MessageType,
    pub payload: serde_json::Value,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash)]
#[serde(rename_all = "snake_case")]
pub enum MessageType {
    Welcome,
    BackToLobby,
    UpdateLobby,
    PlayerStatusUpdate,
    Error,
    ClassicInput,
    ClassicState,
    GameStart,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GameInfo {
    pub name: String,
    pub description: String,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct WelcomeMessage {
    pub client_id: String,
    pub current_games: Vec<GameInfo>,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct PlayerInfo {
    pub in_game: bool,
    pub is_ready: bool,
    pub score: i32,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct LobbyUpdateMessage {
    pub players: std::collections::HashMap<String, PlayerInfo>,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct PlayerStatusUpdatePayload {
    pub is_ready: bool,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct GameStartPayload {
    pub name: String,
    pub description: String,
    pub game_id: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ErrorMessage {
    pub message: String,
}

#[derive(Serialize, Deserialize, Debug, Clone, Copy)]
#[serde(rename_all = "snake_case")]
pub enum PlayerMove {
    Nothing,
    MoveUp,
    MoveRight,
    MoveDown,
    MoveLeft,
    PlaceBomb,
}

impl fmt::Display for PlayerMove {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(
            f,
            "{}",
            match self {
                PlayerMove::Nothing => "nothing",
                PlayerMove::MoveUp => "move_up",
                PlayerMove::MoveRight => "move_right",
                PlayerMove::MoveDown => "move_down",
                PlayerMove::MoveLeft => "move_left",
                PlayerMove::PlaceBomb => "place_bomb",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ClassicInputPayload {
    #[serde(rename = "move")]
    pub a_move: PlayerMove,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PlayerState {
    pub id: String,
    pub pos: Vec2,
    pub health: i32,
    pub score: i32,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct FieldState {
    pub width: usize,
    pub height: usize,
    pub field: Vec<Tile>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct BombState {
    pub pos: Vec2,
    pub fuse: i32,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ClassicStatePayload {
    pub players: Vec<PlayerState>,
    pub field: FieldState,
    pub bombs: Vec<BombState>,
    pub explosions: Vec<Vec2>,
}

#[derive(Serialize, Deserialize, Debug, Clone, Copy)]
pub enum Tile {
    AIR,
    WALL,
    BOX,
}

#[derive(Serialize, Deserialize, Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct Vec2 {
    pub x: i32,
    pub y: i32,
}
