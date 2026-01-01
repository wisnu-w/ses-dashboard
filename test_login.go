package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Test login request
	loginData := map[string]string{
		"username": "admin",
		"password": "password",
	}

	jsonData, _ := json.Marshal(loginData)
	
	resp, err := http.Post("http://localhost:8080/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
	
	if resp.StatusCode == 200 {
		fmt.Println("✅ Login successful!")
	} else {
		fmt.Println("❌ Login failed")
	}
}