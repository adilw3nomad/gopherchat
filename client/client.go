package client

import (
	"fmt"
	"net"
)

// Client Connects to a server and has a channel for receiving message data from the connection
type Client struct {
	Conn net.Conn
	Data chan []byte
}

// receive is a go routine which reads the data from the connection and prints it
func (c *Client) Receive() {
	for {
		message := make([]byte, 4096)
		length, err := c.Conn.Read(message)
		if err != nil {
			c.Conn.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}

func (c *Client) Send() {
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.Data:
			if !ok {
				return
			}
			c.Conn.Write(message)
		}
	}
}
