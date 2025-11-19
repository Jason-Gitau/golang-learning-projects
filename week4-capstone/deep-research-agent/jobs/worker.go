package jobs

import (
	"fmt"
	"log"
	"time"

	"deep-research-agent/agent"
	"deep-research-agent/models"
	ws "deep-research-agent/websocket"
	"gorm.io/gorm"
)

// Worker represents a job processing worker
type Worker struct {
	ID           int
	Queue        *Queue
	ResearchAgent *agent.ResearchAgent
	WebSocketHub *ws.Hub
	DB           *gorm.DB
	stopChan     chan struct{}
}

// NewWorker creates a new worker
func NewWorker(id int, queue *Queue, agent *agent.ResearchAgent, hub *ws.Hub, db *gorm.DB) *Worker {
	return &Worker{
		ID:           id,
		Queue:        queue,
		ResearchAgent: agent,
		WebSocketHub: hub,
		DB:           db,
		stopChan:     make(chan struct{}),
	}
}

// Start begins processing jobs from the queue
func (w *Worker) Start() {
	log.Printf("Worker %d started", w.ID)
	go w.run()
}

// Stop gracefully stops the worker
func (w *Worker) Stop() {
	log.Printf("Worker %d stopping", w.ID)
	close(w.stopChan)
}

// run is the main worker loop
func (w *Worker) run() {
	for {
		select {
		case <-w.stopChan:
			log.Printf("Worker %d stopped", w.ID)
			return
		default:
			// Try to get a job from the queue
			job := w.Queue.Dequeue()
			if job == nil {
				// No jobs available, sleep briefly
				time.Sleep(500 * time.Millisecond)
				continue
			}

			log.Printf("Worker %d processing job %s", w.ID, job.ID)
			w.processJob(job)
		}
	}
}

// processJob processes a single research job
func (w *Worker) processJob(job *models.ResearchJob) {
	// Update job status to running
	now := time.Now()
	job.Status = models.JobStatusRunning
	job.StartedAt = &now
	if err := w.DB.Save(job).Error; err != nil {
		log.Printf("Worker %d: Failed to update job status: %v", w.ID, err)
		return
	}

	// Send research started message via WebSocket
	if w.WebSocketHub != nil {
		msg := ws.NewResearchStartedMessage(job.ID, job.Query, 0)
		w.WebSocketHub.BroadcastToJob(job.ID, msg)
	}

	// Parse query data
	query, err := UnmarshalJobData(job.QueryData)
	if err != nil {
		w.failJob(job, fmt.Sprintf("Failed to parse job data: %v", err))
		return
	}

	// Create progress channel
	progressChan := make(chan models.ResearchProgress, 10)

	// Start progress monitor goroutine
	go w.monitorProgress(job.ID, progressChan)

	// Execute research
	options := models.ResearchOptions{
		SessionID:       job.ID,
		SaveSession:     false, // We're managing the session ourselves
		StreamProgress:  true,
		ProgressChannel: progressChan,
		Timeout:         10 * time.Minute,
	}

	result, err := w.ResearchAgent.Research(*query, options)

	// Close progress channel
	close(progressChan)

	if err != nil {
		w.failJob(job, fmt.Sprintf("Research failed: %v", err))
		return
	}

	// Save result
	resultJSON, err := MarshalJobResult(result)
	if err != nil {
		w.failJob(job, fmt.Sprintf("Failed to marshal result: %v", err))
		return
	}

	// Update job status to completed
	completedAt := time.Now()
	job.Status = models.JobStatusCompleted
	job.Result = resultJSON
	job.Progress = 100
	job.CompletedAt = &completedAt

	if err := w.DB.Save(job).Error; err != nil {
		log.Printf("Worker %d: Failed to save job result: %v", w.ID, err)
		return
	}

	// Send research completed message via WebSocket
	if w.WebSocketHub != nil {
		duration := completedAt.Sub(*job.StartedAt)
		msg := ws.NewResearchCompletedMessage(job.ID, true, duration.String(), result)
		w.WebSocketHub.BroadcastToJob(job.ID, msg)
	}

	log.Printf("Worker %d: Job %s completed successfully", w.ID, job.ID)
}

// monitorProgress monitors research progress and sends WebSocket updates
func (w *Worker) monitorProgress(jobID string, progressChan <-chan models.ResearchProgress) {
	for progress := range progressChan {
		// Update job progress in database
		if err := w.DB.Model(&models.ResearchJob{}).
			Where("id = ?", jobID).
			Updates(map[string]interface{}{
				"progress":     progress.Progress,
				"current_step": progress.CurrentStep,
				"total_steps":  progress.TotalSteps,
				"updated_at":   time.Now(),
			}).Error; err != nil {
			log.Printf("Failed to update job progress: %v", err)
		}

		// Send progress update via WebSocket
		if w.WebSocketHub != nil {
			msg := ws.NewProgressMessage(
				jobID,
				progress.CurrentStep,
				progress.TotalSteps,
				progress.Progress,
				progress.StepDescription,
			)
			w.WebSocketHub.BroadcastToJob(jobID, msg)

			// Also send step-specific messages
			if progress.Status == "step_started" {
				stepMsg := ws.NewStepStartedMessage(
					jobID,
					progress.CurrentStep,
					progress.StepDescription,
					"", // Tool name would need to be added to progress
				)
				w.WebSocketHub.BroadcastToJob(jobID, stepMsg)
			}
		}
	}
}

// failJob marks a job as failed
func (w *Worker) failJob(job *models.ResearchJob, errorMsg string) {
	log.Printf("Worker %d: Job %s failed: %s", w.ID, job.ID, errorMsg)

	completedAt := time.Now()
	job.Status = models.JobStatusFailed
	job.Error = errorMsg
	job.CompletedAt = &completedAt

	if err := w.DB.Save(job).Error; err != nil {
		log.Printf("Worker %d: Failed to save job failure: %v", w.ID, err)
	}

	// Send error message via WebSocket
	if w.WebSocketHub != nil {
		msg := ws.NewErrorMessage(job.ID, errorMsg)
		w.WebSocketHub.BroadcastToJob(job.ID, msg)
	}
}
