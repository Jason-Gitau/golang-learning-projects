package middleware

import (
	"ai-agent-platform/database"
	"ai-agent-platform/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware implements rate limiting per user
func RateLimitMiddleware(requestsPerHour int, windowSize time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := GetUserID(c)
		if !exists {
			// If no user ID, skip rate limiting (for public endpoints)
			c.Next()
			return
		}

		// Check and update rate limit
		allowed, remaining, reset, err := checkRateLimit(userID, requestsPerHour, windowSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check rate limit",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerHour))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", reset.Unix()))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":     "Rate limit exceeded",
				"limit":     requestsPerHour,
				"remaining": 0,
				"reset":     reset.Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if a user has exceeded their rate limit
func checkRateLimit(userID uint, limit int, windowSize time.Duration) (allowed bool, remaining int, reset time.Time, err error) {
	db := database.GetDB()
	now := time.Now()

	var rateLimit models.RateLimit

	// Find or create rate limit record
	result := db.Where("user_id = ?", userID).First(&rateLimit)
	if result.Error != nil {
		// Create new rate limit record
		rateLimit = models.RateLimit{
			UserID:        userID,
			RequestsCount: 0,
			WindowStart:   now,
			WindowEnd:     now.Add(windowSize),
		}
	}

	// Check if window has expired
	if now.After(rateLimit.WindowEnd) {
		// Reset the window
		rateLimit.RequestsCount = 0
		rateLimit.WindowStart = now
		rateLimit.WindowEnd = now.Add(windowSize)
	}

	// Check if limit exceeded
	if rateLimit.RequestsCount >= limit {
		remaining = 0
		reset = rateLimit.WindowEnd
		return false, remaining, reset, nil
	}

	// Increment request count
	rateLimit.RequestsCount++

	// Save to database
	if result.Error != nil {
		// Create new record
		if err := db.Create(&rateLimit).Error; err != nil {
			return false, 0, time.Time{}, err
		}
	} else {
		// Update existing record
		if err := db.Save(&rateLimit).Error; err != nil {
			return false, 0, time.Time{}, err
		}
	}

	remaining = limit - rateLimit.RequestsCount
	reset = rateLimit.WindowEnd

	return true, remaining, reset, nil
}
