package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func send_message(conn net.Conn, node_name string) {
	for {
		// Receives event from generator.py
		reader := bufio.NewReader(os.Stdin)
		generator_event, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", err)
			break
		}

		// Stop the loop if the user types "STOP"
		if strings.TrimSpace(generator_event) == "STOP" {
			fmt.Println("Exiting...")
			break
		} else {
			// Splicing generator message to add node between timestamp and event
			parts := strings.Split(generator_event, " ")
			message := parts[0] + " " + node_name + " " + parts[1]

			// Send message to server
			fmt.Fprintf(conn, message+"\n")
		}
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	node_name := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]

	// Establish connection with server
	CONNECT := host + ":" + port
	c, err := net.Dial("tcp", CONNECT)

	// Check if there was an error creating the connection
	if err != nil {
		fmt.Println(err)
		return
	}

	// First message to be sent is a timestamp with the node name
	currentTime := time.Now().UnixNano()
	connectionMessage := fmt.Sprintf("%.6f - %s Connected\n", float64(currentTime)/1e9, node_name)

	// Print the connection message locally in VM2
	fmt.Println(connectionMessage)

	// Send the connection message to the server
	fmt.Fprintf(c, connectionMessage)

	// Continue sending any received messages
	send_message(c, node_name)

	return
}
