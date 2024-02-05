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

		// Parse the input into nodename and event timestamp
		parts := strings.Split(netData, " ")
		event_timestamp := parts[0]
		nodeName := parts[1]
		messageSize := len(netData)

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

		fmt.Printf("-> %s \n", string(netData))
		logEntry := fmt.Sprintf("%s %s %d \n", timestamp(), event_timestamp, messageSize)

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

		//Create logging file
		var log_err error
		logFile, log_err = os.Create("logFile.txt") // This will create a new file or truncate an existing one
		if log_err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer logFile.Close()

		// Writing a file title
		_, log_title_err := logFile.WriteString("Logging at time " + timestamp() + "\n")
		if log_title_err != nil {
			fmt.Println("Error writing log file tilte:", err)
			return
		}

		// Handle each connection concurrently
		go handleConnection(conn)
	}
}
