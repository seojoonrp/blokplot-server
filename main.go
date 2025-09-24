package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

type User struct {
	Username string `bson:"username" json:"username"`
	Profile interface{} `bson:"profile" json:"profile"`
	Stats interface{} `bson:"stats" json:"stats"`
	Friends []string `bson:"friends" json:"friends"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	mongoClient = client
	log.Println("Successfully connected to MongoDB.")

	http.HandleFunc("/api/users/", getUserHandler)

	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/api/users/"):]

	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	collection := mongoClient.Database("blokplot").Collection("users")

	var result User
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("Error while getting DB:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}