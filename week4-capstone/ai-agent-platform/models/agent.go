package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// AgentStatus represents the status of an agent
type AgentStatus string

const (
	AgentStatusActive   AgentStatus = "active"
	AgentStatusInactive AgentStatus = "inactive"
	AgentStatusArchived AgentStatus = "archived"
)

// ToolConfig represents a tool configuration for an agent
type ToolConfig struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// ToolConfigs is a slice of ToolConfig that implements sql.Scanner and driver.Valuer
type ToolConfigs []ToolConfig

// Value implements driver.Valuer interface for database storage
func (t ToolConfigs) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

// Scan implements sql.Scanner interface for database retrieval
func (t *ToolConfigs) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan ToolConfigs: value is not []byte")
	}

	return json.Unmarshal(bytes, t)
}

// Agent represents an AI agent
type Agent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name         string      `gorm:"not null;index" json:"name" binding:"required"`
	Description  string      `gorm:"type:text" json:"description"`
	SystemPrompt string      `gorm:"type:text;not null" json:"system_prompt" binding:"required"`
	Model        string      `gorm:"not null;default:gpt-4" json:"model" binding:"required"`
	Temperature  float32     `gorm:"default:0.7" json:"temperature" binding:"gte=0,lte=2"`
	MaxTokens    int         `gorm:"default:2000" json:"max_tokens" binding:"gte=0"`
	Status       AgentStatus `gorm:"default:active;index" json:"status"`
	Tools        ToolConfigs `gorm:"type:json" json:"tools"`

	// Foreign key
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// Relationships
	Conversations []Conversation `gorm:"foreignKey:AgentID;constraint:OnDelete:CASCADE" json:"conversations,omitempty"`
	UsageLogs     []UsageLog     `gorm:"foreignKey:AgentID;constraint:OnDelete:SET NULL" json:"-"`
}

// CreateAgentRequest represents a request to create an agent
type CreateAgentRequest struct {
	Name         string      `json:"name" binding:"required"`
	Description  string      `json:"description"`
	SystemPrompt string      `json:"system_prompt" binding:"required"`
	Model        string      `json:"model" binding:"required"`
	Temperature  float32     `json:"temperature" binding:"gte=0,lte=2"`
	MaxTokens    int         `json:"max_tokens" binding:"gte=0"`
	Tools        ToolConfigs `json:"tools"`
}

// UpdateAgentRequest represents a request to update an agent
type UpdateAgentRequest struct {
	Name         *string      `json:"name"`
	Description  *string      `json:"description"`
	SystemPrompt *string      `json:"system_prompt"`
	Model        *string      `json:"model"`
	Temperature  *float32     `json:"temperature" binding:"omitempty,gte=0,lte=2"`
	MaxTokens    *int         `json:"max_tokens" binding:"omitempty,gte=0"`
	Status       *AgentStatus `json:"status"`
	Tools        *ToolConfigs `json:"tools"`
}

// AgentResponse represents an agent response
type AgentResponse struct {
	ID           uint        `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	SystemPrompt string      `json:"system_prompt"`
	Model        string      `json:"model"`
	Temperature  float32     `json:"temperature"`
	MaxTokens    int         `json:"max_tokens"`
	Status       AgentStatus `json:"status"`
	Tools        ToolConfigs `json:"tools"`
	UserID       uint        `json:"user_id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// ToResponse converts an Agent to an AgentResponse
func (a *Agent) ToResponse() AgentResponse {
	return AgentResponse{
		ID:           a.ID,
		Name:         a.Name,
		Description:  a.Description,
		SystemPrompt: a.SystemPrompt,
		Model:        a.Model,
		Temperature:  a.Temperature,
		MaxTokens:    a.MaxTokens,
		Status:       a.Status,
		Tools:        a.Tools,
		UserID:       a.UserID,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
