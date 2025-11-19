# Authentication & Authorization Guide

## Overview

The Deep Research Agent API now includes comprehensive JWT-based authentication, user management, and multi-user session isolation. This guide covers all authentication features and how to use them.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Authentication Flow](#authentication-flow)
3. [API Endpoints](#api-endpoints)
4. [Rate Limiting](#rate-limiting)
5. [Security Best Practices](#security-best-practices)
6. [Examples](#examples)

---

## Architecture Overview

### Components

1. **JWT Authentication** (`auth/jwt.go`)
   - Token generation and validation
   - Token refresh mechanism
   - 24-hour access tokens
   - 7-day refresh window

2. **Password Security** (`auth/password.go`)
   - Bcrypt hashing (cost: 10)
   - Password validation rules
   - Minimum 8 characters
   - Requires letters and numbers

3. **Middleware** (`auth/middleware.go`)
   - `AuthMiddleware` - Protects routes requiring authentication
   - `OptionalAuthMiddleware` - Optional authentication
   - `WebSocketAuthMiddleware` - WebSocket authentication via query params
   - `CORSMiddleware` - CORS handling

4. **Rate Limiting** (`middleware/ratelimit.go`)
   - 100 requests per hour per user
   - 10 concurrent research jobs per user
   - In-memory cache with database persistence

5. **User Models** (`models/user.go`)
   - User accounts with email/password
   - Rate limit tracking
   - Multi-user session isolation

---

## Authentication Flow

### Registration Flow

```
User → POST /api/v1/auth/register
    ↓
Validate email/password
    ↓
Hash password (bcrypt)
    ↓
Create user in database
    ↓
Generate JWT token
    ↓
← Return user + token
```

### Login Flow

```
User → POST /api/v1/auth/login
    ↓
Find user by email
    ↓
Verify password (bcrypt)
    ↓
Generate JWT token
    ↓
← Return user + token
```

### Request Flow (Protected Endpoints)

```
Client → Request with Authorization header
    ↓
Extract JWT token
    ↓
Validate token signature & expiry
    ↓
Extract user ID from claims
    ↓
Add user context to request
    ↓
Check rate limits
    ↓
Process request (user-scoped)
    ↓
← Return response
```

---

## API Endpoints

### Public Endpoints (No Authentication Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |

### Protected Endpoints (Requires Authentication)

#### Authentication Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/me` | Get current user profile |
| PUT | `/api/v1/auth/profile` | Update user profile |
| POST | `/api/v1/auth/refresh` | Refresh JWT token |
| POST | `/api/v1/auth/password` | Change password |

#### Research Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/research/start` | Start new research session |
| GET | `/api/v1/research/:id/status` | Get research status |
| GET | `/api/v1/research/:id/result` | Get research results |
| DELETE | `/api/v1/research/:id` | Cancel research |
| GET | `/api/v1/research/:id/stream` | WebSocket stream (token in query) |

#### Document Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/documents/upload` | Upload document |
| POST | `/api/v1/documents/analyze` | Analyze document |
| GET | `/api/v1/documents` | List user's documents |
| GET | `/api/v1/documents/:id` | Get document details |
| DELETE | `/api/v1/documents/:id` | Delete document |

#### Session Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/sessions` | List user's research sessions |
| GET | `/api/v1/sessions/:id` | Get session details |
| DELETE | `/api/v1/sessions/:id` | Delete session |
| POST | `/api/v1/sessions/:id/export` | Export session |

#### Export & Statistics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/export/:id/download` | Download export |
| GET | `/api/v1/stats` | Get user statistics |

---

## Rate Limiting

### Default Limits

- **Requests**: 100 per hour per user
- **Concurrent Jobs**: 10 per user
- **Window**: 1 hour sliding window

### Rate Limit Headers

All authenticated responses include rate limit headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 2025-01-15T15:30:00Z
```

### Rate Limit Exceeded Response

```json
{
  "error": "Rate limit exceeded",
  "reset_at": "2025-01-15T15:30:00Z",
  "retry_after": 1800
}
```

---

## Security Best Practices

### JWT Secret

**IMPORTANT**: Set a strong JWT secret via environment variable:

```bash
export JWT_SECRET="your-super-secret-key-min-32-chars"
```

**Never** use the default secret in production!

### HTTPS Only

In production, **always** use HTTPS:
- Prevents token interception
- Protects password transmission
- Required for security compliance

### Password Requirements

- Minimum 8 characters
- Must contain at least one letter
- Must contain at least one number
- Hashed with bcrypt (cost: 10)

### Token Storage (Client-Side)

**Recommended**:
- Store tokens in httpOnly cookies (if server supports)
- Or secure localStorage with XSS protection

**Never**:
- Log tokens to console
- Store in URL parameters (except WebSocket)
- Expose in client-side source code

---

## Examples

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "researcher@example.com",
    "password": "SecurePass123",
    "name": "Dr. Jane Smith"
  }'
```

**Response:**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "researcher@example.com",
    "name": "Dr. Jane Smith",
    "created_at": "2025-01-15T14:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "researcher@example.com",
    "password": "SecurePass123"
  }'
```

**Response:**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "researcher@example.com",
    "name": "Dr. Jane Smith",
    "created_at": "2025-01-15T14:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Get Current User Profile

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "researcher@example.com",
  "name": "Dr. Jane Smith",
  "created_at": "2025-01-15T14:00:00Z"
}
```

### 4. Start Research Session (Authenticated)

```bash
curl -X POST http://localhost:8080/api/v1/research/start \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "query": "What are the latest developments in quantum computing?",
    "depth": "medium",
    "use_web": true,
    "use_wiki": true,
    "max_sources": 10
  }'
```

**Response:**
```json
{
  "job_id": "research-123e4567-e89b-12d3-a456-426614174000",
  "session_id": "session-789e0123-e89b-12d3-a456-426614174000",
  "status": "queued",
  "message": "Research job queued successfully"
}
```

### 5. WebSocket Connection (with Token)

```javascript
// JavaScript example
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";
const sessionId = "session-789e0123-e89b-12d3-a456-426614174000";
const ws = new WebSocket(`ws://localhost:8080/api/v1/research/${sessionId}/stream?token=${token}`);

ws.onmessage = (event) => {
  const progress = JSON.parse(event.data);
  console.log(`Progress: ${progress.progress}%`);
};
```

Using `wscat`:
```bash
wscat -c "ws://localhost:8080/api/v1/research/session-123/stream?token=eyJhbGc..."
```

### 6. List User's Documents

```bash
curl -X GET http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 7. Upload Document

```bash
curl -X POST http://localhost:8080/api/v1/documents/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "file=@research-paper.pdf" \
  -F "name=Quantum Computing Research 2025"
```

### 8. Update User Profile

```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dr. Jane Smith-Johnson"
  }'
```

### 9. Change Password

```bash
curl -X POST http://localhost:8080/api/v1/auth/password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "SecurePass123",
    "new_password": "EvenMoreSecure456"
  }'
```

### 10. Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.NEW_TOKEN..."
}
```

---

## Multi-User Isolation

### Database Filtering

All user-specific operations automatically filter by `user_id`:

- **Research Sessions**: Only show sessions created by the authenticated user
- **Documents**: Only show documents uploaded by the authenticated user
- **Jobs**: Only show jobs belonging to the authenticated user

### Authorization Checks

The system ensures users can only:
- View their own sessions and documents
- Delete their own resources
- Export their own research
- Access their own statistics

### Example: Accessing Another User's Session

```bash
# This will return 404 or "access denied" even if the session ID is valid
curl -X GET http://localhost:8080/api/v1/sessions/other-user-session-id \
  -H "Authorization: Bearer your-token"
```

**Response:**
```json
{
  "error": "session not found or access denied: other-user-session-id"
}
```

---

## Configuration

### Environment Variables

```bash
# Required in production
export JWT_SECRET="your-very-long-secure-random-string-min-32-chars"

# Optional (defaults shown)
export API_PORT=8080
export API_ENABLE_AUTH=true
export RATE_LIMIT_ENABLED=true
export RATE_LIMIT_MAX_REQUESTS=100
export RATE_LIMIT_MAX_JOBS=10
```

### Configuration File

Edit `~/.deep-research-agent/config.json`:

```json
{
  "api": {
    "enabled": true,
    "host": "0.0.0.0",
    "port": 8080,
    "enable_cors": true
  },
  "auth": {
    "token_expiry": "24h",
    "refresh_expiry": "168h",
    "enable_auth": true
  },
  "rate_limit": {
    "enabled": true,
    "max_requests_per_hour": 100,
    "max_concurrent_jobs": 10,
    "window_duration": "1h"
  }
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request payload: validation error"
}
```

### 401 Unauthorized
```json
{
  "error": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": "Access denied"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found or access denied"
}
```

### 409 Conflict
```json
{
  "error": "User with this email already exists"
}
```

### 429 Too Many Requests
```json
{
  "error": "Rate limit exceeded",
  "reset_at": "2025-01-15T15:30:00Z",
  "retry_after": 1800
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Troubleshooting

### Token Expired

**Solution**: Use the refresh endpoint to get a new token, or login again.

### Invalid Token Format

**Error**: `Invalid authorization header format. Expected: Bearer <token>`

**Solution**: Ensure header is: `Authorization: Bearer <your-token>`

### Rate Limit Exceeded

**Solution**: Wait for the reset time or contact admin to increase limits.

### CORS Issues

**Solution**: Ensure CORS is enabled in config and proper headers are sent.

---

## Testing

### Running the API Server

```bash
# Set JWT secret
export JWT_SECRET="test-secret-key-for-development-only"

# Start server
cd /home/user/golang-learning-projects/week4-capstone/deep-research-agent
go run cmd/api/main.go
```

### Integration Test Script

```bash
#!/bin/bash
API_URL="http://localhost:8080"

# 1. Register
echo "1. Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST $API_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234","name":"Test User"}')

TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

# 2. Get profile
echo "2. Getting profile..."
curl -s -X GET $API_URL/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq

# 3. Start research
echo "3. Starting research..."
curl -s -X POST $API_URL/api/v1/research/start \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query":"Go programming","depth":"medium"}' | jq

# 4. List sessions
echo "4. Listing sessions..."
curl -s -X GET $API_URL/api/v1/sessions \
  -H "Authorization: Bearer $TOKEN" | jq
```

---

## Summary

The Deep Research Agent API now features:

✅ **JWT-based authentication** with secure token management
✅ **User registration and login** with bcrypt password hashing
✅ **Multi-user session isolation** - users only see their own data
✅ **Rate limiting** - 100 req/hour, 10 concurrent jobs per user
✅ **WebSocket authentication** via query parameter tokens
✅ **Protected routes** with middleware-based access control
✅ **Profile management** - update name, email, password
✅ **Token refresh** mechanism for long-lived sessions

All research operations, document uploads, and sessions are now user-scoped and secure!
