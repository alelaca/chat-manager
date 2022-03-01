package repository

import "github.com/alelaca/chat-manager/src/entities"

type RepositoryMock struct {
	SavePostError error
}

func (r *RepositoryMock) SavePost(post entities.Post) error {
	return r.SavePostError
}
