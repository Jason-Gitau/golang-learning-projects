package models

import (
	"time"

	"gorm.io/gorm"
)

// TaskPriority represents the priority level of a task
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusCancelled  TaskStatus = "cancelled"
)

// Task represents a task in the system
type Task struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Priority    TaskPriority   `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'todo'" json:"status"`
	DueDate     *time.Time     `json:"due_date,omitempty"` // Pointer to allow null values
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"-"` // Don't include full user in task response
}

// CreateTaskRequest represents the payload for creating a new task
type CreateTaskRequest struct {
	Title       string        `json:"title" binding:"required,min=1,max=200"`
	Description string        `json:"description" binding:"max=1000"`
	Priority    *TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status      *TaskStatus   `json:"status" binding:"omitempty,oneof=todo in_progress completed cancelled"`
	DueDate     *time.Time    `json:"due_date"`
}

// UpdateTaskRequest represents the payload for updating an existing task
type UpdateTaskRequest struct {
	Title       *string       `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string       `json:"description" binding:"omitempty,max=1000"`
	Priority    *TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status      *TaskStatus   `json:"status" binding:"omitempty,oneof=todo in_progress completed cancelled"`
	DueDate     *time.Time    `json:"due_date"`
}

// TaskFilterParams represents query parameters for filtering tasks
type TaskFilterParams struct {
	Status   string `form:"status" binding:"omitempty,oneof=todo in_progress completed cancelled"`
	Priority string `form:"priority" binding:"omitempty,oneof=low medium high urgent"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at due_date priority"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// TaskStats represents statistics about tasks
type TaskStats struct {
	TotalTasks      int64            `json:"total_tasks"`
	CompletedTasks  int64            `json:"completed_tasks"`
	PendingTasks    int64            `json:"pending_tasks"`
	OverdueTasks    int64            `json:"overdue_tasks"`
	TasksByPriority map[string]int64 `json:"tasks_by_priority"`
	TasksByStatus   map[string]int64 `json:"tasks_by_status"`
}
