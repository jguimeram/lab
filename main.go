package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var clients []net.Conn

func handleConnection(conn net.Conn) {

	defer conn.Close()

	//Add the new connection to the clients slice
	clients = append(clients, conn)

	//Welcome message
	msg := "Welcome to the server"
	n, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Message not send to %v: %v", conn.RemoteAddr(), err)
		return
	}
	if n < len(msg) {
		log.Printf("Message not successfuly written: %d out of %d", n, len(msg))
	}

	//Read get the data that client send
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer) //in loop to read all the messages is receiveng
		if err != nil {
			log.Printf("Error reading from connection %v: %v", conn.RemoteAddr(), err)
			return
		}
		if n > 0 {
			received := string(buffer[:n])
			fmt.Printf("Received message from %v: %s", conn.RemoteAddr(), received)
		}

	}

}

func listClients() {
	fmt.Println("Clients connected:")
	for i, client := range clients {
		fmt.Printf("%d, %v\n", i, client.RemoteAddr())
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

	for {
		//permite al cliente conectar y devuelve objeto net.Conn
		// Accept va en un loop infinito para poder atender múltiples conexiones. De otro modo, tal y como acepta una, bloquea y ya no recibe más.

		conn, err := ln.Accept() //net.Conn This object represents the connection between the server and a client.
		if err != nil {
			fmt.Println("Connection refused")
		}

		//for each connection a routine handleConnection is created
		go handleConnection(conn)

		//A new goroutine is started using an anonymous function (go func() { ... }()).
		go func() {
			for {
				listClients()
				time.Sleep(10 * time.Second)
			}
		}()

	}
}
