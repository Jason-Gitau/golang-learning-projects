package models

import (
	"time"

	"gorm.io/gorm"
)

// UsageLog represents a log entry for API usage tracking
type UsageLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Token tracking
	InputTokens  int `gorm:"not null;default:0" json:"input_tokens"`
	OutputTokens int `gorm:"not null;default:0" json:"output_tokens"`
	TotalTokens  int `gorm:"not null;default:0" json:"total_tokens"`

	// Cost tracking (in cents)
	Cost float64 `gorm:"not null;default:0" json:"cost"`

	// Model used
	Model string `gorm:"not null;index" json:"model"`

	// Operation type (e.g., "chat", "completion")
	Operation string `gorm:"not null;index" json:"operation"`

	// Foreign keys
	UserID *uint `gorm:"index" json:"user_id"`
	User   *User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`

	AgentID *uint  `gorm:"index" json:"agent_id"`
	Agent   *Agent `gorm:"constraint:OnDelete:SET NULL" json:"agent,omitempty"`
}

// RateLimit represents rate limit tracking for a user
type RateLimit struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`
	User   User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// Rate limit tracking
	RequestsCount int       `gorm:"not null;default:0" json:"requests_count"`
	WindowStart   time.Time `gorm:"not null" json:"window_start"`
	WindowEnd     time.Time `gorm:"not null" json:"window_end"`
}

// UsageStatsRequest represents a request for usage statistics
type UsageStatsRequest struct {
	StartDate string `form:"start_date"` // Format: YYYY-MM-DD
	EndDate   string `form:"end_date"`   // Format: YYYY-MM-DD
	AgentID   *uint  `form:"agent_id"`
}

// UsageStats represents aggregated usage statistics
type UsageStats struct {
	TotalRequests  int     `json:"total_requests"`
	TotalTokens    int     `json:"total_tokens"`
	InputTokens    int     `json:"input_tokens"`
	OutputTokens   int     `json:"output_tokens"`
	TotalCost      float64 `json:"total_cost"`
	PeriodStart    string  `json:"period_start"`
	PeriodEnd      string  `json:"period_end"`
	ByModel        map[string]ModelUsage `json:"by_model"`
	ByAgent        map[string]AgentUsage `json:"by_agent,omitempty"`
}

// ModelUsage represents usage statistics for a specific model
type ModelUsage struct {
	Requests     int     `json:"requests"`
	TotalTokens  int     `json:"total_tokens"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	Cost         float64 `json:"cost"`
}

// AgentUsage represents usage statistics for a specific agent
type AgentUsage struct {
	AgentID      uint    `json:"agent_id"`
	AgentName    string  `json:"agent_name"`
	Requests     int     `json:"requests"`
	TotalTokens  int     `json:"total_tokens"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	Cost         float64 `json:"cost"`
}

// RateLimitInfo represents rate limit information for a user
type RateLimitInfo struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}
