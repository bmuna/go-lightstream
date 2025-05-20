
# go-lightstream

A lightweight, open-source Go library for live audio and video streaming conferencing, optimized for low bandwidth.  
Developers can easily embed this library in their apps or services to handle WebSocket signaling, room management, and peer messaging for WebRTC-based streaming.

---

## Features

- WebSocket-based signaling server
- Room and user connection management
- Broadcast messages within rooms (excluding sender)
- Send direct messages (offers, answers, ICE candidates) to specific peers
- Lightweight and modular — embed easily into your Go apps

---

## Installation

```bash
go get github.com/bmuna/go-lightstream
```

---

## Usage

1. Import the package

```go
import (
  "github.com/bmuna/go-lightstream"
)
```

2. Initialize RoomManager and HTTP WebSocket handler

```go
package main

import (
  "log"
  "net/http"

  "github.com/gorilla/websocket"
  "github.com/yourusername/go-lightstream"
)

var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool { return true },
}

var roomManager = lightstream.NewRoomManager()

func wsHandler(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Print("Upgrade error:", err)
    return
  }
  defer conn.Close()

  // Your WebSocket message handling logic here
  // See example/main.go for full code
}

func main() {
  http.HandleFunc("/ws", wsHandler)
  log.Println("Starting server on :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## WebSocket message format

All messages sent and received via WebSocket should be JSON with the following structure:

```json
{
  "type": "join|message|offer|answer|ice-candidate",
  "senderID": "string",
  "roomID": "string",
  "targetID": "string (optional)",
  "peers": ["string"] (optional),
  "data": "object or string (optional)"
}
```

- **type**: The message category, e.g., "join" to join a room, or "offer" to send a WebRTC offer.
- **senderID**: The unique identifier for the user sending the message.
- **roomID**: The room identifier.
- **targetID**: Optional, the ID of the user to send a direct message to.
- **peers**: Optional, list of peer IDs (used in server responses).
- **data**: Optional, any extra data depending on message type.

---

## How It Works

- When a user connects and sends a `"join"` message, they are added to the specified room.
- The server replies with a list of existing peers in the room.
- The server broadcasts `"user-joined"` and `"user-left"` notifications to all participants.
- `"message"` type messages are broadcast to all other clients in the room.
- `"offer"`, `"answer"`, and `"ice-candidate"` messages are sent directly to the specified target user.

---

## Contributing

Contributions, issues, and feature requests are welcome!  
Feel free to check [issues page](https://github.com/bmuna/go-lightstream/issues) or submit a pull request.

---

## License

[![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)  
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


---

## Contact

Created by [Your Name](https://github.com/bmuna) – feel free to reach out!
