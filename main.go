package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var clients []net.Conn

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

func client(conn net.Conn) {

	defer conn.Close()

	//Add the new connection to the clients slice
	clients = append(clients, conn)

	//Read get the data that client send
	buffer := make([]byte, 1024)
	reader := bufio.NewReader(conn)

	//Write into the conn
	writer := bufio.NewWriter(conn)

	for {
		//n, err := conn.Read(buffer) //in loop to read all the messages is receiveng
		n, err := reader.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection %v: %v", conn.RemoteAddr(), err)
			return
		}
		if n > 0 {
			received := string(buffer[:n])
			fmt.Printf("Received message from %v: %s", conn.RemoteAddr(), received)
		}
		n, err = writer.Write(buffer[:n])
		if err != nil {
			log.Printf("Error writing from connection %v: %v", conn.RemoteAddr(), err)
			return
		}
		if n > 0 {
			sent := string(buffer[:n])
			fmt.Printf("Writing message from %v: %s", conn.RemoteAddr(), sent)
		}

		err = writer.Flush()
		if err != nil {
			fmt.Printf("Error flushing: %v", err)
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

	for {

		conn, err := ln.Accept() //net.Conn This object represents the connection between the server and a client.
		if err != nil {
			fmt.Println("Connection refused")
		}

		go client(conn)

		go func() {
			welcomeMessage(conn)
		}()

	}
}
