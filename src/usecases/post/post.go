package post

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/alelaca/chat-manager/src/entities"
	"github.com/alelaca/chat-manager/src/entities/dtos"
	"github.com/alelaca/chat-manager/src/topics"
	"github.com/alelaca/chat-manager/src/usecases/interfaces"
)

const SystemSender = "System"

type Usecases interface {
	CreatePost(postDTO dtos.PostDTO) (*entities.Post, error)
}

type Handler struct {
	TopicsHandler topics.Usecases
	Repository    interfaces.Repository
}

func InitializePostHandler(topicsHandler topics.Usecases, repository interfaces.Repository) *Handler {
	return &Handler{
		TopicsHandler: topicsHandler,
		Repository:    repository,
	}
}

// Creates a Post and saves it in repository
// Send notification about the message to other APIs
func (h *Handler) CreatePost(postDTO dtos.PostDTO) (*entities.Post, error) {
	post, err := entities.CreatePost(postDTO)
	if err != nil {
		return nil, apperrors.CreateAPIError(http.StatusBadRequest, err.Error())
	}

	err = h.Repository.SavePost(post)
	if err != nil {
		return nil, err
	}

	if isCommandPost(post.Message) {
		err := h.TopicsHandler.NotifyMessage(post)
		if err != nil {
			return nil, err
		}
	}

	return &post, nil
}

// TODO change this method to another file (or maybe not)
func createErrorPost(message string, room string) entities.Post {
	postErr := dtos.PostDTO{
		Message: message,
		Sender:  SystemSender,
		Room:    room,
	}
	post, err := entities.CreatePost(postErr)
	if err != nil {
		fmt.Println("error creating error post")
	}

	return post
}

func isCommandPost(message string) bool {
	return strings.HasPrefix(message, "/")
}
