package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	conn net.Conn
	data chan []byte
}

func main() {
	flagMode := flag.String("mode", "server", "Start in client or server mode")
	flag.Parse()
	if *flagMode == "server" {
		startServerMode()
	} else {
		startClientMode()
	}

}

func (m *ClientManager) startManager() {
	for {
		select {
		case client := <-m.register:
			m.clients[client] = true
			fmt.Printf("A client has joined the server. Connections: %v \n", len(m.clients))
		case client := <-m.unregister:
			delete(m.clients, client)
			fmt.Printf("A client has left the server. Connections: %v \n", len(m.clients))
		case message := <-m.broadcast:
			for client := range m.clients {
				client.data <- message
			}
		}
	}
}

func (m *ClientManager) receive(c *Client) {
	for {
		message := make([]byte, 4096)
		length, err := c.conn.Read(message)
		if err != nil {
			m.unregister <- c
			c.conn.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
			m.broadcast <- message
		}
	}
}

func (c *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := c.conn.Read(message)
		if err != nil {
			c.conn.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}

func (m *ClientManager) send(c *Client) {
	// defer client.socket.Close()
	for {
		select {
		case message, ok := <-c.data:
			if !ok {
				return
			}
			c.conn.Write(message)
		}
	}
}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", ":12345")
	if error != nil {
		fmt.Println(error)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.startManager()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		client := &Client{conn: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}

func startClientMode() {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	client := &Client{conn: connection}
	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(message))
	}
}
