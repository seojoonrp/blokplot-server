// game/models.go

package game

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn *websocket.Conn
	GameIndex int
}

type ClientMessage struct {
	Type string `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ServerMessage struct {
	Type string `json:"type"`
	SenderIndex int `json:"senderIndex"`
	Data json.RawMessage `json:"data"`
}

type GameInitMessage struct {
	PlayerIndex int `json:"playerIndex"`
}

type ChatData struct {
	Message string `json:"message"`
}

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type BlockPlaceData struct {
	Position Vector3 `json:"position"`
	ColorIndex int `json:"colorIndex"`
	BlockType string `json:"blockType"`
}