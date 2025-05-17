package main

type Message struct {
	Type     string `json:"type"`    // "join", "offer", "answer", "ice"
	RoomID   string `json:"roomID"`  // ID of the room
	SenderID string `json:sender_id` // ID of the sender
	TargetID string `json:target_id` //ID of the target(recipient)
	Payload  string `json:"payload"` // Actual offer/answer/ICE data (as JSON string)
}
