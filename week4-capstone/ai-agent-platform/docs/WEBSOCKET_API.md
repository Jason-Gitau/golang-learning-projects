# WebSocket API Documentation

## Overview

The AI Agent Platform provides a real-time WebSocket API for interactive chat with AI agents. The WebSocket connection enables bidirectional streaming communication, allowing for:

- Real-time message exchange
- Streaming AI responses (token-by-token)
- Tool execution notifications
- Connection health monitoring
- Multiple concurrent conversations

## Connection

### Endpoint

```
ws://localhost:8080/ws/chat/{conversation_id}
```

### Authentication

Authentication is performed via user ID passed as a query parameter or header:

**Query Parameter (Development/Testing):**
```
ws://localhost:8080/ws/chat/conv-123?user_id=user-456
```

**Header (Production):**
```
X-User-ID: user-456
```

In production, JWT tokens should be used via the Authorization header and validated by middleware.

### Connection Flow

1. Client initiates WebSocket upgrade request
2. Server validates authentication
3. Connection upgraded to WebSocket
4. Client registered in hub
5. Confirmation message sent to client
6. Connection ready for bidirectional communication

## Message Format

All messages are JSON-encoded with the following base structure:

```json
{
  "type": "message_type",
  "content": "message content",
  "conversation_id": "conv-123",
  "message_id": "msg-456",
  "timestamp": "2024-01-01T12:00:00Z",
  "metadata": {}
}
```

## Message Types

### Client ’ Server

#### User Message

Send a user message to the AI agent.

```json
{
  "type": "user_message",
  "content": "What is Go programming language?",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Fields:**
- `type`: Always `"user_message"`
- `content`: The user's message text (required)
- `timestamp`: ISO 8601 timestamp (auto-set by server if omitted)

#### Ping

Check connection health.

```json
{
  "type": "ping",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Response:** Server responds with `pong` message.

### Server ’ Client

#### Agent Response Chunk

Streaming chunk of AI agent response.

```json
{
  "type": "agent_response_chunk",
  "content": "Go is a statically ",
  "conversation_id": "conv-123",
  "message_id": "msg-456",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Fields:**
- `type`: Always `"agent_response_chunk"`
- `content`: A chunk of the response (word or phrase)
- `conversation_id`: ID of the conversation
- `message_id`: Unique ID for this message
- `timestamp`: When the chunk was generated

#### Agent Response Complete

Indicates the AI response is complete.

```json
{
  "type": "agent_response_complete",
  "conversation_id": "conv-123",
  "message_id": "msg-456",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### Tool Call

Agent is executing a tool.

```json
{
  "type": "tool_call",
  "tool": "calculator",
  "params": {
    "operation": "add",
    "a": 5,
    "b": 3
  },
  "conversation_id": "conv-123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Fields:**
- `tool`: Name of the tool being called
- `params`: Parameters passed to the tool

#### Tool Result

Result from tool execution.

```json
{
  "type": "tool_result",
  "tool": "calculator",
  "result": {
    "operation": "add",
    "result": 8,
    "operands": {"a": 5, "b": 3}
  },
  "conversation_id": "conv-123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Fields:**
- `tool`: Name of the tool that was executed
- `result`: The tool's return value

#### Typing Indicator

Agent is processing a response.

```json
{
  "type": "typing_indicator",
  "content": "Agent is typing...",
  "conversation_id": "conv-123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### Pong

Response to ping message.

```json
{
  "type": "pong",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### Error

Error occurred during processing.

```json
{
  "type": "error",
  "error": "Failed to process message: timeout",
  "conversation_id": "conv-123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Conversation Flow Example

### Simple Chat

```
Client ’ Server: user_message("Hello")
Server ’ Client: typing_indicator
Server ’ Client: agent_response_chunk("Hello!")
Server ’ Client: agent_response_chunk(" I'm")
Server ’ Client: agent_response_chunk(" an")
Server ’ Client: agent_response_chunk(" AI")
Server ’ Client: agent_response_chunk(" assistant.")
Server ’ Client: agent_response_complete
```

### Chat with Tool Usage

```
Client ’ Server: user_message("What is 25 + 17?")
Server ’ Client: typing_indicator
Server ’ Client: tool_call("calculator", {operation: "add", a: 25, b: 17})
Server ’ Client: tool_result("calculator", {result: 42})
Server ’ Client: agent_response_chunk("I'll")
Server ’ Client: agent_response_chunk(" use")
Server ’ Client: agent_response_chunk(" the")
Server ’ Client: agent_response_chunk(" calculator")
Server ’ Client: agent_response_chunk(" tool.")
Server ’ Client: agent_response_chunk(" The")
Server ’ Client: agent_response_chunk(" answer")
Server ’ Client: agent_response_chunk(" is")
Server ’ Client: agent_response_chunk(" 42.")
Server ’ Client: agent_response_complete
```

## Connection Management

### Heartbeat

The server sends periodic ping messages to check connection health. Clients should respond or will be disconnected after timeout.

**Server Configuration:**
- Ping interval: 54 seconds (9/10 of pong wait)
- Pong wait: 60 seconds
- Read deadline extended on each pong received

### Graceful Disconnection

Clients should send a close frame before disconnecting:

```javascript
websocket.close(1000, "Client closing connection");
```

### Reconnection

If disconnected, clients should implement exponential backoff reconnection:

```javascript
let reconnectDelay = 1000; // Start with 1 second

function reconnect() {
  setTimeout(() => {
    connect();
    reconnectDelay = Math.min(reconnectDelay * 2, 30000); // Max 30 seconds
  }, reconnectDelay);
}
```

## Error Handling

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Missing or invalid user_id | Include valid user_id in query param or header |
| 400 Bad Request | Missing conversation_id | Ensure conversation_id is in URL path |
| 1000 Normal Closure | Client closed connection | Normal behavior, reconnect if needed |
| 1001 Going Away | Server shutting down | Reconnect after brief delay |
| 1006 Abnormal Closure | Network issue | Reconnect with exponential backoff |

### Error Message Format

```json
{
  "type": "error",
  "error": "Error description",
  "conversation_id": "conv-123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Best Practices

### Client Implementation

1. **Handle All Message Types:** Implement handlers for all message types
2. **Buffer Management:** Handle rapid streaming chunks efficiently
3. **Reconnection Logic:** Implement exponential backoff
4. **Message Queuing:** Queue messages if connection temporarily down
5. **UI Updates:** Update UI smoothly for streaming responses

### Performance

1. **Connection Pooling:** Reuse connections when possible
2. **Message Batching:** Batch UI updates for streaming chunks
3. **Memory Management:** Clear old messages from memory
4. **Compression:** Enable WebSocket compression for large messages

### Security

1. **Authentication:** Always validate user identity
2. **Rate Limiting:** Implement per-user rate limits
3. **Input Validation:** Sanitize all user input
4. **CORS:** Configure appropriate CORS policies
5. **TLS:** Use WSS (WebSocket Secure) in production

## Example Client (JavaScript)

```javascript
class AIAgentClient {
  constructor(conversationId, userId) {
    this.conversationId = conversationId;
    this.userId = userId;
    this.ws = null;
    this.listeners = {};
  }

  connect() {
    const url = `ws://localhost:8080/ws/chat/${this.conversationId}?user_id=${this.userId}`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log('Connected');
      this.emit('connected');
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.emit('error', error);
    };

    this.ws.onclose = () => {
      console.log('Disconnected');
      this.emit('disconnected');
    };
  }

  handleMessage(message) {
    switch (message.type) {
      case 'agent_response_chunk':
        this.emit('chunk', message.content);
        break;
      case 'agent_response_complete':
        this.emit('complete', message);
        break;
      case 'tool_call':
        this.emit('tool_call', message);
        break;
      case 'tool_result':
        this.emit('tool_result', message);
        break;
      case 'error':
        this.emit('error', message.error);
        break;
    }
  }

  sendMessage(content) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket not connected');
    }

    const message = {
      type: 'user_message',
      content: content,
      timestamp: new Date().toISOString()
    };

    this.ws.send(JSON.stringify(message));
  }

  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = [];
    }
    this.listeners[event].push(callback);
  }

  emit(event, data) {
    if (this.listeners[event]) {
      this.listeners[event].forEach(cb => cb(data));
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}

// Usage
const client = new AIAgentClient('conv-123', 'user-456');

client.on('chunk', (content) => {
  console.log('Received chunk:', content);
  // Update UI with streaming content
});

client.on('complete', (message) => {
  console.log('Response complete');
  // Finalize UI update
});

client.on('tool_call', (message) => {
  console.log('Tool called:', message.tool, message.params);
});

client.connect();
client.sendMessage('Hello, AI!');
```

## Testing

### Manual Testing

Use the built-in test page:
```
http://localhost:8080/
```

### WebSocket Testing Tools

- **websocat:** CLI tool for WebSocket testing
- **wscat:** Node.js WebSocket CLI client
- **Postman:** GUI tool with WebSocket support

### Example with wscat

```bash
# Install wscat
npm install -g wscat

# Connect
wscat -c "ws://localhost:8080/ws/chat/test-conv?user_id=test-user"

# Send message
{"type":"user_message","content":"Hello"}
```

## Monitoring

### Connection Metrics

Monitor these metrics for production:

- Active connections count
- Messages per second
- Average response time
- Tool execution count
- Error rate
- Reconnection rate

### Health Check

```bash
curl http://localhost:8080/health/ws
```

Response:
```json
{
  "status": "healthy",
  "service": "websocket"
}
```

## Troubleshooting

### Connection Refused

**Cause:** Server not running or wrong port

**Solution:**
```bash
# Check if server is running
curl http://localhost:8080/health

# Check port
netstat -an | grep 8080
```

### Messages Not Received

**Cause:** Client not properly handling message events

**Solution:** Add logging to onmessage handler and verify message format

### Frequent Disconnections

**Cause:** Network issues or ping/pong timeout

**Solution:**
- Implement reconnection logic
- Check network stability
- Verify ping/pong handling

### Slow Responses

**Cause:** High server load or slow tool execution

**Solution:**
- Check server metrics
- Increase worker pool size
- Optimize tool implementations

## Advanced Features

### Multiple Conversations

Open multiple WebSocket connections for different conversations:

```javascript
const conv1 = new AIAgentClient('conv-1', 'user-123');
const conv2 = new AIAgentClient('conv-2', 'user-123');

conv1.connect();
conv2.connect();
```

### Message History

Conversation context is maintained server-side. To clear:

```
POST /api/conversations/{conversation_id}/clear
```

### Custom Agents

Different agents can be used per conversation by setting agent_id in metadata:

```json
{
  "type": "user_message",
  "content": "Hello",
  "metadata": {
    "agent_id": "specialized-agent"
  }
}
```

## API Versioning

Current version: **v1**

Future versions will be namespaced:
```
ws://localhost:8080/v2/ws/chat/{conversation_id}
```

## Rate Limits

- **Connections per user:** 10 concurrent
- **Messages per minute:** 60
- **Message size:** 512KB max

Exceeding limits results in error message and possible temporary ban.

## Support

For issues or questions:
- GitHub Issues: [project-repo/issues]
- Documentation: [project-docs]
- Examples: [project-examples]
