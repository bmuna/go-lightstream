package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

//	{
//		"room" : {
//			"connA": "Abel"
//			"connB": "Biruk"
//		}
//	}

type RoomManager struct {
	rooms map[string]map[*websocket.Conn]bool
	users map[string]*websocket.Conn
	lock  sync.RWMutex
}

// Constructor for RoomManager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]map[*websocket.Conn]bool),
		users: make(map[string]*websocket.Conn),
	}
}

// JoinRoom adds a WebSocket connection to the specified room
func (rm *RoomManager) JoinRoom(roomID string, userID string, conn *websocket.Conn) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if _, exists := rm.rooms[roomID]; !exists {
		rm.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	rm.rooms[roomID][conn] = true
	rm.users[userID] = conn
}


// LeaveAllRooms removes a WebSocket connection from all rooms
func (rm *RoomManager) LeaveAllRooms(conn *websocket.Conn) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	for _, clients := range rm.rooms {
		delete(clients, conn)
	}
}

func (rm *RoomManager) LeaveRoom(roomID string, conn *websocket.Conn) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	// Remove from room
	if clients, ok := rm.rooms[roomID]; ok {
		delete(clients, conn)
		if len(clients) == 0 {
			delete(rm.rooms, roomID)
		}
	}

	// Remove from users
	for userID, userConn := range rm.users {
		if userConn == conn {
			delete(rm.users, userID)
			break
		}
	}
}

// BroadcastToRoom messenger except the sender
func (rm *RoomManager) BroadcastToRoom(roomID string, sender *websocket.Conn, message []byte) {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	conns := rm.rooms[roomID]
	for conn := range conns {
		if conn != sender {
			_ = conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (rm *RoomManager) SendToUser(targetID string, message []byte) {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	if conn, found := rm.users[targetID]; found {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
