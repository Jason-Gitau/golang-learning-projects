package main

import (
	"ai-agent-platform/config"
	"ai-agent-platform/database"
	"ai-agent-platform/handlers"
	"ai-agent-platform/middleware"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	if err := database.Initialize(cfg.Database.Path, cfg.Server.Environment == "development"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create router
	router := gin.New()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	agentHandler := handlers.NewAgentHandler()
	conversationHandler := handlers.NewConversationHandler()
	usageHandler := handlers.NewUsageHandler()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		protected.Use(middleware.RateLimitMiddleware(cfg.RateLimit.RequestsPerHour, cfg.RateLimit.WindowSize))
		{
			// User profile
			protected.GET("/auth/me", authHandler.GetProfile)

			// Agent management
			agents := protected.Group("/agents")
			{
				agents.POST("", agentHandler.CreateAgent)
				agents.GET("", agentHandler.ListAgents)
				agents.GET("/:id", agentHandler.GetAgent)
				agents.PUT("/:id", agentHandler.UpdateAgent)
				agents.DELETE("/:id", agentHandler.DeleteAgent)

				// Conversation creation for specific agent
				agents.POST("/:id/conversations", conversationHandler.CreateConversation)
			}

			// Conversation management
			conversations := protected.Group("/conversations")
			{
				conversations.GET("", conversationHandler.ListConversations)
				conversations.GET("/:id", conversationHandler.GetConversation)
				conversations.DELETE("/:id", conversationHandler.DeleteConversation)

				// Message management
				conversations.POST("/:id/messages", conversationHandler.CreateMessage)
				conversations.GET("/:id/messages", conversationHandler.GetMessages)
			}

			// Usage tracking
			usage := protected.Group("/usage")
			{
				usage.GET("", usageHandler.GetUsageStats)
				usage.GET("/rate-limit", usageHandler.GetRateLimitInfo)
				usage.POST("/log", usageHandler.CreateUsageLog)
			}
		}
	}

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server
	log.Printf("Starting AI Agent Platform API on port %s", cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Database: %s", cfg.Database.Path)
	log.Printf("Rate Limit: %d requests per hour", cfg.RateLimit.RequestsPerHour)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
