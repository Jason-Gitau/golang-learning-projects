package agent

import (
	"context"
	"fmt"
	"log"
	"time"
)

// AIService interface for AI providers
type AIService interface {
	StreamResponse(ctx context.Context, messages []AIMessage, systemPrompt string) <-chan AIStreamResponse
}

// AIMessage represents a message for the AI service
type AIMessage struct {
	Role    string
	Content string
}

// AIStreamResponse represents a streaming response from AI
type AIStreamResponse struct {
	Content  string
	Done     bool
	ToolCall *AIToolCall
	Error    error
}

// AIToolCall represents a tool call from AI
type AIToolCall struct {
	Name   string
	Params map[string]interface{}
}

// ToolRegistry interface for tool execution
type ToolRegistry interface {
	Execute(ctx context.Context, name string, params map[string]interface{}) (ToolResult, error)
	GetToolDescriptions() map[string]string
}

// ToolResult represents the result of tool execution
type ToolResult struct {
	Tool    string
	Success bool
	Result  interface{}
	Error   string
}

// AgentConfig holds agent configuration
type AgentConfig struct {
	SystemPrompt string
	Model        string
	Temperature  float64
	MaxTokens    int
	Timeout      time.Duration
}

// DefaultAgentConfig returns default agent configuration
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		SystemPrompt: "You are a helpful AI assistant with access to various tools. Use tools when appropriate to provide accurate information.",
		Model:        "mock-ai",
		Temperature:  0.7,
		MaxTokens:    2000,
		Timeout:      30 * time.Second,
	}
}

// Engine is the core agent execution engine
type Engine struct {
	ai           AIService
	tools        ToolRegistry
	contextMgr   *ContextManager
	config       *AgentConfig
}

// NewEngine creates a new agent engine
func NewEngine(ai AIService, tools ToolRegistry, config *AgentConfig) *Engine {
	if config == nil {
		config = DefaultAgentConfig()
	}

	return &Engine{
		ai:         ai,
		tools:      tools,
		contextMgr: NewContextManager(),
		config:     config,
	}
}

// ProcessMessage processes a user message and generates a response
func (e *Engine) ProcessMessage(ctx context.Context, conversationID, agentID, userMessage string) (<-chan StreamChunk, error) {
	// Get or create conversation context
	convCtx := e.contextMgr.GetOrCreate(conversationID, agentID)

	// Add user message to context
	convCtx.AddUserMessage(userMessage)

	// Create response channel
	responseChan := make(chan StreamChunk, 10)

	// Process in goroutine
	go func() {
		defer close(responseChan)

		// Build messages for AI
		messages := e.buildAIMessages(convCtx)

		// Create context with timeout
		execCtx, cancel := context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()

		// Get streaming response from AI
		aiStream := e.ai.StreamResponse(execCtx, messages, e.config.SystemPrompt)

		// Process AI stream
		fullResponse, toolCalls := e.processAIStream(execCtx, aiStream, responseChan)

		// Add assistant response to context
		convCtx.AddAssistantMessage(fullResponse, toolCalls)

		log.Printf("Message processed for conversation %s: %d messages in context",
			conversationID, convCtx.MessageCount())
	}()

	return responseChan, nil
}

// processAIStream processes the AI stream and handles tool calls
func (e *Engine) processAIStream(ctx context.Context, aiStream <-chan AIStreamResponse, responseChan chan<- StreamChunk) (string, []ToolCallRecord) {
	var fullResponse string
	var toolCalls []ToolCallRecord

	for {
		select {
		case <-ctx.Done():
			responseChan <- StreamChunk{
				Error: ctx.Err(),
				Done:  true,
			}
			return fullResponse, toolCalls

		case chunk, ok := <-aiStream:
			if !ok {
				// AI stream closed
				responseChan <- StreamChunk{Done: true}
				return fullResponse, toolCalls
			}

			// Handle errors
			if chunk.Error != nil {
				responseChan <- StreamChunk{
					Error: chunk.Error,
					Done:  true,
				}
				return fullResponse, toolCalls
			}

			// Handle tool calls
			if chunk.ToolCall != nil {
				log.Printf("AI requested tool: %s", chunk.ToolCall.Name)

				// Broadcast tool call
				responseChan <- StreamChunk{
					ToolCall: &ToolCallInfo{
						Tool:   chunk.ToolCall.Name,
						Params: chunk.ToolCall.Params,
					},
				}

				// Execute tool
				result, err := e.tools.Execute(ctx, chunk.ToolCall.Name, chunk.ToolCall.Params)

				// Record tool call
				toolCall := ToolCallRecord{
					Tool:   chunk.ToolCall.Name,
					Params: chunk.ToolCall.Params,
					Result: result.Result,
				}

				if err != nil || !result.Success {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					} else {
						errMsg = result.Error
					}
					toolCall.Error = errMsg
					log.Printf("Tool execution failed: %s - %s", chunk.ToolCall.Name, errMsg)
				}

				toolCalls = append(toolCalls, toolCall)

				// Broadcast tool result
				responseChan <- StreamChunk{
					ToolResult: &ToolResultInfo{
						Tool:   chunk.ToolCall.Name,
						Result: result.Result,
						Error:  toolCall.Error,
					},
				}

				continue
			}

			// Handle content
			if chunk.Content != "" {
				fullResponse += chunk.Content
				responseChan <- StreamChunk{
					Content: chunk.Content,
				}
			}

			// Handle completion
			if chunk.Done {
				responseChan <- StreamChunk{Done: true}
				return fullResponse, toolCalls
			}
		}
	}
}

// buildAIMessages builds the message list for AI from conversation context
func (e *Engine) buildAIMessages(convCtx *ConversationContext) []AIMessage {
	contextMessages := convCtx.GetMessages()
	aiMessages := make([]AIMessage, 0, len(contextMessages))

	for _, msg := range contextMessages {
		aiMessages = append(aiMessages, AIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return aiMessages
}

// GetConversationContext retrieves a conversation context
func (e *Engine) GetConversationContext(conversationID string) (*ConversationContext, error) {
	ctx, exists := e.contextMgr.Get(conversationID)
	if !exists {
		return nil, fmt.Errorf("conversation %s not found", conversationID)
	}
	return ctx, nil
}

// ClearConversationContext clears a conversation context
func (e *Engine) ClearConversationContext(conversationID string) error {
	ctx, exists := e.contextMgr.Get(conversationID)
	if !exists {
		return fmt.Errorf("conversation %s not found", conversationID)
	}

	ctx.Clear()
	return nil
}

// DeleteConversationContext deletes a conversation context
func (e *Engine) DeleteConversationContext(conversationID string) {
	e.contextMgr.Delete(conversationID)
}

// GetActiveConversationCount returns the number of active conversations
func (e *Engine) GetActiveConversationCount() int {
	return e.contextMgr.Count()
}

// CleanupOldContexts removes contexts older than the specified duration
func (e *Engine) CleanupOldContexts(maxAge time.Duration) int {
	return e.contextMgr.Cleanup(maxAge)
}

// UpdateConfig updates the agent configuration
func (e *Engine) UpdateConfig(config *AgentConfig) {
	e.config = config
}

// GetConfig returns the current agent configuration
func (e *Engine) GetConfig() *AgentConfig {
	return e.config
}
