package websocket

import (
	"fmt"

	"github.com/alelaca/chat-manager/src/entities"
)

type Pool struct {
	Subscribe   chan *client
	Unsubscribe chan *client
	Broadcast   chan entities.Post
	Clients     map[*client]bool
}

func InitializeClientPool() *Pool {
	return &Pool{
		Subscribe:   make(chan *client),
		Unsubscribe: make(chan *client),
		Broadcast:   make(chan entities.Post),
		Clients:     make(map[*client]bool),
	}
}

// Starts a listener for different events from the clients
// - Adds new clients to the pool
// - Remove disconnected clients from the pool
// - Broadcast incoming messages to all clients
func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Subscribe:
			p.Clients[client] = true
			p.broadcast(entities.Post{Message: fmt.Sprintf("Client %p connected", client), Sender: "System"})
			break
		case client := <-p.Unsubscribe:
			delete(p.Clients, client)
			p.broadcast(entities.Post{Message: fmt.Sprintf("Client %p disconnected", client), Sender: "System"})
			client.Connection.Close()
			break
		case message := <-p.Broadcast:
			p.broadcast(message)
			break
		}
	}
}

func (p *Pool) broadcast(post entities.Post) {
	for client := range p.Clients {
		client.sendMessage(post)
	}
}
