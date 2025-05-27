# File Transfer Application using Sockets

## 1. Overview

This project is a socket-based client-server application implemented in Go (Golang). The goal is to provide a file transfer system with the following features:

* Directory listing on server
* Single and multiple file downloads
* Chunked file transmission for files >1MB
* Progress display, speed, and transfer duration on client
* File overwrite check and download into dedicated folder

---

## 2. Features Implemented

### Server Side

* Handles multiple clients concurrently using Goroutines.
* Responds to two commands:

  * `LIST`: Sends list of available files in `shared/` directory.
  * `GET <file1> <file2> ...`: Streams files in 1MB chunks.
* Reports:

  * Client connection log
  * Received command
  * File sending log with size and chunk info

### Client Side

* Menu-driven CLI for:

  * Listing server files
  * Downloading single or multiple files
* Download behavior:

  * Prompts if file exists
  * Saves all received files into `received/` folder
  * Displays progress per chunk:

    * e.g., `Downloading test.zip part 3/5 ... 60%`
  * At end, prints download time and speed

---

## 3. Directory Structure

```txt
file-transfer/
├── client/
│   ├── received/        # all downloaded files go here
│   └── client.go
├── server/
│   ├── shared/          # files available for download
│   └── server.go
└── README.md
```

---

## 4. Protocol Design

* **LIST command**: Sent as `LIST`, server replies with newline-separated filenames.
* **GET command**: Sent as `GET file1 file2`, server responds for each file:

  * `FILE <filename> <filesize>\n`
  * Waits for `READY` or `SKIP`
  * Sends content in chunks until done

---

## 5. Example Output

### Client Console

```bash
Choose an option:
1. List files
2. Download file(s)
0. Exit
Enter a number (0-2): 2
Enter filenames (space-separated): 5MB.zip
File 5MB.zip already exists. Overwrite? (y/n): y
Downloading 5MB.zip part 1/5 ..... 20%
Downloading 5MB.zip part 2/5 ..... 40%
Downloading 5MB.zip part 3/5 ..... 60%
Downloading 5MB.zip part 4/5 ..... 80%
Downloading 5MB.zip part 5/5 ..... 100%
Time: 0.01s | Speed: 731.58 MB/s

File saved at: ./received
```

### Server Console

```bash
Server is listening on port 9000
New connection from 127.0.0.1:54321
Received command: GET 5MB.zip
Sending 5MB.zip (5.00 MB)
...
```

---

## 6. Technologies Used

* Language: Golang
* Socket Library: `net`
* OS: Cross-platform (tested on Windows)

---

## 7. Completeness

* [x] Basic Socket Communication
* [x] Directory Listing
* [x] Chunked File Transfer > 1MB
* [x] Multiple File Download
* [x] Progress, Speed, and Time Reporting
* [x] File overwrite prompt
* [x] Organized file saving

Estimated completeness: **100%**

---

## 8. How to Run

### Using Go

```bash
# Terminal 1 - Run Server
cd server
go run server.go

# Terminal 2 - Run Client
cd client
go run client.go
```

### Using Executable Files

```bash
# Terminal 1 - Run Server
cd release
.\server.exe

# Terminal 2 - Run Client
cd release
.\client.exe
```
