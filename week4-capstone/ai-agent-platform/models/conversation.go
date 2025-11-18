package models

import (
	"time"

	"gorm.io/gorm"
)

// ConversationStatus represents the status of a conversation
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"
	ConversationStatusArchived ConversationStatus = "archived"
	ConversationStatusDeleted  ConversationStatus = "deleted"
)

// Conversation represents a chat conversation
type Conversation struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title  string             `gorm:"not null" json:"title"`
	Status ConversationStatus `gorm:"default:active;index" json:"status"`

	// Foreign keys
	AgentID uint  `gorm:"not null;index" json:"agent_id"`
	Agent   Agent `gorm:"constraint:OnDelete:CASCADE" json:"agent,omitempty"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// Relationships
	Messages []Message `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
}

// CreateConversationRequest represents a request to create a conversation
type CreateConversationRequest struct {
	Title string `json:"title" binding:"required"`
}

// ConversationResponse represents a conversation response
type ConversationResponse struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Status    ConversationStatus `json:"status"`
	AgentID   uint               `json:"agent_id"`
	UserID    uint               `json:"user_id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// ConversationWithMessagesResponse represents a conversation with its messages
type ConversationWithMessagesResponse struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Status    ConversationStatus `json:"status"`
	AgentID   uint               `json:"agent_id"`
	UserID    uint               `json:"user_id"`
	Messages  []MessageResponse  `json:"messages"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// ToResponse converts a Conversation to a ConversationResponse
func (c *Conversation) ToResponse() ConversationResponse {
	return ConversationResponse{
		ID:        c.ID,
		Title:     c.Title,
		Status:    c.Status,
		AgentID:   c.AgentID,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToResponseWithMessages converts a Conversation to a ConversationWithMessagesResponse
func (c *Conversation) ToResponseWithMessages() ConversationWithMessagesResponse {
	messages := make([]MessageResponse, len(c.Messages))
	for i, msg := range c.Messages {
		messages[i] = msg.ToResponse()
	}

	return ConversationWithMessagesResponse{
		ID:        c.ID,
		Title:     c.Title,
		Status:    c.Status,
		AgentID:   c.AgentID,
		UserID:    c.UserID,
		Messages:  messages,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
