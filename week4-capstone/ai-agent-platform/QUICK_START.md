# Quick Start Guide

Get up and running with the AI Agent API Platform in 5 minutes!

## Prerequisites

- Go 1.21+ installed
- Basic understanding of REST APIs
- curl or Postman for testing

## Step 1: Start the Server

```bash
cd /home/user/golang-learning-projects/week4-capstone/ai-agent-platform
go run main.go
```

You should see:
```
Starting AI Agent Platform API on port 8080
Environment: development
Database: ./ai-agent-platform.db
Rate Limit: 100 requests per hour
```

## Step 2: Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "password123",
    "name": "Demo User"
  }'
```

Response:
```json
{
  "user": {
    "id": 1,
    "email": "demo@example.com",
    "name": "Demo User",
    "tier": "free",
    "created_at": "2025-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Save the token - you'll need it for authenticated requests!

## Step 3: Create an AI Agent

Replace `<YOUR_TOKEN>` with the token from Step 2:

```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My First Agent",
    "description": "A helpful assistant",
    "system_prompt": "You are a helpful AI assistant. Be concise and friendly.",
    "model": "gpt-4",
    "temperature": 0.7,
    "max_tokens": 2000
  }'
```

Response:
```json
{
  "id": 1,
  "name": "My First Agent",
  "description": "A helpful assistant",
  "status": "active",
  ...
}
```

## Step 4: Start a Conversation

```bash
curl -X POST http://localhost:8080/api/v1/agents/1/conversations \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My first conversation"
  }'
```

Response:
```json
{
  "id": 1,
  "title": "My first conversation",
  "status": "active",
  "agent_id": 1,
  ...
}
```

## Step 5: Send a Message

```bash
curl -X POST http://localhost:8080/api/v1/conversations/1/messages \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "user",
    "content": "Hello, how are you?"
  }'
```

Response:
```json
{
  "id": 1,
  "role": "user",
  "content": "Hello, how are you?",
  "conversation_id": 1,
  ...
}
```

## Step 6: View Your Conversation

```bash
curl -X GET http://localhost:8080/api/v1/conversations/1 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

Response includes all messages in the conversation.

## Step 7: Check Usage Statistics

```bash
curl -X GET http://localhost:8080/api/v1/usage \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

See token usage, costs, and statistics aggregated by model and agent.

## Quick Test Script

Save this as `test.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

# Register
echo "1. Registering user..."
RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }')

TOKEN=$(echo $RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

# Create agent
echo -e "\n2. Creating agent..."
AGENT=$(curl -s -X POST $BASE_URL/agents \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Agent",
    "system_prompt": "You are helpful",
    "model": "gpt-4"
  }')

AGENT_ID=$(echo $AGENT | jq -r '.id')
echo "Agent ID: $AGENT_ID"

# Start conversation
echo -e "\n3. Starting conversation..."
CONV=$(curl -s -X POST $BASE_URL/agents/$AGENT_ID/conversations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Chat"}')

CONV_ID=$(echo $CONV | jq -r '.id')
echo "Conversation ID: $CONV_ID"

# Send message
echo -e "\n4. Sending message..."
curl -s -X POST $BASE_URL/conversations/$CONV_ID/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "user", "content": "Hello!"}' | jq

# List agents
echo -e "\n5. Listing agents..."
curl -s -X GET $BASE_URL/agents \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n✅ All tests completed!"
```

Make it executable and run:
```bash
chmod +x test.sh
./test.sh
```

## Common Issues

### Port already in use
```bash
# Change the port
PORT=3000 go run main.go
```

### Database locked
```bash
# Remove the database and restart
rm ai-agent-platform.db
go run main.go
```

### Rate limit exceeded
Wait for the window to reset or check your rate limit status:
```bash
curl -X GET http://localhost:8080/api/v1/usage/rate-limit \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

## Environment Configuration

Create a `.env` file (or set environment variables):

```bash
# Server
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_PATH=./ai-agent-platform.db

# JWT
JWT_SECRET=my-super-secret-key
JWT_EXPIRATION=24h

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1h
```

## Next Steps

1. **Explore the API** - Check out `API_DOCUMENTATION.md` for all endpoints
2. **Build a client** - Create a frontend or CLI to interact with the platform
3. **Add AI integration** - Connect to OpenAI or Anthropic APIs
4. **Implement WebSocket** - Add real-time chat capabilities
5. **Add monitoring** - Track performance and usage

## Helpful Commands

```bash
# Check health
curl http://localhost:8080/health

# Get your profile
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <TOKEN>"

# List all conversations
curl http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer <TOKEN>"

# Update an agent
curl -X PUT http://localhost:8080/api/v1/agents/1 \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"status": "inactive"}'

# Delete a conversation
curl -X DELETE http://localhost:8080/api/v1/conversations/1 \
  -H "Authorization: Bearer <TOKEN>"
```

## Support

For detailed API documentation, see `API_DOCUMENTATION.md`.

For architecture and design details, see `README.md`.

Happy coding! 🚀
