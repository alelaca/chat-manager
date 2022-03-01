package rabbitmq

import (
	"encoding/json"
	"strings"

	"github.com/alelaca/chat-manager/src/entities"
	"github.com/alelaca/chat-manager/src/topics"
	"github.com/streadway/amqp"
)

type Handler struct {
	Connection    *amqp.Connection
	messagesTopic topics.Config
}

func InitializeRabbitMQHandler(connection *amqp.Connection) *Handler {
	messagesTopic := topics.Config{
		Name: "messages-topic",
	}

	return &Handler{
		Connection:    connection,
		messagesTopic: messagesTopic,
	}
}

// Sends command results to topic
func (h *Handler) NotifyMessage(post entities.Post) error {
	postBody, err := json.Marshal(post)
	if err != nil {
		return err
	}

	filter := ""
	if strings.HasPrefix(post.Message, "/") {
		filter = strings.Replace(strings.TrimPrefix(post.Message, "/"), "=", ".", 1)
	}

	return h.sendMessage(h.messagesTopic, filter, postBody)
}

// Sends a message to the specified topic
// If filter is set, it applies the filter in the topic. If filter is empty, it sends message to queues without key
func (h *Handler) sendMessage(topic topics.Config, filter string, body []byte) error {
	ch, err := h.Connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.Publish(
		topic.Name, // exchange
		filter,     // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		return err
	}

	return nil
}
