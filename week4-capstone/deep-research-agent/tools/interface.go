package tools

import (
	"context"

	"deep-research-agent/models"
)

// Tool is the interface that all research tools must implement
type Tool interface {
	Name() string
	Description() string
	Parameters() []Parameter
	Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Success  bool                   `json:"success"`
	Data     interface{}            `json:"data"`
	Sources  []models.Source        `json:"sources,omitempty"`
	Error    error                  `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Parameter describes a tool parameter
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // string, int, bool, array, object
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Default     interface{} `json:"default,omitempty"`
}
