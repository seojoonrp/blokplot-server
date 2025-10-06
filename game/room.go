// game/room.go

package game

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender *websocket.Conn
	Content []byte
}

type Room struct {
	players []*websocket.Conn
	broadcast chan Message
}

func NewRoom(players ...*websocket.Conn) *Room {
	return &Room{
		players: players,
		broadcast: make(chan Message),
	}
}

func (r *Room) Run() {
	for _, player := range r.players {
		player.WriteMessage(websocket.TextMessage, []byte("Game matched! Starting game..."))
		go r.readMessages(player)
	}

	for msg := range r.broadcast {
		for _, player := range r.players {
			if msg.Sender != player {
				if err := player.WriteMessage(websocket.TextMessage, msg.Content); err != nil {
					player.WriteMessage(websocket.TextMessage, []byte("Player has disconnected."))
				}
			}
		}
	}

	log.Println("Game room closed.")
}

func (r *Room) readMessages(conn *websocket.Conn) {
	defer func() {
		close(r.broadcast)
		conn.Close()
	}()

	for {
		_, content, err := conn.ReadMessage()

		if err != nil {
			log.Printf("Error while reading message: %v", err)
			close(r.broadcast)
			break
		}

		var baseMessage BaseMessage
		if err := json.Unmarshal(content, &baseMessage); err != nil {
			log.Printf("Invalid JSON message received: %v", err)
			continue
		}

		switch baseMessage.Type {
		case "chat":
			log.Println("Chat message received.")
		case "placeBlock":
			log.Println("Block placement data received.")
		default:
			log.Printf("Unknown type of data received: %s\n", baseMessage.Type)
		}

		r.broadcast <- Message{
			Sender: conn,
			Content: content,
		}
	}
}