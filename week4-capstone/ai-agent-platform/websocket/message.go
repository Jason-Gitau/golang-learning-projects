package websocket

import "time"

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client -> Server
	MessageTypeUserMessage MessageType = "user_message"
	MessageTypePing        MessageType = "ping"

	// Server -> Client
	MessageTypeAgentResponseChunk    MessageType = "agent_response_chunk"
	MessageTypeAgentResponseComplete MessageType = "agent_response_complete"
	MessageTypeToolCall              MessageType = "tool_call"
	MessageTypeToolResult            MessageType = "tool_result"
	MessageTypeTypingIndicator       MessageType = "typing_indicator"
	MessageTypePong                  MessageType = "pong"
	MessageTypeError                 MessageType = "error"
)

// Message represents a WebSocket message
type Message struct {
	Type           MessageType            `json:"type"`
	Content        string                 `json:"content,omitempty"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	MessageID      string                 `json:"message_id,omitempty"`
	Tool           string                 `json:"tool,omitempty"`
	Params         map[string]interface{} `json:"params,omitempty"`
	Result         interface{}            `json:"result,omitempty"`
	Error          string                 `json:"error,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// NewMessage creates a new message with current timestamp
func NewMessage(msgType MessageType, content string) *Message {
	return &Message{
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewErrorMessage creates an error message
func NewErrorMessage(errMsg string) *Message {
	return &Message{
		Type:      MessageTypeError,
		Error:     errMsg,
		Timestamp: time.Now(),
	}
}

// NewToolCallMessage creates a tool call message
func NewToolCallMessage(tool string, params map[string]interface{}) *Message {
	return &Message{
		Type:      MessageTypeToolCall,
		Tool:      tool,
		Params:    params,
		Timestamp: time.Now(),
	}
}

// NewToolResultMessage creates a tool result message
func NewToolResultMessage(tool string, result interface{}) *Message {
	return &Message{
		Type:      MessageTypeToolResult,
		Tool:      tool,
		Result:    result,
		Timestamp: time.Now(),
	}
}
