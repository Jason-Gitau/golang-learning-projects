package ai

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// MockService provides a mock AI service with streaming responses
type MockService struct {
	streamDelay time.Duration
}

// NewMockService creates a new mock AI service
func NewMockService(streamDelay time.Duration) *MockService {
	if streamDelay == 0 {
		streamDelay = 30 * time.Millisecond // Default: 30ms per token
	}

	return &MockService{
		streamDelay: streamDelay,
	}
}

// Message represents a conversation message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	Content    string
	Done       bool
	ToolCall   *ToolCall
	Error      error
}

// ToolCall represents a tool invocation request
type ToolCall struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

// GenerateResponse generates an AI response based on the input
func (s *MockService) GenerateResponse(ctx context.Context, messages []Message, systemPrompt string) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("no messages provided")
	}

	lastMessage := messages[len(messages)-1].Content
	response := s.generateResponseText(lastMessage, systemPrompt)

	return response, nil
}

// StreamResponse generates a streaming AI response
func (s *MockService) StreamResponse(ctx context.Context, messages []Message, systemPrompt string) <-chan StreamResponse {
	responseChan := make(chan StreamResponse, 10)

	go func() {
		defer close(responseChan)

		if len(messages) == 0 {
			responseChan <- StreamResponse{
				Error: fmt.Errorf("no messages provided"),
				Done:  true,
			}
			return
		}

		lastMessage := messages[len(messages)-1].Content

		// Check if we need to call a tool
		toolCall := s.detectToolCall(lastMessage)
		if toolCall != nil {
			// Send tool call notification
			responseChan <- StreamResponse{
				ToolCall: toolCall,
				Done:     false,
			}

			// Wait a bit to simulate processing
			time.Sleep(s.streamDelay * 10)

			// Generate response mentioning tool usage
			response := fmt.Sprintf("I'll use the %s tool to help with that. ", toolCall.Name)
			s.streamText(ctx, response, responseChan)

			// Mark as done
			responseChan <- StreamResponse{Done: true}
			return
		}

		// Generate normal response
		response := s.generateResponseText(lastMessage, systemPrompt)
		s.streamText(ctx, response, responseChan)

		// Send completion signal
		responseChan <- StreamResponse{Done: true}
	}()

	return responseChan
}

// streamText sends text token by token
func (s *MockService) streamText(ctx context.Context, text string, responseChan chan<- StreamResponse) {
	words := strings.Fields(text)

	for i, word := range words {
		select {
		case <-ctx.Done():
			return
		default:
			// Add space between words (except for the first word)
			content := word
			if i > 0 {
				content = " " + word
			}

			responseChan <- StreamResponse{
				Content: content,
				Done:    false,
			}

			// Simulate streaming delay
			time.Sleep(s.streamDelay)
		}
	}
}

// generateResponseText creates a response based on input
func (s *MockService) generateResponseText(input, systemPrompt string) string {
	lowerInput := strings.ToLower(input)

	// Go-related queries
	if strings.Contains(lowerInput, "go") || strings.Contains(lowerInput, "golang") {
		return "Go is a statically typed, compiled programming language designed at Google. It's known for its simplicity, excellent concurrency support through goroutines and channels, and fast compilation times. Go is particularly well-suited for building web servers, command-line tools, and distributed systems."
	}

	// Concurrency queries
	if strings.Contains(lowerInput, "concurrency") || strings.Contains(lowerInput, "goroutine") {
		return "Concurrency in Go is built around goroutines and channels. Goroutines are lightweight threads managed by the Go runtime, and channels provide a way for goroutines to communicate and synchronize. This makes concurrent programming in Go both powerful and relatively simple compared to traditional threading models."
	}

	// API queries
	if strings.Contains(lowerInput, "api") || strings.Contains(lowerInput, "rest") {
		return "A REST API (Representational State Transfer Application Programming Interface) is an architectural style for building web services. It uses HTTP methods like GET, POST, PUT, and DELETE to perform CRUD operations. REST APIs are stateless, scalable, and widely used for building modern web applications."
	}

	// WebSocket queries
	if strings.Contains(lowerInput, "websocket") {
		return "WebSockets provide full-duplex communication channels over a single TCP connection. Unlike HTTP, which follows a request-response pattern, WebSockets allow both the server and client to send messages at any time. This makes them ideal for real-time applications like chat, gaming, and live updates."
	}

	// AI/Agent queries
	if strings.Contains(lowerInput, "ai agent") || strings.Contains(lowerInput, "agent") {
		return "An AI agent is a software program that can perceive its environment, make decisions, and take actions to achieve specific goals. Modern AI agents can use tools, maintain conversation context, and adapt their behavior based on feedback. They're becoming increasingly important in automation and human-computer interaction."
	}

	// Greeting
	if strings.Contains(lowerInput, "hello") || strings.Contains(lowerInput, "hi") {
		return "Hello! I'm an AI assistant. I can help you with questions about Go programming, REST APIs, concurrency, WebSockets, and much more. I can also use various tools like calculator, weather lookup, web search, and datetime operations. How can I assist you today?"
	}

	// Help request
	if strings.Contains(lowerInput, "help") || strings.Contains(lowerInput, "what can you do") {
		return "I can help you with many things! I can answer questions about programming, technology, and general knowledge. I also have access to several tools: calculator for math operations, weather for location-based weather info, search for web queries, and datetime for time-related operations. Just ask me anything!"
	}

	// Default response
	return fmt.Sprintf("I understand you're asking about: '%s'. Based on the information available, I can provide some insights. However, for more specific details, you might want to rephrase your question or use one of my available tools (calculator, weather, search, datetime) for more precise results.", input)
}

// detectToolCall analyzes input to determine if a tool should be called
func (s *MockService) detectToolCall(input string) *ToolCall {
	lowerInput := strings.ToLower(input)

	// Calculator patterns
	calcPatterns := []string{
		`calculate (.+)`,
		`what is (\d+)\s*\+\s*(\d+)`,
		`what is (\d+)\s*-\s*(\d+)`,
		`what is (\d+)\s*\*\s*(\d+)`,
		`what is (\d+)\s*/\s*(\d+)`,
		`(\d+)\s*\+\s*(\d+)`,
		`(\d+)\s*\*\s*(\d+)`,
	}

	for _, pattern := range calcPatterns {
		if matched, _ := regexp.MatchString(pattern, lowerInput); matched {
			// Extract numbers for simple operations
			if nums := extractNumbers(input); len(nums) >= 2 {
				return &ToolCall{
					Name: "calculator",
					Params: map[string]interface{}{
						"operation": detectOperation(input),
						"a":         nums[0],
						"b":         nums[1],
					},
				}
			}
		}
	}

	// Weather patterns
	weatherPatterns := []string{
		`weather in (.+)`,
		`what's the weather in (.+)`,
		`how's the weather`,
		`temperature in (.+)`,
	}

	for _, pattern := range weatherPatterns {
		if matched, _ := regexp.MatchString(pattern, lowerInput); matched {
			location := extractLocation(input)
			if location == "" {
				location = "San Francisco"
			}
			return &ToolCall{
				Name: "weather",
				Params: map[string]interface{}{
					"location": location,
					"units":    "celsius",
				},
			}
		}
	}

	// Search patterns
	searchPatterns := []string{
		`search for (.+)`,
		`look up (.+)`,
		`find information about (.+)`,
		`google (.+)`,
	}

	for _, pattern := range searchPatterns {
		if matched, _ := regexp.MatchString(pattern, lowerInput); matched {
			query := extractSearchQuery(input)
			return &ToolCall{
				Name: "search",
				Params: map[string]interface{}{
					"query":       query,
					"max_results": 5,
				},
			}
		}
	}

	// DateTime patterns
	datetimePatterns := []string{
		`what time is it`,
		`current time`,
		`what's the date`,
		`today's date`,
		`time in (.+)`,
	}

	for _, pattern := range datetimePatterns {
		if matched, _ := regexp.MatchString(pattern, lowerInput); matched {
			timezone := extractTimezone(input)
			if timezone == "" {
				timezone = "UTC"
			}
			return &ToolCall{
				Name: "datetime",
				Params: map[string]interface{}{
					"operation": "current",
					"timezone":  timezone,
				},
			}
		}
	}

	return nil
}

// Helper functions for tool detection

func extractNumbers(input string) []float64 {
	re := regexp.MustCompile(`\d+\.?\d*`)
	matches := re.FindAllString(input, -1)

	numbers := make([]float64, 0)
	for _, match := range matches {
		var num float64
		fmt.Sscanf(match, "%f", &num)
		numbers = append(numbers, num)
	}

	return numbers
}

func detectOperation(input string) string {
	lowerInput := strings.ToLower(input)

	if strings.Contains(lowerInput, "+") || strings.Contains(lowerInput, "plus") || strings.Contains(lowerInput, "add") {
		return "add"
	}
	if strings.Contains(lowerInput, "-") || strings.Contains(lowerInput, "minus") || strings.Contains(lowerInput, "subtract") {
		return "subtract"
	}
	if strings.Contains(lowerInput, "*") || strings.Contains(lowerInput, "times") || strings.Contains(lowerInput, "multiply") {
		return "multiply"
	}
	if strings.Contains(lowerInput, "/") || strings.Contains(lowerInput, "divide") {
		return "divide"
	}

	return "add" // default
}

func extractLocation(input string) string {
	patterns := []string{
		`weather in (.+?)(?:\?|$)`,
		`temperature in (.+?)(?:\?|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(strings.ToLower(input)); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

func extractSearchQuery(input string) string {
	patterns := []string{
		`search for (.+?)(?:\?|$)`,
		`look up (.+?)(?:\?|$)`,
		`find information about (.+?)(?:\?|$)`,
		`google (.+?)(?:\?|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(strings.ToLower(input)); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// If no pattern matches, return the whole input
	return strings.TrimSpace(input)
}

func extractTimezone(input string) string {
	patterns := []string{
		`time in (.+?)(?:\?|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(strings.ToLower(input)); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}
