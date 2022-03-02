package websocket

import (
	"fmt"
	"net/http"

	"github.com/alelaca/chat-manager/src/auth"
	"github.com/alelaca/chat-manager/src/entities"
	"github.com/alelaca/chat-manager/src/entities/dtos"
	"github.com/alelaca/chat-manager/src/usecases/post"
	"github.com/gin-gonic/gin"
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
	AuthHandler auth.Handler
}

type client struct {
	Username    string
	Connection  *websocket.Conn
	Pool        *Pool
	PostHandler post.Usecases
	AuthHandler auth.Handler
}

type AuthenticatedPost struct {
	Token string       `json:"token"`
	Post  dtos.PostDTO `json:"post"`
}

func InitializeWebsocketHandler(postHandler post.Usecases, authHandler auth.Handler) *Handler {
	pool := InitializeClientPool()
	return &Handler{
		Pool:        pool,
		PostHandler: postHandler,
		AuthHandler: authHandler,
	}
}

func (h *Handler) StartPool() {
	go func() {
		h.Pool.Start()
	}()
}

// Handles HTTP requests and enables a WebSocket connection to a client
func (h *Handler) Connect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error upgrading")
	}

	client := &client{
		Connection:  ws,
		Pool:        h.Pool,
		PostHandler: h.PostHandler,
		AuthHandler: h.AuthHandler,
	}

	client.Pool.Subscribe <- client

	client.listen()
}

// Broadcast message to all connected Websockets
func (h *Handler) Broadcast(post entities.Post) {
	h.Pool.Broadcast <- post
}

// Listen to all incoming messages from client
func (c *client) listen() {
	defer func() {
		c.Pool.Unsubscribe <- c
	}()

	for {
		authPost := AuthenticatedPost{}
		err := c.Connection.ReadJSON(&authPost)
		if err != nil {
			fmt.Println("error reading message", err.Error())
			break
		}

		fmt.Println(authPost)

		err = c.AuthHandler.Authenticate(authPost.Token)
		if err != nil {
			fmt.Println("unauth")
			c.Pool.Unsubscribe <- c
		}

		post, err := c.PostHandler.CreatePost(authPost.Post)
		if err != nil {
			fmt.Println(fmt.Sprintf("error creating post: '%s'", err.Error()))
			c.logErrorToClient(err, authPost.Post)
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
	errorMessage := fmt.Sprintf("Error processing your message, please try again. Message not sent: '%s'", postDTO.Message)
	postDTO.Message = errorMessage

	post, err := entities.CreatePost(postDTO)
	if err != nil {
		fmt.Println("error sending error message to user")
		return
	}

	c.sendMessage(post)
}
