package models

import (
	"time"

	"github.com/google/uuid"
)

// AgentState represents the current state of an agent
type AgentState string

const (
	StateIdle  AgentState = "idle"
	StateBusy  AgentState = "busy"
	StateError AgentState = "error"
)

// Request represents a user request to be processed by an agent
type Request struct {
	ID        string                 `json:"id"`
	ToolName  string                 `json:"tool_name"`
	Params    map[string]interface{} `json:"params"`
	CreatedAt time.Time              `json:"created_at"`
	Timeout   time.Duration          `json:"timeout,omitempty"`
}

// Response represents the result of processing a request
type Response struct {
	RequestID  string      `json:"request_id"`
	AgentID    string      `json:"agent_id"`
	Success    bool        `json:"success"`
	Result     interface{} `json:"result,omitempty"`
	Error      string      `json:"error,omitempty"`
	ProcessedAt time.Time  `json:"processed_at"`
	Duration   time.Duration `json:"duration"`
}

// AgentInfo represents information about an agent
type AgentInfo struct {
	ID              string     `json:"id"`
	State           AgentState `json:"state"`
	RequestsHandled int        `json:"requests_handled"`
	LastActive      time.Time  `json:"last_active"`
	CurrentRequest  string     `json:"current_request,omitempty"`
}

// Statistics represents overall system statistics
type Statistics struct {
	TotalRequests     int64         `json:"total_requests"`
	SuccessfulRequests int64        `json:"successful_requests"`
	FailedRequests    int64         `json:"failed_requests"`
	AverageProcessTime time.Duration `json:"average_process_time"`
	ActiveAgents      int           `json:"active_agents"`
	IdleAgents        int           `json:"idle_agents"`
	BusyAgents        int           `json:"busy_agents"`
	ErrorAgents       int           `json:"error_agents"`
	Uptime            time.Duration `json:"uptime"`
}

// NewRequest creates a new request with a unique ID
func NewRequest(toolName string, params map[string]interface{}) *Request {
	return &Request{
		ID:        uuid.New().String(),
		ToolName:  toolName,
		Params:    params,
		CreatedAt: time.Now(),
		Timeout:   30 * time.Second, // Default timeout
	}
}

// NewResponse creates a new response
func NewResponse(requestID, agentID string, success bool, result interface{}, err error, duration time.Duration) *Response {
	resp := &Response{
		RequestID:  requestID,
		AgentID:    agentID,
		Success:    success,
		Result:     result,
		ProcessedAt: time.Now(),
		Duration:   duration,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	return resp
}
