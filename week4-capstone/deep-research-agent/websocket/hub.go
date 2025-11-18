package websocket

import (
	"log"
	"sync"
)

// Hub maintains the set of active WebSocket clients for research jobs
type Hub struct {
	// Registered clients by job ID
	clients map[string]map[*Client]bool

	// All clients by client ID for quick lookup
	clientsByID map[string]*Client

	// Broadcast channel for sending messages to clients
	Broadcast chan *BroadcastMessage

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	mu sync.RWMutex
}

// BroadcastMessage wraps a message with routing information
type BroadcastMessage struct {
	Message        *Message
	JobID          string
	TargetClientID string // If empty, broadcast to all clients watching this job
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]map[*Client]bool),
		clientsByID: make(map[string]*Client),
		Broadcast:   make(chan *BroadcastMessage, 256),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	log.Println("WebSocket Hub started")
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case broadcastMsg := <-h.Broadcast:
			h.handleBroadcast(broadcastMsg)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Add to job clients
	if h.clients[client.JobID] == nil {
		h.clients[client.JobID] = make(map[*Client]bool)
	}
	h.clients[client.JobID][client] = true

	// Add to global client map
	h.clientsByID[client.ID] = client

	log.Printf("WebSocket client registered: %s (job: %s)", client.ID, client.JobID)

	// Send connection confirmation
	confirmMsg := NewConnectionConfirmMessage(client.JobID)
	client.Send <- confirmMsg
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove from job clients
	if clients, ok := h.clients[client.JobID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)

			// Clean up empty job maps
			if len(clients) == 0 {
				delete(h.clients, client.JobID)
			}
		}
	}

	// Remove from global client map
	delete(h.clientsByID, client.ID)

	// Close the client's send channel
	client.Close()

	log.Printf("WebSocket client unregistered: %s (job: %s)", client.ID, client.JobID)
}

// handleBroadcast processes and routes messages
func (h *Hub) handleBroadcast(broadcastMsg *BroadcastMessage) {
	msg := broadcastMsg.Message

	// Handle ping messages
	if msg.Type == MessageTypePing {
		return
	}

	// Broadcast to job or specific client
	if broadcastMsg.TargetClientID != "" {
		h.sendToClient(broadcastMsg.TargetClientID, msg)
	} else {
		h.broadcastToJob(broadcastMsg.JobID, msg)
	}
}

// broadcastToJob sends a message to all clients watching a specific job
func (h *Hub) broadcastToJob(jobID string, msg *Message) {
	h.mu.RLock()
	clients := h.clients[jobID]
	h.mu.RUnlock()

	if len(clients) == 0 {
		return
	}

	for client := range clients {
		select {
		case client.Send <- msg:
		default:
			// Client is slow or disconnected
			log.Printf("Failed to send to client %s, unregistering", client.ID)
			go func(c *Client) {
				h.Unregister <- c
			}(client)
		}
	}
}

// sendToClient sends a message to a specific client
func (h *Hub) sendToClient(clientID string, msg *Message) {
	h.mu.RLock()
	client, exists := h.clientsByID[clientID]
	h.mu.RUnlock()

	if !exists {
		log.Printf("Client %s not found", clientID)
		return
	}

	select {
	case client.Send <- msg:
	default:
		log.Printf("Failed to send to client %s, unregistering", client.ID)
		go func(c *Client) {
			h.Unregister <- c
		}(client)
	}
}

// BroadcastToJob broadcasts a message to all clients watching a specific job
func (h *Hub) BroadcastToJob(jobID string, msg *Message) {
	msg.JobID = jobID
	h.Broadcast <- &BroadcastMessage{
		Message: msg,
		JobID:   jobID,
	}
}

// GetJobClientCount returns the number of clients watching a specific job
func (h *Hub) GetJobClientCount(jobID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.clients[jobID])
}

// GetTotalClientCount returns the total number of connected clients
func (h *Hub) GetTotalClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.clientsByID)
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Println("Shutting down WebSocket hub...")

	// Close all client connections
	for _, client := range h.clientsByID {
		client.Close()
		client.Conn.Close()
	}

	// Clear all maps
	h.clients = make(map[string]map[*Client]bool)
	h.clientsByID = make(map[string]*Client)

	log.Println("WebSocket hub shutdown complete")
}
