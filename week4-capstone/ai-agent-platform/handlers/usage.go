package handlers

import (
	"ai-agent-platform/database"
	"ai-agent-platform/middleware"
	"ai-agent-platform/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UsageHandler struct{}

func NewUsageHandler() *UsageHandler {
	return &UsageHandler{}
}

// GetUsageStats returns usage statistics for the authenticated user
// GET /api/v1/usage
func (h *UsageHandler) GetUsageStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req models.UsageStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Set default date range if not provided (last 30 days)
	var startDate, endDate time.Time
	var err error

	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid start_date format. Use YYYY-MM-DD",
			})
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid end_date format. Use YYYY-MM-DD",
			})
			return
		}
		// Set to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		endDate = time.Now()
	}

	// Build query
	query := db.Model(&models.UsageLog{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate)

	// Filter by agent if provided
	if req.AgentID != nil {
		query = query.Where("agent_id = ?", *req.AgentID)
	}

	// Get total statistics
	var totalLogs []models.UsageLog
	if err := query.Find(&totalLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch usage logs",
		})
		return
	}

	// Calculate aggregated statistics
	stats := models.UsageStats{
		PeriodStart: startDate.Format("2006-01-02"),
		PeriodEnd:   endDate.Format("2006-01-02"),
		ByModel:     make(map[string]models.ModelUsage),
		ByAgent:     make(map[string]models.AgentUsage),
	}

	for _, log := range totalLogs {
		stats.TotalRequests++
		stats.TotalTokens += log.TotalTokens
		stats.InputTokens += log.InputTokens
		stats.OutputTokens += log.OutputTokens
		stats.TotalCost += log.Cost

		// Aggregate by model
		modelUsage := stats.ByModel[log.Model]
		modelUsage.Requests++
		modelUsage.TotalTokens += log.TotalTokens
		modelUsage.InputTokens += log.InputTokens
		modelUsage.OutputTokens += log.OutputTokens
		modelUsage.Cost += log.Cost
		stats.ByModel[log.Model] = modelUsage

		// Aggregate by agent if available
		if log.AgentID != nil {
			agentKey := string(rune(*log.AgentID))
			agentUsage := stats.ByAgent[agentKey]
			agentUsage.AgentID = *log.AgentID
			agentUsage.Requests++
			agentUsage.TotalTokens += log.TotalTokens
			agentUsage.InputTokens += log.InputTokens
			agentUsage.OutputTokens += log.OutputTokens
			agentUsage.Cost += log.Cost
			stats.ByAgent[agentKey] = agentUsage
		}
	}

	// Fetch agent names for the ByAgent stats
	if len(stats.ByAgent) > 0 {
		var agents []models.Agent
		var agentIDs []uint
		for _, usage := range stats.ByAgent {
			agentIDs = append(agentIDs, usage.AgentID)
		}

		if err := db.Where("id IN ?", agentIDs).Find(&agents).Error; err == nil {
			for _, agent := range agents {
				agentKey := string(rune(agent.ID))
				if usage, exists := stats.ByAgent[agentKey]; exists {
					usage.AgentName = agent.Name
					stats.ByAgent[agentKey] = usage
				}
			}
		}
	}

	c.JSON(http.StatusOK, stats)
}

// GetRateLimitInfo returns rate limit information for the authenticated user
// GET /api/v1/usage/rate-limit
func (h *UsageHandler) GetRateLimitInfo(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	db := database.GetDB()

	var rateLimit models.RateLimit
	if err := db.Where("user_id = ?", userID).First(&rateLimit).Error; err != nil {
		// No rate limit record yet, return default
		c.JSON(http.StatusOK, models.RateLimitInfo{
			Limit:     100, // Default from config
			Remaining: 100,
			Reset:     time.Now().Add(1 * time.Hour),
		})
		return
	}

	// Check if window has expired
	now := time.Now()
	if now.After(rateLimit.WindowEnd) {
		// Window expired, reset
		c.JSON(http.StatusOK, models.RateLimitInfo{
			Limit:     100,
			Remaining: 100,
			Reset:     now.Add(1 * time.Hour),
		})
		return
	}

	remaining := 100 - rateLimit.RequestsCount
	if remaining < 0 {
		remaining = 0
	}

	c.JSON(http.StatusOK, models.RateLimitInfo{
		Limit:     100,
		Remaining: remaining,
		Reset:     rateLimit.WindowEnd,
	})
}

// CreateUsageLog creates a new usage log entry
// POST /api/v1/usage/log
// This is typically called internally by the chat system
func (h *UsageHandler) CreateUsageLog(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req struct {
		InputTokens  int    `json:"input_tokens" binding:"required,gte=0"`
		OutputTokens int    `json:"output_tokens" binding:"required,gte=0"`
		Model        string `json:"model" binding:"required"`
		Operation    string `json:"operation" binding:"required"`
		AgentID      *uint  `json:"agent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Calculate cost (this would use the actual model pricing)
	cost := float64(req.InputTokens+req.OutputTokens) * 0.001 // Mock pricing

	usageLog := models.UsageLog{
		UserID:       &userID,
		AgentID:      req.AgentID,
		InputTokens:  req.InputTokens,
		OutputTokens: req.OutputTokens,
		TotalTokens:  req.InputTokens + req.OutputTokens,
		Cost:         cost,
		Model:        req.Model,
		Operation:    req.Operation,
	}

	if err := db.Create(&usageLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create usage log",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usage log created successfully",
		"log_id":  usageLog.ID,
	})
}
