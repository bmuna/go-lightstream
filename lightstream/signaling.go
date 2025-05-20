package lightstream

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type LightstreamServer struct {
	roomManager *RoomManager
	upgrader    websocket.Upgrader
}

func NewServer() *LightstreamServer {
	return &LightstreamServer{
		roomManager: NewRoomManager(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *LightstreamServer) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var currentRoom, currentUser string

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
			currentUser = msg.SenderID
			s.roomManager.JoinRoom(currentRoom, currentUser, conn)
			log.Println("User joined room:", currentRoom)

			existingPeers := s.roomManager.GetPeerIDsInRoomExcept(currentRoom, currentUser)
			response := Message{
				Type:     "peers",
				SenderID: "server",
				RoomID:   currentRoom,
				Peers:    existingPeers,
			}
			peerBytes, _ := json.Marshal(response)
			conn.WriteMessage(websocket.TextMessage, peerBytes)

			notification := Message{
				Type:     "user-joined",
				SenderID: currentUser,
				RoomID:   currentRoom,
			}
			notifyBytes, _ := json.Marshal(notification)
			s.roomManager.BroadcastToRoom(currentRoom, conn, notifyBytes)

		case "message", "offer", "answer", "ice-candidate":
			if msg.TargetID != "" {
				s.roomManager.SendToUser(msg.TargetID, messageBytes)
			} else {
				s.roomManager.BroadcastToRoom(msg.RoomID, conn, messageBytes)
			}

		default:
			log.Println("Unknown message type:", msg.Type)
		}
	}

	if currentRoom != "" && currentUser != "" {
		s.roomManager.LeaveRoom(currentRoom, conn)
		log.Println("User left room:", currentRoom)

		notification := Message{
			Type:     "user-left",
			SenderID: currentUser,
			RoomID:   currentRoom,
		}
		notifyBytes, _ := json.Marshal(notification)
		s.roomManager.BroadcastToRoom(currentRoom, conn, notifyBytes)
	}
}
