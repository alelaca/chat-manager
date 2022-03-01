package queues

type QueueMessageHandler = func(message QueueMessage) error

type QueueMessage struct {
	Message []byte
}

type Config struct {
	Name string
}

type QueuesHandler interface {
	PollPostsMessages(messageHandler QueueMessageHandler) error
}
