use client_rust::{Bomber, BomberBot, ClassicStatePayload, PlayerMove};
use url::Url;

struct Bot;

impl BomberBot for Bot {
    fn calc_next_move(&self, _bot_id: &str, _state: &ClassicStatePayload) -> PlayerMove {
        // Currently a pretty lazy player :(
        PlayerMove::Nothing
    }
}

#[tokio::main]
async fn main() {
    let bot = Bot;
    let url = Url::parse("ws://localhost:8038/ws").unwrap();
    if let Err(e) = Bomber::start(bot, url).await {
        eprintln!("Failed to start Bomber client: {}", e);
    }
}
