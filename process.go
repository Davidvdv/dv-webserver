package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	bufferSize      = 1024
	defaultFilePath = "public/index.html"
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
	buffer := make([]byte, bufferSize)
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
	fmt.Printf("=> Accepted connection from %s %+v\n", remoteAddress, requestDetails)

	if err := serveFileContent(conn); err != nil {
		fmt.Printf("Serve file content error=%s\n", err)
		return
	}

	fmt.Printf("=> Closed connection from %s\n", remoteAddress)
}

func parseRequest(request string) (*RequestDetails, error) {
	fmt.Printf("=> Incoming raw request\n%+v", request)
	lines := strings.Split(request, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
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

func serveFileContent(conn net.Conn) error {
	content, err := os.ReadFile(defaultFilePath)
	if err != nil {
		return fmt.Errorf("ReadFile error=%s\n", err)
	}
	if _, err := io.WriteString(conn, "HTTP/1.1 200 OK\r\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(conn, "Content-Type: text/html\r\n\r\n"); err != nil {
		return err
	}
	_, err = conn.Write(content)
	return err
}
