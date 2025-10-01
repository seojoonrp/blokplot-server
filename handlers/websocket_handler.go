package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/seojoonrp/blokplot-server/matchmaking"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true;
	},
}

type WebSocketHandler struct {
	MatchmakingService *matchmaking.Service
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	log.Println("New player had connected. Adding to the matching queue...")

	h.MatchmakingService.AddPlayer(ws)
}