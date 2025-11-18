# Web UI Integration Guide

Complete guide for integrating the web UI with your Deep Research Agent API.

## Quick Start

### Step 1: Import the web package

```go
import "deep-research-agent/web"
```

### Step 2: Add to your main.go

```go
package main

import (
    "github.com/gin-gonic/gin"
    "deep-research-agent/web"
)

func main() {
    router := gin.Default()

    // Enable CORS for development (optional)
    router.Use(web.CORSMiddleware())

    // Setup your API routes here
    setupAPIRoutes(router)

    // Setup web UI (must be last)
    web.ServeStaticFiles(router)

    router.Run(":8080")
}
```

### Step 3: Access the UI

Open your browser to `http://localhost:8080`

## Complete Integration Example

### Full main.go with API routes

```go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "deep-research-agent/web"
)

func main() {
    // Create Gin router
    router := gin.Default()

    // Enable CORS middleware (for development)
    router.Use(web.CORSMiddleware())

    // API v1 routes
    api := router.Group("/api/v1")
    {
        // Authentication routes
        auth := api.Group("/auth")
        {
            auth.POST("/register", handleRegister)
            auth.POST("/login", handleLogin)
            auth.GET("/profile", authMiddleware, handleProfile)
        }

        // Research routes (protected)
        research := api.Group("/research")
        research.Use(authMiddleware)
        {
            research.POST("/start", handleStartResearch)
            research.GET("/:jobId/status", handleJobStatus)
            research.GET("/:jobId/stream", handleWebSocketStream)
            research.POST("/:jobId/cancel", handleCancelResearch)

            // Sessions
            research.GET("/sessions", handleListSessions)
            research.GET("/sessions/search", handleSearchSessions)
            research.GET("/sessions/:sessionId", handleGetSession)
            research.DELETE("/sessions/:sessionId", handleDeleteSession)
            research.GET("/sessions/:sessionId/export", handleExportSession)
        }

        // Document routes (protected)
        documents := api.Group("/documents")
        documents.Use(authMiddleware)
        {
            documents.POST("/upload", handleUploadDocuments)
            documents.GET("", handleListDocuments)
            documents.GET("/:documentId", handleGetDocument)
            documents.DELETE("/:documentId", handleDeleteDocument)
        }
    }

    // Serve web UI (MUST BE LAST)
    web.ServeStaticFiles(router)

    // Start server
    log.Println("Server starting on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

// Example handler implementations
func handleRegister(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=8"`
        Name     string `json:"name" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Your registration logic here
    // - Hash password
    // - Create user in database
    // - Generate JWT token

    c.JSON(http.StatusOK, gin.H{
        "token": "jwt_token_here",
        "user": gin.H{
            "id":    "user123",
            "email": req.Email,
            "name":  req.Name,
        },
    })
}

func handleLogin(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Your login logic here
    // - Verify credentials
    // - Generate JWT token

    c.JSON(http.StatusOK, gin.H{
        "token": "jwt_token_here",
        "user": gin.H{
            "id":    "user123",
            "email": req.Email,
            "name":  "User Name",
        },
    })
}

func handleProfile(c *gin.Context) {
    // Get user from context (set by authMiddleware)
    user := c.MustGet("user")

    c.JSON(http.StatusOK, gin.H{
        "user": user,
    })
}

func handleStartResearch(c *gin.Context) {
    var req struct {
        Query         string   `json:"query" binding:"required"`
        Depth         string   `json:"depth"`
        DocumentIDs   []string `json:"document_ids"`
        UseWeb        bool     `json:"use_web"`
        UseWikipedia  bool     `json:"use_wikipedia"`
        FactCheck     bool     `json:"fact_check"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Your research start logic here
    // - Create research job
    // - Start background processing
    // - Return job ID

    jobID := "job_" + generateID()

    c.JSON(http.StatusOK, gin.H{
        "job_id": jobID,
        "status": "started",
    })
}

func handleJobStatus(c *gin.Context) {
    jobID := c.Param("jobId")

    // Get job status from your system

    c.JSON(http.StatusOK, gin.H{
        "job_id":       jobID,
        "status":       "running",
        "progress":     60,
        "current_step": "Analyzing Wikipedia article",
        "total_steps":  5,
        "current":      3,
    })
}

func handleWebSocketStream(c *gin.Context) {
    jobID := c.Param("jobId")

    // Upgrade to WebSocket
    // Send real-time updates
    // Example messages:
    // {"type": "progress", "data": {"percentage": 60, "message": "Processing..."}}
    // {"type": "step_started", "data": {"step": "Web Search"}}
    // {"type": "step_completed", "data": {"step": "Web Search", "duration": "2.3s"}}
    // {"type": "research_completed", "data": {"session_id": "123", "results": "..."}}
}

func handleListSessions(c *gin.Context) {
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "20")
    status := c.Query("status")

    // Your session list logic here

    c.JSON(http.StatusOK, gin.H{
        "sessions": []gin.H{
            {
                "id":         "session123",
                "query":      "What are AI trends?",
                "status":     "completed",
                "created_at": "2024-01-15T10:30:00Z",
                "duration":   320,
            },
        },
        "pagination": gin.H{
            "current_page": 1,
            "total_pages":  5,
            "total_items":  100,
            "has_next":     true,
            "has_prev":     false,
        },
    })
}

func handleGetSession(c *gin.Context) {
    sessionID := c.Param("sessionId")

    // Get session details

    c.JSON(http.StatusOK, gin.H{
        "id":         sessionID,
        "query":      "What are AI trends?",
        "status":     "completed",
        "created_at": "2024-01-15T10:30:00Z",
        "duration":   320,
        "results":    "# Research Results\n\n...",
    })
}

func handleDeleteSession(c *gin.Context) {
    sessionID := c.Param("sessionId")

    // Delete session logic

    c.JSON(http.StatusOK, gin.H{
        "message": "Session deleted successfully",
    })
}

func handleExportSession(c *gin.Context) {
    sessionID := c.Param("sessionId")
    format := c.DefaultQuery("format", "markdown")

    // Export session in requested format

    var contentType string
    var filename string
    var content []byte

    switch format {
    case "json":
        contentType = "application/json"
        filename = "research-" + sessionID + ".json"
        // Generate JSON content
    case "markdown":
        contentType = "text/markdown"
        filename = "research-" + sessionID + ".md"
        // Generate Markdown content
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
        return
    }

    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Data(http.StatusOK, contentType, content)
}

func handleUploadDocuments(c *gin.Context) {
    // Parse multipart form
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
        return
    }

    files := form.File["files"]
    var documentIDs []string

    for _, file := range files {
        // Save file and process
        // documentID := saveDocument(file)
        // documentIDs = append(documentIDs, documentID)
    }

    c.JSON(http.StatusOK, gin.H{
        "message":      "Files uploaded successfully",
        "document_ids": documentIDs,
    })
}

func handleListDocuments(c *gin.Context) {
    // List user's documents

    c.JSON(http.StatusOK, gin.H{
        "documents": []gin.H{
            {
                "id":          "doc123",
                "filename":    "research_paper.pdf",
                "size":        1024000,
                "type":        "pdf",
                "uploaded_at": "2024-01-15T10:30:00Z",
            },
        },
    })
}

func handleGetDocument(c *gin.Context) {
    documentID := c.Param("documentId")

    // Get document details

    c.JSON(http.StatusOK, gin.H{
        "id":          documentID,
        "filename":    "research_paper.pdf",
        "size":        1024000,
        "type":        "pdf",
        "uploaded_at": "2024-01-15T10:30:00Z",
    })
}

func handleDeleteDocument(c *gin.Context) {
    documentID := c.Param("documentId")

    // Delete document

    c.JSON(http.StatusOK, gin.H{
        "message": "Document deleted successfully",
    })
}

// Auth middleware
func authMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token"})
        c.Abort()
        return
    }

    // Verify JWT token
    // Extract user from token
    // Set user in context

    c.Set("user", gin.H{
        "id":    "user123",
        "email": "user@example.com",
        "name":  "User Name",
    })

    c.Next()
}

func generateID() string {
    // Generate unique ID
    return "12345"
}
```

## WebSocket Implementation

### Using gorilla/websocket

```go
import (
    "github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Configure properly for production
    },
}

func handleWebSocketStream(c *gin.Context) {
    jobID := c.Param("jobId")

    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }
    defer conn.Close()

    // Get job progress channel
    progressChan := getJobProgressChannel(jobID)

    // Send updates to client
    for update := range progressChan {
        message := map[string]interface{}{
            "type": update.Type,
            "data": update.Data,
        }

        if err := conn.WriteJSON(message); err != nil {
            log.Println("WebSocket write error:", err)
            break
        }
    }
}

// Example progress channel
func getJobProgressChannel(jobID string) <-chan ProgressUpdate {
    ch := make(chan ProgressUpdate)

    go func() {
        defer close(ch)

        // Send progress updates
        ch <- ProgressUpdate{
            Type: "progress",
            Data: map[string]interface{}{
                "percentage": 20,
                "message":    "Starting research...",
            },
        }

        ch <- ProgressUpdate{
            Type: "step_started",
            Data: map[string]interface{}{
                "step": "Web Search",
            },
        }

        // ... more updates

        ch <- ProgressUpdate{
            Type: "research_completed",
            Data: map[string]interface{}{
                "session_id": "session123",
                "results":    "...",
            },
        }
    }()

    return ch
}

type ProgressUpdate struct {
    Type string
    Data map[string]interface{}
}
```

## Testing the Integration

### 1. Start the server

```bash
go run main.go
```

### 2. Access the UI

Open browser to `http://localhost:8080`

### 3. Test authentication

1. Register a new account
2. Login with credentials
3. Check that JWT token is stored in localStorage

### 4. Test research flow

1. Enter a research query
2. Configure options
3. Start research
4. Watch real-time progress
5. View results when complete

### 5. Test WebSocket

Open browser DevTools > Network > WS to see WebSocket messages

## Deployment Checklist

### Development

- [ ] API server running on localhost:8080
- [ ] CORS enabled for development
- [ ] Browser DevTools for debugging
- [ ] Hot reload for Go code changes

### Staging

- [ ] Configure proper CORS origins
- [ ] Use environment variables for config
- [ ] Enable HTTPS/WSS
- [ ] Test with real data

### Production

- [ ] Minify CSS and JavaScript
- [ ] Enable gzip compression
- [ ] Set cache headers
- [ ] Configure CSP headers
- [ ] Use secure WebSocket (wss://)
- [ ] Set httpOnly and secure flags for cookies
- [ ] Rate limiting on API endpoints
- [ ] Monitor WebSocket connections
- [ ] Log errors properly

## Troubleshooting

### Issue: CORS errors

**Solution**: Enable CORS middleware
```go
router.Use(web.CORSMiddleware())
```

### Issue: Static files not loading

**Solution**: Check file paths and ensure server.go is properly configured
```go
web.ServeStaticFiles(router)
```

### Issue: WebSocket not connecting

**Solution**:
1. Check WebSocket endpoint is correct
2. Verify token is passed in query string
3. Check CORS for WebSocket upgrade

### Issue: 404 on page refresh

**Solution**: Ensure NoRoute handler serves index.html for non-API routes

## Advanced Configuration

### Custom static directory

```go
config := &web.ServerConfig{
    StaticDir: "./custom/static/path",
    IndexFile: "./custom/static/path/index.html",
}
web.SetupRoutesWithConfig(router, config)
```

### Custom middleware

```go
// Add logging middleware
router.Use(gin.Logger())

// Add recovery middleware
router.Use(gin.Recovery())

// Custom middleware
router.Use(func(c *gin.Context) {
    // Your middleware logic
    c.Next()
})
```

### Environment configuration

```go
import "os"

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router := gin.Default()
    // ... setup routes
    router.Run(":" + port)
}
```

## Support

For questions or issues:
1. Check this integration guide
2. Review the main README.md
3. Check browser console for errors
4. Verify API endpoints are working
5. Test with curl/Postman first
