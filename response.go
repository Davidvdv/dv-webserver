package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	HttpResponseProtocol = "HTTP/1.1 "
	PublicDir            = "public"
	DefaultFilePath      = "/index.html"
)

func responseHandler(conn net.Conn, requestDetails *RequestDetails) {
	contents, notFoundErr := readFile(requestDetails.Path)
	// TODO: explore http.responseWriter
	sb := strings.Builder{}
	sb.WriteString(HttpResponseProtocol)
	if notFoundErr != nil {
		fmt.Println(notFoundErr)
		sb.WriteString(strconv.Itoa(http.StatusNotFound) + " ")
		sb.WriteString(http.StatusText(http.StatusNotFound) + "\r\n")
	} else {
		sb.WriteString(strconv.Itoa(http.StatusOK) + " ")
		sb.WriteString(http.StatusText(http.StatusOK) + "\r\n")
		sb.WriteString("Content-Type: ")
		sb.WriteString(http.DetectContentType(contents) + "\r\n\r\n")
		sb.Write(contents)
	}
	httpResponseMessage := sb.String()
	fmt.Println(httpResponseMessage)
	_, _ = io.WriteString(conn, httpResponseMessage)
	// TODO: handle errors
}

func readFile(requestPath string) ([]byte, error) {
	if requestPath == "/" || requestPath == "" {
		requestPath = DefaultFilePath
	}
	if contents, err := os.ReadFile(path.Join(PublicDir + requestPath)); err != nil {
		return nil, fmt.Errorf("File not found: %s\n", requestPath)
	} else {
		return contents, nil
	}
}
