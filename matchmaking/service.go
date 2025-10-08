// matchmaking/service.go

package matchmaking

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/seojoonrp/blokplot-server/game"
)

const maxPlayerCount = 2;

type RoomStatusMessage struct {
	Type string `json:"type"`
	Data RoomStatusData `json:"data"`
}

type RoomStatusData struct {
	CurPlayerCount int `json:"curPlayerCount"`
	MaxPlayerCount int `json:"maxPlayerCount"`
}

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
		waitingPlayers = broadcastAndClean(waitingPlayers)

		if len(waitingPlayers) >= maxPlayerCount {
			matchedPlayers := waitingPlayers[:maxPlayerCount];
			waitingPlayers = waitingPlayers[maxPlayerCount:];
			log.Printf("%d players matched. Creating new room...", maxPlayerCount)

			room := game.NewRoom(matchedPlayers...)
			go room.Run()

			if (len(waitingPlayers) > 0) {
				waitingPlayers = broadcastAndClean(waitingPlayers)
			}
		}
	}
}

func broadcastAndClean(players []*websocket.Conn) []*websocket.Conn {
	alivePlayers := make([]*websocket.Conn, 0, len(players))

	for _, player := range players {
		if isConnectionAlive(player) {
			alivePlayers = append(alivePlayers, player)
		} else {
			log.Println("Some player has disconnected. Removing from queue.")
		}
	}

	msg := RoomStatusMessage{
		Type: "roomStatus",
		Data: RoomStatusData{
			CurPlayerCount: len(alivePlayers),
			MaxPlayerCount: maxPlayerCount,
		},
	}
	jsonMsg, _ := json.Marshal(msg)

	for _, player := range alivePlayers {
		player.WriteMessage(websocket.TextMessage, jsonMsg)
	}

	return alivePlayers
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