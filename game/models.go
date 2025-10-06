package game

import "encoding/json"

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type BaseMessage struct {
	Type string `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ChatData struct {
	Message string `json:"message"`
}

type BlockData struct {
	Position Vector3 `json:"position"`
	ColorIndex int `json:"colorIndex"`
	Type string `json:"type"`
}