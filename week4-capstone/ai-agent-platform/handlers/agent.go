package handlers

import (
	"ai-agent-platform/database"
	"ai-agent-platform/middleware"
	"ai-agent-platform/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct{}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{}
}

// CreateAgent creates a new agent
// POST /api/v1/agents
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req models.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set default values if not provided
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2000
	}

	db := database.GetDB()

	agent := models.Agent{
		Name:         req.Name,
		Description:  req.Description,
		SystemPrompt: req.SystemPrompt,
		Model:        req.Model,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
		Status:       models.AgentStatusActive,
		Tools:        req.Tools,
		UserID:       userID,
	}

	if err := db.Create(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create agent",
		})
		return
	}

	c.JSON(http.StatusCreated, agent.ToResponse())
}

// ListAgents returns all agents for the authenticated user
// GET /api/v1/agents
func (h *AgentHandler) ListAgents(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	db := database.GetDB()

	var agents []models.Agent
	query := db.Where("user_id = ?", userID)

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&agents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch agents",
		})
		return
	}

	// Convert to response format
	responses := make([]models.AgentResponse, len(agents))
	for i, agent := range agents {
		responses[i] = agent.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": responses,
		"count":  len(responses),
	})
}

// GetAgent returns a specific agent by ID
// GET /api/v1/agents/:id
func (h *AgentHandler) GetAgent(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	agentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid agent ID",
		})
		return
	}

	db := database.GetDB()

	var agent models.Agent
	if err := db.Where("id = ? AND user_id = ?", agentID, userID).First(&agent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Agent not found",
		})
		return
	}

	c.JSON(http.StatusOK, agent.ToResponse())
}

// UpdateAgent updates an existing agent
// PUT /api/v1/agents/:id
func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	agentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid agent ID",
		})
		return
	}

	var req models.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	db := database.GetDB()

	var agent models.Agent
	if err := db.Where("id = ? AND user_id = ?", agentID, userID).First(&agent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Agent not found",
		})
		return
	}

	// Update fields if provided
	if req.Name != nil {
		agent.Name = *req.Name
	}
	if req.Description != nil {
		agent.Description = *req.Description
	}
	if req.SystemPrompt != nil {
		agent.SystemPrompt = *req.SystemPrompt
	}
	if req.Model != nil {
		agent.Model = *req.Model
	}
	if req.Temperature != nil {
		agent.Temperature = *req.Temperature
	}
	if req.MaxTokens != nil {
		agent.MaxTokens = *req.MaxTokens
	}
	if req.Status != nil {
		agent.Status = *req.Status
	}
	if req.Tools != nil {
		agent.Tools = *req.Tools
	}

	if err := db.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update agent",
		})
		return
	}

	c.JSON(http.StatusOK, agent.ToResponse())
}

// DeleteAgent deletes an agent
// DELETE /api/v1/agents/:id
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	agentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid agent ID",
		})
		return
	}

	db := database.GetDB()

	var agent models.Agent
	if err := db.Where("id = ? AND user_id = ?", agentID, userID).First(&agent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Agent not found",
		})
		return
	}

	// Soft delete
	if err := db.Delete(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete agent",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Agent deleted successfully",
	})
}
