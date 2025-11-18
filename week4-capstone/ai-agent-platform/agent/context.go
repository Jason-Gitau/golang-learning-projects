package agent

import (
	"sync"
	"time"
)

// ConversationContext maintains the context for an ongoing conversation
type ConversationContext struct {
	ConversationID string
	AgentID        string
	Messages       []ContextMessage
	Metadata       map[string]interface{}
	CreatedAt      time.Time
	UpdatedAt      time.Time
	mu             sync.RWMutex
}

// ContextMessage represents a message in the conversation context
type ContextMessage struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	ToolCalls []ToolCallRecord `json:"tool_calls,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// ToolCallRecord records a tool invocation
type ToolCallRecord struct {
	Tool   string                 `json:"tool"`
	Params map[string]interface{} `json:"params"`
	Result interface{}            `json:"result"`
	Error  string                 `json:"error,omitempty"`
}

// NewConversationContext creates a new conversation context
func NewConversationContext(conversationID, agentID string) *ConversationContext {
	return &ConversationContext{
		ConversationID: conversationID,
		AgentID:        agentID,
		Messages:       make([]ContextMessage, 0),
		Metadata:       make(map[string]interface{}),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// AddUserMessage adds a user message to the context
func (c *ConversationContext) AddUserMessage(content string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = append(c.Messages, ContextMessage{
		Role:      "user",
		Content:   content,
		Timestamp: time.Now(),
	})
	c.UpdatedAt = time.Now()
}

// AddAssistantMessage adds an assistant message to the context
func (c *ConversationContext) AddAssistantMessage(content string, toolCalls []ToolCallRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = append(c.Messages, ContextMessage{
		Role:      "assistant",
		Content:   content,
		ToolCalls: toolCalls,
		Timestamp: time.Now(),
	})
	c.UpdatedAt = time.Now()
}

// GetMessages returns a copy of all messages
func (c *ConversationContext) GetMessages() []ContextMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	messages := make([]ContextMessage, len(c.Messages))
	copy(messages, c.Messages)
	return messages
}

// GetRecentMessages returns the last N messages
func (c *ConversationContext) GetRecentMessages(n int) []ContextMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if n >= len(c.Messages) {
		messages := make([]ContextMessage, len(c.Messages))
		copy(messages, c.Messages)
		return messages
	}

	start := len(c.Messages) - n
	messages := make([]ContextMessage, n)
	copy(messages, c.Messages[start:])
	return messages
}

// MessageCount returns the total number of messages
func (c *ConversationContext) MessageCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.Messages)
}

// SetMetadata sets a metadata value
func (c *ConversationContext) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Metadata[key] = value
	c.UpdatedAt = time.Now()
}

// GetMetadata gets a metadata value
func (c *ConversationContext) GetMetadata(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.Metadata[key]
	return value, exists
}

// Clear clears the conversation context
func (c *ConversationContext) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = make([]ContextMessage, 0)
	c.Metadata = make(map[string]interface{})
	c.UpdatedAt = time.Now()
}

// ContextManager manages multiple conversation contexts
type ContextManager struct {
	contexts map[string]*ConversationContext
	mu       sync.RWMutex
}

// NewContextManager creates a new context manager
func NewContextManager() *ContextManager {
	return &ContextManager{
		contexts: make(map[string]*ConversationContext),
	}
}

// GetOrCreate gets an existing context or creates a new one
func (m *ContextManager) GetOrCreate(conversationID, agentID string) *ConversationContext {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ctx, exists := m.contexts[conversationID]; exists {
		return ctx
	}

	ctx := NewConversationContext(conversationID, agentID)
	m.contexts[conversationID] = ctx
	return ctx
}

// Get retrieves a conversation context
func (m *ContextManager) Get(conversationID string) (*ConversationContext, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ctx, exists := m.contexts[conversationID]
	return ctx, exists
}

// Delete removes a conversation context
func (m *ContextManager) Delete(conversationID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.contexts, conversationID)
}

// Count returns the number of active contexts
func (m *ContextManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.contexts)
}

// Cleanup removes contexts older than the specified duration
func (m *ContextManager) Cleanup(maxAge time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	removed := 0

	for id, ctx := range m.contexts {
		ctx.mu.RLock()
		age := now.Sub(ctx.UpdatedAt)
		ctx.mu.RUnlock()

		if age > maxAge {
			delete(m.contexts, id)
			removed++
		}
	}

	return removed
}
