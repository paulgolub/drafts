// Websocket server

// Should be done:
// go mod init websocket
// go get github.com/gorilla/websocket

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading message:", err)
			return
		}

		if messageType == websocket.TextMessage {
			log.Println("Received message:", string(msg))

			// Пример данных, которые будем отправлять
			data := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"value":     "some_data",
			}

			for {
				// Преобразуем данные в формат JSON
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Println("Error while marshalling JSON:", err)
					return
				}

				// Отправляем данные по WebSocket
				err = conn.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					log.Println("Error while writing message:", err)
					return
				}

				// Пауза между отправками
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)

	serverAddr := "localhost:8080"
	log.Printf("WebSocket server started at ws://%s\n", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Error while starting server:", err)
	}
}
