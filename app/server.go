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

	connection, err := listener.Accept()
		if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	request, err := http.ReadRequest(bufio.NewReader(connection))

	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	if request.URL.Path == "/" {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
