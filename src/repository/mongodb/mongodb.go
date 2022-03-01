package mongodb

import (
	"context"

	"github.com/alelaca/chat-manager/src/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseName    = "chatdb"
	PostsCollection = "posts"
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
