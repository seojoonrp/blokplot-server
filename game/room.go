// game/room.go

package game

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender *Player
	Content []byte
}

type Room struct {
	players map[*websocket.Conn]*Player
	broadcast chan Message
	mu sync.Mutex
}

func NewRoom(conns ...*websocket.Conn) *Room {
	players := make(map[*websocket.Conn]*Player)

	for i, conn := range conns {
		players[conn] = &Player{
			Conn: conn,
			GameIndex: i,
		}
	}

	return &Room{
		players: players,
		broadcast: make(chan Message),
	}
}

func (r *Room) removePlayer(player *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.players[player.Conn]; ok {
		delete(r.players, player.Conn)
		log.Printf("Player %d removed from game room.", player.GameIndex)
	}

	if len(r.players) == 0 {
		close(r.broadcast)
	}
}

func (r *Room) Run() {
	for _, player := range r.players {
		gameInitData := GameInitMessage{ PlayerIndex: player.GameIndex }
		dataBytes, _ := json.Marshal(gameInitData)

		initMessage := ServerMessage{
			Type: "gameInit",
			SenderIndex: -1,
			Data: dataBytes,
		}
		finalJson, _ := json.Marshal(initMessage)

		player.Conn.WriteMessage(websocket.TextMessage, finalJson)

		go r.readMessages(player)
	}

	for msg := range r.broadcast {
		var clientMessage ClientMessage
		if err := json.Unmarshal(msg.Content, &clientMessage); err != nil {
			log.Printf("Failed to unmarshal client message: %v", err)
			continue
		}

		serverMessage := ServerMessage{
			Type: clientMessage.Type,
			SenderIndex: msg.Sender.GameIndex,
			Data: clientMessage.Data,
		}

		HandleServerMessage(&serverMessage)
		
		newContent, err := json.Marshal(serverMessage)
		if err != nil {
			log.Printf("Failed to marshal server message: %v", err)
			continue
		}

		for _, player := range r.players {
			if msg.Sender != player {
				player.Conn.WriteMessage(websocket.TextMessage, newContent)
			}
		}
	}

	log.Println("Game room closed.")
}

func (r *Room) readMessages(player *Player) {
	defer func() {
		r.removePlayer(player)
		player.Conn.Close()
	}()

	for {
		_, content, err := player.Conn.ReadMessage()

		if err != nil {
			log.Printf("Error while reading message from player %d: %v", player.GameIndex, err)
			break
		}

		r.broadcast <- Message{
			Sender: player,
			Content: content,
		}
	}
}

func HandleServerMessage(serverMessage *ServerMessage) {
	switch serverMessage.Type {
	case "chat":
		log.Printf("Chat message received from Player #%d", serverMessage.SenderIndex)
	case "blockPlace":
		log.Printf("Block place message received from Player #%d", serverMessage.SenderIndex)
	default:
		log.Printf("Unknown type \"%s\"of message received from Player #%d", serverMessage.Type, serverMessage.SenderIndex)
	}
}