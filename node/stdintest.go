package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
)

func main() {
    // Create a command to execute
    cmd := exec.Command("python3", "-u", "generator.py", "0.1", "5")

    // Create buffers to capture the output and errors
    var out bytes.Buffer
    var stderr bytes.Buffer // Buffer for capturing standard error

    // Assign the buffers to the command's output and error output
    cmd.Stdout = &out
    cmd.Stderr = &stderr // Capture standard error

    // Run the command
    err := cmd.Run()
    if err != nil {
        // Print the error along with the stderr output, if any
        fmt.Printf("Error executing command: %s\n", err)
        if stderr.Len() > 0 {
            // Print the contents of stderr if it's not empty
            fmt.Printf("Command stderr: %s\n", stderr.String())
        }
        return
    } else {
        fmt.Printf("Ran successfully\n")
    }

    // Capture the output
    output := out.String()

    // Open a file for writing
    file, err := os.Create("output.txt")
    if err != nil {
        fmt.Printf("Error creating file: %s\n", err)
        return
    }
    defer file.Close()

    // Write the captured output to the file
    _, err = file.WriteString(output)
    if err != nil {
        fmt.Printf("Error writing to file: %s\n", err)
        return
    }

    fmt.Println("Command output written to output.txt")
}
