package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	http.HandleFunc("/ws", wsHandler)
	log.Println("Starting WebSocket server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Upgrade error", err)
	}

	defer conn.Close()

	for {

		messageType, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("error when reading", err)
			break
		}

		log.Printf("Received: %s", message)

		err = conn.WriteMessage(messageType, message)

		if err != nil {
			log.Println("error when writing", err)
			break
		}

	}
	log.Println("Client disconnected")

}
