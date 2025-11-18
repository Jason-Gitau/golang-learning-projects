package models

import (
	"time"

	"gorm.io/gorm"
)

// ScrapedPage represents a scraped web page stored in the database
type ScrapedPage struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	URL         string         `gorm:"uniqueIndex;not null" json:"url"`
	Title       string         `gorm:"type:text" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	LinkCount   int            `json:"link_count"`
	StatusCode  int            `json:"status_code"`
	Error       string         `gorm:"type:text" json:"error,omitempty"`
	RetryCount  int            `json:"retry_count"`
	ScrapedAt   time.Time      `json:"scraped_at"`
	Duration    int64          `json:"duration_ms"` // Duration in milliseconds
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Link represents a link found on a scraped page
type Link struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	PageID    uint      `gorm:"not null;index" json:"page_id"`
	URL       string    `gorm:"type:text;not null" json:"url"`
	Text      string    `gorm:"type:text" json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// ScrapeJob represents a job to be processed by workers
type ScrapeJob struct {
	URL        string
	RetryCount int
}

// ScrapeResult represents the result of a scraping operation
type ScrapeResult struct {
	URL         string
	Title       string
	Description string
	Links       []LinkData
	StatusCode  int
	Error       error
	Duration    time.Duration
	RetryCount  int
	ScrapedAt   time.Time
}

// LinkData represents a link found during scraping
type LinkData struct {
	URL  string
	Text string
}

// Statistics holds scraping statistics
type Statistics struct {
	TotalPages      int           `json:"total_pages"`
	SuccessfulPages int           `json:"successful_pages"`
	FailedPages     int           `json:"failed_pages"`
	TotalLinks      int           `json:"total_links"`
	TotalDuration   time.Duration `json:"total_duration"`
	AverageDuration time.Duration `json:"average_duration"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
}
