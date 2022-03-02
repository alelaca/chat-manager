package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alelaca/chat-manager/src/adapters/websocket"
	"github.com/alelaca/chat-manager/src/apperrors"
	"github.com/alelaca/chat-manager/src/auth"
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Handler struct {
	WebsocketHandler *websocket.Handler
	AuthHandler      auth.Handler
}

// Initializes an engine to receive HTTP requests
func InitializeRouter(websocket *websocket.Handler, authHandler auth.Handler) *gin.Engine {
	h := Handler{
		WebsocketHandler: websocket,
		AuthHandler:      authHandler,
	}

	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, apperrors.CreateAPIError(http.StatusNotFound, "resource not found"))
	})

	// group created to apply middleware rules
	restGroup := router.Group("/api")
	restGroup.Use(configureCors)

	router.GET("/ws", h.WebsocketHandler.Connect)
	restGroup.POST("/login", h.Login)
	restGroup.POST("/auth", h.Auth)

	return router
}

func (h *Handler) Login(c *gin.Context) {
	credentials := Credentials{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials structure"})
		return
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong credentials structure"})
		return
	}

	token, err := h.AuthHandler.Login(credentials.Username, credentials.Password)
	if err != nil {
		abortWithCustomError(c, http.StatusUnauthorized, err)
		return
	}

	// cant set httpOnly cookies for websockets
	c.JSON(http.StatusOK, token)
}

func (h *Handler) Auth(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "error reading body"})
		return
	}

	err = h.AuthHandler.Authenticate(string(body))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func abortWithCustomError(c *gin.Context, defaultStatus int, err error) {
	status := defaultStatus
	if apiError, ok := err.(apperrors.APIError); ok && apiError.StatusCode != 0 {
		c.AbortWithStatusJSON(apiError.StatusCode, apiError)
		status = apiError.StatusCode
	} else {
		c.AbortWithStatusJSON(defaultStatus, apperrors.CreateAPIError(defaultStatus, err.Error()))
	}

	fmt.Println(fmt.Sprintf("[ERROR][status:%d] %s", status, err.Error()))
}
