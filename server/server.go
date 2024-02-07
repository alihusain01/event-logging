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
	ipNodeMap = make(map[string]string)
	logFile *os.File
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fullErrorMessage := strings.Split(err.Error(), " ")
			ipAddressWithColon := strings.Split(strings.TrimSpace(fullErrorMessage[2]), ">")
			ipAddressTrimmed := strings.Split(strings.TrimSpace(ipAddressWithColon[1]), ":")
			finalIP := ipAddressTrimmed[0]
			fmt.Printf(timestamp() + " - " + ipNodeMap[finalIP] + " disconnected" + "\n")
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
		message := parts[3]
		messageSize := len(netData)

		fmt.Println(eventTimestamp + " " + nodeName + " " + message)

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
		nodeName, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Parse the input into nodename and event timestamp
		connectionMessage := strings.Split(strings.TrimSpace(nodeName), " ")

		eventTimestamp := connectionMessage[0]
		node_name := connectionMessage[2]
		node_ip := connectionMessage[4]

		if _, ok := ipNodeMap[node_ip]; !ok {
			ipNodeMap[node_ip] = node_name
		}

		fmt.Printf(eventTimestamp + " - " + node_name + " connected" + "\n")

		// Lock to safely update the clients map
		mu.Lock()
		clients[conn] = strings.TrimSpace(nodeName)
		mu.Unlock()

		// Handle each connection concurrently
		go handleConnection(conn)
	}
}
