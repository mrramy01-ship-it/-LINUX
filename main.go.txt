package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    fmt.Println("AnonGo v4.0 Public - Linux source")
    file, err := os.Open("sources.list")
    if err != nil {
        fmt.Println("sources.list not found")
        return
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    fmt.Println("Fetching from Linux sources:")
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" || line[0] == '#' {
            continue
        }
        fmt.Println(" -", line)
    }
}
