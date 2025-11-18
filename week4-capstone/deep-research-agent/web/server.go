package web

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ServerConfig holds the configuration for the web server
type ServerConfig struct {
	StaticDir string
	IndexFile string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *ServerConfig {
	return &ServerConfig{
		StaticDir: "./web/static",
		IndexFile: "./web/static/index.html",
	}
}

// SetupRoutes configures the web routes for serving static files
func SetupRoutes(router *gin.Engine) {
	SetupRoutesWithConfig(router, DefaultConfig())
}

// SetupRoutesWithConfig configures the web routes with custom configuration
func SetupRoutesWithConfig(router *gin.Engine, config *ServerConfig) {
	// Serve static files (CSS, JS, etc.)
	router.Static("/static", config.StaticDir)

	// Serve favicon if it exists
	faviconPath := filepath.Join(config.StaticDir, "favicon.ico")
	router.StaticFile("/favicon.ico", faviconPath)

	// Serve index.html for the root path
	router.GET("/", func(c *gin.Context) {
		c.File(config.IndexFile)
	})

	// SPA fallback: serve index.html for any non-API routes
	// This allows client-side routing to work
	router.NoRoute(func(c *gin.Context) {
		// Only serve index.html for non-API routes
		if c.Request.URL.Path[:4] != "/api" && c.Request.URL.Path[:7] != "/static" {
			c.File(config.IndexFile)
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Route not found",
			})
		}
	})
}

// ServeStaticFiles is a convenience function that sets up all necessary routes
// for serving the web UI. It should be called when setting up your API router.
//
// Example usage:
//
//	router := gin.Default()
//
//	// Setup API routes
//	api := router.Group("/api/v1")
//	{
//	    // Your API routes here
//	}
//
//	// Setup web UI routes (should be last)
//	web.ServeStaticFiles(router)
//
func ServeStaticFiles(router *gin.Engine) {
	SetupRoutes(router)
}

// CORS middleware for development
// In production, configure CORS properly based on your domain
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
