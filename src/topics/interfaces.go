package topics

import "github.com/alelaca/chat-manager/src/entities"

type Config struct {
	Name string
}

type Usecases interface {
	NotifyMessage(post entities.Post) error
}
