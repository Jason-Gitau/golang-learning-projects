package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"deep-research-agent/agent"
	"deep-research-agent/auth"
	"deep-research-agent/config"
	"deep-research-agent/jobs"
	"deep-research-agent/middleware"
	"deep-research-agent/models"
	"deep-research-agent/storage"
	"deep-research-agent/tools"
	"deep-research-agent/uploads"
	ws "deep-research-agent/websocket"
)

// Server represents the API server
type Server struct {
	router        *gin.Engine
	httpServer    *http.Server
	wsHub         *ws.Hub
	jobQueue      *jobs.Queue
	db            *gorm.DB
	researchAgent *agent.ResearchAgent
	uploadHandler *uploads.UploadHandler
	jwtManager    *auth.JWTManager
	rateLimiter   *middleware.RateLimiter
	config        ServerConfig
	appConfig     *config.Config
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port          int
	WorkerCount   int
	DBPath        string
	UploadDir     string
	ExportDir     string
	Debug         bool
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:        8080,
		WorkerCount: 10,
		DBPath:      "data/research.db",
		UploadDir:   "uploads/files",
		ExportDir:   "exports",
		Debug:       false,
	}
}

// NewServer creates a new API server
func NewServer(serverConfig ServerConfig) (*Server, error) {
	// Set Gin mode
	if !serverConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create necessary directories
	if err := os.MkdirAll(serverConfig.UploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}
	if err := os.MkdirAll(serverConfig.ExportDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create export directory: %w", err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Load application configuration
	appConfig, err := config.LoadOrCreateConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Validate JWT secret
	if appConfig.Auth.EnableAuth && appConfig.Auth.JWTSecret == "" {
		log.Println("WARNING: JWT_SECRET not set. Using default (INSECURE - DO NOT USE IN PRODUCTION)")
		appConfig.Auth.JWTSecret = "insecure-default-secret-change-me"
	}

	// Initialize database
	db, err := initDatabase(serverConfig.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize storage
	researchDB, err := storage.NewResearchDB(serverConfig.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		appConfig.Auth.JWTSecret,
		appConfig.Auth.TokenExpiry,
		appConfig.Auth.RefreshExpiry,
	)

	// Initialize rate limiter
	rateLimitConfig := &middleware.RateLimitConfig{
		MaxRequestsPerHour: appConfig.RateLimit.MaxRequestsPerHour,
		MaxConcurrentJobs:  appConfig.RateLimit.MaxConcurrentJobs,
		WindowDuration:     appConfig.RateLimit.WindowDuration,
		CleanupInterval:    appConfig.RateLimit.CleanupInterval,
	}
	rateLimiter := middleware.NewRateLimiter(db, rateLimitConfig)

	// Initialize tool registry
	toolRegistry := tools.NewToolRegistry()
	registerAllTools(toolRegistry)

	// Initialize research agent
	agentConfig := models.AgentConfig{
		MaxSteps: 10,
		Timeout:  5 * time.Minute,
	}
	researchAgent := agent.NewResearchAgent(agentConfig, toolRegistry, researchDB)

	// Initialize WebSocket hub
	wsHub := ws.NewHub()

	// Initialize job queue
	jobQueue := jobs.NewQueue(serverConfig.WorkerCount, db, wsHub, researchAgent)

	// Initialize upload handler
	uploadHandler, err := uploads.NewUploadHandler(db, researchAgent, serverConfig.UploadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize upload handler: %w", err)
	}

	// Create Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware())
	router.Use(auth.CORSMiddleware()) // Use auth CORS middleware
	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestIDMiddleware())

	// Create API handler
	apiHandler := NewAPIHandler(jobQueue, uploadHandler, researchAgent, researchDB, wsHub)

	// Setup routes with authentication
	SetupRoutes(router, apiHandler, db, appConfig, jwtManager, rateLimiter)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		router:        router,
		httpServer:    httpServer,
		wsHub:         wsHub,
		jobQueue:      jobQueue,
		db:            db,
		researchAgent: researchAgent,
		uploadHandler: uploadHandler,
		jwtManager:    jwtManager,
		rateLimiter:   rateLimiter,
		config:        serverConfig,
		appConfig:     appConfig,
	}, nil
}

// Start starts the API server
func (s *Server) Start() error {
	// Start WebSocket hub
	go s.wsHub.Run()
	log.Println("WebSocket hub started")

	// Start job queue workers
	if err := s.jobQueue.Start(s.config.WorkerCount); err != nil {
		return fmt.Errorf("failed to start job queue: %w", err)
	}

	// Start HTTP server
	log.Printf("Starting API server on port %d", s.config.Port)
	log.Printf("Server running at http://localhost:%d", s.config.Port)
	log.Printf("API documentation available at http://localhost:%d/health", s.config.Port)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down API server...")

	// Stop job queue
	s.jobQueue.Stop()

	// Shutdown WebSocket hub
	s.wsHub.Shutdown()

	// Close database
	if s.researchAgent != nil {
		s.researchAgent.Close()
	}

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	log.Println("API server shutdown complete")
	return nil
}

// initDatabase initializes the database with all required tables
func initDatabase(dbPath string) (*gorm.DB, error) {
	// Open database directly with GORM
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.ResearchSession{},
		&models.Document{},
		&models.ResearchJob{},
		&models.UploadedFile{},
		&models.ExportJob{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// registerAllTools registers all available tools
func registerAllTools(registry *tools.ToolRegistry) {
	// Register all tools from the tools package
	tools.RegisterAllTools(registry)
	log.Printf("Registered %d tools", len(registry.ListToolNames()))
}
