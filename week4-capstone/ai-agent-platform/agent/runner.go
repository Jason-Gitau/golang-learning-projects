package agent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// WorkerPool manages a pool of workers for concurrent agent processing
type WorkerPool struct {
	engine      *Engine
	broadcaster ResponseBroadcaster
	workers     int
	messageQueue chan *MessageTask
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// MessageTask represents a message processing task
type MessageTask struct {
	ConversationID string
	AgentID        string
	UserMessage    string
	Timestamp      time.Time
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(engine *Engine, broadcaster ResponseBroadcaster, workers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		engine:       engine,
		broadcaster:  broadcaster,
		workers:      workers,
		messageQueue: make(chan *MessageTask, queueSize),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	log.Printf("Starting agent worker pool with %d workers", wp.workers)

	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker processes messages from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	log.Printf("Agent worker %d started", id)

	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("Agent worker %d stopping", id)
			return

		case task, ok := <-wp.messageQueue:
			if !ok {
				log.Printf("Agent worker %d: queue closed", id)
				return
			}

			wp.processTask(id, task)
		}
	}
}

// processTask processes a single message task
func (wp *WorkerPool) processTask(workerID int, task *MessageTask) {
	startTime := time.Now()

	log.Printf("Worker %d processing message for conversation %s", workerID, task.ConversationID)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(wp.ctx, 60*time.Second)
	defer cancel()

	// Process message through engine
	responseChan, err := wp.engine.ProcessMessage(
		ctx,
		task.ConversationID,
		task.AgentID,
		task.UserMessage,
	)

	if err != nil {
		log.Printf("Worker %d: Failed to process message: %v", workerID, err)
		wp.broadcaster.BroadcastError(task.ConversationID, err.Error())
		return
	}

	// Stream the response
	streamHandler := NewStreamHandler(task.ConversationID, wp.broadcaster)
	fullResponse, toolCalls, err := streamHandler.StreamResponse(ctx, responseChan)

	if err != nil {
		log.Printf("Worker %d: Stream error: %v", workerID, err)
		wp.broadcaster.BroadcastError(task.ConversationID, err.Error())
		return
	}

	duration := time.Since(startTime)
	log.Printf("Worker %d completed message for conversation %s in %v (response: %d chars, tools: %d)",
		workerID, task.ConversationID, duration, len(fullResponse), len(toolCalls))
}

// Submit submits a message task to the queue
func (wp *WorkerPool) Submit(task *MessageTask) error {
	select {
	case wp.messageQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return ErrQueueFull
	}
}

// Shutdown gracefully shuts down the worker pool
func (wp *WorkerPool) Shutdown(timeout time.Duration) error {
	log.Println("Shutting down agent worker pool...")

	// Stop accepting new tasks
	close(wp.messageQueue)

	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All agent workers stopped gracefully")
		return nil
	case <-time.After(timeout):
		log.Println("Worker pool shutdown timeout, forcing stop")
		wp.cancel()
		return ErrShutdownTimeout
	}
}

// QueueSize returns the current size of the message queue
func (wp *WorkerPool) QueueSize() int {
	return len(wp.messageQueue)
}

// IsRunning returns whether the worker pool is running
func (wp *WorkerPool) IsRunning() bool {
	select {
	case <-wp.ctx.Done():
		return false
	default:
		return true
	}
}

// Errors
var (
	ErrQueueFull        = fmt.Errorf("message queue is full")
	ErrShutdownTimeout  = fmt.Errorf("shutdown timeout exceeded")
)

// Runner manages the entire agent processing system
type Runner struct {
	engine      *Engine
	broadcaster ResponseBroadcaster
	workerPool  *WorkerPool
}

// NewRunner creates a new agent runner
func NewRunner(engine *Engine, broadcaster ResponseBroadcaster, workers, queueSize int) *Runner {
	return &Runner{
		engine:      engine,
		broadcaster: broadcaster,
		workerPool:  NewWorkerPool(engine, broadcaster, workers, queueSize),
	}
}

// Start starts the runner
func (r *Runner) Start() {
	r.workerPool.Start()
	log.Println("Agent runner started")
}

// ProcessMessage submits a message for processing
func (r *Runner) ProcessMessage(conversationID, agentID, userMessage string) error {
	task := &MessageTask{
		ConversationID: conversationID,
		AgentID:        agentID,
		UserMessage:    userMessage,
		Timestamp:      time.Now(),
	}

	return r.workerPool.Submit(task)
}

// Shutdown gracefully shuts down the runner
func (r *Runner) Shutdown(timeout time.Duration) error {
	return r.workerPool.Shutdown(timeout)
}

// GetStats returns runner statistics
func (r *Runner) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"queue_size":              r.workerPool.QueueSize(),
		"running":                 r.workerPool.IsRunning(),
		"active_conversations":    r.engine.GetActiveConversationCount(),
	}
}
