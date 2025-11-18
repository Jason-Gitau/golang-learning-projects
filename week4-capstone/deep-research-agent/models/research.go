package models

import (
	"time"
)

// ResearchDepth defines how deep the research should go
type ResearchDepth string

const (
	ResearchShallow ResearchDepth = "shallow" // Quick overview
	ResearchMedium  ResearchDepth = "medium"  // Standard research
	ResearchDeep    ResearchDepth = "deep"    // Comprehensive research
)

// ResearchType categorizes the kind of research being conducted
type ResearchType string

const (
	ResearchGeneral  ResearchType = "general"  // General topic research
	ResearchAcademic ResearchType = "academic" // Academic/scientific research
	ResearchDocument ResearchType = "document" // Document analysis
	ResearchMulti    ResearchType = "multi"    // Multi-source research
)

// ResearchQuery represents a research request
type ResearchQuery struct {
	Query      string        `json:"query"`
	Documents  []string      `json:"documents"`   // PDF/DOCX paths
	UseWeb     bool          `json:"use_web"`     // Enable web search
	UseWiki    bool          `json:"use_wiki"`    // Enable Wikipedia
	Depth      ResearchDepth `json:"depth"`       // Research depth
	MaxSources int           `json:"max_sources"` // Maximum sources to use
	MaxSteps   int           `json:"max_steps"`   // Maximum research steps
}

// ResearchStep represents a single step in the research plan
type ResearchStep struct {
	StepNumber  int                    `json:"step_number"`
	Tool        string                 `json:"tool"`        // Tool to use (e.g., "web_search", "pdf_processor")
	Action      string                 `json:"action"`      // Action to perform
	Parameters  map[string]interface{} `json:"parameters"`  // Tool-specific parameters
	Reasoning   string                 `json:"reasoning"`   // Why this step is needed
	DependsOn   []int                  `json:"depends_on"`  // Steps this depends on
	Priority    int                    `json:"priority"`    // Execution priority
	MaxRetries  int                    `json:"max_retries"` // Max retry attempts
}

// StepResult represents the result of executing a research step
type StepResult struct {
	Step      ResearchStep `json:"step"`
	Success   bool         `json:"success"`
	Data      interface{}  `json:"data"`     // Tool-specific result data
	Sources   []Source     `json:"sources"`  // Sources found in this step
	Error     error        `json:"error"`    // Error if failed
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Retries   int          `json:"retries"` // Number of retries attempted
}

// Source represents a citation source
type Source struct {
	ID          string                 `json:"id"`           // Unique source identifier
	Type        string                 `json:"type"`         // "web", "pdf", "docx", "wikipedia"
	URL         string                 `json:"url"`          // Source URL or file path
	Title       string                 `json:"title"`        // Source title
	Snippet     string                 `json:"snippet"`      // Relevant excerpt
	PageNumber  int                    `json:"page_number"`  // For PDFs
	Section     string                 `json:"section"`      // Section/chapter name
	CitationKey string                 `json:"citation_key"` // Unique citation identifier
	Timestamp   time.Time              `json:"timestamp"`    // When retrieved
	Relevance   float64                `json:"relevance"`    // Relevance score (0-1)
	Author      string                 `json:"author"`       // Author name
	Publisher   string                 `json:"publisher"`    // Publisher name
	PublishDate time.Time              `json:"publish_date"` // Publication date
	AccessDate  time.Time              `json:"access_date"`  // Access date
	FilePath    string                 `json:"file_path"`    // File path for local documents
	Content     string                 `json:"content"`      // Full content for documents
	Excerpt     string                 `json:"excerpt"`      // Excerpt from document
	Metadata    map[string]interface{} `json:"metadata"`     // Additional metadata
}

// Finding represents a key research finding
type Finding struct {
	ID         string    `json:"id"`         // Unique finding identifier
	Content    string    `json:"content"`    // The finding text
	Sources    []Source  `json:"sources"`    // Supporting sources
	Source     Source    `json:"source"`     // Primary source (for backward compatibility)
	Confidence float64   `json:"confidence"` // Confidence level (0-1)
	Category   string    `json:"category"`   // Finding category
	Keywords   []string  `json:"keywords"`   // Key terms
	Timestamp  time.Time `json:"timestamp"`  // When found
}

// ResearchResult represents the final research output
type ResearchResult struct {
	Query          string          `json:"query"`
	Summary        string          `json:"summary"`         // Executive summary
	KeyFindings    []Finding       `json:"key_findings"`    // Important findings
	Sources        []Source        `json:"sources"`         // All sources used
	Steps          []StepResult    `json:"steps"`           // All steps executed
	Confidence     float64         `json:"confidence"`      // Overall confidence
	Duration       time.Duration   `json:"duration"`        // Total duration
	TotalSteps     int             `json:"total_steps"`     // Steps executed
	SuccessfulSteps int            `json:"successful_steps"`
	FailedSteps    int             `json:"failed_steps"`
	ResearchType   ResearchType    `json:"research_type"`
	CompletionTime time.Time       `json:"completion_time"`
}

// ResearchSession represents a saved research session
type ResearchSession struct {
	ID         string          `json:"id" gorm:"primaryKey"`
	UserID     string          `json:"user_id" gorm:"type:varchar(36);index;not null"` // Owner of this session
	Query      string          `json:"query" gorm:"type:text"`
	QueryData  string          `json:"query_data" gorm:"type:text"` // JSON serialized ResearchQuery
	Result     string          `json:"result" gorm:"type:text"`     // JSON serialized ResearchResult
	Status     string          `json:"status"`                      // "pending", "in_progress", "completed", "failed"
	Progress   float64         `json:"progress"`                    // 0-100
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	CompletedAt *time.Time     `json:"completed_at"`
}

// Document represents an indexed document
type Document struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id" gorm:"type:varchar(36);index;not null"` // Owner of this document
	Filename    string    `json:"filename"`
	FilePath    string    `json:"file_path"`
	FileType    string    `json:"file_type"` // "pdf", "docx"
	FileSize    int64     `json:"file_size"`
	PageCount   int       `json:"page_count"`
	Indexed     bool      `json:"indexed"`
	IndexData   string    `json:"index_data" gorm:"type:text"` // JSON metadata
	UploadedAt  time.Time `json:"uploaded_at"`
}

// Add composite unique index for user_id + filename
func (Document) TableName() string {
	return "documents"
}

// ResearchProgress represents real-time progress updates
type ResearchProgress struct {
	SessionID      string    `json:"session_id"`
	CurrentStep    int       `json:"current_step"`
	TotalSteps     int       `json:"total_steps"`
	StepDescription string   `json:"step_description"`
	Progress       float64   `json:"progress"` // 0-100
	Status         string    `json:"status"`
	Message        string    `json:"message"`
	Timestamp      time.Time `json:"timestamp"`
}

// ResearchPlan represents the overall research execution plan
type ResearchPlan struct {
	Query         ResearchQuery  `json:"query"`
	Steps         []ResearchStep `json:"steps"`
	EstimatedTime time.Duration  `json:"estimated_time"`
	ResearchType  ResearchType   `json:"research_type"`
	Strategy      string         `json:"strategy"` // Description of research strategy
	CreatedAt     time.Time      `json:"created_at"`
}
