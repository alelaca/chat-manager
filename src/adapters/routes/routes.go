package routes

import (
	"net/http"

	"github.com/alelaca/chat-manager/src/adapters/websocket"
	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	WebsocketHandler websocket.Handler
}

// Initializes an engine to receive HTTP requests
func InitializeRouter(websocket websocket.Handler) *gin.Engine {
	h := Handler{
		WebsocketHandler: websocket,
	}

	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, apperrors.CreateAPIError(http.StatusNotFound, "resource not found"))
	})

	router.GET("/ws", h.WebsocketHandler.Connect)

	return router
}
