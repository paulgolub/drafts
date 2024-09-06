package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

var mySigningKey = []byte("your_secret_key")

const (
	UDP_IP      = "0.0.0.0"
	UDP_PORT    = 8085
	BUFFER_SIZE = 32
	LOG_FILE    = "sensor_data.log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket connection
var wsConn *websocket.Conn

// Database for text
var fake_db = []string{
	"Value 1",
	"Value 2",
  "Value 3",
  "Value 4",
}

func getTokenFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	params := parsedURL.Query()
	token := params.Get("token")
	return token, nil
}

func handleWebSocket(conn *websocket.Conn, rawURL string) {
	token, err := getTokenFromURL(rawURL)
	if err != nil || token == "" {
		conn.Close()
		return
	}

	// Decoding the JWT token
	decodedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		conn.Close()
		return
	}

	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok && decodedToken.Valid {
		expiration := int64(claims["exp"].(float64))

		if expiration < time.Now().Unix() {
			conn.Close()
			return
		}

		fmt.Println("Token valid, user:", claims["username"])
	} else {
		conn.Close()
	}
}

// Key JWT Parcer
func ParseJWT(tokenString string) (*jwt.Token, error) {
    // Разбор токена и проверка подписи
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Проверка метода подписания
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return mySigningKey, nil
    })
    if err != nil {
        return nil, err
    }
    return token, nil
}

// Log data writting
func logData(data string) {
	file, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("%s - %s\n", timestamp, data)
	_, err = file.WriteString(fmt.Sprintf("%s - %s\n", timestamp, data))
	if err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

// UDP data string parcer
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
	// Извлекаем параметры строки запроса
	params := r.URL.Query()
	rawURL := r.URL.Query()

	tokenString := params.Get("token")
	if tokenString == "" {
		http.Error(w, "Authorization parameter is missing", http.StatusBadRequest)
		return
	}

	// if claims, ok := tokenString.Claims.(jwt.MapClaims); ok && tokenString.Valid {
	// 	// Проверка на истечение срока действия токена
	// 	if exp, ok := claims["exp"].(float64); ok {
	// 		if time.Unix(int64(exp), 0).Before(time.Now()) {
	// 			http.Error(w, "Token has expired", http.StatusUnauthorized)
	// 			return
	// 		}
	// 	}
	// 	// Valid token
	// 	// next(w, r)
	// 	fmt.Println("Valid token")
	// } else {
	// 	http.Error(w, "Invalid token", http.StatusUnauthorized)
	// 	fmt.Println("Протухший токен")
	// }
	
	fmt.Println("rawURL:", rawURL)

	var err error
	wsConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}
	defer wsConn.Close()

  //tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU1NTA5MTQsInVzZXJuYW1lIjoiZXhhbXBsZVVzZXIifQ.Ob-3GM5FjYxLYG1qv0_Lxnw15GCGnCKSDi7bAyrVYD0"
	
	//rawURL := "wss://example.com?token=eRJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU1NTA5MTQsInVzZXJuYW1lIjoiZXhhbXBsZVVzZXIifQ.Ob-3GM5FjYxLYG1qv0_Lxnw15GCGnCKSDi7bAyrVYD0"

	//tokenString, err := getTokenFromURL(rawURL)

	token, err := ParseJWT(tokenString)
    if err != nil {
        fmt.Println("Error parsing token:", err)
        return
    }

    // Token validation
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("Claims: %+v\n", claims)
		if exp, ok := claims["exp"].(float64); ok {
			fmt.Printf("Exp time: %+v\n", exp)
			fmt.Printf("Time now: %+v\n", time.Now())
			fmt.Printf("Un time: %+v\n", time.Unix(int64(exp), 0))
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				// http.Error(w, "Token has expired", http.StatusUnauthorized)
				fmt.Println("Token has expired")
				return
			}
		}
        fmt.Println("Token is valid")

    } else {
        fmt.Println("Invalid token")
		defer wsConn.Close()
    }

	log.Println("Client connected")

	// JSON was recived flag
	jsonReceived := false

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Client disconnected")
				return
			}
			log.Println("Error while reading message:", err)
			return
		}

		// If message is JSON
		var jsonMessage map[string]interface{}
		if err := json.Unmarshal(message, &jsonMessage); err == nil {
			// Сообщение является JSON, начинаем отправку данных
			jsonReceived = true
			log.Println("Received JSON from client:", jsonMessage)
		} else {
			log.Println("Received non-JSON message from client, ignoring...")
		}

		// If JSON recived, start sending
		if jsonReceived {
			break
		}
	}

	// sending 
	for {
		// if JSON was recived, send data from UDP
		buffer := make([]byte, BUFFER_SIZE)
		addr := net.UDPAddr{
			Port: UDP_PORT,
			IP:   net.ParseIP(UDP_IP),
		}

		conn, err := net.ListenUDP("udp", &addr)
		if err != nil {
			log.Println("Error creating UDP socket:", err)
			return
		}
		defer conn.Close()

		for {
			n, clientAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Println("Error reading from UDP:", err)
				continue
			}

			dataStr := strings.TrimSpace(string(buffer[:n]))

			logData(dataStr)

			parsedList := parseString(dataStr)

			fmt.Printf("Received message from %s: %s\n", clientAddr, dataStr)
			fmt.Printf("Parsed message: %v\n", parsedList)

			var poiIndex int

			if parsedList[1] == "0" && parsedList[2] == "1" && parsedList[3] == "1" {
				poiIndex = 1
			}
			if parsedList[2] == "0" && parsedList[1] == "1" && parsedList[3] == "1" {
				poiIndex = 2
			}
			if parsedList[3] == "0" && parsedList[2] == "1" && parsedList[1] == "1" {
				poiIndex = 3
			}
			
			if parsedList[3] == "0" && parsedList[2] == "0" && parsedList[1] == "0" {
				poiIndex = 0
			}

			fmt.Println(poiIndex)
		
			// Create JSON struct for sending
			// data := map[string]interface{}{
			// 	"timestamp": time.Now().Format(time.RFC3339),
			// 	"data":      dataStr,
			// 	"parsed":    parsedList,
			// }

			data := map[string]interface{}{
				"POI_index": fmt.Sprintf("%d", poiIndex),
				"media_link": fmt.Sprintf("/stream/%d", poiIndex),
				"audio_text": fake_db[poiIndex],
			}

			// If POI_index = 0. No detected
			if poiIndex == 0 {
				data["POI_index"] = "None"
				data["media_link"] = "Nothing to play"
				data["audio_text"] = fake_db[poiIndex]
			}


			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Println("Error while marshalling JSON:", err)
				continue
			}

			// sending via WebSocket
			err = wsConn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Client disconnected, stopping data transmission")
					return
				}
				log.Println("Error while writing message:", err)
				wsConn = nil
				return
			}
		}
	}
}

func main() {
	// starting WebSocket server
	http.HandleFunc("/ws", handleConnection)

	serverAddr := "0.0.0.0:8765"
	go func() {
		log.Printf("WebSocket server started at ws://%s\n", serverAddr)
		err := http.ListenAndServe(serverAddr, nil)
		if err != nil {
			log.Fatal("Error while starting server:", err)
		}
	}()

	// Main program is waiting websocket connection
	select {}
}
