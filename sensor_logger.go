// This code listen port 1703 and get data from ESP-32 in "#RTNode1=1111" format

package main

import (
    "fmt"
    "net"
    "os"
    "strings"
    "time"
)

const (
    UDP_IP      = "0.0.0.0"
    UDP_PORT    = 1703
    BUFFER_SIZE = 32
    LOG_FILE    = "sensor_data.log"
)

func logData(data string) {
    file, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening log file:", err)
        return
    }
    defer file.Close()

    timestamp := time.Now().Format("2006-01-02 15:04:05")
    _, err = file.WriteString(fmt.Sprintf("%s - %s\n", timestamp, data))
    if err != nil {
        fmt.Println("Error writing to log file:", err)
    }
}

func parseString(inputStr string) []string {
    inputStr = strings.TrimPrefix(inputStr, "#")
    parts := strings.Split(inputStr, "=")
    if len(parts) != 2 {
        return nil
    }

    nodePart := parts[0]
    valuesPart := strings.Split(parts[1], "")
    result := append([]string{nodePart}, valuesPart...)

    return result
}

func main() {
    addr := net.UDPAddr{
        Port: UDP_PORT,
        IP:   net.ParseIP(UDP_IP),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Println("Error creating UDP socket:", err)
        return
    }
    defer conn.Close()

    fmt.Printf("Listening on %s:%d...\n", UDP_IP, UDP_PORT)

    buffer := make([]byte, BUFFER_SIZE)

    for {
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("Error reading from UDP:", err)
            continue
        }

        dataStr := strings.TrimSpace(string(buffer[:n]))

        logData(dataStr)

        parsedList := parseString(dataStr)

        fmt.Printf("Received message from %s: %s\n", clientAddr, dataStr)
        fmt.Printf("Parsed message: %v\n", parsedList)

		for i := 1; i < len(parsedList); i++ {
			fmt.Printf("Received message, item %d: %s\n", i, parsedList[i])
		}
    }
}
