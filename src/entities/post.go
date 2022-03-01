package entities

import (
	"time"

	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/alelaca/chat-manager/src/entities/dtos"
	"github.com/google/uuid"
)

type Post struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	Sender      string    `json:"sender"`
	Room        string    `json:"room"`
	CreatedDate time.Time `json:"created_date"`
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// Creates Post entity
func CreatePost(postDTO dtos.PostDTO) (Post, error) {
	if err := ValidatePost(postDTO); err != nil {
		return Post{}, err
	}

	return Post{
		ID:          uuid.NewString(),
		Message:     postDTO.Message,
		Sender:      postDTO.Sender,
		Room:        postDTO.Room,
		CreatedDate: time.Now(),
	}, nil
}

// Validate Post data
func ValidatePost(postDTO dtos.PostDTO) error {
	if postDTO.Message == "" {
		return apperrors.NewPropertyRequiredError("message", "message")
	}

	if postDTO.Sender == "" {
		return apperrors.NewPropertyRequiredError("message", "sender")
	}

	if postDTO.Room == "" {
		return apperrors.NewPropertyRequiredError("message", "room")
	}

	return nil
}
