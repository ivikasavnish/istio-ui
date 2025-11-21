package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	
	"github.com/ivikasavnish/istio-ui/backend/internal/models"
)

// HandleWebSocket handles WebSocket connections
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	s.wsMutex.Lock()
	s.wsClients[conn] = true
	s.wsMutex.Unlock()

	log.Printf("New WebSocket client connected. Total clients: %d", len(s.wsClients))

	// Send welcome message
	welcomeMsg := models.WebSocketMessage{
		Type: "connected",
		Payload: map[string]string{
			"message": "Connected to MeshControl Center",
		},
	}
	conn.WriteJSON(welcomeMsg)

	// Keep connection alive and handle messages
	for {
		var msg models.WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Handle incoming messages
		s.handleWebSocketMessage(conn, &msg)
	}

	// Remove client on disconnect
	s.wsMutex.Lock()
	delete(s.wsClients, conn)
	s.wsMutex.Unlock()

	log.Printf("WebSocket client disconnected. Total clients: %d", len(s.wsClients))
}

// handleWebSocketMessage handles incoming WebSocket messages
func (s *Server) handleWebSocketMessage(conn *websocket.Conn, msg *models.WebSocketMessage) {
	switch msg.Type {
	case "ping":
		conn.WriteJSON(models.WebSocketMessage{
			Type:    "pong",
			Payload: map[string]string{"status": "alive"},
		})
	case "subscribe":
		// Handle subscription requests
		log.Printf("Client subscribed to: %v", msg.Payload)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// broadcastEvent broadcasts an event to all connected WebSocket clients
func (s *Server) broadcastEvent(eventType string, payload interface{}) {
	msg := models.WebSocketMessage{
		Type:    eventType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v", err)
		return
	}

	s.wsMutex.RLock()
	defer s.wsMutex.RUnlock()

	for conn := range s.wsClients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
		}
	}
}
