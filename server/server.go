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
	clients = make(map[net.Conn]string) // Map to store connections and associated node names
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Lock to safely access the clients map
	mu.Lock()
	nodeName, ok := clients[conn]
	mu.Unlock()

	if !ok {
		fmt.Println("Error: Node name not found for connection.")
		return
	}

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
			fmt.Printf("%s - %s disconnected\n", timestamp(), nodeName)
			// Lock to safely remove the disconnected node from the clients map
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}

		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Printf("%s - %s disconnected\n", timestamp(), nodeName)
			// Lock to safely remove the disconnected node from the clients map
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}

		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte(myTime))
	}
}

func timestamp() string {
	return fmt.Sprintf("%.6f", float64(time.Now().UnixNano())/1e9)
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

		// Get node name from the initial connection message
		nodeName, err := bufio.NewReader(conn).ReadString(' ')
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Lock to safely update the clients map
		mu.Lock()
		clients[conn] = strings.TrimSpace(nodeName)
		mu.Unlock()

		// Print the connected message
		fmt.Printf("%s - %s connected\n", timestamp(), clients[conn])

		// Handle each connection concurrently
		go handleConnection(conn)
	}
}
