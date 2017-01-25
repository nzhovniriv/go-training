package main

import (
	"bufio"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_ROOM string = "global"
const DEFAULT_CLIENT string = "Guest"

type MessageFormat struct{ uuid, message string }
type Message map[string]MessageFormat
type Client chan Message
type MainActivity map[string]Client

var (
	joining           = make(chan MainActivity)            // for client, which joins the chat
	leaving           = make(chan MainActivity)            // for client, which leaves the chat
	messages          = make(chan Message)                 // for all messages
	clients           = make(map[string][]map[Client]bool) // all clients
	number            = make(map[string]int)               // number of messages
	rooms    []string = []string{DEFAULT_ROOM}             // number of rooms
)

func announcer() {
	for {
		select {
		case msg := <-messages:
			var curRoom string
			for room, _ := range msg {
				curRoom = room
				if !strings.HasPrefix(msg[room].message, "--") {
					number[room]++
				}
			}
			// add new room
			if !isValueInSlice(rooms, curRoom) {
				rooms = append(rooms, curRoom)
			}
			// announce the message to all clients
			for _, client := range clients[curRoom] {
				for ch, _ := range client {
					ch <- msg
				}
			}

		case client := <-joining:
			// new client
			for room, ch := range client {
				cl := map[Client]bool{ch: true}
				clients[room] = append(clients[room], cl)
			}

		case client := <-leaving:
			// delete and close when client leaves
			for room, ch := range client {
				for sliceIdx, sliceVal := range clients[room] {
					if _, ok := sliceVal[ch]; ok {
						close(ch)
						clients[room] = removeElement(clients[room], sliceIdx)
					}
				}
			}
		}
	}
}

// handle client connection
func clientConnection(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "\t\t\tWelcome to our tcp chat.\n"+showRooms()+"You can choose the one or enter your own: ")

	reader := bufio.NewReader(conn)

	chatRoom, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	currentRoom := string(chatRoom)
	if len(currentRoom) == 0 {
		currentRoom = DEFAULT_ROOM
	}

	uuid := uuid.NewV4().String()
	ch := make(chan Message)
	go writeMessage(conn, ch, uuid, currentRoom)

	ch <- Message{currentRoom: MessageFormat{"Welcome message " + uuid, "You are in " + currentRoom + " chat room."}}

	client_name, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	name := string(client_name)
	if len(name) == 0 {
		name = DEFAULT_CLIENT
	}

	ch <- Message{currentRoom: MessageFormat{"Initial message " + uuid, "Your name is " + name + ".\n\t\t\tYou can enter a message."}}
	messages <- Message{currentRoom: {uuid, "--" + name + " has joined the " + currentRoom + " chat room at " + time.Now().Format("15:04:05") + "."}}
	joining <- MainActivity{currentRoom: ch}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if text := input.Text(); len(text) > 0 {
			if text == "/show" {
				io.WriteString(conn, showRooms())
			} else if text == "/disconnect" {
				conn.Close()
			} else {
				messages <- Message{currentRoom: {uuid, name + ": " + text}}
			}
		}
	}

	leaving <- MainActivity{currentRoom: ch}
	messages <- Message{currentRoom: {uuid, "--" + name + " has left the " + currentRoom + " chat room at " + time.Now().Format("15:04:05") + "."}}
}

// write message to clients
func writeMessage(conn net.Conn, ch chan Message, uuid string, currentRoom string) {
	for msg := range ch {
		if strings.Contains(msg[currentRoom].uuid, "Welcome message ") {
			io.WriteString(conn, msg[currentRoom].message+"\nEnter your name: ")
		} else if strings.Contains(msg[currentRoom].uuid, "Initial message ") {
			io.WriteString(conn, msg[currentRoom].message+"\n")
			go metricsInfo(conn, clients, number, currentRoom)
		} else if msg[currentRoom].uuid != uuid {
			io.WriteString(conn, msg[currentRoom].message+"\n")
		}
	}
}

// write metrics information
func metricsInfo(conn net.Conn, clients map[string][]map[Client]bool, number map[string]int, currentRoom string) {
	for {
		io.WriteString(conn, "--Number of rooms: "+strconv.Itoa(len(rooms))+". Number of clients: "+strconv.Itoa(len(clients[currentRoom]))+". Number of messages: "+strconv.Itoa(number[currentRoom])+".\n")
		time.Sleep(30 * time.Second)
	}
}

// remove element from slice
func removeElement(slice []map[Client]bool, index int) []map[Client]bool {
	return append(slice[:index], slice[index+1:]...)
}

// check whether a slice contains a certain value
func isValueInSlice(slice []string, value string) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

// show rooms
func showRooms() string {
	return "There are following chat rooms: " + strings.Join(rooms, ", ") + ".\n"
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
			log.Println(err)
			continue
		}
		go clientConnection(conn)
	}
}
