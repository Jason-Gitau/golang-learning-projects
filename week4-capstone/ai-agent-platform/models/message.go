package models

import (
	"time"

	"gorm.io/gorm"
)

// MessageRole represents the role of a message sender
type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)

// Message represents a single message in a conversation
type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Role    MessageRole `gorm:"not null;index" json:"role"`
	Content string      `gorm:"type:text;not null" json:"content"`

	// Token tracking
	TokensUsed int `gorm:"default:0" json:"tokens_used"`

	// Foreign key
	ConversationID uint         `gorm:"not null;index" json:"conversation_id"`
	Conversation   Conversation `gorm:"constraint:OnDelete:CASCADE" json:"conversation,omitempty"`
}

// CreateMessageRequest represents a request to create a message
type CreateMessageRequest struct {
	Role    MessageRole `json:"role" binding:"required,oneof=user assistant system"`
	Content string      `json:"content" binding:"required"`
}

// MessageResponse represents a message response
type MessageResponse struct {
	ID             uint        `json:"id"`
	Role           MessageRole `json:"role"`
	Content        string      `json:"content"`
	TokensUsed     int         `json:"tokens_used"`
	ConversationID uint        `json:"conversation_id"`
	CreatedAt      time.Time   `json:"created_at"`
}

// ToResponse converts a Message to a MessageResponse
func (m *Message) ToResponse() MessageResponse {
	return MessageResponse{
		ID:             m.ID,
		Role:           m.Role,
		Content:        m.Content,
		TokensUsed:     m.TokensUsed,
		ConversationID: m.ConversationID,
		CreatedAt:      m.CreatedAt,
	}
}
