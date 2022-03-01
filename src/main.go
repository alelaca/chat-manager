package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alelaca/chat-manager/src/adapters/queues"
	"github.com/alelaca/chat-manager/src/adapters/routes"
	"github.com/alelaca/chat-manager/src/adapters/websocket"
	"github.com/alelaca/chat-manager/src/auth/jwt"
	rabbitmqqueue "github.com/alelaca/chat-manager/src/queues/rabbitmq"
	"github.com/alelaca/chat-manager/src/repository/mongodb"
	rabbitmqtopics "github.com/alelaca/chat-manager/src/topics/rabbitmq"
	"github.com/alelaca/chat-manager/src/usecases/post"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	authHandler := &jwt.Handler{}
	rabbitmqClient := ConnectRabbitMQ()
	mongodbClient := ConnectMongoDB()

	repository := mongodb.InitializeMongoDB(mongodbClient)

	queuesHandler := rabbitmqqueue.InitializeRabbitMQHandler(rabbitmqClient)
	topicsHandler := rabbitmqtopics.InitializeRabbitMQHandler(rabbitmqClient)

	postHandler := post.InitializePostHandler(topicsHandler, repository)

	websocketHandler := websocket.InitializeWebsocketHandler(postHandler)
	worker := queues.InitializeWorker(queuesHandler, postHandler, websocketHandler)
	router := routes.InitializeRouter(websocketHandler, authHandler)

	websocketHandler.StartPool()
	worker.StartPollingPostsMessages()

	router.Run(":8080")
}

func ConnectRabbitMQ() *amqp.Connection {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(fmt.Sprintf("error initializing rabbitmq connection, log: %s", err.Error()))
	}

	return connection
}

func ConnectMongoDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	return client
}
