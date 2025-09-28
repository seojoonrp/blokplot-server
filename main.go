// main.go

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/seojoonrp/blokplot-server/database"
	"github.com/seojoonrp/blokplot-server/handlers"
)

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

	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}