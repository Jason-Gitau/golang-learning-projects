package tools

import "context"

// Tool represents an executable tool that agents can use
type Tool interface {
	// Name returns the unique name of the tool
	Name() string

	// Description returns a human-readable description of what the tool does
	Description() string

	// Parameters returns a JSON schema describing the tool's parameters
	Parameters() map[string]interface{}

	// Execute runs the tool with the given parameters
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Tool    string      `json:"tool"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Parameter represents a tool parameter definition
type Parameter struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}
