package handlers

import (
	"net/http"
	"strconv"
	"task-management-api/database"
	"task-management-api/middleware"
	"task-management-api/models"
	"time"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles task-related requests
type TaskHandler struct{}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

// CreateTask creates a new task
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body models.CreateTaskRequest true "Task details"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Set default values
	priority := models.PriorityMedium
	if req.Priority != nil {
		priority = *req.Priority
	}

	status := models.StatusTodo
	if req.Status != nil {
		status = *req.Status
	}

	// Create task
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Status:      status,
		DueDate:     req.DueDate,
		UserID:      userID,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasks retrieves all tasks for the authenticated user with optional filtering
// @Summary Get all tasks
// @Description Get all tasks for the authenticated user with optional filtering
// @Tags tasks
// @Produce json
// @Param status query string false "Filter by status" Enums(todo, in_progress, completed, cancelled)
// @Param priority query string false "Filter by priority" Enums(low, medium, high, urgent)
// @Param sort_by query string false "Sort by field" Enums(created_at, due_date, priority)
// @Param order query string false "Sort order" Enums(asc, desc)
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Parse filter parameters
	var filters models.TaskFilterParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Build query
	query := database.DB.Where("user_id = ?", userID)

	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Priority != "" {
		query = query.Where("priority = ?", filters.Priority)
	}

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}

	order := "desc"
	if filters.Order != "" {
		order = filters.Order
	}

	query = query.Order(sortBy + " " + order)

	// Execute query
	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tasks",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask retrieves a single task by ID
// @Summary Get a task
// @Description Get a task by ID (must belong to authenticated user)
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get task ID from URL
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Fetch task
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates an existing task
// @Summary Update a task
// @Description Update a task by ID (must belong to authenticated user)
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body models.UpdateTaskRequest true "Updated task details"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get task ID from URL
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Fetch task
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Bind update request
	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Update fields if provided
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	// Save updates
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task
// @Summary Delete a task
// @Description Delete a task by ID (must belong to authenticated user)
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get task ID from URL
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Delete task
	result := database.DB.Where("id = ? AND user_id = ?", taskID, userID).Delete(&models.Task{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// GetTaskStats returns statistics about the user's tasks
// @Summary Get task statistics
// @Description Get statistics about tasks for the authenticated user
// @Tags tasks
// @Produce json
// @Success 200 {object} models.TaskStats
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /tasks/stats [get]
func (h *TaskHandler) GetTaskStats(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var stats models.TaskStats

	// Total tasks
	database.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&stats.TotalTasks)

	// Completed tasks
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.StatusCompleted).Count(&stats.CompletedTasks)

	// Pending tasks (todo + in_progress)
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status IN ?", userID, []models.TaskStatus{models.StatusTodo, models.StatusInProgress}).Count(&stats.PendingTasks)

	// Overdue tasks
	now := time.Now()
	database.DB.Model(&models.Task{}).
		Where("user_id = ? AND status != ? AND due_date IS NOT NULL AND due_date < ?", userID, models.StatusCompleted, now).
		Count(&stats.OverdueTasks)

	// Tasks by priority
	stats.TasksByPriority = make(map[string]int64)
	for _, priority := range []models.TaskPriority{models.PriorityLow, models.PriorityMedium, models.PriorityHigh, models.PriorityUrgent} {
		var count int64
		database.DB.Model(&models.Task{}).Where("user_id = ? AND priority = ?", userID, priority).Count(&count)
		stats.TasksByPriority[string(priority)] = count
	}

	// Tasks by status
	stats.TasksByStatus = make(map[string]int64)
	for _, status := range []models.TaskStatus{models.StatusTodo, models.StatusInProgress, models.StatusCompleted, models.StatusCancelled} {
		var count int64
		database.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, status).Count(&count)
		stats.TasksByStatus[string(status)] = count
	}

	c.JSON(http.StatusOK, stats)
}
