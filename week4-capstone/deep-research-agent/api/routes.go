package api

import (
	"deep-research-agent/auth"
	"deep-research-agent/config"
	"deep-research-agent/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all API routes with authentication and rate limiting
func SetupRoutes(router *gin.Engine, handler *APIHandler, db *gorm.DB, cfg *config.Config, jwtManager *auth.JWTManager, rateLimiter *middleware.RateLimiter) {
	// Create auth handler
	authHandler := NewAuthHandler(db, jwtManager)

	// Health check endpoint (public)
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public authentication endpoints
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Protected authentication endpoints
		authProtected := v1.Group("/auth")
		authProtected.Use(auth.AuthMiddleware(jwtManager))
		if cfg.RateLimit.Enabled {
			authProtected.Use(rateLimiter.RateLimitMiddleware())
		}
		{
			authProtected.GET("/me", authHandler.GetProfile)
			authProtected.PUT("/profile", authHandler.UpdateProfile)
			authProtected.POST("/refresh", authHandler.RefreshToken)
			authProtected.POST("/password", authHandler.ChangePassword)
		}

		// Research endpoints (protected)
		research := v1.Group("/research")
		research.Use(auth.AuthMiddleware(jwtManager))
		if cfg.RateLimit.Enabled {
			research.Use(rateLimiter.RateLimitMiddleware())
		}
		{
			research.POST("/start", handler.StartResearch)
			research.GET("/:id/status", handler.GetResearchStatus)
			research.GET("/:id/result", handler.GetResearchResult)
			research.DELETE("/:id", handler.CancelResearch)
		}

		// WebSocket for research streaming (protected with query param token)
		wsRoutes := v1.Group("/research")
		wsRoutes.Use(auth.WebSocketAuthMiddleware(jwtManager))
		{
			wsRoutes.GET("/:id/stream", handler.StreamResearch)
		}

		// Document endpoints (protected)
		documents := v1.Group("/documents")
		documents.Use(auth.AuthMiddleware(jwtManager))
		if cfg.RateLimit.Enabled {
			documents.Use(rateLimiter.RateLimitMiddleware())
		}
		{
			documents.POST("/upload", handler.UploadDocument)
			documents.POST("/analyze", handler.AnalyzeDocument)
			documents.GET("", handler.ListDocuments)
			documents.GET("/:id", handler.GetDocument)
			documents.DELETE("/:id", handler.DeleteDocument)
		}

		// Session endpoints (protected)
		sessions := v1.Group("/sessions")
		sessions.Use(auth.AuthMiddleware(jwtManager))
		if cfg.RateLimit.Enabled {
			sessions.Use(rateLimiter.RateLimitMiddleware())
		}
		{
			sessions.GET("", handler.ListSessions)
			sessions.GET("/:id", handler.GetSession)
			sessions.DELETE("/:id", handler.DeleteSession)
			sessions.POST("/:id/export", handler.ExportSession)
		}

		// Export endpoints (protected)
		export := v1.Group("/export")
		export.Use(auth.AuthMiddleware(jwtManager))
		if cfg.RateLimit.Enabled {
			export.Use(rateLimiter.RateLimitMiddleware())
		}
		{
			export.GET("/:id/download", handler.DownloadExport)
		}

		// Statistics (protected)
		stats := v1.Group("/stats")
		stats.Use(auth.AuthMiddleware(jwtManager))
		if cfg.RateLimit.Enabled {
			stats.Use(rateLimiter.RateLimitMiddleware())
		}
		{
			stats.GET("", handler.GetStats)
		}
	}
}

// RouteInfo contains information about an API route
type RouteInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Protected   bool   `json:"protected"`
}

// GetRoutesList returns a list of all available routes
func GetRoutesList() []RouteInfo {
	return []RouteInfo{
		// Health
		{Method: "GET", Path: "/health", Description: "Health check", Protected: false},

		// Auth - Public
		{Method: "POST", Path: "/api/v1/auth/register", Description: "Register new user", Protected: false},
		{Method: "POST", Path: "/api/v1/auth/login", Description: "Login user", Protected: false},

		// Auth - Protected
		{Method: "GET", Path: "/api/v1/auth/me", Description: "Get current user profile", Protected: true},
		{Method: "PUT", Path: "/api/v1/auth/profile", Description: "Update user profile", Protected: true},
		{Method: "POST", Path: "/api/v1/auth/refresh", Description: "Refresh authentication token", Protected: true},
		{Method: "POST", Path: "/api/v1/auth/password", Description: "Change password", Protected: true},

		// Research - Protected
		{Method: "POST", Path: "/api/v1/research/start", Description: "Start new research session", Protected: true},
		{Method: "GET", Path: "/api/v1/research/:id/status", Description: "Get research status", Protected: true},
		{Method: "GET", Path: "/api/v1/research/:id/result", Description: "Get research results", Protected: true},
		{Method: "DELETE", Path: "/api/v1/research/:id", Description: "Cancel research", Protected: true},
		{Method: "GET", Path: "/api/v1/research/:id/stream", Description: "WebSocket stream for research progress", Protected: true},

		// Documents - Protected
		{Method: "POST", Path: "/api/v1/documents/upload", Description: "Upload document", Protected: true},
		{Method: "POST", Path: "/api/v1/documents/analyze", Description: "Analyze document", Protected: true},
		{Method: "GET", Path: "/api/v1/documents", Description: "List user's documents", Protected: true},
		{Method: "GET", Path: "/api/v1/documents/:id", Description: "Get document details", Protected: true},
		{Method: "DELETE", Path: "/api/v1/documents/:id", Description: "Delete document", Protected: true},

		// Sessions - Protected
		{Method: "GET", Path: "/api/v1/sessions", Description: "List user's research sessions", Protected: true},
		{Method: "GET", Path: "/api/v1/sessions/:id", Description: "Get session details", Protected: true},
		{Method: "DELETE", Path: "/api/v1/sessions/:id", Description: "Delete session", Protected: true},
		{Method: "POST", Path: "/api/v1/sessions/:id/export", Description: "Export session", Protected: true},

		// Export - Protected
		{Method: "GET", Path: "/api/v1/export/:id/download", Description: "Download export", Protected: true},

		// Stats - Protected
		{Method: "GET", Path: "/api/v1/stats", Description: "Get user statistics", Protected: true},
	}
}
