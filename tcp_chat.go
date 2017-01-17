package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type Message struct{ ip, message string }
type Client chan<- Message

var (
	joining  = make(chan Client)  // for client, which joins the chat
	leaving  = make(chan Client)  // for client, which leaves the chat
	messages = make(chan Message) // for all messages
)

func announcer() {
	clients := make(map[Client]bool) // all clients
	for {
		select {
		case msg := <-messages:
			// announce the message to all clients
			for client := range clients {
				client <- msg
			}

		case client := <-joining:
			// new client
			clients[client] = true

		case client := <-leaving:
			// delete and close when client leaves
			delete(clients, client)
			close(client)
		}
	}
}

// handle client connection
func clientConnection(conn net.Conn) {
	io.WriteString(conn, "\t\t\tWelcome to our tcp chat.\nPlease enter your name: ")
	address := conn.RemoteAddr().String()
	ch := make(chan Message)
	go writeMessage(conn, ch, address)

	reader := bufio.NewReader(conn)
	line, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	name := string(line)
	ch <- Message{"Initial message " + address, "Your name is " + name + ".\n\t\t\tYou can enter a message."}
	messages <- Message{address, "--" + name + " has joined the chat room at " + time.Now().Format("15:04:05") + "."}
	joining <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- Message{address, name + ": " + input.Text()}
	}

	leaving <- ch
	messages <- Message{address, "--" + name + " has left the chat room at " + time.Now().Format("15:04:05") + "."}
	conn.Close()
}

// write message to clients
func writeMessage(conn net.Conn, ch <-chan Message, ip string) {
	for msg := range ch {
		if strings.Contains(msg.ip, "Initial message") || msg.ip != ip {
			io.WriteString(conn, msg.message+"\n")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go announcer()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go clientConnection(conn)
	}
}
