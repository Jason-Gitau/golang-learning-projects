package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get request details
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Get user ID if available
		userID, exists := GetUserID(c)
		userIDStr := "anonymous"
		if exists {
			userIDStr = string(rune(userID))
		}

		// Log request
		log.Printf("[%s] %s %s | Status: %d | Latency: %v | IP: %s | User: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			latency,
			clientIP,
			userIDStr,
		)
	}
}
