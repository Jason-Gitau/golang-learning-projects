package ai

import (
	"context"
)

// AIServiceAdapter adapts MockService to the agent.AIService interface
type AIServiceAdapter struct {
	mock *MockService
}

// NewAIServiceAdapter creates a new AI service adapter
func NewAIServiceAdapter(mock *MockService) *AIServiceAdapter {
	return &AIServiceAdapter{mock: mock}
}

// AIMessage represents a message for the AI (matches agent.AIMessage)
type AdapterAIMessage struct {
	Role    string
	Content string
}

// AIStreamResponse represents a streaming response (matches agent.AIStreamResponse)
type AdapterAIStreamResponse struct {
	Content  string
	Done     bool
	ToolCall *AdapterAIToolCall
	Error    error
}

// AIToolCall represents a tool call (matches agent.AIToolCall)
type AdapterAIToolCall struct {
	Name   string
	Params map[string]interface{}
}

// StreamResponse streams an AI response
func (a *AIServiceAdapter) StreamResponse(ctx context.Context, messages []AdapterAIMessage, systemPrompt string) <-chan AdapterAIStreamResponse {
	responseChan := make(chan AdapterAIStreamResponse, 10)

	go func() {
		defer close(responseChan)

		// Convert to mock service messages
		mockMessages := make([]Message, len(messages))
		for i, msg := range messages {
			mockMessages[i] = Message{
				Role:    msg.Role,
				Content: msg.Content,
			}
		}

		// Get mock service stream
		mockStream := a.mock.StreamResponse(ctx, mockMessages, systemPrompt)

		// Convert and forward responses
		for mockResp := range mockStream {
			var toolCall *AdapterAIToolCall
			if mockResp.ToolCall != nil {
				toolCall = &AdapterAIToolCall{
					Name:   mockResp.ToolCall.Name,
					Params: mockResp.ToolCall.Params,
				}
			}

			responseChan <- AdapterAIStreamResponse{
				Content:  mockResp.Content,
				Done:     mockResp.Done,
				ToolCall: toolCall,
				Error:    mockResp.Error,
			}
		}
	}()

	return responseChan
}
