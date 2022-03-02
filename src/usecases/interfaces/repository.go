package interfaces

import "github.com/alelaca/chat-manager/src/entities"

type Repository interface {
	SavePost(post entities.Post) error
	AuthenticateUser(username, password string) (*entities.User, error)
}
