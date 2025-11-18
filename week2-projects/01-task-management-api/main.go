package main

import (
	"log"
	"strings"
	"task-management-api/config"
	"task-management-api/database"
	"task-management-api/handlers"
	"task-management-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	if err := database.InitDatabase(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create Gin router
	router := gin.New()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	taskHandler := handlers.NewTaskHandler()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Task Management API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(cfg), authHandler.GetProfile)
		}

		// Task routes (protected)
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware(cfg))
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/stats", taskHandler.GetTaskStats) // Must come before /:id
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}

	// Print available routes
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("🚀 Task Management API Server")
	log.Println(strings.Repeat("=", 60))
	log.Println("Server running on: http://localhost:" + cfg.Server.Port)
	log.Println("Mode:", cfg.Server.Mode)
	log.Println("\nAvailable endpoints:")
	log.Println("  GET    /health                    - Health check")
	log.Println("  POST   /api/v1/auth/register      - Register new user")
	log.Println("  POST   /api/v1/auth/login         - Login user")
	log.Println("  GET    /api/v1/auth/profile       - Get user profile (protected)")
	log.Println("  POST   /api/v1/tasks              - Create task (protected)")
	log.Println("  GET    /api/v1/tasks              - Get all tasks (protected)")
	log.Println("  GET    /api/v1/tasks/stats        - Get task statistics (protected)")
	log.Println("  GET    /api/v1/tasks/:id          - Get task by ID (protected)")
	log.Println("  PUT    /api/v1/tasks/:id          - Update task (protected)")
	log.Println("  DELETE /api/v1/tasks/:id          - Delete task (protected)")
	log.Println(strings.Repeat("=", 60) + "\n")

	// Start server
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
