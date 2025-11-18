package tools

import (
	"context"
	"fmt"
	"sync"
)

// ToolRegistry manages all available tools
type ToolRegistry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}

	r.tools[name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	return tool, nil
}

// List returns all registered tools
func (r *ToolRegistry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// Execute runs a tool by name with given parameters
func (r *ToolRegistry) Execute(ctx context.Context, toolName string, params map[string]interface{}) (*ToolResult, error) {
	tool, err := r.Get(toolName)
	if err != nil {
		return nil, err
	}

	return tool.Execute(ctx, params)
}

// HasTool checks if a tool is registered
func (r *ToolRegistry) HasTool(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.tools[name]
	return exists
}

// ListToolNames returns names of all registered tools
func (r *ToolRegistry) ListToolNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetToolInfo returns information about all registered tools
func (r *ToolRegistry) GetToolInfo() map[string]ToolInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info := make(map[string]ToolInfo)
	for name, tool := range r.tools {
		info[name] = ToolInfo{
			Name:        name,
			Description: tool.Description(),
			Parameters:  tool.Parameters(),
		}
	}

	return info
}

// ToolInfo contains metadata about a tool
type ToolInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters"`
}

// RegisterAllTools registers all available tools
func RegisterAllTools(registry *ToolRegistry) error {
	tools := []Tool{
		NewPDFProcessor(),
		NewDOCXProcessor(),
		NewWebSearch(),
		NewWikipedia(),
		NewURLFetcher(),
		NewSummarizer(),
		NewCitationManager(),
		NewFactChecker(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return fmt.Errorf("failed to register tool %s: %w", tool.Name(), err)
		}
	}

	return nil
}
