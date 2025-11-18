# Quick Start Guide - Agent Orchestrator

Get the Agent Orchestrator up and running in 5 minutes.

## Prerequisites

- Go 1.21 or higher
- curl (for testing)
- Optional: OpenWeatherMap API key for real weather data

## Installation

```bash
# Navigate to project directory
cd /home/user/golang-learning-projects/week3-projects/02-agent-orchestrator

# Download dependencies
go mod download

# Run the application
go run main.go
```

You should see:
```
=== Agent Orchestrator ===
Starting up...
Registered 5 tools: [calculator time random weather text]
Starting Agent Manager with 5 agents
Agent agent-1: Started and ready to process requests
Agent agent-2: Started and ready to process requests
...
Starting HTTP API server on :8080
```

## Quick Test

### 1. Check if the server is running

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "healthy": true,
  "time": "2024-01-15T10:30:00Z"
}
```

### 2. View available tools

```bash
curl http://localhost:8080/api/v1/tools
```

### 3. Submit your first request

**Calculator:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "calculator",
    "params": {
      "operation": "multiply",
      "a": 7,
      "b": 6
    }
  }'
```

Response:
```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "agent_id": "agent-1",
  "success": true,
  "result": {
    "operation": "multiply",
    "a": 7,
    "b": 6,
    "result": 42
  },
  "processed_at": "2024-01-15T10:31:00Z",
  "duration": 150000
}
```

### 4. Check agent status

```bash
curl http://localhost:8080/api/v1/agents
```

### 5. View system statistics

```bash
curl http://localhost:8080/api/v1/stats
```

## Common Requests

### Time Operations

**Current time:**
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

**Timezone conversion:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "time",
    "params": {
      "action": "timezone",
      "timezone": "America/New_York"
    }
  }'
```

**Add duration:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "time",
    "params": {
      "action": "add_duration",
      "duration": "2h30m"
    }
  }'
```

### Random Generation

**Random integer:**
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

**Roll dice:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "random",
    "params": {
      "type": "dice",
      "sides": 6,
      "count": 3
    }
  }'
```

**Random UUID:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "random",
    "params": {
      "type": "uuid"
    }
  }'
```

**Random choice:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "random",
    "params": {
      "type": "choice",
      "choices": ["red", "green", "blue", "yellow"]
    }
  }'
```

### Text Manipulation

**Convert to uppercase:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "text",
    "params": {
      "operation": "uppercase",
      "text": "hello world"
    }
  }'
```

**Reverse text:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "text",
    "params": {
      "operation": "reverse",
      "text": "golang"
    }
  }'
```

**Word count:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "text",
    "params": {
      "operation": "word_count",
      "text": "The quick brown fox jumps over the lazy dog"
    }
  }'
```

**Replace text:**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "text",
    "params": {
      "operation": "replace",
      "text": "Hello World, Hello Universe",
      "find": "Hello",
      "replace_with": "Hi"
    }
  }'
```

### Weather

**Get weather (mock data by default):**
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "weather",
    "params": {
      "location": "London",
      "units": "metric"
    }
  }'
```

**Get weather (real data - requires API key):**
```bash
# Set environment variable
export WEATHER_API_KEY="your_api_key"

# Restart the application
go run main.go

# Make request
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "weather",
    "params": {
      "location": "Paris",
      "units": "metric"
    }
  }'
```

## Testing Concurrency

Run multiple requests simultaneously to see agents working concurrently:

```bash
# Terminal 1
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/request \
    -H "Content-Type: application/json" \
    -d "{\"tool_name\":\"calculator\",\"params\":{\"operation\":\"add\",\"a\":$i,\"b\":$i}}" &
done
wait

# Check statistics
curl http://localhost:8080/api/v1/stats
```

## Configuration

### Change number of agents

Edit `config/config.go` or modify `main.go`:

```go
cfg := config.Default().WithAgents(10) // 10 agents instead of 5
```

### Change port

```bash
PORT=3000 go run main.go
```

Or use environment variable:
```bash
export PORT=3000
go run main.go
```

### Adjust timeout

Modify request timeout:
```bash
curl -X POST http://localhost:8080/api/v1/request \
  -H "Content-Type: application/json" \
  -d '{
    "tool_name": "calculator",
    "params": {
      "operation": "add",
      "a": 10,
      "b": 5
    },
    "timeout": 60
  }'
```

## Monitoring

### View all agents
```bash
curl http://localhost:8080/api/v1/agents | jq
```

### View specific agent
```bash
curl http://localhost:8080/api/v1/agents/agent-1 | jq
```

### View statistics
```bash
curl http://localhost:8080/api/v1/stats | jq
```

### Watch statistics in real-time
```bash
watch -n 1 'curl -s http://localhost:8080/api/v1/stats | jq'
```

## Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o agent-orchestrator

# Run binary
./agent-orchestrator

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o agent-orchestrator-linux
GOOS=darwin GOARCH=amd64 go build -o agent-orchestrator-mac
GOOS=windows GOARCH=amd64 go build -o agent-orchestrator.exe
```

## Troubleshooting

### Port already in use
```bash
# Use different port
PORT=8081 go run main.go
```

### Dependencies not found
```bash
go mod tidy
go mod download
```

### Application not responding
```bash
# Check if running
curl http://localhost:8080/health

# Check logs for errors
# Restart the application
```

## Next Steps

1. Read [START_HERE.md](START_HERE.md) for in-depth learning guide
2. Read [README.md](README.md) for complete documentation
3. Try creating your own custom tool
4. Experiment with different agent pool sizes
5. Add monitoring and metrics
6. Implement request prioritization

## Stopping the Application

Press `Ctrl+C` to gracefully shut down:
```
Shutdown signal received
Stopping HTTP API server...
Stopping Agent Manager...
Agent agent-1: Shutting down
Agent agent-2: Shutting down
...
Agent Manager stopped
Agent Orchestrator shut down successfully
```

## Getting Help

- Check the [README.md](README.md) for detailed information
- Review [START_HERE.md](START_HERE.md) for learning concepts
- Examine the code comments for implementation details
- Test each tool individually to understand behavior

## Summary

You now have a running agent orchestrator! You can:
- Submit requests to various tools
- Monitor agent status
- View system statistics
- Test concurrent operations
- Extend with custom tools

Happy orchestrating!
