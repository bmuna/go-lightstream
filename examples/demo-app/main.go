package main

import (
	"github/com/bmuna/go-lightstream/lightstream"
	"log"
	"net/http"
)

func main() {
	server := lightstream.NewServer()
	http.HandleFunc("/ws", server.HandleWS)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
