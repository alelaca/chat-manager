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
	"go.mongodb.org/mongo-driver/bson"
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

	if profile == ScopeIntgDocker || profile == ScopeLocal {
		InitializeMongoData(mongodbClient)
		InitializeRabbitMQInfra(rabbitmqClient, envConfig)
	}

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

// Methods only for local and docker compose!! because docker files couldn't be loaded directly from docker
func InitializeMongoData(mongo *mongo.Client) {
	mongo.Database("chatdb").Collection("users").InsertMany(context.Background(), []interface{}{
		bson.M{"_id": "1221a07d-6ee9-4d60-92e7-46850e7565a2", "username": "Bob", "password": "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"},
		bson.M{"_id": "3f68696f-cf68-4d8e-a440-05230f3e1a4f", "username": "Alice", "password": "4cc8f4d609b717356701c57a03e737e5ac8fe885da8c7163d3de47e01849c635"},
		bson.M{"_id": "1cd8e021-1921-4d85-abe9-323acf326656", "username": "Jane", "password": "68487dc295052aa79c530e283ce698b8c6bb1b42ff0944252e1910dbecdc5425"},
		bson.M{"_id": "f3690232-27c0-425c-9962-8828d32ee4eb", "username": "Justin", "password": "69f7f7a7f8bca9970fa6f9c0b8dad06901d3ef23fd599d3213aa5eee5621c3e3"},
	})
}

// Methods only for local and docker compose!! because docker files couldn't be loaded directly from docker
func InitializeRabbitMQInfra(rabbit *amqp.Connection, envConfig config.Config) {
	ch, _ := rabbit.Channel()
	defer ch.Close()

	ch.ExchangeDeclare(
		envConfig.RabbitMQ.MessagesTopic.Name, // name
		"topic",                               // type
		true,                                  // durable
		false,                                 // auto-deleted
		false,                                 // internal
		false,                                 // no-wait
		nil,                                   // arguments
	)

	ch.QueueDeclare(
		envConfig.RabbitMQ.PostsQueue.Name, // name
		true,                               // durable
		false,                              // delete when unused
		false,                              // exclusive
		false,                              // no-wait
		nil,                                // arguments
	)

	ch.QueueDeclare(
		envConfig.RabbitMQ.CommandQueue.Name, // name
		true,                                 // durable
		false,                                // delete when unused
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)

	ch.QueueBind(
		envConfig.RabbitMQ.CommandQueue.Name,  // queue name
		"stock.#",                             // routing key
		envConfig.RabbitMQ.MessagesTopic.Name, // exchange
		false,
		nil)
}
