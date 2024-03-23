package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var clients = make(map[string]net.Conn)
var leaving = make(chan message)
var messages = make(chan message)

type message struct {
	text    string
	address string
}

func newMessage(msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    addr + msg,
		address: addr,
	}
}

func handleConnnection(conn net.Conn) {

	defer conn.Close()
	clients[conn.RemoteAddr().String()] = conn

	messages <- newMessage(" joined", conn)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- newMessage(input.Text(), conn)
	}

	//built in method that deletes an item from a map
	delete(clients, conn.RemoteAddr().String())

	leaving <- newMessage(" has left", conn)
}

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			for _, conn := range clients {
				if msg.address == conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(conn, msg.text) // NOTE: ignoring network errors
			}

		case msg := <-leaving:
			for _, conn := range clients {
				fmt.Fprintln(conn, msg.text) // NOTE: ignoring network errors
			}

		}
	}
}

func main() {

	//abre la conexión para poder conectar desde un cliente
	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}

	defer ln.Close()

	fmt.Println("Listening on port 3000")

	go broadcaster()

	for {

		conn, err := ln.Accept() //net.Conn This object represents the connection between the server and a client.
		if err != nil {
			fmt.Println("Connection refused")
			continue
		}

		go func() {
			welcomeMessage(conn)
		}()

		// buffer := make([]byte, 1024)

		go handleConnnection(conn)
	}
}

func welcomeMessage(conn net.Conn) {
	//Welcome message

	msg := "Welcome to the server\n"
	n, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Message not send to %v: %v", conn.RemoteAddr(), err)
		return
	}
	if n < len(msg) {
		log.Printf("Message not successfuly written: %d out of %d", n, len(msg))
	}
}
