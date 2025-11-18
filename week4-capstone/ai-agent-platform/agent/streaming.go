package agent

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
)

// StreamHandler manages the streaming of agent responses
type StreamHandler struct {
	conversationID string
	messageID      string
	broadcaster    ResponseBroadcaster
}

// ResponseBroadcaster is an interface for broadcasting responses
type ResponseBroadcaster interface {
	BroadcastChunk(conversationID, messageID, content string) error
	BroadcastToolCall(conversationID, tool string, params map[string]interface{}) error
	BroadcastToolResult(conversationID, tool string, result interface{}) error
	BroadcastComplete(conversationID, messageID string) error
	BroadcastError(conversationID, errorMsg string) error
}

// NewStreamHandler creates a new stream handler
func NewStreamHandler(conversationID string, broadcaster ResponseBroadcaster) *StreamHandler {
	return &StreamHandler{
		conversationID: conversationID,
		messageID:      uuid.New().String(),
		broadcaster:    broadcaster,
	}
}

// StreamResponse streams a response token by token
func (h *StreamHandler) StreamResponse(ctx context.Context, responseChan <-chan StreamChunk) (string, []ToolCallRecord, error) {
	var fullResponse strings.Builder
	var toolCalls []ToolCallRecord

	for {
		select {
		case <-ctx.Done():
			return fullResponse.String(), toolCalls, ctx.Err()

		case chunk, ok := <-responseChan:
			if !ok {
				// Channel closed
				return fullResponse.String(), toolCalls, nil
			}

			// Handle errors
			if chunk.Error != nil {
				log.Printf("Stream error for conversation %s: %v", h.conversationID, chunk.Error)
				h.broadcaster.BroadcastError(h.conversationID, chunk.Error.Error())
				return fullResponse.String(), toolCalls, chunk.Error
			}

			// Handle tool calls
			if chunk.ToolCall != nil {
				log.Printf("Tool call in conversation %s: %s", h.conversationID, chunk.ToolCall.Tool)

				err := h.broadcaster.BroadcastToolCall(
					h.conversationID,
					chunk.ToolCall.Tool,
					chunk.ToolCall.Params,
				)
				if err != nil {
					log.Printf("Failed to broadcast tool call: %v", err)
				}

				// Record tool call
				toolCalls = append(toolCalls, ToolCallRecord{
					Tool:   chunk.ToolCall.Tool,
					Params: chunk.ToolCall.Params,
				})

				continue
			}

			// Handle tool results
			if chunk.ToolResult != nil {
				log.Printf("Tool result in conversation %s: %s", h.conversationID, chunk.ToolResult.Tool)

				err := h.broadcaster.BroadcastToolResult(
					h.conversationID,
					chunk.ToolResult.Tool,
					chunk.ToolResult.Result,
				)
				if err != nil {
					log.Printf("Failed to broadcast tool result: %v", err)
				}

				// Update tool call record with result
				for i := range toolCalls {
					if toolCalls[i].Tool == chunk.ToolResult.Tool {
						toolCalls[i].Result = chunk.ToolResult.Result
						if chunk.ToolResult.Error != "" {
							toolCalls[i].Error = chunk.ToolResult.Error
						}
						break
					}
				}

				continue
			}

			// Handle content chunks
			if chunk.Content != "" {
				fullResponse.WriteString(chunk.Content)

				err := h.broadcaster.BroadcastChunk(
					h.conversationID,
					h.messageID,
					chunk.Content,
				)
				if err != nil {
					log.Printf("Failed to broadcast chunk: %v", err)
				}
			}

			// Handle completion
			if chunk.Done {
				log.Printf("Stream complete for conversation %s, message %s", h.conversationID, h.messageID)

				err := h.broadcaster.BroadcastComplete(h.conversationID, h.messageID)
				if err != nil {
					log.Printf("Failed to broadcast completion: %v", err)
				}

				return fullResponse.String(), toolCalls, nil
			}
		}
	}
}

// GetMessageID returns the message ID for this stream
func (h *StreamHandler) GetMessageID() string {
	return h.messageID
}

// StreamChunk represents a chunk of streamed content
type StreamChunk struct {
	Content    string
	ToolCall   *ToolCallInfo
	ToolResult *ToolResultInfo
	Done       bool
	Error      error
}

// ToolCallInfo contains information about a tool call
type ToolCallInfo struct {
	Tool   string
	Params map[string]interface{}
}

// ToolResultInfo contains information about a tool result
type ToolResultInfo struct {
	Tool   string
	Result interface{}
	Error  string
}
