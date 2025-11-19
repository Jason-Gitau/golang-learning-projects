package middleware

import (
	"net/http"
	"sync"
	"time"

	"deep-research-agent/auth"
	"deep-research-agent/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MaxRequestsPerHour int           // Maximum requests per hour per user
	MaxConcurrentJobs  int           // Maximum concurrent jobs per user
	WindowDuration     time.Duration // Rate limit window duration
	CleanupInterval    time.Duration // How often to cleanup old rate limit records
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		MaxRequestsPerHour: 100,
		MaxConcurrentJobs:  10,
		WindowDuration:     1 * time.Hour,
		CleanupInterval:    10 * time.Minute,
	}
}

// RateLimiter handles rate limiting logic
type RateLimiter struct {
	db     *gorm.DB
	config *RateLimitConfig
	mu     sync.RWMutex
	cache  map[string]*models.RateLimit // In-memory cache for performance
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(db *gorm.DB, config *RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		db:     db,
		config: config,
		cache:  make(map[string]*models.RateLimit),
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (must be called after auth middleware)
		userID, exists := auth.GetUserIDFromContext(c)
		if !exists {
			// If no user authenticated, skip rate limiting (for public endpoints)
			c.Next()
			return
		}

		// Check rate limit
		allowed, remaining, resetTime, err := rl.checkRateLimit(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check rate limit",
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", string(rune(rl.config.MaxRequestsPerHour)))
		c.Header("X-RateLimit-Remaining", string(rune(remaining)))
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":      "Rate limit exceeded",
				"reset_at":   resetTime.Format(time.RFC3339),
				"retry_after": int(time.Until(resetTime).Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if a user has exceeded their rate limit
func (rl *RateLimiter) checkRateLimit(userID string) (allowed bool, remaining int, resetTime time.Time, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Try to get from cache first
	limit, exists := rl.cache[userID]

	if !exists || time.Since(limit.WindowStart) >= rl.config.WindowDuration {
		// Create new window or load from database
		var dbLimit models.RateLimit
		err := rl.db.Where("user_id = ?", userID).First(&dbLimit).Error

		if err == gorm.ErrRecordNotFound || time.Since(dbLimit.WindowStart) >= rl.config.WindowDuration {
			// Create new rate limit window
			limit = &models.RateLimit{
				ID:            uuid.New().String(),
				UserID:        userID,
				RequestCount:  0,
				WindowStart:   now,
				ActiveJobs:    0,
				LastRequestAt: now,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
		} else if err != nil {
			return false, 0, now, err
		} else {
			limit = &dbLimit
		}

		rl.cache[userID] = limit
	}

	// Check if limit exceeded
	if limit.RequestCount >= rl.config.MaxRequestsPerHour {
		resetTime = limit.WindowStart.Add(rl.config.WindowDuration)
		return false, 0, resetTime, nil
	}

	// Increment request count
	limit.RequestCount++
	limit.LastRequestAt = now
	limit.UpdatedAt = now

	// Update database (async to avoid blocking)
	go func() {
		dbLimit := *limit
		if err := rl.db.Save(&dbLimit).Error; err != nil {
			// Log error but don't fail the request
			// In production, use proper logging
		}
	}()

	remaining = rl.config.MaxRequestsPerHour - limit.RequestCount
	resetTime = limit.WindowStart.Add(rl.config.WindowDuration)

	return true, remaining, resetTime, nil
}

// CheckConcurrentJobs checks if a user can start a new job
func (rl *RateLimiter) CheckConcurrentJobs(userID string) (bool, error) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	limit, exists := rl.cache[userID]
	if !exists {
		// Load from database
		var dbLimit models.RateLimit
		err := rl.db.Where("user_id = ?", userID).First(&dbLimit).Error
		if err == gorm.ErrRecordNotFound {
			return true, nil // No limits yet, allow
		} else if err != nil {
			return false, err
		}
		limit = &dbLimit
	}

	return limit.ActiveJobs < rl.config.MaxConcurrentJobs, nil
}

// IncrementActiveJobs increments the active job count for a user
func (rl *RateLimiter) IncrementActiveJobs(userID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit, exists := rl.cache[userID]
	if !exists {
		var dbLimit models.RateLimit
		err := rl.db.Where("user_id = ?", userID).First(&dbLimit).Error
		if err == gorm.ErrRecordNotFound {
			// Create new limit record
			limit = &models.RateLimit{
				ID:            uuid.New().String(),
				UserID:        userID,
				RequestCount:  0,
				WindowStart:   time.Now(),
				ActiveJobs:    0,
				LastRequestAt: time.Now(),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
		} else if err != nil {
			return err
		} else {
			limit = &dbLimit
		}
		rl.cache[userID] = limit
	}

	limit.ActiveJobs++
	limit.UpdatedAt = time.Now()

	return rl.db.Save(limit).Error
}

// DecrementActiveJobs decrements the active job count for a user
func (rl *RateLimiter) DecrementActiveJobs(userID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit, exists := rl.cache[userID]
	if !exists {
		var dbLimit models.RateLimit
		err := rl.db.Where("user_id = ?", userID).First(&dbLimit).Error
		if err != nil {
			return err
		}
		limit = &dbLimit
		rl.cache[userID] = limit
	}

	if limit.ActiveJobs > 0 {
		limit.ActiveJobs--
		limit.UpdatedAt = time.Now()
		return rl.db.Save(limit).Error
	}

	return nil
}

// GetUserStats returns rate limit statistics for a user
func (rl *RateLimiter) GetUserStats(userID string) (*models.RateLimit, error) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	limit, exists := rl.cache[userID]
	if exists {
		return limit, nil
	}

	var dbLimit models.RateLimit
	err := rl.db.Where("user_id = ?", userID).First(&dbLimit).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &dbLimit, err
}

// cleanupLoop periodically cleans up old rate limit records
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup removes expired rate limit records
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.config.WindowDuration * 2)

	// Clean cache
	for userID, limit := range rl.cache {
		if limit.WindowStart.Before(cutoff) {
			delete(rl.cache, userID)
		}
	}

	// Clean database (async)
	go func() {
		rl.db.Where("window_start < ?", cutoff).Delete(&models.RateLimit{})
	}()
}

// ResetUserLimit resets the rate limit for a specific user (admin function)
func (rl *RateLimiter) ResetUserLimit(userID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.cache, userID)
	return rl.db.Where("user_id = ?", userID).Delete(&models.RateLimit{}).Error
}
