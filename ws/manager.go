/*
WebSocket manager logic and client registry:

The clients map and mu mutex (to track all connected clients and protect concurrent access).

The broadcast channel.

The handleWSConnection function (upgrades HTTP to WebSocket, creates a new Client, starts read/write goroutines, adds/removes clients to/from the registry).

The broadcaster function (reads from broadcast and sends to all clients).
*/
package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	Clients   map[*Client]bool
	Broadcast chan []byte
	Mu        sync.Mutex
	Upgrader  websocket.Upgrader
}

func NewManager() *Manager { 
	return &Manager{
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan []byte),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
		},
	}
}

func (m *Manager) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	m.Mu.Lock()
	if len(m.Clients) >= 4 { // Limit to 4 clients
		m.Mu.Unlock()
		http.Error(w, "Room is full (max 4 players)", http.StatusTooManyRequests)
		return
	}
	m.Mu.Unlock()
	
	conn, err := m.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{Conn: conn, Send: make(chan []byte)}
	m.Mu.Lock()
	m.Clients[client] = true
	m.Mu.Unlock()

	go m.readMessages(client)
	go m.writeMessages(client)
}

func (m *Manager) Broadcaster() {
	for {
		msg := <-m.Broadcast
		m.Mu.Lock()
		for client := range m.Clients {
			select {
			case client.Send <- msg: // Send message to each client
			default:
				close(client.Send)      // Close channel if client is not ready
				delete(m.Clients, client) // Remove client from map
			}
		}
		m.Mu.Unlock()
	}
}
