package main

import (
	"fmt"
	"github/com/bmuna/go-lightstream/lightstream"
	"github/com/bmuna/go-lightstream/lightstream/abr"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("bandwidth_data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write header if file is empty
	info, _ := file.Stat()
	if info.Size() == 0 {
		file.WriteString("timestamp,time_kbps\n")
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		kbps := abr.MeasureBandwidth()
		fmt.Printf("📶 Bandwidth: %.2f kbps\n", kbps)
		abr.LogToFile(kbps, file)
	}

	server := lightstream.NewServer()
	http.HandleFunc("/ws", server.HandleWS)        // WebSocket endpoint
	http.HandleFunc("/health", server.HealthCheck) // Health check

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
