package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var startTime time.Time

type PingResponse struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Uptime  string    `json:"uptime"`
	AppName string    `json:"app_name"`
	Env     string    `json:"env"`
}

func main() {
	startTime = time.Now()

	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/", handleRoot)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Current Version: 1.0.3")

	fmt.Printf("Server started on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	uptime := time.Since(startTime).String()

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "your_ip"
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "playground"
	}

	response := PingResponse{
		Message: "pong",
		Time:    currentTime,
		Uptime:  uptime,
		AppName: appName,
		Env:     env,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Check if the X-Forwarded-For header is provided (typical with proxies)
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		// Directly fetch the public IP without using r.RemoteAddr as a fallback
		publicIP, err := getPublicIP()
		if err != nil {
			clientIP = "unable to determine your IP"
		} else {
			clientIP = publicIP
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Your IP Address : " + clientIP))
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}
