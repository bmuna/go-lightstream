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
		log.Println("Upgrade error:", err)
		return
	}
	defer func() {
		roomManager.LeaveAllRooms(conn)
		conn.Close()
		log.Println("Client disconnected")
	}()

	// Extract roomID from query (e.g., ws://localhost:8080/ws?room=room1)
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		log.Println("Missing room ID")
		return
	}

	roomManager.JoinRoom(roomID, conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received from room %s: %s", roomID, message)

		// Broadcast to others in the room
		roomManager.BroadcastToRoom(roomID, conn, message)

		// Optional: echo back to sender
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
