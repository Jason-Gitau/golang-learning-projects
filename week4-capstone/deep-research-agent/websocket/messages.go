package websocket

import "time"

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Research progress message types
	MessageTypeResearchStarted    MessageType = "research_started"
	MessageTypeStepStarted        MessageType = "step_started"
	MessageTypeStepCompleted      MessageType = "step_completed"
	MessageTypeProgress           MessageType = "progress"
	MessageTypeToolExecution      MessageType = "tool_execution"
	MessageTypeResearchCompleted  MessageType = "research_completed"
	MessageTypeError              MessageType = "error"
	MessageTypePing               MessageType = "ping"
	MessageTypePong               MessageType = "pong"
	MessageTypeConnectionConfirm  MessageType = "connection_confirm"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType            `json:"type"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	JobID     string                 `json:"job_id,omitempty"`
}

// NewMessage creates a new message with current timestamp
func NewMessage(msgType MessageType, jobID string) *Message {
	return &Message{
		Type:      msgType,
		JobID:     jobID,
		Timestamp: time.Now(),
		Data:      make(map[string]interface{}),
	}
}

// NewResearchStartedMessage creates a research started message
func NewResearchStartedMessage(jobID, query string, totalSteps int) *Message {
	msg := NewMessage(MessageTypeResearchStarted, jobID)
	msg.Data = map[string]interface{}{
		"query":       query,
		"total_steps": totalSteps,
		"status":      "started",
	}
	return msg
}

// NewStepStartedMessage creates a step started message
func NewStepStartedMessage(jobID string, stepNumber int, stepDescription, tool string) *Message {
	msg := NewMessage(MessageTypeStepStarted, jobID)
	msg.Data = map[string]interface{}{
		"step":        stepNumber,
		"description": stepDescription,
		"tool":        tool,
		"status":      "running",
	}
	return msg
}

// NewStepCompletedMessage creates a step completed message
func NewStepCompletedMessage(jobID string, stepNumber int, success bool, duration string) *Message {
	msg := NewMessage(MessageTypeStepCompleted, jobID)
	msg.Data = map[string]interface{}{
		"step":     stepNumber,
		"success":  success,
		"duration": duration,
		"status":   "completed",
	}
	return msg
}

// NewProgressMessage creates a progress update message
func NewProgressMessage(jobID string, current, total int, percentage float64, message string) *Message {
	msg := NewMessage(MessageTypeProgress, jobID)
	msg.Data = map[string]interface{}{
		"current":    current,
		"total":      total,
		"percentage": percentage,
		"message":    message,
	}
	return msg
}

// NewToolExecutionMessage creates a tool execution message
func NewToolExecutionMessage(jobID, tool, status string, params map[string]interface{}) *Message {
	msg := NewMessage(MessageTypeToolExecution, jobID)
	msg.Data = map[string]interface{}{
		"tool":   tool,
		"status": status,
		"params": params,
	}
	return msg
}

// NewResearchCompletedMessage creates a research completed message
func NewResearchCompletedMessage(jobID string, success bool, duration string, result interface{}) *Message {
	msg := NewMessage(MessageTypeResearchCompleted, jobID)
	msg.Data = map[string]interface{}{
		"success":  success,
		"duration": duration,
		"result":   result,
		"status":   "completed",
	}
	return msg
}

// NewErrorMessage creates an error message
func NewErrorMessage(jobID, errorMsg string) *Message {
	msg := NewMessage(MessageTypeError, jobID)
	msg.Data = map[string]interface{}{
		"error":   errorMsg,
		"status":  "error",
	}
	return msg
}

// NewPongMessage creates a pong message (response to ping)
func NewPongMessage() *Message {
	return NewMessage(MessageTypePong, "")
}

// NewConnectionConfirmMessage creates a connection confirmation message
func NewConnectionConfirmMessage(jobID string) *Message {
	msg := NewMessage(MessageTypeConnectionConfirm, jobID)
	msg.Data = map[string]interface{}{
		"status":  "connected",
		"message": "WebSocket connection established",
	}
	return msg
}
