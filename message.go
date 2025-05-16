package main

type Message struct {
	Type    string `json:"type"`    // "join", "offer", "answer", "ice"
	RoomID  string `json:"roomID"`  // ID of the room
	Payload string `json:"payload"` // Actual offer/answer/ICE data (as JSON string)
}
