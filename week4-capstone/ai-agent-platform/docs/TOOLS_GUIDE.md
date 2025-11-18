# AI Agent Tools Development Guide

## Overview

The AI Agent Platform provides a flexible, extensible tool system that allows AI agents to perform actions beyond text generation. Tools enable agents to:

- Perform calculations
- Fetch real-time data
- Execute searches
- Interact with external APIs
- Access system information

## Architecture

### Tool Interface

All tools implement the `Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() map[string]interface{}
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}
```

### Components

1. **Tool Interface** - Defines the contract for all tools
2. **Tool Registry** - Manages tool registration and execution
3. **Tool Implementations** - Individual tool implementations
4. **AI Integration** - Tool calling detection and execution

## Built-in Tools

### 1. Calculator Tool

Performs mathematical operations.

**Operations:**
- `add` - Addition
- `subtract` - Subtraction
- `multiply` - Multiplication
- `divide` - Division
- `power` - Exponentiation
- `sqrt` - Square root
- `mod` - Modulo

**Example Usage:**

```json
{
  "tool": "calculator",
  "params": {
    "operation": "add",
    "a": 25,
    "b": 17
  }
}
```

**Response:**

```json
{
  "tool": "calculator",
  "success": true,
  "result": {
    "operation": "add",
    "result": 42,
    "operands": {"a": 25, "b": 17}
  }
}
```

**Natural Language Triggers:**
- "calculate 5 + 3"
- "what is 10 * 7"
- "25 divided by 5"

### 2. Weather Tool

Retrieves weather information for a location (mock data).

**Parameters:**
- `location` (string, required) - City or location name
- `units` (string, optional) - "celsius" or "fahrenheit" (default: "celsius")

**Example Usage:**

```json
{
  "tool": "weather",
  "params": {
    "location": "San Francisco",
    "units": "celsius"
  }
}
```

**Response:**

```json
{
  "tool": "weather",
  "success": true,
  "result": {
    "location": "San Francisco",
    "temperature": 18.5,
    "units": "celsius",
    "condition": "Partly Cloudy",
    "humidity": 65,
    "wind_speed": 15,
    "wind_unit": "km/h",
    "precipitation_chance": 20,
    "description": "Partly Cloudy with a temperature of 18.5°C"
  }
}
```

**Natural Language Triggers:**
- "what's the weather in London"
- "weather in Paris"
- "temperature in Tokyo"

### 3. Search Tool

Performs web searches and returns relevant results (mock data).

**Parameters:**
- `query` (string, required) - Search query
- `max_results` (number, optional) - Maximum results (default: 5)

**Example Usage:**

```json
{
  "tool": "search",
  "params": {
    "query": "golang concurrency",
    "max_results": 3
  }
}
```

**Response:**

```json
{
  "tool": "search",
  "success": true,
  "result": {
    "query": "golang concurrency",
    "total_results": 3,
    "results": [
      {
        "title": "Concurrency in Go",
        "url": "https://example.com/go-concurrency",
        "snippet": "Learn about goroutines and channels...",
        "relevance": 0.95
      }
    ]
  }
}
```

**Natural Language Triggers:**
- "search for golang tutorials"
- "look up AI agents"
- "find information about REST APIs"

### 4. DateTime Tool

Provides date and time operations.

**Operations:**
- `current` - Get current date/time
- `timezone` - Convert timezone
- `format` - Format date/time
- `parse` - Parse date string
- `add` - Add duration
- `diff` - Calculate difference

**Example Usage:**

```json
{
  "tool": "datetime",
  "params": {
    "operation": "current",
    "timezone": "UTC"
  }
}
```

**Response:**

```json
{
  "tool": "datetime",
  "success": true,
  "result": {
    "timestamp": 1704110400,
    "iso8601": "2024-01-01T12:00:00Z",
    "readable": "Monday, January 1, 2024 at 12:00:00 PM UTC",
    "timezone": "UTC",
    "date": "2024-01-01",
    "time": "12:00:00",
    "day_of_week": "Monday",
    "day_of_year": 1,
    "week_of_year": 1
  }
}
```

**Natural Language Triggers:**
- "what time is it"
- "current time"
- "what's today's date"

## Creating Custom Tools

### Step 1: Implement the Tool Interface

```go
package tools

import (
    "context"
    "fmt"
)

type MyCustomTool struct{}

func NewMyCustomTool() *MyCustomTool {
    return &MyCustomTool{}
}

func (t *MyCustomTool) Name() string {
    return "my_custom_tool"
}

func (t *MyCustomTool) Description() string {
    return "A custom tool that does something useful"
}

func (t *MyCustomTool) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "param1": map[string]interface{}{
                "type":        "string",
                "description": "First parameter",
            },
            "param2": map[string]interface{}{
                "type":        "number",
                "description": "Second parameter",
            },
        },
        "required": []string{"param1"},
    }
}

func (t *MyCustomTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Extract parameters
    param1, ok := params["param1"].(string)
    if !ok {
        return nil, fmt.Errorf("param1 must be a string")
    }

    // Perform tool logic
    result := fmt.Sprintf("Processed: %s", param1)

    // Return result
    return map[string]interface{}{
        "success": true,
        "output": result,
    }, nil
}
```

### Step 2: Register the Tool

```go
// In main.go or initialization code
toolsRegistry := tools.NewRegistry(10 * time.Second)
toolsRegistry.Register(tools.NewMyCustomTool())
```

### Step 3: Test the Tool

```go
func TestMyCustomTool(t *testing.T) {
    tool := NewMyCustomTool()

    params := map[string]interface{}{
        "param1": "test value",
        "param2": 42,
    }

    ctx := context.Background()
    result, err := tool.Execute(ctx, params)

    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }

    // Verify result
    // ...
}
```

## Best Practices

### 1. Input Validation

Always validate input parameters:

```go
func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Validate required parameters
    value, ok := params["required_param"]
    if !ok {
        return nil, fmt.Errorf("required_param is required")
    }

    // Validate type
    strValue, ok := value.(string)
    if !ok {
        return nil, fmt.Errorf("required_param must be a string")
    }

    // Validate range/format
    if len(strValue) == 0 {
        return nil, fmt.Errorf("required_param cannot be empty")
    }

    // ... rest of logic
}
```

### 2. Context Handling

Respect context cancellation and timeouts:

```go
func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Check context before expensive operations
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // For long operations, check periodically
    for i := 0; i < 100; i++ {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            // Do work
        }
    }

    return result, nil
}
```

### 3. Error Handling

Provide clear, actionable error messages:

```go
func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    value, ok := params["url"].(string)
    if !ok {
        return nil, fmt.Errorf("url parameter must be a string, got %T", params["url"])
    }

    resp, err := http.Get(value)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch URL %s: %w", value, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("URL returned status %d: %s", resp.StatusCode, resp.Status)
    }

    // ... rest of logic
}
```

### 4. Structured Results

Return consistent, structured results:

```go
type ToolResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // ... perform operation

    return ToolResponse{
        Success: true,
        Data: resultData,
        Metadata: map[string]interface{}{
            "execution_time_ms": executionTime,
            "version": "1.0",
        },
    }, nil
}
```

### 5. Documentation

Use clear JSON schema for parameters:

```go
func (t *MyTool) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "query": map[string]interface{}{
                "type":        "string",
                "description": "The search query to execute",
                "minLength":   1,
                "maxLength":   200,
            },
            "limit": map[string]interface{}{
                "type":        "integer",
                "description": "Maximum number of results to return",
                "minimum":     1,
                "maximum":     100,
                "default":     10,
            },
        },
        "required": []string{"query"},
    }
}
```

## Advanced Features

### Rate Limiting

Implement per-tool rate limiting:

```go
type RateLimitedTool struct {
    baseTool Tool
    limiter  *rate.Limiter
}

func (t *RateLimitedTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    if err := t.limiter.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limit exceeded: %w", err)
    }

    return t.baseTool.Execute(ctx, params)
}
```

### Caching

Cache tool results:

```go
type CachedTool struct {
    baseTool Tool
    cache    map[string]interface{}
    mu       sync.RWMutex
}

func (t *CachedTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    cacheKey := fmt.Sprintf("%v", params)

    // Check cache
    t.mu.RLock()
    if result, ok := t.cache[cacheKey]; ok {
        t.mu.RUnlock()
        return result, nil
    }
    t.mu.RUnlock()

    // Execute tool
    result, err := t.baseTool.Execute(ctx, params)
    if err != nil {
        return nil, err
    }

    // Store in cache
    t.mu.Lock()
    t.cache[cacheKey] = result
    t.mu.Unlock()

    return result, nil
}
```

### Async Execution

Execute tools asynchronously:

```go
type AsyncResult struct {
    Result interface{}
    Error  error
}

func ExecuteAsync(ctx context.Context, tool Tool, params map[string]interface{}) <-chan AsyncResult {
    resultChan := make(chan AsyncResult, 1)

    go func() {
        defer close(resultChan)

        result, err := tool.Execute(ctx, params)
        resultChan <- AsyncResult{
            Result: result,
            Error:  err,
        }
    }()

    return resultChan
}
```

## Tool Registry API

### Register a Tool

```go
registry := tools.NewRegistry(10 * time.Second)
err := registry.Register(myTool)
```

### Execute a Tool

```go
result, err := registry.Execute(ctx, "my_tool", params)
```

### List Available Tools

```go
toolList := registry.List()
for _, tool := range toolList {
    fmt.Printf("%s: %s\n", tool.Name(), tool.Description())
}
```

### Get Tool Descriptions

```go
descriptions := registry.GetToolDescriptions()
// Returns: map[string]string{"calculator": "Performs math...", ...}
```

### Unregister a Tool

```go
registry.Unregister("my_tool")
```

## Integration with AI Agents

### Tool Detection

The AI service automatically detects when tools should be called based on user input:

```go
// In ai/mock.go
func (s *MockService) detectToolCall(input string) *ToolCall {
    lowerInput := strings.ToLower(input)

    // Calculator pattern
    if strings.Contains(lowerInput, "calculate") {
        return &ToolCall{
            Name: "calculator",
            Params: extractParams(input),
        }
    }

    // ... other patterns
}
```

### Execution Flow

1. User sends message: "What is 5 + 3?"
2. AI detects calculator tool needed
3. AI sends tool_call message
4. Engine executes tool
5. Engine sends tool_result message
6. AI incorporates result in response
7. AI streams final answer to user

## Testing Tools

### Unit Tests

```go
func TestCalculatorTool(t *testing.T) {
    calc := NewCalculatorTool()

    tests := []struct {
        name    string
        params  map[string]interface{}
        want    float64
        wantErr bool
    }{
        {
            name: "addition",
            params: map[string]interface{}{
                "operation": "add",
                "a": 5.0,
                "b": 3.0,
            },
            want: 8.0,
        },
        {
            name: "division by zero",
            params: map[string]interface{}{
                "operation": "divide",
                "a": 5.0,
                "b": 0.0,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := calc.Execute(context.Background(), tt.params)

            if tt.wantErr {
                if err == nil {
                    t.Error("Expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Errorf("Unexpected error: %v", err)
                return
            }

            // Verify result
            // ...
        })
    }
}
```

### Integration Tests

```go
func TestToolRegistry(t *testing.T) {
    registry := tools.NewRegistry(5 * time.Second)

    // Register tools
    registry.Register(tools.NewCalculatorTool())

    // Execute tool
    result, err := registry.Execute(
        context.Background(),
        "calculator",
        map[string]interface{}{
            "operation": "add",
            "a": 10,
            "b": 5,
        },
    )

    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    // Verify result
    // ...
}
```

## Security Considerations

### 1. Input Sanitization

Always sanitize user input:

```go
func sanitizeInput(input string) string {
    // Remove dangerous characters
    // Validate format
    // Limit length
    return cleaned
}
```

### 2. Resource Limits

Limit resource usage:

```go
const (
    maxExecutionTime = 10 * time.Second
    maxMemoryUsage = 100 * 1024 * 1024 // 100MB
)
```

### 3. Sandboxing

Execute tools in isolated environments when possible.

### 4. Audit Logging

Log all tool executions:

```go
func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    log.Printf("Tool %s executed with params: %v", t.Name(), params)

    result, err := t.performOperation(ctx, params)

    if err != nil {
        log.Printf("Tool %s failed: %v", t.Name(), err)
    }

    return result, err
}
```

## Performance Optimization

### 1. Connection Pooling

Reuse connections for external APIs:

```go
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 100,
        MaxIdleConnsPerHost: 10,
    },
}
```

### 2. Parallel Execution

Execute independent tools in parallel:

```go
type ToolExecution struct {
    Tool string
    Result interface{}
    Error error
}

func ExecuteParallel(ctx context.Context, registry *Registry, executions []ToolRequest) []ToolExecution {
    results := make([]ToolExecution, len(executions))
    var wg sync.WaitGroup

    for i, exec := range executions {
        wg.Add(1)
        go func(idx int, req ToolRequest) {
            defer wg.Done()

            result, err := registry.Execute(ctx, req.Tool, req.Params)
            results[idx] = ToolExecution{
                Tool: req.Tool,
                Result: result.Result,
                Error: err,
            }
        }(i, exec)
    }

    wg.Wait()
    return results
}
```

### 3. Result Streaming

Stream large results:

```go
type StreamingTool interface {
    Tool
    ExecuteStream(ctx context.Context, params map[string]interface{}) <-chan StreamChunk
}
```

## Troubleshooting

### Tool Not Found

**Error:** `tool my_tool not found`

**Solution:** Ensure tool is registered before use

### Timeout Errors

**Error:** `context deadline exceeded`

**Solution:** Increase timeout in registry or optimize tool implementation

### Parameter Errors

**Error:** `param must be a string, got float64`

**Solution:** Ensure correct parameter types in tool calls

## Future Enhancements

- **Dynamic Tool Loading:** Load tools from plugins
- **Tool Marketplace:** Share and discover tools
- **Version Management:** Support multiple tool versions
- **Tool Composition:** Chain tools together
- **Visual Tool Builder:** GUI for creating tools

## Resources

- [Tool Interface Documentation](./tools/interface.go)
- [Example Tools](./tools/)
- [Tool Registry](./tools/registry.go)
- [AI Integration](./ai/mock.go)

## Contributing

To contribute a new tool:

1. Implement the Tool interface
2. Add comprehensive tests
3. Document parameters using JSON schema
4. Update this guide with usage examples
5. Submit a pull request

For questions or issues, please open a GitHub issue.
