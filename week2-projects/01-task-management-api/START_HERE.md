# 🎓 START HERE - Learning Guide

Welcome to Week 2 of your Go learning journey! This guide will help you understand the Task Management API project and extract maximum learning value.

## 📚 What You'll Learn

By studying and building this project, you'll master:

1. **Web Development with Gin**
2. **Database Operations with GORM**
3. **JWT Authentication**
4. **REST API Design**
5. **Middleware Patterns**
6. **Project Structure and Organization**

## 🗺️ Learning Path

### Step 1: Understand the Big Picture (30 minutes)

Before diving into code, understand what this API does:

1. **Run the application:**
   ```bash
   go run main.go
   ```

2. **Test the endpoints** (see QUICK_START.md for examples)
   - Register a user
   - Login
   - Create some tasks
   - Get tasks with filters
   - Check statistics

3. **Questions to answer:**
   - What happens when I register?
   - Where is my data stored?
   - How does authentication work?
   - What can I do with tasks?

### Step 2: Study the Code Structure (45 minutes)

Read files in this order for maximum comprehension:

#### Phase 1: Data Layer
1. **models/user.go** (10 min)
   - Understanding structs and tags
   - GORM model definition
   - JSON serialization
   - Request/Response DTOs

2. **models/task.go** (10 min)
   - Enums in Go (TaskPriority, TaskStatus)
   - Pointer fields for optional values
   - Data validation tags
   - Relationships between models

3. **database/database.go** (5 min)
   - Database connection
   - Auto-migrations
   - Error handling

#### Phase 2: Business Logic
4. **utils/password.go** (5 min)
   - Password hashing with bcrypt
   - Security best practices

5. **utils/jwt.go** (10 min)
   - JWT token structure
   - Claims and expiration
   - Token generation and validation

6. **config/config.go** (5 min)
   - Configuration management
   - Environment variables
   - Sensible defaults

#### Phase 3: HTTP Layer
7. **middleware/auth.go** (10 min)
   - Middleware concept
   - Authorization header parsing
   - Context usage in Gin
   - Request interception

8. **middleware/logger.go** (5 min)
   - Request logging
   - Performance tracking

9. **middleware/error.go** (10 min)
   - Centralized error handling
   - Panic recovery
   - CORS configuration

10. **handlers/auth.go** (15 min)
    - Registration logic
    - Login flow
    - Profile retrieval
    - HTTP status codes
    - JSON responses

11. **handlers/task.go** (20 min)
    - CRUD operations
    - Query parameters
    - Filtering and sorting
    - Statistics aggregation
    - User ownership verification

12. **main.go** (15 min)
    - Application initialization
    - Route setup
    - Middleware ordering
    - Server configuration

**Total time: ~2 hours**

### Step 3: Key Concepts Deep Dive

#### Concept 1: HTTP Middleware

**What is it?**
Functions that execute before/after your handlers, allowing you to add cross-cutting functionality.

**Example from the project:**
```go
// middleware/auth.go
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before handler: validate token
        token := c.GetHeader("Authorization")
        // ... validation logic ...

        c.Next() // Execute the actual handler

        // After handler: could add cleanup logic here
    }
}
```

**Why it matters:**
- DRY principle (Don't Repeat Yourself)
- Separation of concerns
- Easy to add/remove functionality

**Try this:**
Create a new middleware that counts API requests.

#### Concept 2: JWT Authentication

**The Flow:**
```
1. User registers → Server hashes password, stores user
2. User logs in → Server validates password, generates JWT
3. Client stores JWT
4. Client sends JWT with each request
5. Server validates JWT, extracts user info
```

**Why JWT?**
- Stateless (no session storage)
- Contains user info (no database lookup)
- Can't be forged (cryptographically signed)

**Try this:**
Modify token expiration time and test what happens.

#### Concept 3: GORM Model Relationships

**One-to-Many:**
```go
type User struct {
    ID    uint
    Tasks []Task `gorm:"foreignKey:UserID"`
}

type Task struct {
    ID     uint
    UserID uint
    User   User
}
```

**GORM does:**
- Creates foreign key constraints
- Enables cascade operations
- Allows eager/lazy loading

**Try this:**
Add `Preload("Tasks")` when fetching users to include their tasks.

#### Concept 4: Input Validation

**Using Gin Binding:**
```go
type CreateTaskRequest struct {
    Title string `json:"title" binding:"required,min=1,max=200"`
}
```

**Validation tags:**
- `required` - Must be present
- `min=1` - Minimum length
- `email` - Valid email format
- `oneof=todo in_progress` - Must be one of values

**Try this:**
Add validation for due_date to ensure it's in the future.

#### Concept 5: RESTful Design

**HTTP Methods:**
- `POST` - Create resource
- `GET` - Read resource(s)
- `PUT`/`PATCH` - Update resource
- `DELETE` - Delete resource

**Status Codes:**
- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Not authenticated
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Server error

**Resource URLs:**
```
/api/v1/tasks           → List/Create
/api/v1/tasks/:id       → Get/Update/Delete
/api/v1/tasks/stats     → Special operation
```

### Step 4: Hands-On Exercises

#### Exercise 1: Add Task Categories (Easy)
**Goal:** Add categories to tasks

1. Add `Category` field to Task model
2. Update CreateTaskRequest
3. Test creating tasks with categories

<details>
<summary>Solution Hint</summary>

```go
// In models/task.go
type Task struct {
    // ... existing fields ...
    Category string `json:"category"`
}

// In CreateTaskRequest
type CreateTaskRequest struct {
    // ... existing fields ...
    Category string `json:"category" binding:"max=50"`
}
```
</details>

#### Exercise 2: Add Pagination (Medium)
**Goal:** Add pagination to task list

1. Add `page` and `limit` query parameters
2. Use GORM's `Offset` and `Limit`
3. Return total count with results

<details>
<summary>Solution Hint</summary>

```go
// Parse pagination params
page := c.DefaultQuery("page", "1")
limit := c.DefaultQuery("limit", "10")

// Convert to integers and apply
query.Offset((page - 1) * limit).Limit(limit)

// Get total count
var total int64
database.DB.Model(&models.Task{}).Where(...).Count(&total)
```
</details>

#### Exercise 3: Add Email Validation (Hard)
**Goal:** Send welcome email on registration

1. Create a utils/email.go file
2. Implement SendWelcomeEmail function
3. Call it after user registration
4. Handle errors gracefully

#### Exercise 4: Add Task Search (Hard)
**Goal:** Full-text search in title/description

1. Add `search` query parameter
2. Use GORM's `Where` with `LIKE`
3. Search both title and description

<details>
<summary>Solution Hint</summary>

```go
search := c.Query("search")
if search != "" {
    query = query.Where(
        "title LIKE ? OR description LIKE ?",
        "%"+search+"%",
        "%"+search+"%",
    )
}
```
</details>

### Step 5: Testing Your Knowledge

Answer these questions without looking at code:

1. **What happens when a user logs in?**
   - Where is the password checked?
   - What is returned?
   - How long is the token valid?

2. **How does AuthMiddleware work?**
   - What header does it check?
   - What happens if token is invalid?
   - What does it add to the context?

3. **How are tasks filtered?**
   - What query parameters are supported?
   - How does GORM build the query?
   - What's the default sort order?

4. **What security measures are in place?**
   - How are passwords stored?
   - How is authentication handled?
   - Can users access other users' tasks?

## 🎯 Learning Objectives Checklist

After completing this guide, you should be able to:

### REST API Basics
- [ ] Explain HTTP methods and status codes
- [ ] Design RESTful endpoints
- [ ] Handle JSON requests/responses
- [ ] Implement proper error handling

### Gin Framework
- [ ] Set up routes and route groups
- [ ] Create and use middleware
- [ ] Bind and validate requests
- [ ] Return JSON responses

### Database with GORM
- [ ] Define models with tags
- [ ] Perform CRUD operations
- [ ] Use query builders
- [ ] Handle relationships
- [ ] Run migrations

### Authentication & Security
- [ ] Hash passwords with bcrypt
- [ ] Generate JWT tokens
- [ ] Validate tokens
- [ ] Implement auth middleware
- [ ] Secure endpoints

### Project Organization
- [ ] Structure Go projects
- [ ] Separate concerns into packages
- [ ] Manage configuration
- [ ] Handle errors consistently

## 🚀 Next Steps

### Immediate (This Week)
1. Complete all exercises above
2. Add custom features you'd like
3. Write tests for at least one handler
4. Deploy to a cloud platform (Railway, Render, Fly.io)

### Short-term (This Month)
1. Add PostgreSQL support
2. Implement rate limiting
3. Add Swagger documentation
4. Create a simple frontend
5. Add WebSocket support for real-time updates

### Long-term (Next Month)
1. Migrate to microservices architecture
2. Add caching with Redis
3. Implement event-driven architecture
4. Add monitoring and observability
5. Learn about Kubernetes deployment

## 📚 Additional Resources

### Official Documentation
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM](https://gorm.io/docs/)
- [JWT Specification](https://jwt.io/introduction)
- [Go net/http](https://pkg.go.dev/net/http)

### Tutorials
- [Let's Go by Alex Edwards](https://lets-go.alexedwards.net/)
- [Go Web Examples](https://gowebexamples.com/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)

### Best Practices
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## 💡 Tips for Success

1. **Read the code, don't just run it** - Understanding beats copying
2. **Make mistakes** - Break things, see what happens, fix them
3. **Add logging** - See what your code is doing
4. **Use a debugger** - Step through the code
5. **Test everything** - Try edge cases and error scenarios
6. **Refactor** - Come back and improve your code
7. **Ask questions** - Why does it work this way?

## 🎉 Congratulations!

You've built a production-ready REST API with:
- ✅ User authentication
- ✅ Database integration
- ✅ Input validation
- ✅ Error handling
- ✅ Clean architecture

This is a significant achievement! You now have the skills to build real-world web applications with Go.

**Keep coding, keep learning, and most importantly - have fun!** 🚀

---

**Questions? Stuck? Tips:**
1. Re-read the relevant section above
2. Check the inline code comments
3. Review the official documentation
4. Search for specific error messages
5. Join the Go community on Reddit/Slack

**Ready to code?** Head to [QUICK_START.md](QUICK_START.md) to run the project!
