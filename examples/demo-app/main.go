package main

import (
	"github/com/bmuna/go-lightstream/lightstream"
	"log"
	"net/http"
)

func main() {
	server := lightstream.NewServer()
	http.HandleFunc("/ws", server.HandleWS)        // WebSocket endpoint
	http.HandleFunc("/health", server.HealthCheck) // Health check

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
