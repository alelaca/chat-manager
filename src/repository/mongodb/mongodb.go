package mongodb

import (
	"context"

	"github.com/alelaca/chat-manager/src/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseName    = "chatdb"
	PostsCollection = "posts"
	UsersCollection = "users"
)

type Handler struct {
	Client          *mongo.Client
	database        *mongo.Database
	postsCollection *mongo.Collection
}

func InitializeMongoDB(client *mongo.Client) *Handler {
	return &Handler{
		Client: client,
	}
}

func (h *Handler) SavePost(post entities.Post) error {
	_, err := h.Client.Database(DatabaseName).Collection(PostsCollection).InsertOne(context.Background(), post)
	return err
}

func (h *Handler) AuthenticateUser(username, password string) (*entities.User, error) {
	ctx := context.Background()
	user := entities.User{}

	result, err := h.Client.Database(DatabaseName).Collection(UsersCollection).Find(ctx, bson.M{"username": username, "password": password})
	if err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	if !result.Next(ctx) {
		return nil, nil // not found
	}

	if err = result.Decode(&user); err != nil {
		return nil, err
	}

	return &user, err
}
