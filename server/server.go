package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	clients = make(map[net.Conn]bool)
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Add the new connection to the clients map
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	defer func() {
		// Remove the connection from the clients map when the function exits
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
	}()

	// Read and print the initial connection message
	netData, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("-> ", string(netData))

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Connection closed:", conn.RemoteAddr())
			return
		}

		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte(myTime))
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fmt.Println("Server listening on", PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Handle each connection concurrently
		go handleConnection(conn)
	}
}
