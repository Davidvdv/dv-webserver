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
	publicDir       = "public"
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

	if err := serveFileContent(conn, requestDetails.Path); err != nil {
		fmt.Printf("Serve file content error=%s\n", err)
		return
	}

	fmt.Printf("=> Closed connection from %s\n", remoteAddress)
}

func parseRequest(request string) (*RequestDetails, error) {
	requestLines := strings.Split(request, "\n")
	for i := range requestLines {
		requestLines[i] = strings.TrimRight(requestLines[i], "\r")
	}
	if len(requestLines) >= 4 {
		firstLineParts := strings.Split(requestLines[0], " ")
		if len(firstLineParts) >= 3 {
			return &RequestDetails{
				HttpVersion: firstLineParts[2],
				Method:      firstLineParts[0],
				Path:        firstLineParts[1],
				Host:        requestLines[1],
				UserAgent:   requestLines[2],
				Accept:      requestLines[3],
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid request %s", request[0:100]+" ...")
}

func serveFileContent(conn net.Conn, path string) error {
	content, err := os.ReadFile(publicDir + path)
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
