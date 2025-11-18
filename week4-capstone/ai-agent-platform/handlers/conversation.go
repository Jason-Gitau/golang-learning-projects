package handlers

import (
	"ai-agent-platform/database"
	"ai-agent-platform/middleware"
	"ai-agent-platform/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct{}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{}
}

// CreateConversation creates a new conversation for an agent
// POST /api/v1/agents/:id/conversations
func (h *ConversationHandler) CreateConversation(c *gin.Context) {
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

	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Verify agent exists and belongs to user
	var agent models.Agent
	if err := db.Where("id = ? AND user_id = ?", agentID, userID).First(&agent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Agent not found",
		})
		return
	}

	// Create conversation
	conversation := models.Conversation{
		Title:   req.Title,
		AgentID: uint(agentID),
		UserID:  userID,
		Status:  models.ConversationStatusActive,
	}

	if err := db.Create(&conversation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create conversation",
		})
		return
	}

	c.JSON(http.StatusCreated, conversation.ToResponse())
}

// ListConversations returns all conversations for the authenticated user
// GET /api/v1/conversations
func (h *ConversationHandler) ListConversations(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	db := database.GetDB()

	var conversations []models.Conversation
	query := db.Where("user_id = ?", userID)

	// Filter by agent if provided
	if agentID := c.Query("agent_id"); agentID != "" {
		query = query.Where("agent_id = ?", agentID)
	}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("updated_at DESC").Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch conversations",
		})
		return
	}

	// Convert to response format
	responses := make([]models.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		responses[i] = conv.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": responses,
		"count":         len(responses),
	})
}

// GetConversation returns a specific conversation with its messages
// GET /api/v1/conversations/:id
func (h *ConversationHandler) GetConversation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid conversation ID",
		})
		return
	}

	db := database.GetDB()

	var conversation models.Conversation
	if err := db.Preload("Messages").Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Conversation not found",
		})
		return
	}

	c.JSON(http.StatusOK, conversation.ToResponseWithMessages())
}

// DeleteConversation deletes a conversation
// DELETE /api/v1/conversations/:id
func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid conversation ID",
		})
		return
	}

	db := database.GetDB()

	var conversation models.Conversation
	if err := db.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Conversation not found",
		})
		return
	}

	// Soft delete
	if err := db.Delete(&conversation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete conversation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Conversation deleted successfully",
	})
}

// CreateMessage creates a new message in a conversation
// POST /api/v1/conversations/:id/messages
func (h *ConversationHandler) CreateMessage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid conversation ID",
		})
		return
	}

	var req models.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	db := database.GetDB()

	// Verify conversation exists and belongs to user
	var conversation models.Conversation
	if err := db.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Conversation not found",
		})
		return
	}

	// Create message
	message := models.Message{
		ConversationID: uint(conversationID),
		Role:           req.Role,
		Content:        req.Content,
		TokensUsed:     0, // Will be calculated by the real-time agent
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create message",
		})
		return
	}

	c.JSON(http.StatusCreated, message.ToResponse())
}

// GetMessages returns all messages for a conversation
// GET /api/v1/conversations/:id/messages
func (h *ConversationHandler) GetMessages(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid conversation ID",
		})
		return
	}

	db := database.GetDB()

	// Verify conversation exists and belongs to user
	var conversation models.Conversation
	if err := db.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Conversation not found",
		})
		return
	}

	var messages []models.Message
	if err := db.Where("conversation_id = ?", conversationID).Order("created_at ASC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch messages",
		})
		return
	}

	// Convert to response format
	responses := make([]models.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = msg.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": responses,
		"count":    len(responses),
	})
}
