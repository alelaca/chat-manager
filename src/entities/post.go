package entities

import (
	"time"

	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/google/uuid"
)

type Post struct {
	ID          string
	Message     string
	Sender      string
	Room        string
	CreatedDate time.Time
}

type Room struct {
	ID       string
	Name     string
	Messages []Post
}

type User struct {
	ID       string
	Nickname string
	Password string
}

func CreatePost(message string, sender string, room string) (Post, error) {
	if err := ValidatePost(message, sender, room); err != nil {
		return Post{}, err
	}

	return Post{
		ID:          uuid.NewString(),
		Message:     message,
		Sender:      sender,
		Room:        room,
		CreatedDate: time.Now(),
	}, nil
}

func ValidatePost(message string, sender string, room string) error {
	if message == "" {
		return apperrors.NewPropertyRequiredError("message", message)
	}

	if sender == "" {
		return apperrors.NewPropertyRequiredError("message", "sender")
	}

	if room == "" {
		return apperrors.NewPropertyRequiredError("message", "room")
	}

	return nil
}
