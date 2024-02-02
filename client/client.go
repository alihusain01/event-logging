package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Fprintf(os.Stderr, "Usage: %s host port\n", os.Args[0])
        os.Exit(1)
    }
    host := os.Args[1]
    port := os.Args[2]
    conn, err := net.Dial("tcp", host+":"+port)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
        os.Exit(1)
    }
    defer conn.Close()
    fmt.Println("Connected to server")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        if !scanner.Scan() {
            break
        }
        message := scanner.Text()
        fmt.Fprintf(conn, message+"\n")
    }
}