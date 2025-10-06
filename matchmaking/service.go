// matchmaking/service.go

package matchmaking

import (
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/seojoonrp/blokplot-server/game"
)

type Service struct {
	queue chan *websocket.Conn
}

func NewService() *Service {
	return &Service{
		queue: make(chan *websocket.Conn),
	}
}

func (s *Service) AddPlayer(conn *websocket.Conn) {
	s.queue <- conn
}

func (s *Service) Run() {
	log.Println("Matchmaking service started.")

	var waitingPlayers []*websocket.Conn

	for {
		player := <- s.queue
		waitingPlayers = append(waitingPlayers, player)

		if len(waitingPlayers) >= 2 {
			player1 := waitingPlayers[0]
			player2 := waitingPlayers[1]

			if !isConnectionAlive(player1) {
				log.Println("Player 1 disconnected while waiting. Removing from queue...")
				waitingPlayers = waitingPlayers[1:]
				continue
			}
			if !isConnectionAlive(player2) {
				log.Println("Player 2 disconnected while waiting. Removing from queue...")
				waitingPlayers = append(waitingPlayers[:1], waitingPlayers[2:]...)
				continue
			}

			waitingPlayers = waitingPlayers[2:]
			log.Println("2 players matched. Creating game room...")

			room := game.NewRoom(player1, player2)
			go room.Run()
		}
	}
}

func isConnectionAlive(conn *websocket.Conn) bool {
	err := conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return false
	}

	err = conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		conn.SetWriteDeadline(time.Time{})
		return false
	}

	conn.SetWriteDeadline(time.Time{})
	return true
}