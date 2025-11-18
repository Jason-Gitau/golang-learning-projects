package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper CORS checking
		// For now, allow all origins
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	Hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		Hub: hub,
	}
}

// HandleWebSocket handles WebSocket connection requests
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get conversation ID from URL
	vars := mux.Vars(r)
	conversationID := vars["conversation_id"]

	if conversationID == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	// For now, we'll extract from query param or header as fallback
	userID := h.extractUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	clientID := uuid.New().String()
	client := NewClient(clientID, userID, conversationID, h.Hub, conn)

	// Register client with hub
	h.Hub.Register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket connection established: client=%s, user=%s, conversation=%s",
		clientID, userID, conversationID)
}

// extractUserID extracts user ID from request
// In production, this would come from JWT token validation
func (h *Handler) extractUserID(r *http.Request) string {
	// Try to get from context (set by auth middleware)
	if userID := r.Context().Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}

	// Fallback: Try query parameter (for development/testing)
	if userID := r.URL.Query().Get("user_id"); userID != "" {
		return userID
	}

	// Fallback: Try header
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return userID
	}

	return ""
}

// HandleHealthCheck handles health check requests
func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"status":         "healthy",
		"total_clients":  h.Hub.GetTotalClientCount(),
		"service":        "websocket",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Simple JSON response without encoding package dependency
	w.Write([]byte(`{"status":"healthy","service":"websocket"}`))
}
