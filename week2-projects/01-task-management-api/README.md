# 📋 Task Management REST API

A production-ready REST API built with Go and Gin framework for managing tasks with user authentication, JWT tokens, and full CRUD operations.

## 🎯 Week 2 Project Overview

This project is part of the **1-Month Go Learning Path - Week 2: Web Development & APIs**

### Learning Objectives

- ✅ Build REST APIs with Gin framework
- ✅ Implement JWT authentication
- ✅ Work with databases using GORM
- ✅ Create middleware for authentication and logging
- ✅ Handle HTTP requests and responses
- ✅ Validate input data
- ✅ Structure a real-world Go web application

## ✨ Features

### Core Features
- 🔐 **User Authentication** - Register and login with JWT tokens
- 📝 **Task Management** - Full CRUD operations for tasks
- 👤 **User-Specific Tasks** - Each user has their own task list
- 🎨 **Task Metadata** - Priority levels, status tracking, due dates
- 🔍 **Filtering & Sorting** - Filter tasks by status/priority, sort by various fields
- 📊 **Statistics** - Get insights about your tasks

### Technical Features
- ✅ RESTful API design
- ✅ JWT-based authentication
- ✅ Input validation with Gin binding
- ✅ Error handling middleware
- ✅ Request logging
- ✅ CORS support
- ✅ SQLite database with GORM
- ✅ Password hashing with bcrypt
- ✅ Clean architecture with separate packages

## 🏗️ Project Structure

```
01-task-management-api/
├── main.go                  # Application entry point, routes setup
├── go.mod                   # Go module dependencies
├── QUICK_START.md          # Get started in 2 minutes
├── README.md               # This file
├── START_HERE.md           # Learning guide
│
├── config/                 # Configuration management
│   └── config.go          # App configuration, environment variables
│
├── models/                # Data models and DTOs
│   ├── user.go           # User model and request/response types
│   └── task.go           # Task model, enums, and request types
│
├── database/              # Database connection and setup
│   └── database.go       # Database initialization, migrations
│
├── handlers/              # HTTP request handlers
│   ├── auth.go          # Authentication handlers (register, login, profile)
│   └── task.go          # Task CRUD handlers
│
├── middleware/            # HTTP middleware
│   ├── auth.go          # JWT authentication middleware
│   ├── logger.go        # Request logging middleware
│   └── error.go         # Error handling and CORS middleware
│
├── utils/                 # Utility functions
│   ├── jwt.go           # JWT token generation and validation
│   └── password.go      # Password hashing and verification
│
└── data/                  # Database storage (created automatically)
    └── tasks.db          # SQLite database file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- curl or Postman for testing (optional)

### Installation

1. **Navigate to the project directory:**
   ```bash
   cd week2-projects/01-task-management-api
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

### Building

To build an executable:

```bash
# Linux/Mac
go build -o task-api main.go

# Windows
go build -o task-api.exe main.go
```

Then run:
```bash
./task-api        # Linux/Mac
task-api.exe      # Windows
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

All task endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

### Endpoints

#### Health Check
```http
GET /health
```

#### Authentication

##### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "john",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

##### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

##### Get Profile
```http
GET /api/v1/auth/profile
Authorization: Bearer YOUR_TOKEN
```

#### Tasks

##### Create Task
```http
POST /api/v1/tasks
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json

{
  "title": "Complete project",
  "description": "Finish the REST API implementation",
  "priority": "high",
  "status": "todo",
  "due_date": "2024-12-31T23:59:59Z"
}
```

**Priority Options:** `low`, `medium`, `high`, `urgent`
**Status Options:** `todo`, `in_progress`, `completed`, `cancelled`

##### Get All Tasks
```http
GET /api/v1/tasks
Authorization: Bearer YOUR_TOKEN

# With filters
GET /api/v1/tasks?status=todo&priority=high&sort_by=due_date&order=asc
```

**Query Parameters:**
- `status` - Filter by status
- `priority` - Filter by priority
- `sort_by` - Sort by field (`created_at`, `due_date`, `priority`)
- `order` - Sort order (`asc`, `desc`)

##### Get Single Task
```http
GET /api/v1/tasks/:id
Authorization: Bearer YOUR_TOKEN
```

##### Update Task
```http
PUT /api/v1/tasks/:id
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json

{
  "title": "Updated title",
  "status": "in_progress"
}
```

All fields are optional - only send what you want to update.

##### Delete Task
```http
DELETE /api/v1/tasks/:id
Authorization: Bearer YOUR_TOKEN
```

##### Get Task Statistics
```http
GET /api/v1/tasks/stats
Authorization: Bearer YOUR_TOKEN
```

**Response:**
```json
{
  "total_tasks": 10,
  "completed_tasks": 3,
  "pending_tasks": 5,
  "overdue_tasks": 2,
  "tasks_by_priority": {
    "high": 4,
    "low": 2,
    "medium": 3,
    "urgent": 1
  },
  "tasks_by_status": {
    "cancelled": 0,
    "completed": 3,
    "in_progress": 2,
    "todo": 5
  }
}
```

## 🔧 Configuration

The application can be configured using environment variables:

```bash
# Server configuration
export SERVER_PORT=8080        # Server port (default: 8080)
export GIN_MODE=release        # Gin mode: debug or release (default: debug)

# Database configuration
export DB_PATH=./data/tasks.db # Database file path (default: ./data/tasks.db)

# JWT configuration
export JWT_SECRET=your-secret-key  # JWT signing secret (CHANGE IN PRODUCTION!)
```

## 🧪 Testing Examples

### Using curl

**1. Register and save the token:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123"
  }' | jq -r '.token' > token.txt

TOKEN=$(cat token.txt)
```

**2. Create a task:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Learn Go Web Development",
    "description": "Complete Week 2 project",
    "priority": "high",
    "status": "in_progress"
  }'
```

**3. Get all tasks:**
```bash
curl http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer $TOKEN"
```

**4. Get statistics:**
```bash
curl http://localhost:8080/api/v1/tasks/stats \
  -H "Authorization: Bearer $TOKEN"
```

**5. Filter tasks:**
```bash
# Get all high-priority tasks
curl "http://localhost:8080/api/v1/tasks?priority=high" \
  -H "Authorization: Bearer $TOKEN"

# Get all completed tasks sorted by date
curl "http://localhost:8080/api/v1/tasks?status=completed&sort_by=created_at&order=desc" \
  -H "Authorization: Bearer $TOKEN"
```

## 🏛️ Architecture & Design Patterns

### Clean Architecture

The project follows clean architecture principles with clear separation of concerns:

- **Models**: Data structures and business logic
- **Handlers**: HTTP request/response handling
- **Middleware**: Cross-cutting concerns (auth, logging, errors)
- **Database**: Data persistence layer
- **Utils**: Reusable helper functions

### Key Concepts Demonstrated

1. **RESTful Design**: Proper use of HTTP methods and status codes
2. **Middleware Pattern**: Reusable request processing logic
3. **Dependency Injection**: Passing configuration to handlers
4. **DTO Pattern**: Separate request/response models from database models
5. **Repository Pattern**: Database abstraction with GORM
6. **Error Handling**: Centralized error handling middleware

## 🔒 Security Features

- ✅ **Password Hashing**: Bcrypt for secure password storage
- ✅ **JWT Tokens**: Stateless authentication
- ✅ **Input Validation**: Gin binding with validation tags
- ✅ **SQL Injection Prevention**: GORM parameterized queries
- ✅ **CORS**: Configurable cross-origin resource sharing
- ✅ **Authorization**: User-specific resource access

## 📖 Learning Notes

### What You'll Learn

#### 1. **Gin Framework**
   - Routing and route groups
   - Middleware usage
   - Request binding and validation
   - JSON responses

#### 2. **Database with GORM**
   - Model definition
   - Auto migrations
   - CRUD operations
   - Relationships (User has many Tasks)
   - Query building

#### 3. **Authentication**
   - JWT token generation
   - Token validation
   - Middleware-based auth
   - Context sharing between middleware and handlers

#### 4. **Best Practices**
   - Project structure
   - Error handling
   - Configuration management
   - Logging
   - API versioning

## 🎓 Next Steps

After completing this project, you should be able to:

- ✅ Build REST APIs from scratch with Go
- ✅ Implement authentication systems
- ✅ Work with databases using ORMs
- ✅ Structure larger Go applications
- ✅ Handle HTTP requests and responses
- ✅ Validate and sanitize user input
- ✅ Implement middleware for cross-cutting concerns

### Enhancements to Try

1. **Add more features:**
   - Task categories/tags
   - Task sharing between users
   - File attachments
   - Comments on tasks

2. **Improve the API:**
   - Pagination for task lists
   - Full-text search
   - Bulk operations
   - Webhooks for task updates

3. **Add testing:**
   - Unit tests for handlers
   - Integration tests for API endpoints
   - Test coverage reporting

4. **Production readiness:**
   - Add PostgreSQL support
   - Implement rate limiting
   - Add API documentation (Swagger)
   - Docker containerization
   - CI/CD pipeline

## 🐛 Common Issues

### Database Locked
If you get "database is locked" errors:
- Make sure only one instance of the app is running
- Close any DB browser tools

### Token Expired
Tokens expire after 24 hours. Solution:
- Login again to get a new token

### Port Already in Use
If port 8080 is busy:
```bash
export SERVER_PORT=9000
go run main.go
```

## 📝 License

This is a learning project. Feel free to use it for educational purposes.

## 🙏 Acknowledgments

- Built as part of the 1-Month Go Learning Path
- Gin framework: https://gin-gonic.com/
- GORM: https://gorm.io/
- JWT-Go: https://github.com/golang-jwt/jwt

---

**Happy Coding!** 🚀

For a quick start, see [QUICK_START.md](QUICK_START.md)
For learning guidance, see [START_HERE.md](START_HERE.md)
