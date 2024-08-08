package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/wilsonle/excalidraw-collab/constants"
)

func generateSessionId() string {
	return uuid.New().String()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var connections = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Retrieve the username from the context
	_, ok := r.Context().Value(constants.USERNAME_CONTEXT_KEY).(string)
	if !ok {
		log.Println("Unable to retrieve username from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate a unique session ID for this connection
	sessionId := generateSessionId()

	// Send the session ID to the client
	err = conn.WriteMessage(websocket.TextMessage, []byte("session_id_"+sessionId))
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	// Add connection to the global slice
	mu.Lock()
	connections[conn] = true
	mu.Unlock()

	// Remove connection from the global slice when done
	defer func() {
		mu.Lock()
		delete(connections, conn)
		mu.Unlock()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Embed the session ID in the message
		messageWithId := map[string]interface{}{
			"sessionId":      sessionId,
			"excalidrawData": string(message),
		}
		messageWithIdBytes, err := json.Marshal(messageWithId)
		if err != nil {
			log.Println("JSON marshal error:", err)
			continue
		}

		// Broadcast the message to all connected clients
		broadcastMessage(messageType, messageWithIdBytes)
	}
}

func broadcastMessage(messageType int, message []byte) {
	mu.Lock()
	defer mu.Unlock()
	for conn := range connections {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			conn.Close()
			delete(connections, conn)
		}
	}
}
