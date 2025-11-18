package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-learning/agent-orchestrator/agent"
	"github.com/golang-learning/agent-orchestrator/config"
	"github.com/golang-learning/agent-orchestrator/models"
)

// Server represents the HTTP API server
type Server struct {
	router  *gin.Engine
	manager *agent.Manager
	config  *config.Config
	server  *http.Server
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, manager *agent.Manager) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	s := &Server{
		router:  router,
		manager: manager,
		config:  cfg,
	}

	s.setupRoutes()

	return s
}

// setupRoutes sets up all API routes
func (s *Server) setupRoutes() {
	// Enable CORS if configured
	if s.config.EnableCORS {
		s.router.Use(corsMiddleware())
	}

	// Health check
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Request endpoints
		v1.POST("/request", s.submitRequest)

		// Agent endpoints
		v1.GET("/agents", s.getAllAgents)
		v1.GET("/agents/:id", s.getAgent)

		// Tool endpoints
		v1.GET("/tools", s.getTools)

		// Statistics endpoints
		v1.GET("/stats", s.getStatistics)
	}

	// Root endpoint
	s.router.GET("/", s.index)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := ":" + s.config.APIPort

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Printf("Starting HTTP API server on %s", addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping HTTP API server...")
	return s.server.Shutdown(ctx)
}

// index handles the root endpoint
func (s *Server) index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "Agent Orchestrator",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": map[string]string{
			"health":     "GET /health",
			"request":    "POST /api/v1/request",
			"agents":     "GET /api/v1/agents",
			"agent":      "GET /api/v1/agents/:id",
			"tools":      "GET /api/v1/tools",
			"statistics": "GET /api/v1/stats",
		},
	})
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	healthy := s.manager.HealthCheck()

	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"healthy": healthy,
		"time":    time.Now().Format(time.RFC3339),
	})
}

// submitRequest handles request submissions
func (s *Server) submitRequest(c *gin.Context) {
	var reqInput struct {
		ToolName string                 `json:"tool_name" binding:"required"`
		Params   map[string]interface{} `json:"params"`
		Timeout  int                    `json:"timeout"` // in seconds
	}

	if err := c.ShouldBindJSON(&reqInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Create request
	req := models.NewRequest(reqInput.ToolName, reqInput.Params)

	// Set custom timeout if provided
	if reqInput.Timeout > 0 {
		req.Timeout = time.Duration(reqInput.Timeout) * time.Second
	}

	// Submit request
	resp, err := s.manager.SubmitRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Request failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// getAllAgents returns information about all agents
func (s *Server) getAllAgents(c *gin.Context) {
	agents := s.manager.GetAllAgentInfo()
	c.JSON(http.StatusOK, gin.H{
		"count":  len(agents),
		"agents": agents,
	})
}

// getAgent returns information about a specific agent
func (s *Server) getAgent(c *gin.Context) {
	agentID := c.Param("id")

	info, err := s.manager.GetAgentInfo(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// getTools returns a list of available tools
func (s *Server) getTools(c *gin.Context) {
	registry := s.manager.GetToolRegistry()
	tools := registry.List()

	toolInfos := make([]map[string]string, len(tools))
	for i, tool := range tools {
		toolInfos[i] = map[string]string{
			"name":        tool.Name(),
			"description": tool.Description(),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(toolInfos),
		"tools": toolInfos,
	})
}

// getStatistics returns system statistics
func (s *Server) getStatistics(c *gin.Context) {
	stats := s.manager.GetStatistics()
	c.JSON(http.StatusOK, stats)
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
