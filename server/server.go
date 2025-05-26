package main

import (
	"fmt"
	"net"
)

const (
	PORT       = ":9000"
	CHUNK_SIZE = 1024 * 1024
	SHARE_DIR  = "./shared"
)

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		fmt.Printf("New connection from %s\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
}
