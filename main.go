package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

const network = "tcp"
const defaultPort = ":8080"

func main() {
	fmt.Println("Start dv-webserver")
	defer fmt.Println("Stopped dv-webserver")

	address := getAddress()
	listener, err := net.Listen(network, address)
	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("Listener close error=%s\n", err)
		}
	}()
	if err != nil {
		log.Fatalf("ListenIP error=%s\n", err)
	}

	fmt.Printf("=> Listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error=%s\n", err)
			continue
		}
		go processConn(conn)
	}
}

func getAddress() string {
	if len(os.Args) < 2 {
		return defaultPort
	}
	port := os.Args[1]

	if portNum, err := strconv.Atoi(port); err != nil || portNum < 1 || portNum > 65535 {
		return defaultPort
	}
	return ":" + port
}
