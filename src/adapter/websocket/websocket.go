package websocket

import (
	"fmt"
	"net/http"

	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/alelaca/chat-manager/src/entities"
	"github.com/alelaca/chat-manager/src/entities/dtos"
	"github.com/alelaca/chat-manager/src/usecases/post"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const SenderSystem = "System"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // authentication coming in next versions
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	Pool        *Pool
	PostHandler post.Usecases
}

type client struct {
	ID          string
	Connection  *websocket.Conn
	Pool        *Pool
	PostHandler post.Usecases
}

func InitializeWebsocketHandler(postHandler post.Usecases) *Handler {
	pool := InitializeClientPool()
	return &Handler{
		Pool:        pool,
		PostHandler: postHandler,
	}
}

func (h *Handler) StartPool() {
	h.Pool.Start()
}

// Handles HTTP requests and enables a WebSocket connection to a client
func (h *Handler) Connect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error upgrading") // TODO change this to an error somewhere
	}

	client := &client{
		ID:          uuid.NewString(),
		Connection:  ws,
		Pool:        h.Pool,
		PostHandler: h.PostHandler,
	}

	client.Pool.Subscribe <- client

	client.listen()
}

// Listen to all incoming messages from client
func (c *client) listen() {
	defer func() {
		c.Pool.Unsubscribe <- c
	}()

	for {
		postDTO := dtos.PostDTO{}
		err := c.Connection.ReadJSON(&postDTO)
		if err != nil {
			fmt.Println("error reading message", err.Error())
			break
		}

		post, err := c.PostHandler.CreatePost(postDTO)
		if err != nil {
			fmt.Println(fmt.Sprintf("error creating post: '%s'", err.Error()))
			c.logErrorToClient(err, postDTO)
			continue
		}

		c.Pool.Broadcast <- *post
	}
}

// Send message to client
func (c *client) sendMessage(post entities.Post) {
	err := c.Connection.WriteJSON(post)
	if err != nil {
		fmt.Println("error sending message")
	}
}

// Send error messages to client
func (c *client) logErrorToClient(err error, postDTO dtos.PostDTO) {
	apiError, ok := err.(apperrors.APIError)
	if !ok || apiError.StatusCode < 500 {
		return
	}

	errorMessage := fmt.Sprintf("Error processing your message, please try again. Message not sent: '%s'", postDTO.Message)
	postDTO.Message = errorMessage

	post, err := entities.CreatePost(postDTO)
	if err != nil {
		fmt.Println("error sending error message to user")
		return
	}

	c.sendMessage(post)
}