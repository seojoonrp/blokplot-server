// main.go

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/websocket"
	"github.com/seojoonrp/blokplot-server/database"
	"github.com/seojoonrp/blokplot-server/handlers"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var matchmakingQueue = make(chan *websocket.Conn)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	log.Println("Successfully connected to MongoDB.")

	userStore := database.NewUserStore(client.Database("blokplot").Collection("users"))
	userHandler := &handlers.UserHandler{Store: userStore}

	http.HandleFunc("/api/login", userHandler.LoginOrRegisterHandler)

	go runMatchmaker()

	http.HandleFunc("/ws", handleConnections)

	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	log.Println("New client connected.")
	matchmakingQueue <- ws
}

func runMatchmaker() {
	var waitingPlayers []*websocket.Conn

	for {
		player := <- matchmakingQueue
		waitingPlayers = append(waitingPlayers, player)
		
		if len(waitingPlayers) >= 2 {
			player1 := waitingPlayers[0]
			player2 := waitingPlayers[1]
			waitingPlayers = waitingPlayers[2:]

			log.Println("2 Players matched. Creating game room...")
			go handleGameRoom(player1, player2)
		}
	}
}

func handleGameRoom(p1, p2 *websocket.Conn) {
	broadcast := make(chan []byte)

	go readMessages(p1, broadcast)
	go readMessages(p2, broadcast)

	p1.WriteMessage(websocket.TextMessage, []byte("Game Matched! Starting game..."))
	p2.WriteMessage(websocket.TextMessage, []byte("Game Matched! Starting game..."))

	for {
		msg := <- broadcast

		p1.WriteMessage(websocket.TextMessage, msg)
		p2.WriteMessage(websocket.TextMessage, msg)
	}
}

func readMessages(conn *websocket.Conn, broadcast chan<- []byte) {
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			break
		}
		broadcast <- msg
	}
}