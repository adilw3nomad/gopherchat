package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/adilw3nomad/gopherchat/client"
	"github.com/adilw3nomad/gopherchat/clientManager"
)

func main() {
	flagMode := flag.String("mode", "server", "Start in client or server mode")
	flag.Parse()
	if *flagMode == "server" {
		startServerMode()
	} else {
		startClientMode()
	}

}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", ":12345")
	if error != nil {
		fmt.Println(error)
	}
	manager := clientManager.NewClientManager()
	go manager.Start()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		client := &client.Client{Conn: connection, Data: make(chan []byte)}
		manager.Register <- client
		go manager.Receive(client)
		go client.Send()
	}
}

func startClientMode() {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	client := &client.Client{Conn: connection}
	go client.Receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(message))
	}
}
