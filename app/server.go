package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	request, err := http.ReadRequest(bufio.NewReader(connection))
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	path := request.URL.Path
	if path == "/" {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	} else if path[0:6] == "/echo/" {
		connection.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(path[6:]), path[6:])))
		return
	} else if path[0:7] == "/files/" {
		file, err := os.ReadFile(os.Args[2] + path[7:])
		if err != nil {
			fmt.Println("Error opening file: ", err.Error())
			connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			connection.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(file), file)))
		}
			return
	}

	agent := request.UserAgent()
	if len(agent) > 0 {
		connection.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(agent), agent)))
		return
	}

	connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
