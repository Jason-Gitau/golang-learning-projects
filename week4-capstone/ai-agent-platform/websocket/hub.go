package websocket

import (
	"log"
	"sync"
)

// BroadcastMessage wraps a message with routing information
type BroadcastMessage struct {
	Message        *Message
	ConversationID string
	TargetClientID string // If empty, broadcast to all clients in conversation
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by conversation ID
	clients map[string]map[*Client]bool

	// All clients by client ID for quick lookup
	clientsByID map[string]*Client

	// Inbound messages from clients
	Broadcast chan *BroadcastMessage

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Agent message queue for processing
	AgentQueue chan *Message

	mu sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub(agentQueue chan *Message) *Hub {
	return &Hub{
		clients:     make(map[string]map[*Client]bool),
		clientsByID: make(map[string]*Client),
		Broadcast:   make(chan *BroadcastMessage, 256),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		AgentQueue:  agentQueue,
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
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

	// Add to conversation clients
	if h.clients[client.ConversationID] == nil {
		h.clients[client.ConversationID] = make(map[*Client]bool)
	}
	h.clients[client.ConversationID][client] = true

	// Add to global client map
	h.clientsByID[client.ID] = client

	log.Printf("Client registered: %s (user: %s, conversation: %s)",
		client.ID, client.UserID, client.ConversationID)

	// Send connection confirmation
	confirmMsg := NewMessage(MessageTypePong, "Connected successfully")
	confirmMsg.ConversationID = client.ConversationID
	client.Send <- confirmMsg
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove from conversation clients
	if clients, ok := h.clients[client.ConversationID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)

			// Clean up empty conversation maps
			if len(clients) == 0 {
				delete(h.clients, client.ConversationID)
			}
		}
	}

	// Remove from global client map
	delete(h.clientsByID, client.ID)

	// Close the client's send channel
	client.Close()

	log.Printf("Client unregistered: %s (user: %s, conversation: %s)",
		client.ID, client.UserID, client.ConversationID)
}

// handleBroadcast processes and routes messages
func (h *Hub) handleBroadcast(broadcastMsg *BroadcastMessage) {
	msg := broadcastMsg.Message

	// If it's a user message, send to agent queue for processing
	if msg.Type == MessageTypeUserMessage {
		select {
		case h.AgentQueue <- msg:
			log.Printf("User message queued for processing: conversation=%s",
				broadcastMsg.ConversationID)

			// Send typing indicator to all clients in conversation
			typingMsg := NewMessage(MessageTypeTypingIndicator, "Agent is typing...")
			typingMsg.ConversationID = broadcastMsg.ConversationID
			h.broadcastToConversation(broadcastMsg.ConversationID, typingMsg)
		default:
			// Queue is full
			errMsg := NewErrorMessage("System is busy, please try again")
			errMsg.ConversationID = broadcastMsg.ConversationID
			h.broadcastToConversation(broadcastMsg.ConversationID, errMsg)
		}
		return
	}

	// For other message types, broadcast to conversation or specific client
	if broadcastMsg.TargetClientID != "" {
		h.sendToClient(broadcastMsg.TargetClientID, msg)
	} else {
		h.broadcastToConversation(broadcastMsg.ConversationID, msg)
	}
}

// broadcastToConversation sends a message to all clients in a conversation
func (h *Hub) broadcastToConversation(conversationID string, msg *Message) {
	h.mu.RLock()
	clients := h.clients[conversationID]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.Send <- msg:
		default:
			// Client is slow or disconnected
			log.Printf("Failed to send to client %s, unregistering", client.ID)
			h.Unregister <- client
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
		h.Unregister <- client
	}
}

// BroadcastToConversation broadcasts a message to all clients in a conversation
func (h *Hub) BroadcastToConversation(conversationID string, msg *Message) {
	msg.ConversationID = conversationID
	h.Broadcast <- &BroadcastMessage{
		Message:        msg,
		ConversationID: conversationID,
	}
}

// GetConversationClientCount returns the number of clients in a conversation
func (h *Hub) GetConversationClientCount(conversationID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.clients[conversationID])
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
