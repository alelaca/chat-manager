package main

import (
	"github.com/alelaca/chat-manager/src/adapter/routes"
	"github.com/alelaca/chat-manager/src/adapter/websocket"
	"github.com/alelaca/chat-manager/src/infrastructure/repository/local"
	"github.com/alelaca/chat-manager/src/usecases/post"
)

func main() {
	repository := local.CreateLocalMemory()
	postHandler := post.InitializePostHandler(repository)
	websocketHandler := websocket.InitializeWebsocketHandler(postHandler)
	router := routes.InitializeRouter(*websocketHandler)

	go websocketHandler.StartPool()
	router.Run(":8080")
}
