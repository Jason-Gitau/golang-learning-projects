package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"deep-research-agent/agent"
	"deep-research-agent/jobs"
	"deep-research-agent/models"
	"deep-research-agent/storage"
	"deep-research-agent/uploads"
	ws "deep-research-agent/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// APIHandler handles all API requests
type APIHandler struct {
	jobQueue      *jobs.Queue
	uploadHandler *uploads.UploadHandler
	researchAgent *agent.ResearchAgent
	storage       *storage.ResearchDB
	wsHub         *ws.Hub
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(
	jobQueue *jobs.Queue,
	uploadHandler *uploads.UploadHandler,
	researchAgent *agent.ResearchAgent,
	storage *storage.ResearchDB,
	wsHub *ws.Hub,
) *APIHandler {
	return &APIHandler{
		jobQueue:      jobQueue,
		uploadHandler: uploadHandler,
		researchAgent: researchAgent,
		storage:       storage,
		wsHub:         wsHub,
	}
}

// ============================================================================
// Research Operations
// ============================================================================

// StartResearch handles POST /api/v1/research/start
func (h *APIHandler) StartResearch(c *gin.Context) {
	var request jobs.ResearchJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate request
	if request.Query == "" {
		c.JSON(400, gin.H{"error": "Query is required"})
		return
	}

	// Set defaults
	if request.Depth == "" {
		request.Depth = models.ResearchMedium
	}

	// Enqueue job
	job, err := h.jobQueue.Enqueue(request)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to enqueue research job", "details": err.Error()})
		return
	}

	c.JSON(202, gin.H{
		"job_id":  job.ID,
		"status":  job.Status,
		"message": "Research job queued successfully",
		"query":   job.Query,
	})
}

// GetResearchStatus handles GET /api/v1/research/:id/status
func (h *APIHandler) GetResearchStatus(c *gin.Context) {
	jobID := c.Param("id")

	job, err := h.jobQueue.GetJob(jobID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Job not found", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"job_id":       job.ID,
		"query":        job.Query,
		"status":       job.Status,
		"progress":     job.Progress,
		"current_step": job.CurrentStep,
		"total_steps":  job.TotalSteps,
		"created_at":   job.CreatedAt,
		"started_at":   job.StartedAt,
		"completed_at": job.CompletedAt,
		"error":        job.Error,
	})
}

// GetResearchResult handles GET /api/v1/research/:id/result
func (h *APIHandler) GetResearchResult(c *gin.Context) {
	jobID := c.Param("id")

	job, err := h.jobQueue.GetJob(jobID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Job not found", "details": err.Error()})
		return
	}

	if job.Status != models.JobStatusCompleted {
		c.JSON(400, gin.H{
			"error":   "Job not completed",
			"status":  job.Status,
			"message": "Research job has not completed yet",
		})
		return
	}

	// Parse result
	var result models.ResearchResult
	if err := json.Unmarshal([]byte(job.Result), &result); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse result", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"job_id": job.ID,
		"query":  job.Query,
		"result": result,
	})
}

// CancelResearch handles DELETE /api/v1/research/:id
func (h *APIHandler) CancelResearch(c *gin.Context) {
	jobID := c.Param("id")

	if err := h.jobQueue.CancelJob(jobID); err != nil {
		c.JSON(400, gin.H{"error": "Failed to cancel job", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Research job cancelled successfully",
		"job_id":  jobID,
	})
}

// ============================================================================
// WebSocket Stream
// ============================================================================

// StreamResearch handles WS /api/v1/research/:id/stream
func (h *APIHandler) StreamResearch(c *gin.Context) {
	jobID := c.Param("id")

	// Verify job exists
	_, err := h.jobQueue.GetJob(jobID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	// Create client
	clientID := uuid.New().String()
	client := ws.NewClient(clientID, jobID, h.wsHub, conn)

	// Register client with hub
	h.wsHub.Register <- client

	// Start read/write pumps
	go client.WritePump()
	go client.ReadPump()
}

// ============================================================================
// Document Operations
// ============================================================================

// UploadDocument handles POST /api/v1/documents/upload
func (h *APIHandler) UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded", "details": err.Error()})
		return
	}

	// Handle upload
	uploadedFile, err := h.uploadHandler.HandleUpload(file)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to upload file", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"file_id":       uploadedFile.ID,
		"filename":      uploadedFile.OriginalName,
		"size":          uploadedFile.Size,
		"status":        uploadedFile.Status,
		"uploaded_at":   uploadedFile.UploadedAt,
		"message":       "File uploaded successfully",
	})
}

// AnalyzeDocument handles POST /api/v1/documents/analyze
func (h *APIHandler) AnalyzeDocument(c *gin.Context) {
	var request struct {
		FileID string `json:"file_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Process file asynchronously
	go func() {
		if err := h.uploadHandler.ProcessFile(request.FileID); err != nil {
			log.Printf("Failed to process file %s: %v", request.FileID, err)
		}
	}()

	c.JSON(202, gin.H{
		"message": "Document analysis started",
		"file_id": request.FileID,
	})
}

// GetDocument handles GET /api/v1/documents/:id
func (h *APIHandler) GetDocument(c *gin.Context) {
	fileID := c.Param("id")

	file, err := h.uploadHandler.GetFile(fileID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Document not found", "details": err.Error()})
		return
	}

	c.JSON(200, file)
}

// DeleteDocument handles DELETE /api/v1/documents/:id
func (h *APIHandler) DeleteDocument(c *gin.Context) {
	fileID := c.Param("id")

	if err := h.uploadHandler.DeleteFile(fileID); err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete document", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Document deleted successfully",
		"file_id": fileID,
	})
}

// ListDocuments handles GET /api/v1/documents
func (h *APIHandler) ListDocuments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	files, err := h.uploadHandler.ListFiles(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list documents", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"documents": files,
		"limit":     limit,
		"offset":    offset,
		"count":     len(files),
	})
}

// ============================================================================
// Session Management
// ============================================================================

// ListSessions handles GET /api/v1/sessions
func (h *APIHandler) ListSessions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	sessions, err := h.storage.ListSessions(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list sessions", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"sessions": sessions,
		"limit":    limit,
		"offset":   offset,
		"count":    len(sessions),
	})
}

// GetSession handles GET /api/v1/sessions/:id
func (h *APIHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")

	session, err := h.storage.GetSession(sessionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found", "details": err.Error()})
		return
	}

	// Parse result if exists
	var result *models.ResearchResult
	if session.Result != "" {
		if err := json.Unmarshal([]byte(session.Result), &result); err != nil {
			log.Printf("Warning: Failed to parse session result: %v", err)
		}
	}

	c.JSON(200, gin.H{
		"session": session,
		"result":  result,
	})
}

// DeleteSession handles DELETE /api/v1/sessions/:id
func (h *APIHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")

	if err := h.storage.DeleteSession(sessionID); err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete session", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":    "Session deleted successfully",
		"session_id": sessionID,
	})
}

// ============================================================================
// Export Operations
// ============================================================================

// ExportSession handles POST /api/v1/sessions/:id/export
func (h *APIHandler) ExportSession(c *gin.Context) {
	sessionID := c.Param("id")

	var request struct {
		Format string `json:"format" binding:"required"` // markdown, json, pdf
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Get session
	session, err := h.storage.GetSession(sessionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found", "details": err.Error()})
		return
	}

	// Parse result (for future use in export)
	_, err = storage.ParseSessionResult(session)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse session result", "details": err.Error()})
		return
	}

	// Generate export file
	var exportPath string
	switch request.Format {
	case "markdown", "json", "pdf":
		exportPath = fmt.Sprintf("exports/%s.%s", sessionID, request.Format)
		// Export functionality stub - would use reporting package
		// For now, just create a placeholder file
		os.WriteFile(exportPath, []byte(fmt.Sprintf("Export for session %s", sessionID)), 0644)
	default:
		c.JSON(400, gin.H{"error": "Unsupported format", "supported": []string{"markdown", "json", "pdf"}})
		return
	}

	c.JSON(200, gin.H{
		"message":     "Export generated successfully",
		"session_id":  sessionID,
		"format":      request.Format,
		"export_path": exportPath,
		"download_url": fmt.Sprintf("/api/v1/export/%s/download", sessionID),
	})
}

// DownloadExport handles GET /api/v1/export/:id/download
func (h *APIHandler) DownloadExport(c *gin.Context) {
	sessionID := c.Param("id")
	format := c.DefaultQuery("format", "markdown")

	var exportPath string
	var contentType string

	switch format {
	case "markdown":
		exportPath = fmt.Sprintf("exports/%s.md", sessionID)
		contentType = "text/markdown"
	case "json":
		exportPath = fmt.Sprintf("exports/%s.json", sessionID)
		contentType = "application/json"
	case "pdf":
		exportPath = fmt.Sprintf("exports/%s.pdf", sessionID)
		contentType = "application/pdf"
	default:
		c.JSON(400, gin.H{"error": "Unsupported format"})
		return
	}

	c.Header("Content-Type", contentType)
	c.File(exportPath)
}

// ============================================================================
// Statistics
// ============================================================================

// GetStats handles GET /api/v1/stats
func (h *APIHandler) GetStats(c *gin.Context) {
	// Get session count
	sessionCount, err := h.storage.GetSessionCount()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get session count", "details": err.Error()})
		return
	}

	// Get document count
	docCount, err := h.storage.GetDocumentCount()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get document count", "details": err.Error()})
		return
	}

	// Get job queue stats
	queueStats, err := h.jobQueue.GetQueueStats()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get queue stats", "details": err.Error()})
		return
	}

	// Get agent metrics
	agentMetrics := h.researchAgent.GetMetrics()

	c.JSON(200, gin.H{
		"sessions":      sessionCount,
		"documents":     docCount,
		"queue_stats":   queueStats,
		"agent_metrics": agentMetrics,
		"websocket": gin.H{
			"total_clients": h.wsHub.GetTotalClientCount(),
		},
	})
}

// ============================================================================
// Health Check
// ============================================================================

// HealthCheck handles GET /health
func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "deep-research-agent-api",
		"version": "1.0.0",
	})
}
