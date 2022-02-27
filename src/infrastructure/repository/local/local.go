package local

import (
	"github.com/alelaca/chat-manager/src/entities"
)

type Handler struct {
	Posts map[string][]entities.Post // map rooms to list of posts
}

func CreateLocalMemory() *Handler {
	return &Handler{
		Posts: make(map[string][]entities.Post),
	}
}

func (h *Handler) SavePost(post entities.Post) error {
	if _, ok := h.Posts[post.Room]; !ok {
		h.Posts[post.Room] = []entities.Post{}
	}

	h.Posts[post.Room] = append(h.Posts[post.Room], post)

	return nil
}
