package clientManager

import (
	"fmt"

	"github.com/adilw3nomad/gopherchat/client"
)

// ClientManager keeps track of connected clients and broadcasts any messages that it receives
type ClientManager struct {
	clients    map[*client.Client]bool
	broadcast  chan []byte
	Register   chan *client.Client
	unregister chan *client.Client
}

// start is a go routine that deals with receiving values from the chainManager channels
func (m *ClientManager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.clients[client] = true
			fmt.Printf("A client has joined the server. Connections: %v \n", len(m.clients))
		case client := <-m.unregister:
			delete(m.clients, client)
			fmt.Printf("A client has left the server. Connections: %v \n", len(m.clients))
		case message := <-m.broadcast:
			for client := range m.clients {
				client.Data <- message
			}
		}
	}
}

func (m *ClientManager) Receive(c *client.Client) {
	for {
		message := make([]byte, 4096)
		length, err := c.Conn.Read(message)
		if err != nil {
			m.unregister <- c
			c.Conn.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
			m.broadcast <- message
		}
	}
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*client.Client]bool),
		broadcast:  make(chan []byte),
		Register:   make(chan *client.Client),
		unregister: make(chan *client.Client),
	}
}
