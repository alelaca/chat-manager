package topics

import "github.com/alelaca/chat-manager/src/entities"

type TopicMock struct {
	NotifyMessageError error
}

func (m *TopicMock) NotifyMessage(post entities.Post) error {
	return m.NotifyMessageError
}
