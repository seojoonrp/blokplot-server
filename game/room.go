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
	player1 *websocket.Conn
	player2 *websocket.Conn
	broadcast chan Message
}

func NewRoom(p1, p2 *websocket.Conn) *Room {
	log.Println("Creating new game room.")

	return &Room{
		player1: p1,
		player2: p2,
		broadcast: make(chan Message),
	}
}

func (r *Room) Run() {
	go r.readMessages(r.player1)
	go r.readMessages(r.player2)

	r.player1.WriteMessage(websocket.TextMessage, []byte("Game matched! Starting game..."))
	r.player2.WriteMessage(websocket.TextMessage, []byte("Game matched! Starting game..."))

	for msg := range r.broadcast {
		if msg.Sender != r.player1 {
			if err := r.player1.WriteMessage(websocket.TextMessage, msg.Content); err != nil {
				r.player2.WriteMessage(websocket.TextMessage, []byte("Opponent has disconnected."))
			}
		}

		if msg.Sender != r.player2 {
			if err := r.player2.WriteMessage(websocket.TextMessage, msg.Content); err != nil {
				r.player1.WriteMessage(websocket.TextMessage, []byte("Opponent has disconnected."))
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