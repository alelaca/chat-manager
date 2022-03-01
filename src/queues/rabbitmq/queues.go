package rabbitmq

import (
	"fmt"

	"github.com/alelaca/chat-manager/src/queues"
	"github.com/streadway/amqp"
)

type Handler struct {
	Connection *amqp.Connection
	postsQueue queues.Config
}

func InitializeRabbitMQHandler(connection *amqp.Connection) *Handler {
	postsQueue := queues.Config{
		Name: "posts-queue",
	}

	return &Handler{
		Connection: connection,
		postsQueue: postsQueue,
	}
}

// Poll messages from posts queue
func (h *Handler) PollPostsMessages(messageHandler queues.QueueMessageHandler) error {
	return h.receiveMessages(h.postsQueue, messageHandler)
}

// Receives messages from specified queue
// Handle any messages received with parameter function
func (h *Handler) receiveMessages(queue queues.Config, messageHandler queues.QueueMessageHandler) error {
	ch, err := h.Connection.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			queueMessage := queues.QueueMessage{
				Message: msg.Body,
			}

			err := messageHandler(queueMessage)
			if err != nil {
				fmt.Println(err.Error())
				msg.Reject(true)
				continue
			}

			msg.Ack(false)
		}
	}()

	<-forever

	return nil
}
