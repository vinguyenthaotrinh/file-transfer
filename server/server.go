package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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

	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	command := strings.TrimSpace(string(buffer[:n]))

	fmt.Printf("Received command: %s\n", command)

	if command == "LIST" {
		handleList(conn)
	} else if strings.HasPrefix(command, "GET ") {
		filename := strings.TrimSpace(strings.TrimPrefix(command, "GET "))
		handleGet(conn, filename)
	} else {
		conn.Write([]byte("ERR: Unknown command\n"))
	}
}

func handleList(conn net.Conn) {
	fmt.Println("Sending file list to client...")

	files, err := os.ReadDir(SHARE_DIR)
	if err != nil {
		conn.Write([]byte("ERR: Cannot read directory\n"))
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	if len(fileList) == 0 {
		conn.Write([]byte("No files available\n"))
	} else {
		conn.Write([]byte(strings.Join(fileList, "\n")))
	}
}

func handleGet(conn net.Conn, filename string) {
	fullPath := SHARE_DIR + "/" + filename
	file, err := os.Open(fullPath)
	if err != nil {
		conn.Write([]byte("ERR: Cannot open file\n"))
		return
	}
	defer file.Close()

	info, _ := file.Stat()
	filesize := info.Size()

	fmt.Printf("Sending file '%s' (%.2f MB)\n", filename, float64(filesize)/(1024*1024))

	conn.Write([]byte(fmt.Sprintf("SIZE %d\n", filesize)))

	ack := make([]byte, 16)
	conn.Read(ack)

	buffer := make([]byte, CHUNK_SIZE)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			conn.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}
}
