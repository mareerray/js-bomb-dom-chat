package main

import (
	"fmt"
	"js-bomb-dom-chat/ws"
	"log"
	"net/http"
)

// the server setup and wiring.
func main() {
	manager := ws.NewManager()

	http.Handle("/", http.FileServer(http.Dir("./static"))) // Serve static files from ./static directory
	http.HandleFunc("/ws", manager.HandleWSConnection)      // Handle WebSocket connections

	go manager.Broadcaster() // Start the broadcaster in a separate goroutine

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
