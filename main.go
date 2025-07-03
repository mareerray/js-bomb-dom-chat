package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	// playerID string // Optional: to identify the player
}

var (
	clients   = make(map[*Client]bool) // Connected clients
	broadcast = make(chan []byte)      // Broadcast channel for messages
	mu        sync.Mutex               // Mutex to protect the client map
)

func handleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	// playerID := uuid.New().String() // Generate a unique ID for the player
	client := &Client{conn: conn, send: make(chan []byte)}
	// initMsg := map[string]interface{}{"type": "init", "playerID": playerID}
	// msgBytes, _ := json.Marshal(initMsg)
	// client.send <- msgBytes

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	go readMessages(client)
	go writeMessages(client)
}

func readMessages(client *Client) {
	defer func() { /* ...disconnect logic... */
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
		client.conn.Close()
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err == nil {
			if data["type"] == "chat" {
				broadcast <- msg // Broadcast chat messages
			}
		}
		// msgType, _ := data["type"].(string)
		// switch msgType {
		// 	case "move":
		// 		// Handle move messages (e.g., game actions)
		// 	case "chat":
		// 		// Handle chat messages
		// 	case "bomb":
		// 		// Handle bomb messages
		// }
		// broadcast <- msg // Broadcast the received message
	}
}

func writeMessages(client *Client) {
	for msg := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

func broadcaster() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			select {
			case client.send <- msg: // Send message to each client
			default:
				close(client.send)      // Close channel if client is not ready
				delete(clients, client) // Remove client from map
			}
		}
		mu.Unlock()
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Serve static files from ./static directory
	http.HandleFunc("/ws", handleWSConnection)              // Handle WebSocket connections

	go broadcaster() // Start the broadcaster in a separate goroutine

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

/*
This server:

- Serves your static files.
- Handles WebSocket connections at /ws.
- Broadcasts any message received from any client to all connected clients
(good for both chat and game actions) */
