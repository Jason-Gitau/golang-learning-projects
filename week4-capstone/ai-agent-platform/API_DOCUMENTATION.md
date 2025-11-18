# API Documentation

Comprehensive API documentation for the AI Agent API Platform.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

### Success Response
```json
{
  "field1": "value1",
  "field2": "value2"
}
```

### Error Response
```json
{
  "error": "Error message",
  "details": "Optional additional details"
}
```

## HTTP Status Codes

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., duplicate email)
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

---

## Authentication Endpoints

### Register User

Create a new user account.

**Endpoint:** `POST /auth/register`

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"
}
```

**Validation:**
- `email`: Required, valid email format
- `password`: Required, minimum 6 characters
- `name`: Required

**Response:** `201 Created`
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "tier": "free",
    "created_at": "2025-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `400` - Invalid request body
- `409` - Email already exists

---

### Login

Authenticate a user and receive a JWT token.

**Endpoint:** `POST /auth/login`

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "tier": "free",
    "created_at": "2025-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid email or password

---

### Get Profile

Get the current authenticated user's profile.

**Endpoint:** `GET /auth/me`

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "tier": "free",
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Errors:**
- `401` - Not authenticated
- `404` - User not found

---

## Agent Endpoints

### Create Agent

Create a new AI agent.

**Endpoint:** `POST /agents`

**Authentication:** Required

**Request Body:**
```json
{
  "name": "Customer Support Agent",
  "description": "Handles customer inquiries and support tickets",
  "system_prompt": "You are a helpful customer support agent. Be polite and professional.",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 2000,
  "tools": [
    {
      "name": "search_kb",
      "description": "Search the knowledge base",
      "enabled": true,
      "config": {
        "index": "support_docs"
      }
    }
  ]
}
```

**Validation:**
- `name`: Required
- `system_prompt`: Required
- `model`: Required
- `temperature`: Optional, 0.0-2.0
- `max_tokens`: Optional, >= 0

**Response:** `201 Created`
```json
{
  "id": 1,
  "name": "Customer Support Agent",
  "description": "Handles customer inquiries and support tickets",
  "system_prompt": "You are a helpful customer support agent...",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 2000,
  "status": "active",
  "tools": [...],
  "user_id": 1,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Not authenticated

---

### List Agents

Get all agents for the authenticated user.

**Endpoint:** `GET /agents`

**Authentication:** Required

**Query Parameters:**
- `status` (optional): Filter by status (`active`, `inactive`, `archived`)

**Example:** `GET /agents?status=active`

**Response:** `200 OK`
```json
{
  "agents": [
    {
      "id": 1,
      "name": "Customer Support Agent",
      "description": "Handles customer inquiries",
      "system_prompt": "You are helpful...",
      "model": "gpt-4",
      "temperature": 0.7,
      "max_tokens": 2000,
      "status": "active",
      "tools": [...],
      "user_id": 1,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**Errors:**
- `401` - Not authenticated

---

### Get Agent

Get details of a specific agent.

**Endpoint:** `GET /agents/:id`

**Authentication:** Required

**Path Parameters:**
- `id`: Agent ID

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Customer Support Agent",
  "description": "Handles customer inquiries",
  "system_prompt": "You are helpful...",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 2000,
  "status": "active",
  "tools": [...],
  "user_id": 1,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Errors:**
- `400` - Invalid agent ID
- `401` - Not authenticated
- `404` - Agent not found

---

### Update Agent

Update an existing agent.

**Endpoint:** `PUT /agents/:id`

**Authentication:** Required

**Path Parameters:**
- `id`: Agent ID

**Request Body:** (all fields optional)
```json
{
  "name": "Updated Agent Name",
  "description": "Updated description",
  "system_prompt": "Updated prompt",
  "model": "gpt-4-turbo",
  "temperature": 0.8,
  "max_tokens": 3000,
  "status": "inactive",
  "tools": [...]
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Updated Agent Name",
  ...
}
```

**Errors:**
- `400` - Invalid request body or agent ID
- `401` - Not authenticated
- `404` - Agent not found

---

### Delete Agent

Delete an agent (soft delete).

**Endpoint:** `DELETE /agents/:id`

**Authentication:** Required

**Path Parameters:**
- `id`: Agent ID

**Response:** `200 OK`
```json
{
  "message": "Agent deleted successfully"
}
```

**Errors:**
- `400` - Invalid agent ID
- `401` - Not authenticated
- `404` - Agent not found

---

## Conversation Endpoints

### Create Conversation

Start a new conversation with an agent.

**Endpoint:** `POST /agents/:id/conversations`

**Authentication:** Required

**Path Parameters:**
- `id`: Agent ID

**Request Body:**
```json
{
  "title": "Product inquiry conversation"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "title": "Product inquiry conversation",
  "status": "active",
  "agent_id": 1,
  "user_id": 1,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Errors:**
- `400` - Invalid request body or agent ID
- `401` - Not authenticated
- `404` - Agent not found

---

### List Conversations

Get all conversations for the authenticated user.

**Endpoint:** `GET /conversations`

**Authentication:** Required

**Query Parameters:**
- `agent_id` (optional): Filter by agent
- `status` (optional): Filter by status (`active`, `archived`, `deleted`)

**Example:** `GET /conversations?agent_id=1&status=active`

**Response:** `200 OK`
```json
{
  "conversations": [
    {
      "id": 1,
      "title": "Product inquiry",
      "status": "active",
      "agent_id": 1,
      "user_id": 1,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:35:00Z"
    }
  ],
  "count": 1
}
```

**Errors:**
- `401` - Not authenticated

---

### Get Conversation

Get a specific conversation with all its messages.

**Endpoint:** `GET /conversations/:id`

**Authentication:** Required

**Path Parameters:**
- `id`: Conversation ID

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Product inquiry",
  "status": "active",
  "agent_id": 1,
  "user_id": 1,
  "messages": [
    {
      "id": 1,
      "role": "user",
      "content": "What features are available?",
      "tokens_used": 5,
      "conversation_id": 1,
      "created_at": "2025-01-15T10:31:00Z"
    },
    {
      "id": 2,
      "role": "assistant",
      "content": "We offer the following features...",
      "tokens_used": 15,
      "conversation_id": 1,
      "created_at": "2025-01-15T10:31:05Z"
    }
  ],
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:31:05Z"
}
```

**Errors:**
- `400` - Invalid conversation ID
- `401` - Not authenticated
- `404` - Conversation not found

---

### Delete Conversation

Delete a conversation (soft delete).

**Endpoint:** `DELETE /conversations/:id`

**Authentication:** Required

**Path Parameters:**
- `id`: Conversation ID

**Response:** `200 OK`
```json
{
  "message": "Conversation deleted successfully"
}
```

**Errors:**
- `400` - Invalid conversation ID
- `401` - Not authenticated
- `404` - Conversation not found

---

### Create Message

Add a message to a conversation.

**Endpoint:** `POST /conversations/:id/messages`

**Authentication:** Required

**Path Parameters:**
- `id`: Conversation ID

**Request Body:**
```json
{
  "role": "user",
  "content": "What are the pricing options?"
}
```

**Validation:**
- `role`: Required, one of: `user`, `assistant`, `system`
- `content`: Required

**Response:** `201 Created`
```json
{
  "id": 3,
  "role": "user",
  "content": "What are the pricing options?",
  "tokens_used": 0,
  "conversation_id": 1,
  "created_at": "2025-01-15T10:32:00Z"
}
```

**Errors:**
- `400` - Invalid request body or conversation ID
- `401` - Not authenticated
- `404` - Conversation not found

---

### Get Messages

Get all messages in a conversation.

**Endpoint:** `GET /conversations/:id/messages`

**Authentication:** Required

**Path Parameters:**
- `id`: Conversation ID

**Response:** `200 OK`
```json
{
  "messages": [
    {
      "id": 1,
      "role": "user",
      "content": "Hello",
      "tokens_used": 1,
      "conversation_id": 1,
      "created_at": "2025-01-15T10:31:00Z"
    }
  ],
  "count": 1
}
```

**Errors:**
- `400` - Invalid conversation ID
- `401` - Not authenticated
- `404` - Conversation not found

---

## Usage Endpoints

### Get Usage Statistics

Get aggregated usage statistics for the authenticated user.

**Endpoint:** `GET /usage`

**Authentication:** Required

**Query Parameters:**
- `start_date` (optional): Start date in YYYY-MM-DD format
- `end_date` (optional): End date in YYYY-MM-DD format
- `agent_id` (optional): Filter by specific agent

**Example:** `GET /usage?start_date=2025-01-01&end_date=2025-01-31`

**Response:** `200 OK`
```json
{
  "total_requests": 150,
  "total_tokens": 45000,
  "input_tokens": 20000,
  "output_tokens": 25000,
  "total_cost": 2.35,
  "period_start": "2025-01-01",
  "period_end": "2025-01-31",
  "by_model": {
    "gpt-4": {
      "requests": 100,
      "total_tokens": 30000,
      "input_tokens": 13000,
      "output_tokens": 17000,
      "cost": 1.80
    },
    "gpt-3.5-turbo": {
      "requests": 50,
      "total_tokens": 15000,
      "input_tokens": 7000,
      "output_tokens": 8000,
      "cost": 0.55
    }
  },
  "by_agent": {
    "1": {
      "agent_id": 1,
      "agent_name": "Customer Support Agent",
      "requests": 150,
      "total_tokens": 45000,
      "input_tokens": 20000,
      "output_tokens": 25000,
      "cost": 2.35
    }
  }
}
```

**Errors:**
- `400` - Invalid date format
- `401` - Not authenticated

---

### Get Rate Limit Info

Get current rate limit status for the authenticated user.

**Endpoint:** `GET /usage/rate-limit`

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "limit": 100,
  "remaining": 75,
  "reset": "2025-01-15T11:30:00Z"
}
```

**Errors:**
- `401` - Not authenticated

---

### Create Usage Log

Manually create a usage log entry (typically called by the chat system).

**Endpoint:** `POST /usage/log`

**Authentication:** Required

**Request Body:**
```json
{
  "input_tokens": 100,
  "output_tokens": 150,
  "model": "gpt-4",
  "operation": "chat",
  "agent_id": 1
}
```

**Validation:**
- `input_tokens`: Required, >= 0
- `output_tokens`: Required, >= 0
- `model`: Required
- `operation`: Required
- `agent_id`: Optional

**Response:** `201 Created`
```json
{
  "message": "Usage log created successfully",
  "log_id": 1
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Not authenticated

---

## Rate Limiting

All protected endpoints are subject to rate limiting. The default limit is 100 requests per hour per user.

### Rate Limit Headers

Every response includes rate limit information:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 75
X-RateLimit-Reset: 1705318200
```

### Rate Limit Exceeded Response

When the rate limit is exceeded, the API returns:

**Status:** `429 Too Many Requests`

```json
{
  "error": "Rate limit exceeded",
  "limit": 100,
  "remaining": 0,
  "reset": 1705318200
}
```

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "Human-readable error message",
  "details": "Optional additional details (validation errors, etc.)"
}
```

### Common Error Scenarios

**Missing Authorization:**
```json
{
  "error": "Authorization header required"
}
```

**Invalid Token:**
```json
{
  "error": "Invalid or expired token"
}
```

**Validation Error:**
```json
{
  "error": "Invalid request body",
  "details": "Key: 'CreateAgentRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

**Resource Not Found:**
```json
{
  "error": "Agent not found"
}
```

---

## Health Check

### Check API Health

**Endpoint:** `GET /health`

**Authentication:** Not required

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "time": "2025-01-15T10:30:00Z"
}
```

---

## Model Options

Supported AI models:
- `gpt-4`
- `gpt-4-turbo`
- `gpt-3.5-turbo`
- `claude-3-opus`
- `claude-3-sonnet`
- `claude-3-haiku`

Each model has different pricing (see `utils/tokens.go` for details).

---

## Best Practices

1. **Always handle rate limits** - Check the rate limit headers and implement exponential backoff
2. **Store JWT tokens securely** - Never expose tokens in client-side code
3. **Use HTTPS in production** - Never send tokens over unencrypted connections
4. **Validate input** - The API validates input, but client-side validation improves UX
5. **Monitor usage** - Regularly check usage statistics to avoid unexpected costs
6. **Handle errors gracefully** - Implement proper error handling for all API calls

---

## Postman Collection

Import this base configuration to test the API:

```json
{
  "info": {
    "name": "AI Agent API Platform",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api/v1"
    },
    {
      "key": "token",
      "value": ""
    }
  ]
}
```

Set the `token` variable after login to automatically include it in authenticated requests.
