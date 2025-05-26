package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	SERVER_ADDR = "localhost:9000"
	CHUNK_SIZE  = 1024 * 1024
)

func main() {
	for {
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. List files")
		fmt.Println("2. Download file")
		fmt.Println("0. Exit")
		fmt.Print("Choose: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			listFiles()
		case 2:
			fmt.Print("Enter filename to download: ")
			var filename string
			fmt.Scanln(&filename)
			downloadFile(filename)
		case 0:
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}

func listFiles() {
	conn, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("LIST\n"))

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	fmt.Println("\nFiles on server:")
	fmt.Println(string(buffer[:n]))
}

func downloadFile(filename string) {
	conn, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("GET " + filename + "\n"))

	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading SIZE:", err)
		return
	}

	header := string(buf[:n])
	if strings.HasPrefix(header, "ERR") {
		fmt.Println(header)
		return
	}

	var filesize int64
	fmt.Sscanf(header, "SIZE %d", &filesize)

	conn.Write([]byte("READY"))

	out, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer out.Close()

	received := int64(0)
	buffer := make([]byte, CHUNK_SIZE)

	for received < filesize {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("\nError during download:", err)
			return
		}
		out.Write(buffer[:n])
		received += int64(n)

		percent := float64(received) / float64(filesize) * 100
		fmt.Printf("\rDownloading %s ... %.0f%%", filename, percent)
	}
	fmt.Println("\nDownload complete.")
	fmt.Printf("File saved at: %s\n", filename)
}
