package repository

import "github.com/alelaca/chat-manager/src/entities"

type RepositoryMock struct {
	SavePostError         error
	AuthenticateUserError error
	AuthenticatedUser     entities.User
}

func (r *RepositoryMock) SavePost(post entities.Post) error {
	return r.SavePostError
}

func (r *RepositoryMock) AuthenticateUser(username, password string) (*entities.User, error) {
	return &r.AuthenticatedUser, r.AuthenticateUserError
}
