package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s port\n", os.Args[0])
        os.Exit(1)
    }
    port := os.Args[1]
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
        os.Exit(1)
    }
    defer listener.Close()
    fmt.Printf("Listening on port %s...\n", port)
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
            continue
        }
        fmt.Printf("New node connected: %s\n", conn.RemoteAddr().String())
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Printf("%s sending the message\n", conn.RemoteAddr().String())
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Printf("Node disconnected: %s\n", conn.RemoteAddr().String())
            return
        }
        fmt.Printf("%s: %s\n", conn.RemoteAddr().String(), string(buf[:n]))
    }
}