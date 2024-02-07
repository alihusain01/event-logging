package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Unknown"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "Unknown"
}

func send_message(conn net.Conn, node_name, ip_address string) {
	for {
		// Receives event from generator.py
		reader := bufio.NewReader(os.Stdin)
		generator_event, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", err)
			break
		}

		if strings.TrimSpace(generator_event) == "STOP" {
			fmt.Println("Exiting...")
			break
		} else {
			parts := strings.Split(generator_event, " ")
			message := parts[0] + " " + node_name + " " + ip_address + " " + parts[1]
			fmt.Fprintf(conn, message+"\n")
		}
	}
}

func main() {
	arguments := os.Args
	if len(arguments) != 4 {
		fmt.Println("Please provide node name, host, and port.")
		return
	}

	node_name := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]

	ip_address := getLocalIP()

	CONNECT := host + ":" + port
	c, err := net.Dial("tcp", CONNECT)

	if err != nil {
		fmt.Println(err)
		return
	}

	// First message to be sent is a timestamp with the node name and IP address
	currentTime := time.Now().UnixNano()
	connectionMessage := fmt.Sprintf("%.6f - %s - %s Connected\n", float64(currentTime)/1e9, node_name, ip_address)

	// Print the connection message locally and to server
	fmt.Println(connectionMessage)
	fmt.Fprintf(c, connectionMessage)

	// Continue sending any received messages
	send_message(c, node_name, ip_address)

	return
}
