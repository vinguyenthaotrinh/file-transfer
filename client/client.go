package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	SERVER_ADDR  = "localhost:9000"
	CHUNK_SIZE   = 1024 * 1024
	RECEIVED_DIR = "./received"
)

func main() {
	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. List files")
		fmt.Println("2. Download file(s)")
		fmt.Println("0. Exit")
		fmt.Print("Enter a number (0-2): ")

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
			fmt.Println("\nInvalid choice.")
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

	for _, name := range filenames {
		// Nhận header: FILE <filename> <size>
		headerBuf := make([]byte, 128)
		n, err := conn.Read(headerBuf)
		if err != nil {
			fmt.Println("Error receiving header:", err)
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

		// Nếu file đã tồn tại, hỏi người dùng
		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("File %s already exists. Overwrite? (y/n): ", filename)
			var ans string
			fmt.Scanln(&ans)
			if strings.ToLower(ans) != "y" {
				fmt.Printf("Skipped %s\n", filename)
				conn.Write([]byte("SKIP"))
				continue
			}
		}

		conn.Write([]byte("READY"))

		os.MkdirAll(RECEIVED_DIR, os.ModePerm)
		filepath := RECEIVED_DIR + "/" + filename
		out, err := os.Create(filepath)

		if err != nil {
			fmt.Println("Failed to create file:", filename)
			continue
		}
		defer out.Close()

		start := time.Now()
		received := int64(0)
		buffer := make([]byte, CHUNK_SIZE)

		totalParts := (filesize + CHUNK_SIZE - 1) / CHUNK_SIZE
		partNum := 1

		for received < filesize {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("\nError downloading %s\n", filename)
				break
			}
			out.Write(buffer[:n])
			received += int64(n)

			percent := float64(received) / float64(filesize) * 100
			fmt.Printf("\rDownloading %s part %d/%d .... %.0f%%\n", filename, partNum, totalParts, percent)

			partNum = partNum + n/CHUNK_SIZE
		}
		duration := time.Since(start).Seconds()
		speed := float64(received) / (1024 * 1024) / duration

		fmt.Printf("Time: %.2fs | Speed: %.2f MB/s\n\n", duration, speed)
	}
	fmt.Printf("File saved at: %s\n", RECEIVED_DIR)
}
