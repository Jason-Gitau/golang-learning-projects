package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB
)

// Client represents a WebSocket client connection
type Client struct {
	ID             string
	UserID         string
	ConversationID string
	Hub            *Hub
	Conn           *websocket.Conn
	Send           chan *Message
	mu             sync.Mutex
}

// NewClient creates a new WebSocket client
func NewClient(id, userID, conversationID string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ID:             id,
		UserID:         userID,
		ConversationID: conversationID,
		Hub:            hub,
		Conn:           conn,
		Send:           make(chan *Message, 256),
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %s: %v", c.ID, err)
			}
			break
		}

		// Handle ping messages
		if msg.Type == MessageTypePing {
			pongMsg := NewMessage(MessageTypePong, "")
			c.Send <- pongMsg
			continue
		}

		// Set metadata
		if msg.Metadata == nil {
			msg.Metadata = make(map[string]interface{})
		}
		msg.Metadata["client_id"] = c.ID
		msg.Metadata["user_id"] = c.UserID
		msg.ConversationID = c.ConversationID
		msg.Timestamp = time.Now()

		// Forward message to hub for processing
		c.Hub.Broadcast <- &BroadcastMessage{
			Message:        &msg,
			ConversationID: c.ConversationID,
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error writing message to client %s: %v", c.ID, err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage safely sends a message to the client
func (c *Client) SendMessage(msg *Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.Send <- msg:
		return nil
	default:
		// Channel full, client may be slow or disconnected
		return websocket.ErrCloseSent
	}
}

// Close gracefully closes the client connection
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.Send:
	default:
		close(c.Send)
	}
}

// MarshalJSON implements custom JSON marshaling for logging
func (c *Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID             string `json:"id"`
		UserID         string `json:"user_id"`
		ConversationID string `json:"conversation_id"`
	}{
		ID:             c.ID,
		UserID:         c.UserID,
		ConversationID: c.ConversationID,
	})
}
