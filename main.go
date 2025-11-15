package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type RequestDetails struct {
	HttpVersion string
	Method      string
	Path        string
	Host        string
	UserAgent   string
	Accept      string
}

func main() {
	fmt.Println("# dv-webserver started")
	defer fmt.Println("# dv-webserver stopped")

	listener, err := net.Listen("tcp", getAddress())
	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("Listener close error=%s\n", err)
		}
	}()

	if err != nil {
		log.Fatalf("ListenIP error=%s\n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error=%s\n", err)
		}
		go acceptConn(conn)
	}
}

func acceptConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Connection close error=%s\n", err)
		}
	}(conn)

	remoteAddress := conn.RemoteAddr().String()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Read error=%s\n", err)
		return
	}

	requestDetails, err := parseRequest(string(buffer[:n]))
	if err != nil {
		fmt.Printf("Parse error=%s\n", err)
		return
	}
	fmt.Printf("=> Accepted connection from %s %s\n", remoteAddress, requestDetails)

	_, err = io.WriteString(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\nRequested path: %s\r\nMethod: %s\r\n", requestDetails.Path, requestDetails.Method))
	if err != nil {
		fmt.Printf("WriteString error=%s\n", err)
		return
	}

	fmt.Printf("=> Closed connection from %s\n", remoteAddress)
}

func parseRequest(request string) (*RequestDetails, error) {
	fmt.Printf("=> Incoming raw request\n%s", request)
	lines := strings.Split(request, "\n")
	if len(lines) >= 4 {
		parts := strings.Split(lines[0], " ")
		if len(parts) >= 3 {
			return &RequestDetails{
				HttpVersion: parts[2],
				Method:      parts[0],
				Path:        parts[1],
				Host:        lines[1],
				UserAgent:   lines[2],
				Accept:      lines[3],
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid request %s", request[0:100]+" ...")
}

func getAddress() string {
	if len(os.Args) < 2 {
		return ":8080"
	}
	port := os.Args[1]

	if portNum, err := strconv.Atoi(port); err != nil || portNum < 1 || portNum > 65535 {
		return ":8080"
	}
	return ":" + port
}
