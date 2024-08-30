package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	UDP_IP      = "0.0.0.0"
	UDP_PORT    = 1703
	BUFFER_SIZE = 32
	LOG_FILE    = "sensor_data.log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket соединение
var wsConn *websocket.Conn

// Функция для записи данных в файл
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

// Функция для парсинга строки
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

func handleConnection(w http.ResponseWriter, r *http.Request) {
	var err error
	wsConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}
	defer wsConn.Close()

	log.Println("Client connected")

	// Бесконечный цикл для удержания соединения открытым
	for {
		if _, _, err := wsConn.ReadMessage(); err != nil {
			log.Println("Error while reading message:", err)
			return
		}
	}
}

func main() {
	// Запуск WebSocket-сервера
	http.HandleFunc("/ws", handleConnection)

	serverAddr := "0.0.0.0:8080"
	go func() {
		log.Printf("WebSocket server started at ws://%s\n", serverAddr)
		err := http.ListenAndServe(serverAddr, nil)
		if err != nil {
			log.Fatal("Error while starting server:", err)
		}
	}()

	// Запуск UDP-сервера
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

		// Логирование данных
		logData(dataStr)

		// Парсинг строки
		parsedList := parseString(dataStr)

		fmt.Printf("Received message from %s: %s\n", clientAddr, dataStr)
		fmt.Printf("Parsed message: %v\n", parsedList)

		for i := 1; i < len(parsedList); i++ {
			fmt.Printf("Received message, item %d: %s\n", i, parsedList[i])
		}

		// Отправка данных через WebSocket в формате JSON
		if wsConn != nil {
			/*data1 := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"data":      dataStr,
				"parsed":    parsedList,
			}*/

			data := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"data":      dataStr,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Println("Error while marshalling JSON:", err)
				continue
			}

			err = wsConn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Error while writing message:", err)
				wsConn = nil
			}
		}
	}
}
