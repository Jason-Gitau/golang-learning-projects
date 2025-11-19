package jobs

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"deep-research-agent/agent"
	"deep-research-agent/models"
	ws "deep-research-agent/websocket"
)

// Queue manages the research job queue
type Queue struct {
	jobs      chan *models.ResearchJob
	workers   []*Worker
	db        *gorm.DB
	wsHub     *ws.Hub
	agent     *agent.ResearchAgent
	mu        sync.RWMutex
	isRunning bool
}

// NewQueue creates a new job queue
func NewQueue(workerCount int, db *gorm.DB, wsHub *ws.Hub, agent *agent.ResearchAgent) *Queue {
	return &Queue{
		jobs:      make(chan *models.ResearchJob, 100),
		workers:   make([]*Worker, 0, workerCount),
		db:        db,
		wsHub:     wsHub,
		agent:     agent,
		isRunning: false,
	}
}

// Start initializes and starts the worker pool
func (q *Queue) Start(workerCount int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isRunning {
		return fmt.Errorf("queue is already running")
	}

	log.Printf("Starting job queue with %d workers", workerCount)

	// Create and start workers
	for i := 0; i < workerCount; i++ {
		worker := NewWorker(i+1, q, q.agent, q.wsHub, q.db)
		q.workers = append(q.workers, worker)
		worker.Start()
	}

	q.isRunning = true

	// Load pending jobs from database
	go q.loadPendingJobs()

	log.Printf("Job queue started with %d workers", workerCount)
	return nil
}

// Stop gracefully stops the job queue
func (q *Queue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.isRunning {
		return
	}

	log.Println("Stopping job queue...")

	// Stop all workers
	for _, worker := range q.workers {
		worker.Stop()
	}

	q.workers = nil
	q.isRunning = false

	log.Println("Job queue stopped")
}

// Enqueue adds a new research job to the queue
func (q *Queue) Enqueue(request ResearchJobRequest) (*models.ResearchJob, error) {
	// Create research query
	query := request.ToResearchQuery()

	// Serialize query data
	queryData, err := MarshalJobData(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Create job record
	job := &models.ResearchJob{
		ID:          uuid.New().String(),
		Query:       request.Query,
		QueryData:   queryData,
		Depth:       string(request.Depth),
		Status:      models.JobStatusPending,
		Progress:    0,
		CurrentStep: 0,
		TotalSteps:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := q.db.Create(job).Error; err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Add to queue
	select {
	case q.jobs <- job:
		log.Printf("Job %s enqueued: %s", job.ID, job.Query)
		return job, nil
	default:
		// Queue is full, update status to failed
		job.Status = models.JobStatusFailed
		job.Error = "Job queue is full"
		q.db.Save(job)
		return nil, fmt.Errorf("job queue is full")
	}
}

// Dequeue retrieves the next job from the queue
func (q *Queue) Dequeue() *models.ResearchJob {
	select {
	case job := <-q.jobs:
		return job
	default:
		return nil
	}
}

// GetJob retrieves a job by ID from the database
func (q *Queue) GetJob(jobID string) (*models.ResearchJob, error) {
	var job models.ResearchJob
	if err := q.db.Where("id = ?", jobID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job not found: %s", jobID)
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return &job, nil
}

// CancelJob cancels a pending or running job
func (q *Queue) CancelJob(jobID string) error {
	var job models.ResearchJob
	if err := q.db.Where("id = ?", jobID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("job not found: %s", jobID)
		}
		return fmt.Errorf("failed to get job: %w", err)
	}

	// Can only cancel pending or running jobs
	if job.Status != models.JobStatusPending && job.Status != models.JobStatusRunning {
		return fmt.Errorf("job is not pending or running (status: %s)", job.Status)
	}

	// Update status to cancelled
	now := time.Now()
	job.Status = models.JobStatusCancelled
	job.CompletedAt = &now

	if err := q.db.Save(&job).Error; err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}

	// Send cancellation message via WebSocket
	if q.wsHub != nil {
		msg := ws.NewErrorMessage(jobID, "Job cancelled by user")
		q.wsHub.BroadcastToJob(jobID, msg)
	}

	log.Printf("Job %s cancelled", jobID)
	return nil
}

// GetQueueStats returns statistics about the job queue
func (q *Queue) GetQueueStats() (map[string]interface{}, error) {
	var stats struct {
		Total     int64
		Pending   int64
		Running   int64
		Completed int64
		Failed    int64
		Cancelled int64
	}

	// Count total jobs
	if err := q.db.Model(&models.ResearchJob{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Count by status
	if err := q.db.Model(&models.ResearchJob{}).Where("status = ?", models.JobStatusPending).Count(&stats.Pending).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.ResearchJob{}).Where("status = ?", models.JobStatusRunning).Count(&stats.Running).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.ResearchJob{}).Where("status = ?", models.JobStatusCompleted).Count(&stats.Completed).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.ResearchJob{}).Where("status = ?", models.JobStatusFailed).Count(&stats.Failed).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.ResearchJob{}).Where("status = ?", models.JobStatusCancelled).Count(&stats.Cancelled).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":        stats.Total,
		"pending":      stats.Pending,
		"running":      stats.Running,
		"completed":    stats.Completed,
		"failed":       stats.Failed,
		"cancelled":    stats.Cancelled,
		"queue_length": len(q.jobs),
		"workers":      len(q.workers),
	}, nil
}

// loadPendingJobs loads pending jobs from the database on startup
func (q *Queue) loadPendingJobs() {
	var pendingJobs []models.ResearchJob
	if err := q.db.Where("status = ?", models.JobStatusPending).
		Order("created_at ASC").
		Find(&pendingJobs).Error; err != nil {
		log.Printf("Failed to load pending jobs: %v", err)
		return
	}

	log.Printf("Loading %d pending jobs from database", len(pendingJobs))

	for i := range pendingJobs {
		job := &pendingJobs[i]
		select {
		case q.jobs <- job:
			log.Printf("Loaded pending job %s", job.ID)
		default:
			log.Printf("Queue full, skipping job %s", job.ID)
		}
	}
}
