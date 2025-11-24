package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const network = "tcp"
const defaultPort = ":8080"

func main() {
	fmt.Println("=> Start dv-webserver")
	defer fmt.Println("=> Stopped dv-webserver")

	listener := start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go handleShutdownSignal(stop, listener)

	doListen(listener)
}

func doListen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Printf("Accept error=%s\n", err)
			break
		}
		go processConn(conn)
	}
}

func start() net.Listener {
	address := getAddress()
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("ListenIP error=%s\n", err)
	}
	fmt.Printf("=> Listening on %s\n", address)
	return listener
}

func handleShutdownSignal(stop chan os.Signal, listener net.Listener) {
	<-stop
	fmt.Println("=> Signal received, shutting down...")
	if err := listener.Close(); err != nil {
		fmt.Printf("Listener close error=%s\n", err)
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
