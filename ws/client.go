/*
Client struct and logic for each connection:

The Client struct (holds the connection, send channel, and any player-specific data like playerID).

The readMessages(client *Client) function (reads messages from the client, handles disconnect).

The writeMessages(client *Client) function (writes messages to the client).
*/
package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

func (m *Manager) readMessages(client *Client) {
	defer func() { /* ...disconnect logic... */
		client.Conn.Close()
		m.Mu.Lock()
		delete(m.Clients, client)
		m.Mu.Unlock()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err == nil {
			switch data["type"] {
			case "chat":
				m.Broadcast <- msg // Broadcast chat messages
			default:
				log.Println("Unknown message type:", data["type"])
			}
		}
	}
}

func (m *Manager) writeMessages(client *Client) {
	defer client.Conn.Close()
	for msg := range client.Send {
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
