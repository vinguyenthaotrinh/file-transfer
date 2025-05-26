package main

import (
	"bufio"
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
		fmt.Println("2. Download file(s)")
		fmt.Println("0. Exit")
		fmt.Print("Choose: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			listFiles()
		case 2:
			fmt.Print("Enter filenames (space-separated): ")
			var input string
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input = scanner.Text()
			}
			files := strings.Fields(input)
			downloadFiles(files)
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

func downloadFiles(filenames []string) {
	conn, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("GET " + strings.Join(filenames, " ") + "\n"))

	for i, name := range filenames {
		// Nhận dòng FILE <name> <size>
		headerBuf := make([]byte, 128)
		n, err := conn.Read(headerBuf)
		if err != nil {
			fmt.Println("\nError receiving header:", err)
			return
		}

		header := string(headerBuf[:n])
		if strings.HasPrefix(header, "ERR") {
			fmt.Printf("Server error: cannot send %s\n", name)
			continue
		}

		var filename string
		var filesize int64
		fmt.Sscanf(header, "FILE %s %d", &filename, &filesize)

		conn.Write([]byte("READY"))

		out, err := os.Create(filename)
		if err != nil {
			fmt.Println("Failed to create file:", filename)
			continue
		}
		defer out.Close()

		received := int64(0)
		buffer := make([]byte, CHUNK_SIZE)

		part := i + 1
		for received < filesize {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("\nError downloading %s\n", filename)
				break
			}
			out.Write(buffer[:n])
			received += int64(n)

			percent := float64(received) / float64(filesize) * 100
			fmt.Printf("\rDownloading %s part %d ... %.0f%%", filename, part, percent)
		}
		fmt.Printf("\nDownload of %s complete. Saved at: %s\n", filename, filename)
	}
}
