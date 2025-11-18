package tools

import (
	"context"
)

// RegistryAdapter adapts Registry to the agent.ToolRegistry interface
type RegistryAdapter struct {
	registry *Registry
}

// NewRegistryAdapter creates a new registry adapter
func NewRegistryAdapter(registry *Registry) *RegistryAdapter {
	return &RegistryAdapter{registry: registry}
}

// AdapterToolResult matches agent.ToolResult
type AdapterToolResult struct {
	Tool    string
	Success bool
	Result  interface{}
	Error   string
}

// Execute executes a tool and returns adapted result
func (a *RegistryAdapter) Execute(ctx context.Context, name string, params map[string]interface{}) (AdapterToolResult, error) {
	result, err := a.registry.Execute(ctx, name, params)

	return AdapterToolResult{
		Tool:    result.Tool,
		Success: result.Success,
		Result:  result.Result,
		Error:   result.Error,
	}, err
}

// GetToolDescriptions returns tool descriptions
func (a *RegistryAdapter) GetToolDescriptions() map[string]string {
	return a.registry.GetToolDescriptions()
}
