package main

import (
	"fmt"

	"github.com/alelaca/chat-manager/src/adapters/queues"
	"github.com/alelaca/chat-manager/src/adapters/routes"
	"github.com/alelaca/chat-manager/src/adapters/websocket"
	rabbitmqqueue "github.com/alelaca/chat-manager/src/queues/rabbitmq"
	"github.com/alelaca/chat-manager/src/repository/local"
	rabbitmqtopics "github.com/alelaca/chat-manager/src/topics/rabbitmq"
	"github.com/alelaca/chat-manager/src/usecases/post"
	"github.com/streadway/amqp"
)

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(fmt.Sprintf("error initializing rabbitmq connection, log: %s", err.Error()))
	}

	repository := local.CreateLocalMemory()

	queuesHandler := rabbitmqqueue.InitializeRabbitMQHandler(connection)
	topicsHandler := rabbitmqtopics.InitializeRabbitMQHandler(connection)

	postHandler := post.InitializePostHandler(topicsHandler, repository)

	websocketHandler := websocket.InitializeWebsocketHandler(postHandler)
	worker := queues.InitializeWorker(queuesHandler, postHandler, websocketHandler)
	router := routes.InitializeRouter(*websocketHandler)

	websocketHandler.StartPool()
	worker.StartPollingPostsMessages()

	router.Run(":8080")
}
