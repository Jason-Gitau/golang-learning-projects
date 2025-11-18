# 🎓 START HERE - Learning Guide

Welcome to Week 3 of your Go learning journey! This guide will help you master concurrency in Go through building a production-ready web scraper.

## 📚 What You'll Learn

By studying and building this project, you'll master:

1. **Goroutines and Channels**
2. **Worker Pool Pattern**
3. **Synchronization Primitives** (WaitGroup, Mutex)
4. **Context for Cancellation**
5. **Rate Limiting with time.Ticker**
6. **Concurrent Database Operations**
7. **HTML Parsing**
8. **Production-Ready CLI Applications**

## 🗺️ Learning Path

### Step 1: Understand the Big Picture (30 minutes)

Before diving into code, understand what this scraper does:

1. **Run the application:**
   ```bash
   go run main.go
   ```

2. **Observe the output:**
   - Workers starting up
   - Concurrent scraping
   - Real-time progress
   - Final statistics

3. **Try different configurations:**
   ```bash
   # More workers
   go run main.go -workers 10

   # Different rate limit
   go run main.go -rate 5.0

   # View statistics
   go run main.go -stats
   ```

4. **Questions to answer:**
   - How many workers are running simultaneously?
   - What happens when you press Ctrl+C?
   - How does rate limiting affect performance?
   - Where is the scraped data stored?

### Step 2: Study the Code Structure (2 hours)

Read files in this order for maximum comprehension:

#### Phase 1: Data Layer (20 minutes)

**1. models/models.go** (15 min)
- Database models (ScrapedPage, Link)
- DTOs (ScrapeJob, ScrapeResult)
- Statistics struct

**Key Concepts:**
```go
// GORM model with tags
type ScrapedPage struct {
    ID    uint   `gorm:"primarykey"`
    URL   string `gorm:"uniqueIndex;not null"`
}

// DTO for passing data between components
type ScrapeResult struct {
    URL         string
    Title       string
    Error       error
    Duration    time.Duration
}
```

**2. config/config.go** (5 min)
- Configuration struct
- Validation logic
- File loading/saving

**Questions:**
- Why separate models from DTOs?
- How are GORM tags used?
- What validation is applied?

#### Phase 2: Supporting Components (30 minutes)

**3. ratelimiter/limiter.go** (15 min)
- Token bucket algorithm
- Channel-based implementation
- Context integration

**Key Pattern:**
```go
// Create rate limiter
limiter := NewRateLimiter(ctx, 2.0) // 2 req/s

// Wait for token
if err := limiter.Wait(ctx); err != nil {
    return err
}
```

**Why it matters:**
- Prevents overwhelming servers
- Smooth traffic distribution
- Respects rate limits

**4. scraper/parser.go** (10 min)
- HTML parsing with goquery
- Title/description extraction
- Link collection

**Key Techniques:**
```go
// Find elements
title := doc.Find("title").First().Text()

// Get attributes
href, exists := s.Attr("href")

// Iterate over elements
doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
    // Extract link
})
```

**5. storage/database.go** (5 min)
- GORM setup
- Concurrent-safe operations
- Statistics queries

**Concurrency Safety:**
```go
type Database struct {
    db *gorm.DB
    mu sync.RWMutex  // Protects database access
}

func (d *Database) SavePage(result *ScrapeResult) error {
    d.mu.Lock()         // Exclusive lock
    defer d.mu.Unlock()
    // Database operations
}
```

#### Phase 3: Core Scraping Logic (45 minutes)

**6. scraper/scraper.go** (15 min)
- HTTP client configuration
- Request creation
- Response handling
- Error handling

**HTTP Best Practices:**
```go
// Create request with context (for cancellation)
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

// Set appropriate headers
req.Header.Set("User-Agent", userAgent)

// Handle timeouts
client := &http.Client{
    Timeout: 30 * time.Second,
}
```

**7. scraper/worker.go** (30 min) ⭐ MOST IMPORTANT
- Worker pool pattern
- Channel-based communication
- Statistics tracking
- Retry logic

**Worker Pool Pattern:**
```go
// Create channels
jobs := make(chan ScrapeJob, bufferSize)
results := make(chan *ScrapeResult, bufferSize)

// Start workers
for i := 0; i < workerCount; i++ {
    wg.Add(1)
    go worker(i)
}

// Worker function
func worker(id int) {
    defer wg.Done()
    for job := range jobs {
        // Process job
        result := processJob(job)
        results <- result
    }
}
```

**Why buffered channels?**
- Prevents blocking
- Improves throughput
- Decouples producers/consumers

#### Phase 4: Application Entry Point (25 minutes)

**8. main.go** (25 min)
- CLI flag parsing
- Component initialization
- Signal handling
- Coordination

**Signal Handling:**
```go
// Create cancellable context
ctx, cancel := context.WithCancel(context.Background())

// Handle Ctrl+C
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
go func() {
    <-sigChan
    cancel() // Cancel everything
}()
```

**Total Reading Time: ~2 hours**

### Step 3: Key Concepts Deep Dive

#### Concept 1: Goroutines and Channels

**What are goroutines?**
Lightweight threads managed by Go runtime.

```go
// Start a goroutine
go func() {
    // Runs concurrently
    fmt.Println("Hello from goroutine")
}()

// Start multiple goroutines
for i := 0; i < 10; i++ {
    go doWork(i)
}
```

**What are channels?**
Pipes for communication between goroutines.

```go
// Create channel
ch := make(chan string)

// Send (blocks until received)
ch <- "message"

// Receive (blocks until sent)
msg := <-ch

// Buffered channel (doesn't block if buffer not full)
ch := make(chan string, 10)
```

**In this project:**
- `jobs` channel - Sends URLs to workers
- `results` channel - Receives scraping results
- Both are buffered for efficiency

**Try this:**
Modify worker count and observe behavior:
```bash
go run main.go -workers 1   # Sequential
go run main.go -workers 10  # Highly concurrent
```

#### Concept 2: Worker Pool Pattern

**What is it?**
A pattern for limiting concurrency by having a fixed pool of workers.

**Why use it?**
- Limit resource usage (memory, connections)
- Control concurrency level
- Efficient job distribution

**Implementation:**
```go
// 1. Create job queue
jobs := make(chan Job, bufferSize)

// 2. Start fixed number of workers
for i := 0; i < workerCount; i++ {
    go worker(jobs)
}

// 3. Send jobs to queue
for _, job := range allJobs {
    jobs <- job
}

// 4. Close queue when done
close(jobs)
```

**In this project:**
- Workers: Goroutines that scrape pages
- Job queue: Channel of URLs to scrape
- Result queue: Channel of scraping results

**Analogy:**
Think of a restaurant:
- Workers = Chefs
- Job queue = Order tickets
- Result queue = Finished dishes

**Try this:**
Add logging to see which worker processes which URL:
```go
log.Printf("Worker %d: processing %s", id, job.URL)
```

#### Concept 3: Synchronization with WaitGroup

**What is sync.WaitGroup?**
A counter for waiting for goroutines to finish.

```go
var wg sync.WaitGroup

// Add to counter
wg.Add(1)

// Start goroutine
go func() {
    defer wg.Done() // Decrement counter
    // Do work
}()

// Wait for counter to reach zero
wg.Wait()
```

**In this project:**
```go
// Start workers
for i := 0; i < workerCount; i++ {
    wp.wg.Add(1)
    go wp.worker(i)
}

// Wait for all workers to finish
wp.wg.Wait()
```

**Common mistake:**
```go
// WRONG - Add inside goroutine (race condition)
go func() {
    wg.Add(1)
    defer wg.Done()
}()

// RIGHT - Add before starting goroutine
wg.Add(1)
go func() {
    defer wg.Done()
}()
```

#### Concept 4: Thread-Safe Operations with Mutex

**Why needed?**
Multiple goroutines accessing shared data can cause race conditions.

```go
// NOT THREAD-SAFE
type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++ // Race condition!
}

// THREAD-SAFE
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

**In this project:**
```go
type Statistics struct {
    mu              sync.Mutex
    CompletedJobs   int
    SuccessfulJobs  int
}

func (s *Statistics) RecordSuccess() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.CompletedJobs++
    s.SuccessfulJobs++
}
```

**Read vs Write Locks:**
```go
// RWMutex allows multiple readers OR one writer
var mu sync.RWMutex

// Read lock (multiple allowed)
mu.RLock()
defer mu.RUnlock()
// Read data

// Write lock (exclusive)
mu.Lock()
defer mu.Unlock()
// Modify data
```

**Try this:**
Run with race detector to find race conditions:
```bash
go run -race main.go
```

#### Concept 5: Context for Cancellation

**What is context.Context?**
A way to carry deadlines, cancellation signals, and values across API boundaries.

```go
// Create cancellable context
ctx, cancel := context.WithCancel(context.Background())

// Cancel everything
defer cancel()

// Check if cancelled
select {
case <-ctx.Done():
    return ctx.Err() // context.Canceled
default:
    // Continue
}
```

**Common patterns:**
```go
// 1. With timeout
ctx, cancel := context.WithTimeout(parent, 30*time.Second)

// 2. With deadline
deadline := time.Now().Add(1 * time.Hour)
ctx, cancel := context.WithDeadline(parent, deadline)

// 3. Pass to functions
func doWork(ctx context.Context) error {
    // Check cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Do work
    }
}
```

**In this project:**
- HTTP requests use context for cancellation
- Workers check context before processing
- Signal handler cancels context on Ctrl+C

**Flow:**
```
User presses Ctrl+C
    → signal.Notify triggers
    → cancel() called
    → ctx.Done() closes
    → All goroutines receive signal
    → Clean shutdown
```

#### Concept 6: Rate Limiting

**What is rate limiting?**
Controlling the rate of operations (e.g., requests per second).

**Token Bucket Algorithm:**
```
1. Tokens are added to bucket at fixed rate
2. Each request consumes one token
3. If no tokens available, request waits
4. Bucket has maximum capacity
```

**Implementation:**
```go
ticker := time.NewTicker(interval)
tokens := make(chan struct{}, capacity)

// Generate tokens
go func() {
    for range ticker.C {
        select {
        case tokens <- struct{}{}:
        default: // Bucket full
        }
    }
}()

// Consume token
<-tokens // Blocks if empty
```

**In this project:**
```go
// Create limiter: 2 requests per second
limiter := ratelimiter.NewRateLimiter(ctx, 2.0)

// Wait for token before request
limiter.Wait(ctx)
makeRequest()
```

**Try this:**
Compare performance with different rates:
```bash
# Slow
time go run main.go -rate 0.5

# Fast
time go run main.go -rate 10.0
```

### Step 4: Hands-On Exercises

#### Exercise 1: Add Request Counter (Easy)

**Goal:** Count total HTTP requests made (including retries)

**Steps:**
1. Add `TotalRequests` field to Statistics
2. Increment in worker before making request
3. Display in statistics output

<details>
<summary>Solution</summary>

```go
// In scraper/worker.go
type Statistics struct {
    mu            sync.Mutex
    TotalRequests int
    // ... other fields
}

// In worker function
wp.stats.mu.Lock()
wp.stats.TotalRequests++
wp.stats.mu.Unlock()

// In PrintStatistics
fmt.Printf("Total Requests:   %d\n", stats.TotalRequests)
```
</details>

#### Exercise 2: Add Timeout Handling (Medium)

**Goal:** Add per-request timeout that's separate from HTTP timeout

**Steps:**
1. Add timeout to worker processing
2. Use context.WithTimeout
3. Cancel if worker takes too long

<details>
<summary>Solution Hint</summary>

```go
// In worker function
ctx, cancel := context.WithTimeout(wp.ctx, 60*time.Second)
defer cancel()

result := wp.scraper.ScrapeURL(ctx, job.URL)
```
</details>

#### Exercise 3: Add Progress Bar (Medium)

**Goal:** Show progress percentage during scraping

**Steps:**
1. Track total jobs vs completed jobs
2. Calculate percentage
3. Print progress updates

<details>
<summary>Solution Hint</summary>

```go
// In resultProcessor or separate goroutine
ticker := time.NewTicker(1 * time.Second)
go func() {
    for range ticker.C {
        stats := wp.GetStatistics()
        progress := float64(stats.CompletedJobs) / float64(stats.TotalJobs) * 100
        fmt.Printf("\rProgress: %.1f%% (%d/%d)", progress, stats.CompletedJobs, stats.TotalJobs)
    }
}()
```
</details>

#### Exercise 4: Add Domain-Based Rate Limiting (Hard)

**Goal:** Different rate limits for different domains

**Steps:**
1. Create map of domain → rate limiter
2. Extract domain from URL
3. Use appropriate limiter per domain

<details>
<summary>Solution Hint</summary>

```go
type DomainRateLimiter struct {
    mu       sync.Mutex
    limiters map[string]*ratelimiter.RateLimiter
}

func (d *DomainRateLimiter) Wait(ctx context.Context, url string) error {
    domain := extractDomain(url)

    d.mu.Lock()
    limiter, exists := d.limiters[domain]
    if !exists {
        limiter = ratelimiter.NewRateLimiter(ctx, 2.0)
        d.limiters[domain] = limiter
    }
    d.mu.Unlock()

    return limiter.Wait(ctx)
}
```
</details>

#### Exercise 5: Add Link Following (Hard)

**Goal:** Scrape pages and also scrape the links found

**Steps:**
1. Add max depth limit
2. Track visited URLs to avoid duplicates
3. Add found links as new jobs
4. Respect domain boundaries (optional)

<details>
<summary>Solution Hint</summary>

```go
type Job struct {
    URL   string
    Depth int
}

// In result processor
if result.Error == nil && job.Depth < maxDepth {
    for _, link := range result.Links {
        if !visited[link.URL] {
            wp.AddJob(Job{
                URL:   link.URL,
                Depth: job.Depth + 1,
            })
        }
    }
}
```
</details>

### Step 5: Testing Your Knowledge

Answer these questions without looking at code:

#### Concurrency Questions

1. **What is a goroutine and how is it different from a thread?**
   - Lightweight, managed by Go runtime
   - Can have millions of goroutines
   - Cheaper to create than OS threads

2. **When should you use buffered vs unbuffered channels?**
   - Unbuffered: Synchronous communication
   - Buffered: Decouple sender/receiver, improve throughput

3. **What happens if you forget to call wg.Wait()?**
   - Program may exit before goroutines finish
   - Work may be incomplete

4. **Why do we need mutex for statistics?**
   - Multiple goroutines update same data
   - Without mutex: race conditions, data corruption

#### Architecture Questions

5. **How does the worker pool pattern work?**
   - Fixed number of workers
   - Workers pull from shared job queue
   - Limits concurrency

6. **How does rate limiting prevent server overload?**
   - Controls request rate
   - Spreads requests over time
   - Tokens limit concurrent operations

7. **What is the purpose of context in this project?**
   - Graceful shutdown on Ctrl+C
   - Request timeouts
   - Cancellation propagation

8. **How does graceful shutdown work?**
   - Signal handler catches Ctrl+C
   - Cancels context
   - Workers finish current jobs
   - Clean resource cleanup

## 🎯 Learning Objectives Checklist

After completing this guide, you should be able to:

### Concurrency Basics
- [ ] Create and manage goroutines
- [ ] Use channels for communication
- [ ] Understand buffered vs unbuffered channels
- [ ] Prevent goroutine leaks

### Synchronization
- [ ] Use sync.WaitGroup correctly
- [ ] Protect shared data with mutex
- [ ] Choose between Mutex and RWMutex
- [ ] Detect and fix race conditions

### Advanced Patterns
- [ ] Implement worker pool pattern
- [ ] Build rate limiters
- [ ] Use context for cancellation
- [ ] Handle graceful shutdown

### Go Specifics
- [ ] Parse HTML with goquery
- [ ] Make HTTP requests with context
- [ ] Use GORM for database operations
- [ ] Build CLI applications with flags
- [ ] Handle OS signals

## 🚀 Next Steps

### Immediate (This Week)

1. **Complete all exercises** (4-6 hours)
   - Start with easy ones
   - Build up to hard ones
   - Test thoroughly

2. **Add custom features** (2-4 hours)
   - Export to CSV
   - Add filtering options
   - Improve error handling

3. **Profile and optimize** (2-3 hours)
   ```bash
   # CPU profile
   go test -cpuprofile=cpu.prof

   # Memory profile
   go test -memprofile=mem.prof

   # View profiles
   go tool pprof cpu.prof
   ```

### Short-term (This Month)

1. **Add Testing**
   - Unit tests for parser
   - Integration tests for scraper
   - Mock HTTP responses
   - Test concurrency with race detector

2. **Improve Reliability**
   - Better error handling
   - More sophisticated retry logic
   - Circuit breaker pattern
   - Health checks

3. **Add Features**
   - Robots.txt support
   - Sitemap parsing
   - JavaScript rendering
   - Image downloading

### Long-term (Next Month)

1. **Distributed Scraping**
   - Use message queue (RabbitMQ, Kafka)
   - Distribute across multiple machines
   - Centralized job management

2. **Production Deployment**
   - Docker containerization
   - Kubernetes deployment
   - Monitoring with Prometheus
   - Logging with ELK stack

3. **Advanced Features**
   - GraphQL API for querying
   - Real-time dashboard
   - Machine learning for content extraction
   - Distributed storage

## 📚 Additional Resources

### Official Documentation
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go Blog - Concurrency Patterns](https://go.dev/blog/pipelines)
- [Go Memory Model](https://go.dev/ref/mem)

### Books
- "Concurrency in Go" by Katherine Cox-Buday
- "The Go Programming Language" by Donovan & Kernighan
- "Go in Action" by Kennedy, Ketelsen & Martin

### Video Courses
- [JustForFunc - Concurrency](https://www.youtube.com/c/JustForFunc)
- [GopherCon Talks](https://www.youtube.com/c/GopherAcademy)

### Tools
- Race detector: `go run -race`
- Profiler: `go tool pprof`
- Benchmarks: `go test -bench`

## 💡 Tips for Success

### 1. Start Simple
Don't try to understand everything at once. Start with:
- main.go - See the flow
- scraper/worker.go - See the pattern
- Gradually explore other files

### 2. Draw Diagrams
Visualize:
- Goroutine communication
- Channel flow
- Component relationships

### 3. Add Logging
See what's happening:
```go
log.Printf("Worker %d: starting job %s", id, job.URL)
log.Printf("Worker %d: completed in %s", id, duration)
```

### 4. Use the Race Detector
Find concurrency bugs:
```bash
go run -race main.go
```

### 5. Experiment
- Change worker counts
- Modify buffer sizes
- Add artificial delays
- Break things and fix them

### 6. Read Error Messages
Go's error messages are helpful:
- "fatal error: all goroutines are asleep - deadlock!"
- "panic: send on closed channel"
- Learn what they mean

### 7. Use the Debugger
Step through concurrent code:
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug main.go
```

## 🎉 Congratulations!

You've built a production-ready concurrent web scraper with:
- ✅ Worker pool pattern
- ✅ Rate limiting
- ✅ Graceful shutdown
- ✅ Database persistence
- ✅ Error handling and retries
- ✅ Real-time statistics

This is a significant achievement! You now understand:
- How to write concurrent Go code
- How to coordinate goroutines
- How to build scalable systems
- How to handle real-world challenges

**Keep coding, keep learning, and most importantly - have fun!** 🚀

---

**Questions? Stuck? Tips:**
1. Re-read the relevant section above
2. Check the inline code comments
3. Run with -race to find concurrency bugs
4. Add logging to see what's happening
5. Join the Go community on Reddit/Slack

**Ready to code?** Head to [QUICK_START.md](QUICK_START.md) to run the project!
