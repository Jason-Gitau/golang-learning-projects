package ratelimiter

import (
	"context"
	"time"
)

// RateLimiter controls the rate of operations using a token bucket algorithm
type RateLimiter struct {
	ticker   *time.Ticker
	tokens   chan struct{}
	rate     float64
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewRateLimiter creates a new rate limiter
// rate is requests per second (e.g., 2.0 = 2 requests per second)
func NewRateLimiter(ctx context.Context, rate float64) *RateLimiter {
	if rate <= 0 {
		rate = 1.0
	}

	// Calculate interval between tokens
	interval := time.Duration(float64(time.Second) / rate)

	limiterCtx, cancel := context.WithCancel(ctx)

	rl := &RateLimiter{
		ticker:   time.NewTicker(interval),
		tokens:   make(chan struct{}, int(rate)+1), // Buffer for burst capacity
		rate:     rate,
		interval: interval,
		ctx:      limiterCtx,
		cancel:   cancel,
	}

	// Pre-fill tokens for immediate start
	for i := 0; i < int(rate)+1; i++ {
		select {
		case rl.tokens <- struct{}{}:
		default:
			break
		}
	}

	// Start token generation
	go rl.generateTokens()

	return rl
}

// generateTokens continuously generates tokens at the specified rate
func (rl *RateLimiter) generateTokens() {
	for {
		select {
		case <-rl.ctx.Done():
			rl.ticker.Stop()
			return
		case <-rl.ticker.C:
			// Try to add a token (non-blocking)
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Token bucket is full, skip
			}
		}
	}
}

// Wait blocks until a token is available or context is cancelled
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.ctx.Done():
		return rl.ctx.Err()
	case <-rl.tokens:
		return nil
	}
}

// TryWait attempts to acquire a token without blocking
// Returns true if a token was acquired, false otherwise
func (rl *RateLimiter) TryWait() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Stop stops the rate limiter and releases resources
func (rl *RateLimiter) Stop() {
	rl.cancel()
}

// GetRate returns the current rate limit (requests per second)
func (rl *RateLimiter) GetRate() float64 {
	return rl.rate
}

// GetInterval returns the interval between tokens
func (rl *RateLimiter) GetInterval() time.Duration {
	return rl.interval
}
