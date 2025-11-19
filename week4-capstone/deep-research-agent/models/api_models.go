package models

import (
	"time"
)

// ResearchJob represents an asynchronous research job
type ResearchJob struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Query       string         `json:"query" gorm:"type:text"`
	QueryData   string         `json:"query_data" gorm:"type:text"` // JSON serialized ResearchQuery
	Depth       string         `json:"depth"`
	Status      string         `json:"status"` // pending, running, completed, failed, cancelled
	Progress    float64        `json:"progress"` // 0-100
	CurrentStep int            `json:"current_step"`
	TotalSteps  int            `json:"total_steps"`
	Result      string         `json:"result" gorm:"type:text"` // JSON serialized ResearchResult
	Error       string         `json:"error" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	StartedAt   *time.Time     `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// UploadedFile represents a file uploaded for analysis
type UploadedFile struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	Path         string    `json:"path"`
	Status       string    `json:"status"` // uploaded, processing, indexed, failed
	Error        string    `json:"error" gorm:"type:text"`
	UploadedAt   time.Time `json:"uploaded_at"`
	ProcessedAt  *time.Time `json:"processed_at"`
}

// ExportJob represents a session export job
type ExportJob struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	SessionID   string     `json:"session_id"`
	Format      string     `json:"format"` // markdown, json, pdf
	Status      string     `json:"status"` // pending, processing, completed, failed
	FilePath    string     `json:"file_path"`
	Error       string     `json:"error" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// JobStatus constants
const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
	JobStatusCancelled = "cancelled"
)

// FileStatus constants
const (
	FileStatusUploaded   = "uploaded"
	FileStatusProcessing = "processing"
	FileStatusIndexed    = "indexed"
	FileStatusFailed     = "failed"
)

// ExportFormat constants
const (
	ExportFormatMarkdown = "markdown"
	ExportFormatJSON     = "json"
	ExportFormatPDF      = "pdf"
)
