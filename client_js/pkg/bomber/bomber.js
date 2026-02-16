const WebSocket = require("ws");
const { MessageType, PlayerMove } = require("./message");
const { Tile } = require("./tile");
const logger = require("./logger");

class Bomber {
  constructor(bot) {
    this.bomberID = "";
    this.conn = null;
    this.bot = bot;
  }

  send(msgType, payload) {
    const msg = {
      type: msgType,
      payload: payload,
    };
    if (this.conn.readyState === WebSocket.OPEN) {
      this.conn.send(JSON.stringify(msg));
    }
  }

  start(url) {
    logger.info("Trying to connect to %s...", url);

    try {
      this.conn = new WebSocket(url);
    } catch (e) {
      logger.error("Error while trying to connect");
      console.error(e);
      process.exit(1);
    }

    this.conn.on("open", () => {
      logger.info("Connection established");
      const authToken = process.env.BOMBERMAN_CLIENT_AUTH_TOKEN || "";
      const payload = {
        isReady: true,
        authToken: authToken,
      };
      this.send(MessageType.PlayerStatusUpdate, payload);
    });

    this.conn.on("message", (data) => {
      this.handleMessage(data);
    });

    this.conn.on("close", () => {
      logger.info("The server closed the connection");
      process.exit(0);
    });

    this.conn.on("error", (err) => {
      logger.error("WebSocket error: %s", err.message);
      process.exit(1);
    });

    process.on("SIGINT", () => {
      logger.info("Detected interrupt, closing connection");
      if (this.conn) {
        this.conn.close(1000, "Normal Closure");
      }
      process.exit(0);
    });
  }

  handleMessage(data) {
    let msg;
    try {
      msg = JSON.parse(data);
    } catch (e) {
      logger.error("Error parsing message: %s", data);
      return;
    }

    switch (msg.type) {
      case MessageType.Welcome: {
        const payload = msg.payload;
        this.bomberID = payload.clientId;
        logger.success(
          "You connected to the bomberman server: %s",
          this.bomberID,
        );
        logger.info("Available Games:");
        payload.currentGames.forEach((gameInfo) => {
          logger.info("- %s: %s", gameInfo.name, gameInfo.description);
        });
        break;
      }
      case MessageType.UpdateLobby: {
        const payload = msg.payload;
        logger.info("Lobby:");
        for (const pID in payload.players) {
          const playerInfo = payload.players[pID];
          if (this.bomberID && this.bomberID === pID) {
            logger.info("You (%s):", pID.slice(-4));
          } else {
            logger.info("Player ...%s:", pID.slice(-4));
          }
          const isReady = playerInfo.isReady
            ? logger.green("READY")
            : logger.red("NOT READY");
          const isInGame = playerInfo.inGame
            ? logger.red("IS IN A GAME")
            : logger.green("IS AVAILABLE");
          logger.info("- Score: %d", playerInfo.score);
          logger.info("- State:");
          logger.info("  - %s", isReady);
          logger.info("  - %s", isInGame);
        }
        break;
      }
      case MessageType.Error: {
        const errorPayload = msg.payload;
        logger.error("Server Error: %s", errorPayload.message);
        break;
      }
      case MessageType.GameStart: {
        const gameStartPayload = msg.payload;
        logger.info("A new %s has started", gameStartPayload.name);
        break;
      }
      case MessageType.ClassicState: {
        const classicState = msg.payload;
        const nextMove = this.bot.calcNextMove(this.bomberID, classicState);
        const newPayload = {
          move: nextMove,
        };
        this.send(MessageType.ClassicInput, newPayload);
        printClassicState(classicState, this.bomberID);
        break;
      }
      case MessageType.BackToLobby: {
        logger.info("You are back inside the lobby");
        const payload = {
          isReady: true,
          authToken: "",
        };
        this.send(MessageType.PlayerStatusUpdate, payload);
        break;
      }
      default:
        logger.debug("Received unknown message type: %s", msg.type);
    }
  }
}

function printClassicState(s, ownID) {
  const width = s.field.width;
  const height = s.field.height;

  const playerIcons = {};
  const otherPlayers = [];
  const ownPlayerIcon = "🤖";
  const otherPlayerIcons = ["🏃", "🚶", "💃", "🕺"];

  s.players.forEach((p) => {
    if (p.id === ownID) {
      playerIcons[p.id] = ownPlayerIcon;
    } else {
      otherPlayers.push(p);
    }
  });

  otherPlayers.sort((a, b) => a.id.localeCompare(b.id));

  otherPlayers.forEach((p, i) => {
    playerIcons[p.id] = otherPlayerIcons[i % otherPlayerIcons.length];
  });

  const grid = Array(height)
    .fill(null)
    .map(() => Array(width).fill("  "));

  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      const idx = y * width + x;
      if (idx < s.field.field.length) {
        const tile = s.field.field[idx];
        switch (tile) {
          case Tile.WALL:
            grid[y][x] = "🧱";
            break;
          case Tile.BOX:
            grid[y][x] = "📦";
            break;
        }
      }
    }
  }

  s.explosions.forEach((exp) => {
    if (exp.y >= 0 && exp.y < height && exp.x >= 0 && exp.x < width) {
      grid[exp.y][exp.x] = "💥";
    }
  });

  s.bombs.forEach((bomb) => {
    if (
      bomb.pos.y >= 0 &&
      bomb.pos.y < height &&
      bomb.pos.x >= 0 &&
      bomb.pos.x < width
    ) {
      grid[bomb.pos.y][bomb.pos.x] = "💣";
    }
  });

  s.players.forEach((player) => {
    if (
      player.pos.y >= 0 &&
      player.pos.y < height &&
      player.pos.x >= 0 &&
      player.pos.x < width
    ) {
      grid[player.pos.y][player.pos.x] = playerIcons[player.id];
    }
  });

  let sb = "";

  console.clear();

  sb += "╔" + "══".repeat(width) + "╗\n";

  for (let y = 0; y < height; y++) {
    sb += "║" + grid[y].join("") + "║\n";
  }

  sb += "╚" + "══".repeat(width) + "╝\n";

  sb += "--- PLAYERS ---\n";
  s.players
    .sort((a, b) => a.id.localeCompare(b.id))
    .forEach((p) => {
      sb += `${playerIcons[p.id]} Player ...${p.id.slice(-4)} | Health: ${p.health}, Score: ${p.score}\n`;
    });

  if (s.bombs.length > 0) {
    sb += "--- BOMBS ---\n";
    s.bombs.forEach((b) => {
      sb += `💣 at (${b.pos.x},${b.pos.y}) | Fuse: ${b.fuse}\n`;
    });
  }

  process.stdout.write(sb);
}

module.exports = { Bomber, PlayerMove };
