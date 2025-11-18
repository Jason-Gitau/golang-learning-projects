# AI Agent API Platform - Core Platform

A production-ready REST API platform for managing AI agents, conversations, and tracking usage. This is the foundational backend component of the AI Agent API Platform capstone project.

## Features

### Core Functionality
- **User Authentication** - JWT-based authentication with bcrypt password hashing
- **Agent Management** - Full CRUD operations for AI agents with customizable configurations
- **Conversation Management** - Create and manage chat conversations with agents
- **Message Handling** - Store and retrieve conversation messages
- **Usage Tracking** - Track API usage, tokens, and costs per user and agent
- **Rate Limiting** - Sliding window rate limiting to prevent abuse

### Technical Highlights
- **RESTful API** - Clean, well-structured REST endpoints
- **Database** - SQLite with GORM ORM, proper migrations and indexes
- **Middleware Stack** - Authentication, rate limiting, logging, CORS, error handling
- **Configuration** - Environment-based configuration with sensible defaults
- **Production Ready** - Proper error handling, timeouts, and security measures

## Project Structure

```
ai-agent-platform/
├── main.go                  # Application entry point
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
│
├── config/
│   └── config.go           # Configuration management
│
├── models/
│   ├── user.go             # User model and DTOs
│   ├── agent.go            # Agent model and DTOs
│   ├── conversation.go     # Conversation model and DTOs
│   ├── message.go          # Message model and DTOs
│   └── usage.go            # Usage tracking models
│
├── database/
│   └── database.go         # Database initialization and migrations
│
├── handlers/
│   ├── auth.go             # Authentication endpoints
│   ├── agent.go            # Agent CRUD endpoints
│   ├── conversation.go     # Conversation management endpoints
│   └── usage.go            # Usage statistics endpoints
│
├── middleware/
│   ├── auth.go             # JWT authentication middleware
│   ├── ratelimit.go        # Rate limiting middleware
│   ├── logger.go           # Request logging middleware
│   └── error.go            # Error handling and CORS middleware
│
└── utils/
    ├── jwt.go              # JWT token utilities
    ├── password.go         # Password hashing utilities
    └── tokens.go           # Token counting and pricing utilities
```

## Quick Start

### Prerequisites
- Go 1.21 or higher
- SQLite (included with Go)

### Installation

1. Navigate to the project directory:
```bash
cd /home/user/golang-learning-projects/week4-capstone/ai-agent-platform
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

### Environment Variables

You can configure the application using environment variables:

```bash
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration
DATABASE_PATH=./ai-agent-platform.db

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# Rate Limiting Configuration
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1h
```

## API Endpoints

### Authentication

#### Register a new user
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Get current user profile
```bash
GET /api/v1/auth/me
Authorization: Bearer <token>
```

### Agent Management

#### Create an agent
```bash
POST /api/v1/agents
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Customer Support Agent",
  "description": "Handles customer inquiries",
  "system_prompt": "You are a helpful customer support agent.",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 2000,
  "tools": [
    {
      "name": "search",
      "description": "Search knowledge base",
      "enabled": true
    }
  ]
}
```

#### List all agents
```bash
GET /api/v1/agents
Authorization: Bearer <token>

# Optional query parameters
?status=active
```

#### Get specific agent
```bash
GET /api/v1/agents/:id
Authorization: Bearer <token>
```

#### Update an agent
```bash
PUT /api/v1/agents/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Agent Name",
  "status": "inactive"
}
```

#### Delete an agent
```bash
DELETE /api/v1/agents/:id
Authorization: Bearer <token>
```

### Conversation Management

#### Start a new conversation
```bash
POST /api/v1/agents/:id/conversations
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Help with product features"
}
```

#### List all conversations
```bash
GET /api/v1/conversations
Authorization: Bearer <token>

# Optional query parameters
?agent_id=1
?status=active
```

#### Get conversation with messages
```bash
GET /api/v1/conversations/:id
Authorization: Bearer <token>
```

#### Delete a conversation
```bash
DELETE /api/v1/conversations/:id
Authorization: Bearer <token>
```

#### Create a message
```bash
POST /api/v1/conversations/:id/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "role": "user",
  "content": "What features are available?"
}
```

#### Get all messages
```bash
GET /api/v1/conversations/:id/messages
Authorization: Bearer <token>
```

### Usage Tracking

#### Get usage statistics
```bash
GET /api/v1/usage
Authorization: Bearer <token>

# Optional query parameters
?start_date=2025-01-01
?end_date=2025-01-31
?agent_id=1
```

#### Get rate limit information
```bash
GET /api/v1/usage/rate-limit
Authorization: Bearer <token>
```

## Database Models

### User
- ID, Email, Password (hashed), Name, Tier
- Relationships: Agents, Conversations, UsageLogs

### Agent
- ID, Name, Description, SystemPrompt, Model, Temperature, MaxTokens, Status, Tools
- UserID (foreign key)
- Relationships: Conversations, UsageLogs

### Conversation
- ID, Title, Status
- AgentID, UserID (foreign keys)
- Relationships: Messages

### Message
- ID, Role, Content, TokensUsed
- ConversationID (foreign key)

### UsageLog
- ID, InputTokens, OutputTokens, TotalTokens, Cost, Model, Operation
- UserID, AgentID (foreign keys)

### RateLimit
- ID, RequestsCount, WindowStart, WindowEnd
- UserID (foreign key)

## Rate Limiting

The platform implements sliding window rate limiting:
- Default: 100 requests per hour per user
- Rate limit headers included in responses:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Unix timestamp when the window resets

When rate limit is exceeded, returns `429 Too Many Requests`.

## Usage Tracking

All agent interactions are logged with:
- Token usage (input, output, total)
- Cost calculation based on model pricing
- Aggregation by model and agent
- Time-based filtering

Token counting uses a simple estimation algorithm (can be replaced with tiktoken in production).

## Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT Authentication**: HS256 signing with configurable expiration
- **Rate Limiting**: Per-user request throttling
- **CORS**: Configured for cross-origin requests
- **Input Validation**: Request body validation with Gin binding
- **Error Handling**: Consistent error responses, no sensitive data leakage

## Testing the API

### Using curl

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.token')

# Create agent
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Agent","system_prompt":"You are helpful","model":"gpt-4"}'

# List agents
curl -X GET http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer $TOKEN"
```

## Development

### Running in Development Mode
```bash
ENVIRONMENT=development go run main.go
```

### Building for Production
```bash
go build -o ai-agent-platform main.go
```

### Running in Production
```bash
ENVIRONMENT=production \
JWT_SECRET=your-secure-secret \
PORT=8080 \
./ai-agent-platform
```

## Next Steps

This core platform provides the foundation for:
- **Real-time Chat** - WebSocket support for live agent interactions
- **AI Integration** - Connect to OpenAI, Anthropic, or other LLM providers
- **Advanced Features** - Tool calling, streaming responses, multi-agent orchestration

## Architecture Notes

This component focuses on:
- REST API endpoints
- Database models and migrations
- Authentication and authorization
- Rate limiting and usage tracking
- Core business logic

The real-time chat features (WebSocket, AI streaming) will be handled by a separate component that builds on this foundation.

## License

MIT License - Part of the Golang Learning Projects Week 4 Capstone
