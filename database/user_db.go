// databse/user_db.go

package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/seojoonrp/blokplot-server/models"
)

type UserStore struct {
	Collection *mongo.Collection
}

func NewUserStore(collection *mongo.Collection) *UserStore {
	return &UserStore{Collection: collection}
}

func (store *UserStore) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	filter := bson.M{"username": username}
	err := store.Collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (store *UserStore) CreateUser(ctx context.Context, user models.User) error {
	_, err := store.Collection.InsertOne(ctx, user)

	return err
}