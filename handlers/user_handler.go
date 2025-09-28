// handlers/user_handler.go

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/seojoonrp/blokplot-server/database"
	"github.com/seojoonrp/blokplot-server/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	Store *database.UserStore
}

func (handler *UserHandler) LoginOrRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		http.Error(w, "Username is required!", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	curUser, err := handler.Store.GetUserByUsername(ctx, req.Username)

	if err == mongo.ErrNoDocuments {
		log.Printf("User '%s' not found. Creating new user...", req.Username)

		newUser := models.User{
			Username: req.Username,
			Profile: models.Profile{
				AvatarUrl: "/avatars/default.png",
				BannerUrl: "/banners/default.png",
			},
			Stats: models.Stats{
				Wins:     0,
				LiarWins: 0,
				Trophies: 0,
				Coins:    0,
			},
			Friends: []string{},
		}

		insertErr := handler.Store.CreateUser(ctx, newUser)

		if insertErr != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
		return
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

	log.Printf("User '%s' found. Logging in...", curUser.Username)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(curUser)
}