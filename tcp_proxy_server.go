package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		makeRequest(conn)
	}
}

func makeRequest(c net.Conn) {
	defer c.Close()
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	input := bufio.NewScanner(conn)
	if input.Scan() {
		io.WriteString(c, input.Text())
	}
}
