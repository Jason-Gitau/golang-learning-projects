package models

import (
	"sync"
	"time"
)

// AgentStatus represents the current status of the research agent
type AgentStatus string

const (
	AgentIdle       AgentStatus = "idle"
	AgentPlanning   AgentStatus = "planning"
	AgentResearching AgentStatus = "researching"
	AgentAggregating AgentStatus = "aggregating"
	AgentCompleted  AgentStatus = "completed"
	AgentFailed     AgentStatus = "failed"
)

// AgentConfig contains configuration for the research agent
type AgentConfig struct {
	MaxSteps        int           `json:"max_steps"`         // Maximum research steps
	Timeout         time.Duration `json:"timeout"`           // Overall timeout
	StepTimeout     time.Duration `json:"step_timeout"`      // Per-step timeout
	MaxSources      int           `json:"max_sources"`       // Maximum sources per step
	ConcurrentTools int           `json:"concurrent_tools"`  // Max concurrent tool executions
	RetryAttempts   int           `json:"retry_attempts"`    // Retry failed steps
	RetryDelay      time.Duration `json:"retry_delay"`       // Delay between retries
	EnableCaching   bool          `json:"enable_caching"`    // Cache tool results
	CacheTTL        time.Duration `json:"cache_ttl"`         // Cache time-to-live
	MinConfidence   float64       `json:"min_confidence"`    // Minimum confidence threshold
}

// DefaultAgentConfig returns sensible default configuration
func DefaultAgentConfig() AgentConfig {
	return AgentConfig{
		MaxSteps:        10,
		Timeout:         5 * time.Minute,
		StepTimeout:     30 * time.Second,
		MaxSources:      10,
		ConcurrentTools: 3,
		RetryAttempts:   2,
		RetryDelay:      2 * time.Second,
		EnableCaching:   true,
		CacheTTL:        30 * time.Minute,
		MinConfidence:   0.3,
	}
}

// AgentMetrics tracks agent performance metrics
type AgentMetrics struct {
	TotalResearches     int           `json:"total_researches"`
	SuccessfulResearches int          `json:"successful_researches"`
	FailedResearches    int           `json:"failed_researches"`
	AverageSteps        float64       `json:"average_steps"`
	AverageDuration     time.Duration `json:"average_duration"`
	TotalStepsExecuted  int           `json:"total_steps_executed"`
	ToolUsageCount      map[string]int `json:"tool_usage_count"`
	CacheHits           int           `json:"cache_hits"`
	CacheMisses         int           `json:"cache_misses"`
	LastUpdated         time.Time     `json:"last_updated"`
	mu                  sync.RWMutex
}

// NewAgentMetrics creates a new metrics tracker
func NewAgentMetrics() *AgentMetrics {
	return &AgentMetrics{
		ToolUsageCount: make(map[string]int),
		LastUpdated:    time.Now(),
	}
}

// RecordResearch records metrics for a completed research
func (m *AgentMetrics) RecordResearch(result *ResearchResult, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalResearches++
	if success {
		m.SuccessfulResearches++
	} else {
		m.FailedResearches++
	}

	m.TotalStepsExecuted += result.TotalSteps
	m.AverageSteps = float64(m.TotalStepsExecuted) / float64(m.TotalResearches)

	// Update average duration
	totalDuration := m.AverageDuration * time.Duration(m.TotalResearches-1)
	m.AverageDuration = (totalDuration + result.Duration) / time.Duration(m.TotalResearches)

	// Update tool usage
	for _, step := range result.Steps {
		if step.Success {
			m.ToolUsageCount[step.Step.Tool]++
		}
	}

	m.LastUpdated = time.Now()
}

// RecordCacheHit records a cache hit
func (m *AgentMetrics) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheHits++
}

// RecordCacheMiss records a cache miss
func (m *AgentMetrics) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheMisses++
}

// GetMetrics returns a copy of current metrics
func (m *AgentMetrics) GetMetrics() AgentMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Deep copy tool usage
	toolUsage := make(map[string]int)
	for k, v := range m.ToolUsageCount {
		toolUsage[k] = v
	}

	return AgentMetrics{
		TotalResearches:      m.TotalResearches,
		SuccessfulResearches: m.SuccessfulResearches,
		FailedResearches:     m.FailedResearches,
		AverageSteps:         m.AverageSteps,
		AverageDuration:      m.AverageDuration,
		TotalStepsExecuted:   m.TotalStepsExecuted,
		ToolUsageCount:       toolUsage,
		CacheHits:            m.CacheHits,
		CacheMisses:          m.CacheMisses,
		LastUpdated:          m.LastUpdated,
	}
}

// ResearchOptions contains options for research execution
type ResearchOptions struct {
	SessionID       string        `json:"session_id"`       // Optional session ID
	SaveSession     bool          `json:"save_session"`     // Save to database
	StreamProgress  bool          `json:"stream_progress"`  // Stream progress updates
	ProgressChannel chan ResearchProgress `json:"-"`       // Channel for progress updates
	Timeout         time.Duration `json:"timeout"`          // Override default timeout
	MaxSteps        int           `json:"max_steps"`        // Override max steps
}

// ToolResult represents the result from a tool execution
type ToolResult struct {
	ToolName  string        `json:"tool_name"`
	Success   bool          `json:"success"`
	Data      interface{}   `json:"data"`
	Error     error         `json:"error"`
	Duration  time.Duration `json:"duration"`
	Sources   []Source      `json:"sources"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// CacheEntry represents a cached tool result
type CacheEntry struct {
	Key        string     `json:"key"`
	Result     ToolResult `json:"result"`
	Timestamp  time.Time  `json:"timestamp"`
	TTL        time.Duration `json:"ttl"`
}

// IsExpired checks if the cache entry has expired
func (c *CacheEntry) IsExpired() bool {
	return time.Since(c.Timestamp) > c.TTL
}
