package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RequestData is the structure for the request body
type RequestData struct {
	Username string `json:"username"`
}

// ResponseData is the structure for the response body
type ResponseData struct {
	Token string `json:"token"`
}

func main() {
	// Prepare the request data
	requestData := RequestData{
		Username: "testuser",
	}

	// Convert the request data to JSON
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling request data:", err)
		return
	}

	// Send the POST request
	resp, err := http.Post("http://172.16.0.195/token", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var responseData ResponseData
		err := json.Unmarshal(responseBody, &responseData)
		if err != nil {
			fmt.Println("Error unmarshaling response data:", err)
			return
		}
		fmt.Println("Received token:", responseData.Token)
	} else {
		fmt.Println("Error response:", string(responseBody))
	}
}
