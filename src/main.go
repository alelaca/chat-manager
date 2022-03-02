package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alelaca/chat-manager/src/adapters/queues"
	"github.com/alelaca/chat-manager/src/adapters/routes"
	"github.com/alelaca/chat-manager/src/adapters/websocket"
	"github.com/alelaca/chat-manager/src/auth/jwt"
	"github.com/alelaca/chat-manager/src/config"
	rabbitmqqueue "github.com/alelaca/chat-manager/src/queues/rabbitmq"
	"github.com/alelaca/chat-manager/src/repository/mongodb"
	rabbitmqtopics "github.com/alelaca/chat-manager/src/topics/rabbitmq"
	"github.com/alelaca/chat-manager/src/usecases/post"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ScopeLocal      = "local"
	ScopeIntgDocker = "intg_docker"
)

func main() {
	profile := getProfile()

	configFileName := getConfigFileName(profile)
	envConfig, err := config.GetConfig(configFileName)
	if err != nil {
		log.Panicf(err.Error())
	}

	rabbitmqClient := ConnectRabbitMQ(envConfig)
	mongodbClient := ConnectMongoDB(envConfig)

	repository := mongodb.InitializeMongoDB(mongodbClient)

	authHandler := &jwt.Handler{
		Repository: repository,
	}

	queuesHandler := rabbitmqqueue.InitializeRabbitMQHandler(rabbitmqClient)
	topicsHandler := rabbitmqtopics.InitializeRabbitMQHandler(rabbitmqClient)

	postHandler := post.InitializePostHandler(topicsHandler, repository)

	websocketHandler := websocket.InitializeWebsocketHandler(postHandler, authHandler)
	worker := queues.InitializeWorker(queuesHandler, postHandler, websocketHandler)
	router := routes.InitializeRouter(websocketHandler, authHandler)

	websocketHandler.StartPool()
	worker.StartPollingPostsMessages()

	router.Run(":8080")
}

func getConfigFileName(profile string) string {
	return fmt.Sprintf("src/configfiles/config-%s.yml", strings.ToLower(profile))
}

func getProfile() string {
	scope := strings.ToLower(os.Getenv("SCOPE"))

	if scope != ScopeLocal && scope != ScopeIntgDocker {
		log.Panicf("Wrong env var SCOPE defined. It should be %s or %s", ScopeLocal, ScopeIntgDocker)
	}

	return scope
}

func ConnectRabbitMQ(envConfig config.Config) *amqp.Connection {
	connection, err := amqp.Dial(envConfig.RabbitMQ.URL)
	if err != nil {
		panic(fmt.Sprintf("error initializing rabbitmq connection, log: %s", err.Error()))
	}

	return connection
}

func ConnectMongoDB(envConfig config.Config) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(envConfig.MongoDB.URL))
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
