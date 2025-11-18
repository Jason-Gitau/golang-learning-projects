# Quick Start Guide - Real-time & AI Engine

This guide will get you up and running with the Real-time & AI Engine in 5 minutes.

## Prerequisites

- Go 1.21 or higher
- Terminal/Command prompt
- Web browser

## Step 1: Navigate to Project

```bash
cd /home/user/golang-learning-projects/week4-capstone/ai-agent-platform
```

## Step 2: Install Dependencies

```bash
go mod download
```

## Step 3: Start the Server

```bash
go run main.go
```

You should see:
```
Starting AI Agent API Platform...
Registered 4 tools
AI service initialized (mock mode)
Agent engine initialized
WebSocket hub started
Agent runner started with 10 workers
Server starting on port 8080
```

## Step 4: Test with Web Browser

Open your browser and go to:
```
http://localhost:8080/
```

You'll see the WebSocket test page with:
- Connection status indicator
- Message input box
- Send button
- Connect/Disconnect controls

## Step 5: Try Some Messages

The system will auto-connect. Try these example messages:

### 1. Simple Greeting
```
Hello!
```

**Expected:** AI responds with greeting and explains capabilities

### 2. Calculator Tool
```
What is 25 + 17?
```

**Expected:**
- Tool call notification: `calculator(operation: "add", a: 25, b: 17)`
- Tool result: `42`
- AI response incorporating the result

### 3. Weather Tool
```
What's the weather in San Francisco?
```

**Expected:**
- Tool call: `weather(location: "San Francisco")`
- Mock weather data
- AI description of weather

### 4. Search Tool
```
Search for golang concurrency
```

**Expected:**
- Tool call: `search(query: "golang concurrency")`
- Mock search results
- AI summary

### 5. DateTime Tool
```
What time is it?
```

**Expected:**
- Tool call: `datetime(operation: "current")`
- Current date/time information
- Formatted response

### 6. General Question
```
What is Go programming language?
```

**Expected:**
- No tool call (knowledge-based)
- Streaming response about Go
- No tool execution

## Understanding the UI

### Connection Status
- **Green "Connected":** WebSocket is active
- **Red "Disconnected":** WebSocket is closed

### Message Types

#### User Messages (Blue, Right-aligned)
Your input messages

#### Agent Messages (Green, Left-aligned)
AI responses, streamed word-by-word

#### Tool Messages (Orange, Small text)
- Tool calls with parameters
- Tool results

#### Errors (Red)
Any errors that occur

## Advanced Testing with wscat

If you have Node.js installed:

```bash
# Install wscat
npm install -g wscat

# Connect
wscat -c "ws://localhost:8080/ws/chat/test-conversation?user_id=test-user"

# Send a message (paste and press Enter)
{"type":"user_message","content":"What is 5 + 3?"}

# You'll see streaming responses
{"type":"typing_indicator","content":"Agent is typing..."}
{"type":"tool_call","tool":"calculator","params":{"operation":"add","a":5,"b":3}}
{"type":"tool_result","tool":"calculator","result":{"operation":"add","result":8}}
{"type":"agent_response_chunk","content":"I'll"}
{"type":"agent_response_chunk","content":" use"}
...
{"type":"agent_response_complete"}
```

## Testing Different Scenarios

### Multiple Conversations

Open multiple browser tabs to `http://localhost:8080/` - each will create a separate conversation with its own context.

### Connection Resilience

1. Send a message
2. Click "Disconnect"
3. Click "Connect"
4. Send another message - context is maintained

### Tool Detection Patterns

Try these phrases to trigger different tools:

**Calculator:**
- "calculate 10 * 5"
- "what is 100 divided by 4"
- "15 + 23"

**Weather:**
- "weather in London"
- "temperature in Tokyo"
- "how's the weather"

**Search:**
- "search for REST APIs"
- "look up AI agents"
- "find information about WebSockets"

**DateTime:**
- "current time"
- "what's today's date"
- "what time is it"

## Health Checks

Check if the server is running:

```bash
# Overall health
curl http://localhost:8080/health

# WebSocket health
curl http://localhost:8080/health/ws
```

## Monitoring Server Logs

Watch the terminal where you ran `go run main.go` to see:

```
[WebSocket] Client registered: abc-123 (user: test-user, conversation: test-conversation)
[Agent] Worker 3 processing message for conversation test-conversation
[AI] Tool call detected: calculator
[Tools] Tool executed: calculator (add, 5+3)
[Agent] Message processed for conversation test-conversation: 2 messages in context
[WebSocket] Client unregistered: abc-123
```

## Stopping the Server

Press `Ctrl+C` in the terminal. You'll see graceful shutdown:

```
Shutting down server...
Shutting down agent worker pool...
All agent workers stopped gracefully
Shutting down WebSocket hub...
WebSocket hub shutdown complete
Server stopped gracefully
```

## Troubleshooting

### Port Already in Use

If you see `bind: address already in use`:

```bash
# Find and kill the process
lsof -ti:8080 | xargs kill -9

# Or use a different port
PORT=8081 go run main.go
```

### Dependencies Not Found

```bash
# Clean and reinstall
go clean -modcache
go mod download
```

### WebSocket Won't Connect

1. Check server is running: `curl http://localhost:8080/health`
2. Check browser console for errors (F12)
3. Try different browser
4. Check firewall settings

## Next Steps

Once you've tested the basic functionality:

1. **Read the Documentation:**
   - `docs/WEBSOCKET_API.md` - Complete WebSocket API
   - `docs/TOOLS_GUIDE.md` - Tool development guide
   - `docs/REALTIME_ENGINE_README.md` - Architecture details

2. **Create a Custom Tool:**
   - Follow the guide in `docs/TOOLS_GUIDE.md`
   - Add your tool to `tools/` directory
   - Register it in `tools/registry.go`

3. **Integrate Real AI:**
   - Add OpenAI or Anthropic API key
   - Replace mock AI service
   - See integration guide in README

4. **Scale the System:**
   - Increase worker pool size
   - Add Redis for pub/sub
   - Deploy to production

## Example Complete Session

```
1. Open http://localhost:8080/
2. See "Connected" status
3. Type: "Hello!"
4. See: AI greeting with capability description
5. Type: "What is 50 * 4?"
6. See: Tool call ’ calculator ’ result: 200 ’ AI response
7. Type: "Weather in Paris"
8. See: Tool call ’ weather ’ mock data ’ AI description
9. Type: "What is Go?"
10. See: Knowledge-based response about Go (no tool)
11. Click "Disconnect"
12. Click "Connect"
13. Type: "Tell me more about Go"
14. See: Response using conversation context
```

## Quick Reference

| Action | Command/URL |
|--------|-------------|
| Start Server | `go run main.go` |
| Test Page | `http://localhost:8080/` |
| Health Check | `curl http://localhost:8080/health` |
| WebSocket URL | `ws://localhost:8080/ws/chat/{conv_id}?user_id={user_id}` |
| Stop Server | `Ctrl+C` |

## Support

For issues or questions:
- Check `REALTIME_ENGINE_SUMMARY.md` for overview
- Read detailed docs in `docs/` directory
- Review code comments in source files

Happy testing!
