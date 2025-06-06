package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	HOST       = "127.0.0.1"
	PORT       = "9000"
	CHUNK_SIZE = 1024 * 1024
	SHARED_DIR = "./shared"
)

func main() {
	listener, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port", PORT)

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

	buffer := make([]byte, 2048)
	n, _ := conn.Read(buffer)
	command := strings.TrimSpace(string(buffer[:n]))

	fmt.Printf("Received command: %s\n", command)

	if command == "LIST" {
		handleList(conn)
	} else if strings.HasPrefix(command, "GET ") {
		files := strings.Fields(strings.TrimPrefix(command, "GET "))
		handleMultiGet(conn, files)
	} else {
		conn.Write([]byte("ERR: Unknown command\n"))
	}
}

func handleList(conn net.Conn) {
	fmt.Println("Sending file list to client...")

	files, err := os.ReadDir(SHARED_DIR)
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

func handleMultiGet(conn net.Conn, filenames []string) {
	for _, name := range filenames {
		fullPath := SHARED_DIR + "/" + name
		file, err := os.Open(fullPath)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("ERR %s\n", name)))
			continue
		}

		info, _ := file.Stat()
		filesize := info.Size()

		fmt.Printf("Sending file %s (%.2f MB)\n", name, float64(filesize)/(1024*1024))
		conn.Write([]byte(fmt.Sprintf("FILE %s %d\n", name, filesize)))

		ack := make([]byte, 16)
		conn.Read(ack)
		if strings.TrimSpace(string(ack)) == "SKIP" {
			fmt.Printf("Client skipped %s\n", name)
			file.Close()
			continue
		}

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
		file.Close()
	}
}
