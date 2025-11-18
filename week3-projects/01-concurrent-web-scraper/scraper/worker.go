package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/user/web-scraper/models"
	"github.com/user/web-scraper/ratelimiter"
	"github.com/user/web-scraper/storage"
)

// WorkerPool manages a pool of workers for concurrent scraping
type WorkerPool struct {
	workerCount int
	jobs        chan models.ScrapeJob
	results     chan *models.ScrapeResult
	scraper     *Scraper
	db          *storage.Database
	rateLimiter *ratelimiter.RateLimiter
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	stats       *Statistics
}

// Statistics tracks scraping progress
type Statistics struct {
	mu              sync.Mutex
	TotalJobs       int
	CompletedJobs   int
	SuccessfulJobs  int
	FailedJobs      int
	InProgressJobs  int
	TotalRetries    int
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(ctx context.Context, workerCount int, scraper *Scraper, db *storage.Database, rateLimiter *ratelimiter.RateLimiter) *WorkerPool {
	workerCtx, cancel := context.WithCancel(ctx)

	return &WorkerPool{
		workerCount: workerCount,
		jobs:        make(chan models.ScrapeJob, workerCount*2), // Buffer for efficiency
		results:     make(chan *models.ScrapeResult, workerCount*2),
		scraper:     scraper,
		db:          db,
		rateLimiter: rateLimiter,
		ctx:         workerCtx,
		cancel:      cancel,
		stats:       &Statistics{},
	}
}

// Start starts all workers
func (wp *WorkerPool) Start() {
	log.Printf("Starting worker pool with %d workers", wp.workerCount)

	// Start workers
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i + 1)
	}

	// Start result processor
	wp.wg.Add(1)
	go wp.resultProcessor()
}

// worker is the worker goroutine that processes jobs
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("Worker %d stopped", id)
			return

		case job, ok := <-wp.jobs:
			if !ok {
				log.Printf("Worker %d: jobs channel closed", id)
				return
			}

			// Update statistics
			wp.stats.mu.Lock()
			wp.stats.InProgressJobs++
			wp.stats.mu.Unlock()

			// Wait for rate limiter
			if err := wp.rateLimiter.Wait(wp.ctx); err != nil {
				log.Printf("Worker %d: rate limiter cancelled", id)
				wp.stats.mu.Lock()
				wp.stats.InProgressJobs--
				wp.stats.mu.Unlock()
				return
			}

			// Process the job
			log.Printf("Worker %d: scraping %s (attempt %d)", id, job.URL, job.RetryCount+1)
			result := wp.scraper.ScrapeURL(wp.ctx, job.URL)
			result.RetryCount = job.RetryCount

			// Send result to processor
			select {
			case wp.results <- result:
			case <-wp.ctx.Done():
				wp.stats.mu.Lock()
				wp.stats.InProgressJobs--
				wp.stats.mu.Unlock()
				return
			}

			// Update statistics
			wp.stats.mu.Lock()
			wp.stats.InProgressJobs--
			wp.stats.CompletedJobs++
			wp.stats.mu.Unlock()
		}
	}
}

// resultProcessor processes scraping results
func (wp *WorkerPool) resultProcessor() {
	defer wp.wg.Done()

	log.Println("Result processor started")

	for {
		select {
		case <-wp.ctx.Done():
			log.Println("Result processor stopped")
			return

		case result, ok := <-wp.results:
			if !ok {
				log.Println("Result processor: results channel closed")
				return
			}

			// Save to database
			if err := wp.db.SavePage(result); err != nil {
				log.Printf("Failed to save page %s: %v", result.URL, err)
			}

			// Update statistics
			wp.stats.mu.Lock()
			if result.Error != nil {
				wp.stats.FailedJobs++
				log.Printf("✗ Failed: %s - %v", result.URL, result.Error)

				// Retry logic
				if result.RetryCount < wp.scraper.maxRetries {
					wp.stats.TotalRetries++
					wp.stats.mu.Unlock()

					// Re-queue with incremented retry count
					retryJob := models.ScrapeJob{
						URL:        result.URL,
						RetryCount: result.RetryCount + 1,
					}
					select {
					case wp.jobs <- retryJob:
						log.Printf("↻ Retrying: %s (attempt %d)", result.URL, retryJob.RetryCount+1)
					case <-wp.ctx.Done():
						return
					}
					continue
				} else {
					log.Printf("✗ Max retries reached for: %s", result.URL)
				}
			} else {
				wp.stats.SuccessfulJobs++
				log.Printf("✓ Success: %s (status: %d, title: %s, links: %d, duration: %s)",
					result.URL, result.StatusCode, truncate(result.Title, 50),
					len(result.Links), result.Duration)
			}
			wp.stats.mu.Unlock()
		}
	}
}

// AddJob adds a job to the queue
func (wp *WorkerPool) AddJob(url string) {
	wp.stats.mu.Lock()
	wp.stats.TotalJobs++
	wp.stats.mu.Unlock()

	job := models.ScrapeJob{
		URL:        url,
		RetryCount: 0,
	}

	select {
	case wp.jobs <- job:
	case <-wp.ctx.Done():
		log.Printf("Cannot add job for %s: context cancelled", url)
	}
}

// AddJobs adds multiple jobs to the queue
func (wp *WorkerPool) AddJobs(urls []string) {
	for _, url := range urls {
		wp.AddJob(url)
	}
}

// Wait waits for all jobs to complete
func (wp *WorkerPool) Wait() {
	// Close jobs channel to signal workers to finish
	close(wp.jobs)

	// Wait for all workers to finish
	wp.wg.Wait()

	// Close results channel
	close(wp.results)
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	wp.cancel()
	wp.Wait()
}

// GetStatistics returns current statistics
func (wp *WorkerPool) GetStatistics() Statistics {
	wp.stats.mu.Lock()
	defer wp.stats.mu.Unlock()
	return *wp.stats
}

// PrintStatistics prints current statistics
func (wp *WorkerPool) PrintStatistics() {
	stats := wp.GetStatistics()
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Scraping Statistics")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Jobs:       %d\n", stats.TotalJobs)
	fmt.Printf("Completed:        %d\n", stats.CompletedJobs)
	fmt.Printf("Successful:       %d\n", stats.SuccessfulJobs)
	fmt.Printf("Failed:           %d\n", stats.FailedJobs)
	fmt.Printf("In Progress:      %d\n", stats.InProgressJobs)
	fmt.Printf("Total Retries:    %d\n", stats.TotalRetries)
	fmt.Println(strings.Repeat("=", 60))
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
