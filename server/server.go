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
	fileMu  sync.Mutex
	clients = make(map[net.Conn]string) // Map to store connections and associated node names
	logFile *os.File
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Printf("Error reading from connection: %s\n", err)
			// Lock to safely remove the disconnected node from the clients map
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}

		// Parse the input into nodename and event timestamp
		parts := strings.Split(strings.TrimSpace(netData), " ")

		// Ensure that there are at least two parts before accessing them
		if len(parts) < 2 {
			fmt.Printf("Received invalid data from connection: %s\n", netData)
			continue
		}

		eventTimestamp := parts[0]
		nodeName := parts[1]
		messageSize := len(netData)

		// Handle STOP command
		if nodeName == "STOP" {
			fmt.Printf("%s - %s disconnected\n", timestamp(), nodeName)
			// Lock to safely remove the disconnected node from the clients map
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}

		fmt.Printf("-> %s \n", netData)
		logEntry := fmt.Sprintf("%s %s %d \n", timestamp(), eventTimestamp, messageSize)

		fileMu.Lock()
		_, err = logFile.WriteString(logEntry)
		if err != nil {
			fmt.Println("Error writing to log file:", err)
		}
		fileMu.Unlock()

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

	// Create log file once
	var logErr error
	logFile, logErr = os.Create("logFile.txt") // This will create a new file or truncate an existing one
	if logErr != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer logFile.Close()

	// Writing a file title
	_, logTitleErr := logFile.WriteString("Program start: " + timestamp() + "\n")
	if logTitleErr != nil {
		fmt.Println("Error writing log file title:", err)
		return
	}

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
		fmt.Printf("%s - %s connected\n", timestamp(), nodeName)

		// Handle each connection concurrently
		go handleConnection(conn)
	}
}
