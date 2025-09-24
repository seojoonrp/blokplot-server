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

type Profile struct {
	AvatarUrl string `bson:"avatarUrl" json:"avatarUrl"`
	BannerUrl string `bson:"bannerUrl" json:"bannerUrl"`
}

type Stats struct {
	Wins int `bson:"wins" json:"wins"`
	LiarWins int `bson:"liar-wins" json:"liar-wins"`
	Trophies int `bson:"trophies" json:"trophies"`
	Coins int `bson:"coins" json:"coins"`
}

type User struct {
	Username string `bson:"username" json:"username"`
	Profile Profile `bson:"profile" json:"profile"`
	Stats Stats `bson:"stats" json:"stats"`
	Friends []string `bson:"friends" json:"friends"`
}

type LoginRequest struct {
	Username string `json:"username"`
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

	http.HandleFunc("/api/login", loginOrRegisterHandler)

	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginOrRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req);

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	collection := mongoClient.Database("blokplot").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	var existingUser User
	filter := bson.M{"username": req.Username}
	err = collection.FindOne(ctx, filter).Decode(&existingUser)

	if err == mongo.ErrNoDocuments {
		log.Printf("User '%s' not found. Creating new user...", req.Username)

		newUser := User{
			Username: req.Username,
			Profile: Profile{
				AvatarUrl: "/avatars/default.png",
				BannerUrl: "/banners/default.png",
			},
			Stats: Stats{
				Wins:     0,
				LiarWins: 0,
				Trophies: 0,
				Coins:    0,
			},
			Friends: []string{},
		}

		_, insertErr := collection.InsertOne(ctx, newUser)
		if insertErr != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
		return
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("User '%s' found. Logging in...", existingUser.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingUser)
}