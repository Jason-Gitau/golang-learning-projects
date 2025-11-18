package messaging

import (
	"github.com/golang-learning/agent-orchestrator/models"
)

// MessageType represents different types of messages in the system
type MessageType string

const (
	MessageTypeRequest  MessageType = "request"
	MessageTypeResponse MessageType = "response"
	MessageTypeShutdown MessageType = "shutdown"
	MessageTypeHealth   MessageType = "health"
)

// Message represents a message passed between components
type Message struct {
	Type     MessageType      `json:"type"`
	Request  *models.Request  `json:"request,omitempty"`
	Response *models.Response `json:"response,omitempty"`
}

// NewRequestMessage creates a new request message
func NewRequestMessage(req *models.Request) *Message {
	return &Message{
		Type:    MessageTypeRequest,
		Request: req,
	}
}

// NewResponseMessage creates a new response message
func NewResponseMessage(resp *models.Response) *Message {
	return &Message{
		Type:     MessageTypeResponse,
		Response: resp,
	}
}

// NewShutdownMessage creates a new shutdown message
func NewShutdownMessage() *Message {
	return &Message{
		Type: MessageTypeShutdown,
	}
}

// NewHealthMessage creates a new health check message
func NewHealthMessage() *Message {
	return &Message{
		Type: MessageTypeHealth,
	}
}
