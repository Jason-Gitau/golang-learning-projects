# Agent Orchestrator

A sophisticated concurrent agent orchestration system built in Go that demonstrates advanced concurrency patterns, channel-based message passing, and tool management.

## Overview

The Agent Orchestrator manages a pool of concurrent "agents" (goroutines) that can handle user requests independently. Each agent can execute various tools (weather API, calculator, time operations, etc.) and communicate through a message-passing system using Go channels.

## Key Features

- **Concurrent Agent Pool**: Configurable number of agents running as goroutines
- **Message-Based Communication**: Channel-based message passing between agents
- **Tool Registry Pattern**: Extensible tool system with easy registration
- **HTTP API**: RESTful API built with Gin framework
- **State Management**: Thread-safe agent state tracking (idle, busy, error)
- **Request Routing**: Smart routing of requests to available agents
- **Context & Timeouts**: Proper cancellation and timeout handling
- **Graceful Shutdown**: Clean shutdown with context cancellation
- **Real-time Statistics**: Monitor agent performance and system health

## Architecture

```
┌─────────────┐
│  HTTP API   │
│  (Gin)      │
└──────┬──────┘
       │
       ▼
┌─────────────┐      ┌──────────────┐
│   Router    │◄────►│ Tool Registry│
│  (Channels) │      └──────────────┘
└──────┬──────┘
       │
       ├───────┬───────┬───────┐
       ▼       ▼       ▼       ▼
   ┌──────┐┌──────┐┌──────┐┌──────┐
   │Agent1││Agent2││Agent3││Agent4│
   └──┬───┘└──┬───┘└──┬───┘└──┬───┘
      │       │       │       │
      └───────┴───────┴───────┘
              │
         ┌────▼─────┐
         │  State   │
         │ Manager  │
         └──────────┘
```

## Project Structure

```
02-agent-orchestrator/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── agent/
│   ├── agent.go           # Agent implementation
│   ├── manager.go         # Agent pool manager
│   └── state.go           # State management
├── tools/
│   ├── registry.go        # Tool registry
│   ├── calculator.go      # Math operations
│   ├── time.go            # Time/date operations
│   ├── random.go          # Random generation
│   ├── weather.go         # Weather API
│   └── text.go            # Text manipulation
├── messaging/
│   ├── message.go         # Message types
│   └── router.go          # Message routing
├── api/
│   └── server.go          # HTTP API server
├── models/
│   └── models.go          # Data models
├── config/
│   └── config.go          # Configuration
├── README.md              # This file
├── QUICK_START.md         # Quick start guide
└── START_HERE.md          # Learning guide
```

## Available Tools

### 1. Calculator
Performs mathematical operations:
- **Operations**: add, subtract, multiply, divide
- **Parameters**: `operation`, `a`, `b`

### 2. Time
Provides time and date information:
- **Actions**: current, timezone, format, add_duration, unix
- **Parameters**: `action`, `timezone`, `format`, `duration`

### 3. Random
Generates random data:
- **Types**: int, float, bool, string, uuid, choice, dice
- **Parameters**: `type`, `min`, `max`, `length`, `choices`, `sides`, `count`

### 4. Weather
Fetches weather information:
- **Parameters**: `location`, `units` (metric/imperial)
- **Note**: Configure `WEATHER_API_KEY` for real data

### 5. Text
Manipulates text strings:
- **Operations**: uppercase, lowercase, title, reverse, count, trim, replace, split, contains, word_count
- **Parameters**: `operation`, `text`, `find`, `replace_with`, `delimiter`, `search`

## API Endpoints

### Core Endpoints

- `GET /` - Service information
- `GET /health` - Health check
- `POST /api/v1/request` - Submit a request
- `GET /api/v1/agents` - List all agents
- `GET /api/v1/agents/:id` - Get agent details
- `GET /api/v1/tools` - List available tools
- `GET /api/v1/stats` - Get system statistics

## Configuration

Default configuration:
- **Agents**: 5 concurrent agents
- **Queue Size**: 100 requests
- **Port**: 8080
- **Timeout**: 30 seconds

Environment variables:
- `PORT` - API server port
- `WEATHER_API_KEY` - OpenWeatherMap API key (optional)

## Go Concurrency Concepts Demonstrated

1. **Goroutines**: Each agent runs in its own goroutine
2. **Channels**: Message passing between components
3. **Select Statements**: Multiplexing channel operations
4. **Context**: Cancellation and timeout propagation
5. **Mutex**: Thread-safe state management
6. **WaitGroup**: Coordinating goroutine shutdown
7. **Buffered Channels**: Request queuing
8. **Channel Closing**: Graceful shutdown signaling

## Example Usage

### Submit a Calculator Request
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "calculator",
    "params": {
      "operation": "add",
      "a": 10,
      "b": 5
    }
  }'
```

### Get Current Time
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "time",
    "params": {
      "action": "current"
    }
  }'
```

### Generate Random Number
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "random",
    "params": {
      "type": "int",
      "min": 1,
      "max": 100
    }
  }'
```

### Check Agent Status
```bash
curl http://localhost:8080/api/v1/agents
```

### View Statistics
```bash
curl http://localhost:8080/api/v1/stats
```

## Building and Running

```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# Build binary
go build -o agent-orchestrator

# Run binary
./agent-orchestrator
```

## Extending the System

### Adding a New Tool

1. Create a new file in `tools/` directory:

```go
package tools

import "context"

type MyTool struct{}

func NewMyTool() *MyTool {
    return &MyTool{}
}

func (t *MyTool) Name() string {
    return "mytool"
}

func (t *MyTool) Description() string {
    return "Description of what my tool does"
}

func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Implementation
    return result, nil
}
```

2. Register in `main.go`:

```go
if err := registry.Register(tools.NewMyTool()); err != nil {
    return err
}
```

## Performance Considerations

- **Agent Pool Size**: Adjust based on workload and CPU cores
- **Queue Size**: Larger queues handle bursts better but use more memory
- **Timeouts**: Set appropriate timeouts for different tool types
- **Mutex Contention**: State manager uses RWMutex for better read performance

## Error Handling

The system handles various error scenarios:
- Tool not found
- Invalid parameters
- Request timeout
- Agent unavailability
- Graceful shutdown during processing

## Monitoring and Observability

- Real-time agent state tracking
- Request success/failure statistics
- Average processing time
- Agent utilization metrics
- Health check endpoint for load balancers

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./agent
go test ./tools
```

## Best Practices Demonstrated

1. **Separation of Concerns**: Clear package boundaries
2. **Interface-Based Design**: Tool interface for extensibility
3. **Thread Safety**: Proper mutex usage
4. **Resource Management**: Proper cleanup and shutdown
5. **Context Propagation**: Timeout and cancellation handling
6. **Error Handling**: Comprehensive error checking
7. **Logging**: Structured logging throughout

## Common Issues and Solutions

### High CPU Usage
- Reduce number of agents
- Increase request timeout
- Check for infinite loops in tools

### Memory Growth
- Review queue size
- Check for goroutine leaks
- Monitor agent cleanup

### Slow Response Times
- Increase agent pool size
- Optimize tool execution
- Check for blocking operations

## Contributing

To add new features:
1. Follow existing code structure
2. Implement proper error handling
3. Add appropriate tests
4. Update documentation

## License

MIT License - See LICENSE file for details

## Further Reading

- [QUICK_START.md](QUICK_START.md) - Get started quickly
- [START_HERE.md](START_HERE.md) - Detailed learning guide
- Go Concurrency Patterns: https://go.dev/blog/pipelines
- Effective Go: https://go.dev/doc/effective_go

## Acknowledgments

Built as part of a Go learning path focusing on:
- Concurrent programming
- Channel-based communication
- System architecture
- Production-ready patterns
