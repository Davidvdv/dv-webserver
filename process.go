package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

func processConn(conn net.Conn) {
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

	content, err := os.ReadFile("public/index.html")
	if err != nil {
		fmt.Printf("ReadFile error=%s\n", err)
		return
	}
	_, _ = io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	_, _ = io.WriteString(conn, "Content-Type: text/html\r\n\r\n")
	_, _ = conn.Write(content)
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
