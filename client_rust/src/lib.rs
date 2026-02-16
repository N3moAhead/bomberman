use futures_util::{stream::SplitStream, SinkExt, StreamExt};
use serde_json::json;
use std::{collections::HashMap, sync::Arc};
use tokio::{
    net::TcpStream,
    sync::{mpsc, Mutex},
};
use tokio_tungstenite::{
    connect_async,
    tungstenite::{self, protocol::Message as TungsteniteMessage},
    MaybeTlsStream, WebSocketStream,
};
use url::Url;

pub mod types;
pub use types::*;

pub trait BomberBot: Send + Sync + 'static {
    fn calc_next_move(&self, bot_id: &str, state: &ClassicStatePayload) -> PlayerMove;
}

type BomberId = String;

pub struct Bomber;

impl Bomber {
    pub async fn start<B: BomberBot>(bot: B, url: Url) -> Result<(), tungstenite::Error> {
        println!("Trying to connect to {}...", url);
        let (ws_stream, _) = connect_async(url.as_str()).await?;
        let (mut write, read) = ws_stream.split();
        println!("Successfully connected!");

        let (send_tx, mut send_rx) = mpsc::unbounded_channel();

        // Write pump
        tokio::spawn(async move {
            while let Some(msg) = send_rx.recv().await {
                if write.send(msg).await.is_err() {
                    eprintln!("Error sending message");
                    break;
                }
            }
        });

        let bot = Arc::new(bot);
        let bomber_id = Arc::new(Mutex::new(None));

        // Initial ready status
        let auth_token =
            std::env::var("BOMBERMAN_CLIENT_AUTH_TOKEN").unwrap_or_else(|_| "".to_string());
        let payload = PlayerStatusUpdatePayload {
            is_ready: true,
            auth_token,
        };
        let ready_payload = json!({
            "type": MessageType::PlayerStatusUpdate,
            "payload": payload
        });
        if let Ok(msg_str) = serde_json::to_string(&ready_payload) {
            if send_tx.send(TungsteniteMessage::Text(msg_str)).is_err() {
                eprintln!("Failed to send initial ready status");
            }
        }

        Self::start_reading(read, bot, bomber_id, send_tx).await;

        Ok(())
    }

    async fn start_reading<B: BomberBot>(
        mut read: SplitStream<WebSocketStream<MaybeTlsStream<TcpStream>>>,
        bot: Arc<B>,
        bomber_id: Arc<Mutex<Option<BomberId>>>,
        send_tx: mpsc::UnboundedSender<TungsteniteMessage>,
    ) {
        while let Some(Ok(msg)) = read.next().await {
            if let TungsteniteMessage::Text(text) = msg {
                if let Ok(parsed_message) = serde_json::from_str::<Message>(&text) {
                    let bot_clone = bot.clone();
                    let bomber_id_clone = bomber_id.clone();
                    let send_tx_clone = send_tx.clone();
                    tokio::spawn(async move {
                        Self::handle_message(
                            parsed_message,
                            bot_clone,
                            bomber_id_clone,
                            send_tx_clone,
                        )
                        .await;
                    });
                }
            }
        }
    }

    async fn handle_message<B: BomberBot>(
        msg: Message,
        bot: Arc<B>,
        bomber_id: Arc<Mutex<Option<BomberId>>>,
        tx: mpsc::UnboundedSender<TungsteniteMessage>,
    ) {
        match msg.type_of {
            MessageType::Welcome => {
                if let Ok(payload) = serde_json::from_value::<WelcomeMessage>(msg.payload) {
                    println!(
                        "You connected to the bomberman server: {}",
                        payload.client_id
                    );
                    *bomber_id.lock().await = Some(payload.client_id);
                    println!("Available Games:");
                    for game_info in payload.current_games {
                        println!("- {}: {}", game_info.name, game_info.description);
                    }
                }
            }
            MessageType::UpdateLobby => {
                if let Ok(payload) = serde_json::from_value::<LobbyUpdateMessage>(msg.payload) {
                    println!("Lobby:");
                    let id_lock = bomber_id.lock().await;
                    for (p_id, player_info) in payload.players {
                        if let Some(id) = &*id_lock {
                            if id == &p_id {
                                println!("You ({}):", p_id);
                            } else {
                                println!("{}:", p_id);
                            }
                        } else {
                            println!("{}:", p_id);
                        }
                        let is_ready = if player_info.is_ready {
                            "READY"
                        } else {
                            "NOT READY"
                        };
                        let is_in_game = if player_info.in_game {
                            "IS IN A GAME"
                        } else {
                            "IS AVAILABLE"
                        };
                        println!("- Score: {}", player_info.score);
                        println!("- State:");
                        println!("  - {}", is_ready);
                        println!("  - {}", is_in_game);
                    }
                }
            }
            MessageType::Error => {
                if let Ok(error_payload) = serde_json::from_value::<ErrorMessage>(msg.payload) {
                    eprintln!("Server Error: {}", error_payload.message);
                }
            }
            MessageType::GameStart => {
                if let Ok(game_start_payload) =
                    serde_json::from_value::<GameStartPayload>(msg.payload)
                {
                    println!("A new {} has started", game_start_payload.name);
                }
            }
            MessageType::ClassicState => {
                if let Ok(classic_state) =
                    serde_json::from_value::<ClassicStatePayload>(msg.payload)
                {
                    let id_lock = bomber_id.lock().await;
                    if let Some(id) = &*id_lock {
                        let next_move = bot.calc_next_move(id, &classic_state);
                        let payload = json!({
                            "type": MessageType::ClassicInput,
                            "payload": { "move": next_move.to_string() }
                        });
                        if let Ok(msg_str) = serde_json::to_string(&payload) {
                            if tx.send(TungsteniteMessage::Text(msg_str)).is_err() {
                                eprintln!("Failed to send move");
                            }
                        }

                        print_classic_state(&classic_state, id);
                    }
                }
            }
            MessageType::BackToLobby => {
                println!("You are back inside the lobby");
                let payload = PlayerStatusUpdatePayload {
                    is_ready: true,
                    auth_token: "".to_string(),
                };
                let msg = json!({
                    "type": MessageType::PlayerStatusUpdate,
                    "payload": payload
                });
                if let Ok(msg_str) = serde_json::to_string(&msg) {
                    if tx.send(TungsteniteMessage::Text(msg_str)).is_err() {
                        eprintln!("Failed to send ready status");
                    }
                }
            }
            _ => {
                println!("Received: {:?}", msg);
            }
        }
    }
}

fn print_classic_state(s: &ClassicStatePayload, own_id: &str) {
    let width = s.field.width;
    let height = s.field.height;

    let mut player_icons: HashMap<String, &str> = HashMap::new();
    let mut other_players: Vec<&PlayerState> = Vec::new();
    let own_player_icon = "🤖";
    let other_player_icons = vec!["🏃", "🚶", "💃", "🕺"];

    for p in &s.players {
        if p.id == own_id {
            player_icons.insert(p.id.clone(), own_player_icon);
        } else {
            other_players.push(p);
        }
    }

    other_players.sort_by(|a, b| a.id.cmp(&b.id));

    for (i, p) in other_players.iter().enumerate() {
        player_icons.insert(
            p.id.clone(),
            other_player_icons[i % other_player_icons.len()],
        );
    }

    let mut grid = vec![vec!["  "; width]; height];

    for y in 0..height {
        for x in 0..width {
            let idx = y * width + x;
            if idx < s.field.field.len() {
                let tile = &s.field.field[idx];
                match tile {
                    Tile::WALL => grid[y][x] = "🧱",
                    Tile::BOX => grid[y][x] = "📦",
                    _ => grid[y][x] = "  ",
                }
            }
        }
    }

    for exp in &s.explosions {
        if exp.y >= 0 && (exp.y as usize) < height && exp.x >= 0 && (exp.x as usize) < width {
            grid[exp.y as usize][exp.x as usize] = "💥";
        }
    }

    for bomb in &s.bombs {
        if bomb.pos.y >= 0
            && (bomb.pos.y as usize) < height
            && bomb.pos.x >= 0
            && (bomb.pos.x as usize) < width
        {
            grid[bomb.pos.y as usize][bomb.pos.x as usize] = "💣";
        }
    }

    for player in &s.players {
        if player.pos.y >= 0
            && (player.pos.y as usize) < height
            && player.pos.x >= 0
            && (player.pos.x as usize) < width
        {
            grid[player.pos.y as usize][player.pos.x as usize] =
                player_icons.get(&player.id).unwrap_or(&" ");
        }
    }

    let mut sb = String::new();

    sb.push_str("\x1B[H\x1B[2J");

    sb.push('╔');
    sb.push_str(&"══".repeat(width));
    sb.push_str("╗\n");

    for y in 0..height {
        sb.push('║');
        sb.push_str(&grid[y].join(""));
        sb.push_str("║\n");
    }

    sb.push('╚');
    sb.push_str(&"══".repeat(width));
    sb.push_str("╝\n");

    sb.push_str("-- PLAYERS --\n");
    for p in &s.players {
        sb.push_str(&format!(
            "{} Player ...{} | Health: {}, Score: {}\n",
            player_icons.get(&p.id).unwrap_or(&" "),
            &p.id[p.id.len() - 4..],
            p.health,
            p.score
        ));
    }

    if !s.bombs.is_empty() {
        sb.push_str("-- BOMBS --\n");
        for b in &s.bombs {
            sb.push_str(&format!(
                "💣 at ({},{}) | Fuse: {}\n",
                b.pos.x, b.pos.y, b.fuse
            ));
        }
    }

    print!("{}", sb);
}
