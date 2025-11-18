package tools

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Registry manages all available tools
type Registry struct {
	tools   map[string]Tool
	mu      sync.RWMutex
	timeout time.Duration
}

// NewRegistry creates a new tool registry
func NewRegistry(timeout time.Duration) *Registry {
	return &Registry{
		tools:   make(map[string]Tool),
		timeout: timeout,
	}
}

// Register adds a tool to the registry
func (r *Registry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s is already registered", name)
	}

	r.tools[name] = tool
	log.Printf("Tool registered: %s - %s", name, tool.Description())
	return nil
}

// Unregister removes a tool from the registry
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tools, name)
	log.Printf("Tool unregistered: %s", name)
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	return tool, nil
}

// List returns all registered tools
func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// Execute runs a tool with the given parameters
func (r *Registry) Execute(ctx context.Context, name string, params map[string]interface{}) (*ToolResult, error) {
	tool, err := r.Get(name)
	if err != nil {
		return &ToolResult{
			Tool:    name,
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Execute tool
	result, err := tool.Execute(execCtx, params)
	if err != nil {
		return &ToolResult{
			Tool:    name,
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &ToolResult{
		Tool:    name,
		Success: true,
		Result:  result,
	}, nil
}

// GetToolDescriptions returns a map of tool names to descriptions
func (r *Registry) GetToolDescriptions() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	descriptions := make(map[string]string)
	for name, tool := range r.tools {
		descriptions[name] = tool.Description()
	}

	return descriptions
}

// GetToolParameters returns the parameters for a specific tool
func (r *Registry) GetToolParameters(name string) (map[string]interface{}, error) {
	tool, err := r.Get(name)
	if err != nil {
		return nil, err
	}

	return tool.Parameters(), nil
}

// Count returns the number of registered tools
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.tools)
}

// RegisterDefaultTools registers all default tools
func (r *Registry) RegisterDefaultTools() error {
	defaultTools := []Tool{
		NewCalculatorTool(),
		NewWeatherTool(),
		NewSearchTool(),
		NewDateTimeTool(),
	}

	for _, tool := range defaultTools {
		if err := r.Register(tool); err != nil {
			return err
		}
	}

	return nil
}
