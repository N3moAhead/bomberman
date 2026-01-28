```
    ____                  __
   / __ )____  ____ ___  / /_  ___  _________ ___  ____ _____
  / __  / __ \/ __ `__ \/ __ \/ _ \/ ___/ __ `__ \/ __ `/ __ \
 / /_/ / /_/ / / / / / / /_/ /  __/ /  / / / / / / /_/ / / / /
/_____/\____/_/ /_/ /_/_.___/\___/_/  /_/ /_/ /_/\__,_/_/ /_/
```

![alt text](images/game_example_image.png)

This repository contains a client-server implementation of the classic game Bomberman. The server is written in Go, and there are example clients in Go, JavaScript, and Rust that you can use as a starting point to develop your own bot.

# Getting Started

The game consists of a central server and multiple clients (your bots). You need to run the server first, and then you can connect your clients to it.

## 1. Run the Server

The server is located in the `server/` directory.

```bash
cd server/
make run
```

The server will start and listen for connections on port `:8038`.

## 2. Run a Client

You can run one of the provided clients or create your own.

### Go Client

```bash
cd client_go/
make run
```

### JavaScript Client

```bash
cd client_js/
npm install
npm start
```

### Rust Client

```bash
cd client_rust/
cargo run
```

# How to Write a Bot

Your bot will connect to the game server via WebSockets (`ws://localhost:8038/ws`). The core loop for a bot is simple: receive the game state from the server, decide on an action, and send that action back to the server.

Bots are automatically set to a "ready" state upon connecting, so you don't need to implement that manually.

## Game State (`classic_state`)

On every game tick, the server broadcasts the current game state to all clients. This message has the type `classic_state` and its payload contains the following information:

-   `players` ([]PlayerState): A list of all players currently in the game.
    -   `id` (string): The unique ID of the player.
    -   `pos` (Vec2): The `{x, y}` coordinates of the player on the map.
    -   `health` (int): The current health of the player.
    -   `score` (int): The current score of the player.
-   `field` (FieldState): The game map itself.
    -   `width` (int): The width of the map.
    -   `height` (int): The height of the map.
    -   `field` ([]Tile): A 1D array representing the 2D map. The type of each tile can be `air`, `wall`, or `box`.
-   `bombs` ([]BombState): A list of active bombs on the map.
    -   `pos` (Vec2): The `{x, y}` coordinates of the bomb.
    -   `fuse` (int): The number of ticks until the bomb explodes.
-   `explosions` ([]Vec2): A list of `{x, y}` coordinates where explosions are currently active.

## Player Actions (`classic_input`)

To control your bot, you send a message with the type `classic_input`. The payload must contain your desired move.

-   `move` (PlayerMove): The action to perform on the next tick.

Possible values for `move` are:
-   `nothing`: Do nothing for one turn.
-   `move_up`: Move one tile up.
-   `move_down`: Move one tile down.
-   `move_left`: Move one tile left.
-   `move_right`: Move one tile right.
-   `place_bomb`: Place a bomb at the player's current position.

You can study the `message.go` file in `client_go/pkg/bomber/` for the exact data structures.

# Legacy C Game

The original version of this project was a SDL-based game written entirely in C. You can still find it in the `c_game/` directory. It has its own set of rules and instructions for writing bots. If you're interested, please refer to the `README.md` inside that directory for more information.
