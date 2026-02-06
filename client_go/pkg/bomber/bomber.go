package bomber

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"

	"github.com/gorilla/websocket"
)

type BomberBot interface {
	CalcNextMove(bomberId string, state ClassicStatePayload) PlayerMove
}

type Bomber struct {
	bomberID    string
	conn        *websocket.Conn
	done        chan struct{}
	interrupt   chan os.Signal
	sendChannel chan Message
	bot         BomberBot
}

func NewBomber(bot BomberBot) *Bomber {
	return &Bomber{
		done:        make(chan struct{}),
		interrupt:   make(chan os.Signal, 1),
		sendChannel: make(chan Message),
		bot:         bot,
	}
}

func (b *Bomber) send(msgType MessageType, payload any) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		error("Error marshalling payload for send: %v", err)
		return
	}

	msg := Message{
		Type:    msgType,
		Payload: payloadBytes,
	}

	b.sendChannel <- msg
}

func (b *Bomber) writePump() {
	for {
		select {
		case message := <-b.sendChannel:
			messageBytes, err := json.Marshal(message)
			if err != nil {
				error("Error marshalling message for write: %v", err)
				continue
			}
			if err := b.conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				error("Write Error: %v", err)
				return
			}
		case <-b.interrupt:
			// Main loop will handle the clean close message
			return
		case <-b.done:
			// Read pump closed, so we should also close
			return
		}
	}
}

func (b *Bomber) Start(u url.URL) {
	info("Trying to connect to %s...", u.String())

	authToken := os.Getenv("BOMBERMAN_CLIENT_AUTH_TOKEN")

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		error("Error while trying to connect")
		log.Fatal(err)
	}
	b.conn = conn

	go b.writePump()
	go b.ReadMessages()

	// Let's notify the server that we are ready for a game
	payload := PlayerStatusUpdatePayload{
		IsReady:   true,
		AuthToken: authToken,
	}
	b.send(PlayerStatusUpdate, payload)

	signal.Notify(b.interrupt, os.Interrupt)
	select {
	case <-b.interrupt:
		info("Detected interrupt closing connection")
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			error("Error while closing connection: %v", err)
		}
	case <-b.done:
		info("The server closed the connection")
	}
}

func (b *Bomber) ReadMessages() {
	defer close(b.done)
	for {
		_, messageBytes, err := b.conn.ReadMessage()
		if err != nil {
			error("Read Error %v", err)
			return
		}
		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			continue
		}
		switch msg.Type {
		case Welcome:
			var payload WelcomeMessage
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				error("Error unmarshalling WelcomeMessage: %v", err)
				continue
			}
			b.bomberID = payload.ClientID
			success("You connected to the bomberman server: %s", b.bomberID)
			info("Available Games:")
			for _, gameInfo := range payload.CurrentGames {
				info("- %s: %s", gameInfo.Name, gameInfo.Description)
			}
		case UpdateLobby:
			var payload LobbyUpdateMessage
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				error("Error unmarshalling LobbyUpdatePayload payload from: %v", err)
				continue
			}
			info("Lobby:")
			for pID, playerInfo := range payload.Players {
				if b.bomberID != "" && b.bomberID == pID {
					info("You (%s):", pID)
				} else {
					info("%s", pID)
				}
				isReady := red("NOT READY")
				if playerInfo.IsReady {
					isReady = green("READY")
				}
				isInGame := green("IS AVAILABLE")
				if playerInfo.InGame {
					isInGame = red("IS IN A GAME")
				}
				info("- Score: %d", playerInfo.Score)
				info("- State:")
				info("  - %s", isReady)
				info("  - %s", isInGame)
			}
		case Error:
			var errorPaylaod ErrorMessage
			err := json.Unmarshal(msg.Payload, &errorPaylaod)
			if err != nil {
				error("Error while trying to unmarshal ErrorMessage: %v", err)
			}
			error("Server Error: %s", errorPaylaod.Message)
		case GameStart:
			var gameStartPayload GameStartPayload
			err := json.Unmarshal(msg.Payload, &gameStartPayload)
			if err != nil {
				error("Error while trying to unmarshal GameStart message: %v", err)
			}
			info("A new %s has started", gameStartPayload.Name)
		case ClassicState:
			var classicState ClassicStatePayload
			err := json.Unmarshal(msg.Payload, &classicState)
			if err != nil {
				error("Error while trying to unmarshal ClassicStatePayload: %v", err)
			}
			nextMove := b.bot.CalcNextMove(b.bomberID, classicState)
			newPayload := ClassicInputPayload{
				Move: nextMove,
			}
			b.send(ClassicInput, newPayload)

			printClassicState(classicState, b.bomberID)
		case BackToLobby:
			info("Your back inside the lobby")
			payload := PlayerStatusUpdatePayload{
				IsReady: true,
			}
			b.send(PlayerStatusUpdate, payload)
		default:
			info("Received: %s\n", messageBytes)
		}
	}
}

func printClassicState(s ClassicStatePayload, ownID string) {
	width := s.Field.Width
	height := s.Field.Height

	playerIcons := make(map[string]string)
	otherPlayers := []PlayerState{}
	ownPlayerIcon := "🤖"
	otherPlayerIcons := []string{"🏃", "🚶", "💃", "🕺"}

	for _, p := range s.Players {
		if p.ID == ownID {
			playerIcons[p.ID] = ownPlayerIcon
		} else {
			otherPlayers = append(otherPlayers, p)
		}
	}

	// Sort other players by ID for consistent icon assignment
	sort.Slice(otherPlayers, func(i, j int) bool {
		return otherPlayers[i].ID < otherPlayers[j].ID
	})

	for i, p := range otherPlayers {
		playerIcons[p.ID] = otherPlayerIcons[i%len(otherPlayerIcons)]
	}

	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
	}

	for y := range height {
		for x := range width {
			idx := y*width + x
			if idx < len(s.Field.Field) {
				tile := s.Field.Field[idx]
				switch tile {
				case WALL:
					grid[y][x] = "🧱"
				case BOX:
					grid[y][x] = "📦"
				default:
					grid[y][x] = "  "
				}
			}
		}
	}

	for _, exp := range s.Explosions {
		if exp.Y >= 0 && exp.Y < height && exp.X >= 0 && exp.X < width {
			grid[exp.Y][exp.X] = "💥"
		}
	}

	for _, bomb := range s.Bombs {
		if bomb.Pos.Y >= 0 && bomb.Pos.Y < height && bomb.Pos.X >= 0 && bomb.Pos.X < width {
			grid[bomb.Pos.Y][bomb.Pos.X] = "💣"
		}
	}

	for _, player := range s.Players {
		if player.Pos.Y >= 0 && player.Pos.Y < height && player.Pos.X >= 0 && player.Pos.X < width {
			grid[player.Pos.Y][player.Pos.X] = playerIcons[player.ID]
		}
	}

	var sb strings.Builder

	sb.WriteString("\033[H\033[2J")

	sb.WriteString("╔")
	sb.WriteString(strings.Repeat("══", width))
	sb.WriteString("╗\n")

	for y := range height {
		sb.WriteString("║")
		sb.WriteString(strings.Join(grid[y], ""))
		sb.WriteString("║\n")
	}

	sb.WriteString("╚")
	sb.WriteString(strings.Repeat("══", width))
	sb.WriteString("╝\n")

	sb.WriteString("--- PLAYERS ---\n")
	for _, p := range s.Players {
		fmt.Fprintf(&sb, "%s Player ...%s | Health: %d, Score: %d\n", playerIcons[p.ID], p.ID[len(p.ID)-4:], p.Health, p.Score)
	}

	if len(s.Bombs) > 0 {
		sb.WriteString("--- BOMBS ---\n")
		for _, b := range s.Bombs {
			fmt.Fprintf(&sb, "💣 at (%d,%d) | Fuse: %d\n", b.Pos.X, b.Pos.Y, b.Fuse)
		}
	}

	fmt.Print(sb.String())
}
