package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	bufferSize = 1024
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

	responseHandler(conn, requestDetails)

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
