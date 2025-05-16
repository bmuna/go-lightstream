package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var roomManager = NewRoomManager()

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
		log.Print("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var currentRoom string

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch msg.Type {
		case "join":
			currentRoom = msg.RoomID
			roomManager.JoinRoom(msg.RoomID, conn)
			log.Println("User joined room:", msg.RoomID)

		case "message":
			// Broadcast message to everyone else in the same room
			roomManager.BroadcastToRoom(msg.RoomID, conn, messageBytes)

		default:
			log.Println("Unknown message type:", msg.Type)
		}
	}

	// Cleanup on disconnect
	if currentRoom != "" {
		roomManager.LeaveRoom(currentRoom, conn)
		log.Println("User left room:", currentRoom)
	}
}
