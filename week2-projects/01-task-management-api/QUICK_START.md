# 🚀 Quick Start Guide

Get the Task Management API running in **2 minutes**!

## Prerequisites

- Go 1.21 or higher installed
- Terminal/Command prompt

## Steps

### 1. Navigate to Project Directory

```bash
cd week2-projects/01-task-management-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Run the Server

```bash
go run main.go
```

You should see:

```
============================================================
🚀 Task Management API Server
============================================================
Server running on: http://localhost:8080
Mode: debug

Available endpoints:
  GET    /health                    - Health check
  POST   /api/v1/auth/register      - Register new user
  POST   /api/v1/auth/login         - Login user
  ...
============================================================
```

### 4. Test the API

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Register a User:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Create a Task (use the token from registration):**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Learn Go",
    "description": "Complete Week 2 of Go learning path",
    "priority": "high",
    "status": "todo"
  }'
```

## What's Next?

- Read the [README.md](README.md) for full documentation
- Read the [START_HERE.md](START_HERE.md) for learning guide
- Check out the [example requests](#example-requests) below

## Example Requests

### 1. Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@test.com", "password": "secure123"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@test.com", "password": "secure123"}'
```

Copy the `token` from the response!

### 3. Get Tasks
```bash
curl http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Create Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Build REST API",
    "description": "Complete the task management API",
    "priority": "urgent",
    "status": "in_progress",
    "due_date": "2024-12-31T23:59:59Z"
  }'
```

### 5. Get Statistics
```bash
curl http://localhost:8080/api/v1/tasks/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Tips

- Save your token after login/register - you'll need it for all task operations
- The database file is created at `./data/tasks.db`
- Use `Ctrl+C` to stop the server
- Check the terminal logs to see all requests

## Troubleshooting

**"Failed to connect to database"**
- The `data` directory is created automatically
- Make sure you have write permissions

**"Invalid token"**
- Make sure you're using `Bearer YOUR_TOKEN` in the Authorization header
- Token expires after 24 hours - login again

**"Port already in use"**
- Change the port: `export SERVER_PORT=9000` (Linux/Mac) or `set SERVER_PORT=9000` (Windows)

---

**You're ready to go!** 🎉

For complete documentation, see [README.md](README.md)
