package post

import (
	"net/http"

	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/alelaca/chat-manager/src/entities"
	"github.com/alelaca/chat-manager/src/entities/dtos"
	"github.com/alelaca/chat-manager/src/usecases/interfaces"
)

type Usecases interface {
	CreatePost(postDTO dtos.PostDTO) (*entities.Post, error)
}

type Handler struct {
	Repository interfaces.Repository
}

func InitializePostHandler(repository interfaces.Repository) *Handler {
	return &Handler{
		Repository: repository,
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

	// TODO send to queue

	return &post, nil
}
