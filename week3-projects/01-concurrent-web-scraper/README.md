# 🌐 Concurrent Web Scraper

A high-performance concurrent web scraper built with Go featuring worker pools, rate limiting, and persistent storage.

## 🎯 Week 3 Project Overview

This project is part of the **1-Month Go Learning Path - Week 3: Concurrency & Advanced Patterns**

### Learning Objectives

- ✅ Master goroutines and channels
- ✅ Implement worker pool pattern
- ✅ Use sync primitives (WaitGroup, Mutex)
- ✅ Implement rate limiting with time.Ticker
- ✅ Handle context for cancellation
- ✅ Work with concurrent database operations
- ✅ Parse HTML with goquery
- ✅ Build production-ready CLI applications

## ✨ Features

### Core Features
- 🔄 **Concurrent Scraping** - Multiple goroutines scrape websites simultaneously
- 👷 **Worker Pool Pattern** - Configurable number of workers to control concurrency
- ⏱️ **Rate Limiting** - Token bucket algorithm to control request rate
- 💾 **SQLite Storage** - Persistent storage with GORM ORM
- 🔁 **Retry Logic** - Automatic retry for failed requests
- 📊 **Progress Tracking** - Real-time statistics and progress updates
- 🛑 **Graceful Shutdown** - Clean shutdown on Ctrl+C
- 🎯 **HTML Parsing** - Extract title, description, and links

### Technical Features
- ✅ Worker pool with buffered channels
- ✅ Context-based cancellation
- ✅ sync.WaitGroup for coordination
- ✅ Mutex-protected statistics
- ✅ Rate limiter with time.Ticker
- ✅ Concurrent-safe database operations
- ✅ Signal handling (SIGINT, SIGTERM)
- ✅ Comprehensive error handling
- ✅ CLI with flags for configuration

## 🏗️ Project Structure

```
01-concurrent-web-scraper/
├── main.go                  # Application entry point, CLI interface
├── go.mod                   # Go module dependencies
├── urls.json                # Example URLs to scrape
├── QUICK_START.md          # Get started in 2 minutes
├── README.md               # This file
├── START_HERE.md           # Learning guide
│
├── config/                 # Configuration management
│   └── config.go          # Config struct, loading, validation
│
├── models/                # Data models
│   └── models.go          # Database models and DTOs
│
├── scraper/               # Core scraping logic
│   ├── scraper.go        # HTTP client and scraping operations
│   ├── worker.go         # Worker pool implementation
│   └── parser.go         # HTML parsing with goquery
│
├── storage/               # Data persistence
│   └── database.go       # SQLite database operations with GORM
│
├── ratelimiter/          # Rate limiting
│   └── limiter.go        # Token bucket rate limiter
│
└── data/                 # Database storage (created automatically)
    └── scraper.db        # SQLite database file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Internet connection (for scraping)

### Installation

1. **Navigate to the project directory:**
   ```bash
   cd week3-projects/01-concurrent-web-scraper
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the scraper:**
   ```bash
   go run main.go
   ```

### Building

To build an executable:

```bash
# Linux/Mac
go build -o scraper main.go

# Windows
go build -o scraper.exe main.go
```

Then run:
```bash
./scraper        # Linux/Mac
scraper.exe      # Windows
```

## 📚 Usage

### Basic Usage

Scrape URLs from `urls.json`:
```bash
go run main.go
```

### Command Line Flags

```bash
# Custom configuration
go run main.go -workers 10 -rate 5.0 -retries 5

# Use custom URLs file
go run main.go -urls custom-urls.json

# Custom database path
go run main.go -db ./data/my-scraper.db

# Full options
go run main.go \
  -workers 10 \
  -rate 5.0 \
  -retries 3 \
  -timeout 30 \
  -db ./data/scraper.db \
  -urls ./urls.json \
  -user-agent "MyBot/1.0"
```

### Available Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-workers` | 5 | Number of concurrent workers |
| `-rate` | 2.0 | Requests per second (rate limit) |
| `-retries` | 3 | Maximum retry attempts per URL |
| `-timeout` | 30 | Request timeout in seconds |
| `-db` | ./data/scraper.db | Database file path |
| `-urls` | ./urls.json | URLs file path |
| `-user-agent` | GoWebScraper/1.0 | User agent string |
| `-stats` | false | Show database statistics and exit |
| `-clear` | false | Clear all data from database |

### View Statistics

```bash
go run main.go -stats
```

Output:
```
============================================================
Total Pages:      10
Successful:       8
Failed:           2
Total Links:      247
Average Duration: 523ms
Total Duration:   4.8s
First Scraped:    2024-01-15 10:30:45
Last Scraped:     2024-01-15 10:31:12
============================================================
```

### Clear Database

```bash
go run main.go -clear
```

## 📋 URLs File Format

Create a JSON file with URLs to scrape:

```json
{
  "urls": [
    "https://example.com",
    "https://another-site.com",
    "https://go.dev/"
  ]
}
```

## 🔧 Configuration

### Using Flags (Recommended)

Pass configuration via command-line flags:
```bash
go run main.go -workers 10 -rate 3.0
```

### Using Config File

Create a config file (`config.json`):
```json
{
  "worker_count": 10,
  "rate_limit": 3.0,
  "max_retries": 5,
  "request_timeout": 30,
  "database_path": "./data/scraper.db",
  "urls_file": "./urls.json",
  "user_agent": "MyBot/1.0",
  "follow_redirects": true
}
```

Use it:
```bash
go run main.go -config config.json
```

### Configuration Limits

The configuration is validated with sensible limits:
- Workers: 1-100
- Rate limit: 0.1-100 req/s
- Max retries: 0-10
- Timeout: 1-300 seconds

## 🎯 How It Works

### Architecture Overview

```
                     ┌─────────────────┐
                     │   Main Goroutine│
                     │   (CLI Interface)│
                     └────────┬────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
            ┌───────▼────────┐   ┌──────▼───────┐
            │  Worker Pool   │   │ Rate Limiter │
            │   (Manages)    │   │ (Controls)   │
            └───────┬────────┘   └──────┬───────┘
                    │                   │
        ┌───────────┼───────────────────┼──────────┐
        │           │                   │          │
   ┌────▼────┐ ┌────▼────┐ ┌────▼────┐ ┌───▼────┐
   │Worker 1 │ │Worker 2 │ │Worker 3 │ │   ...  │
   └────┬────┘ └────┬────┘ └────┬────┘ └───┬────┘
        │           │           │           │
        └───────────┴───────────┴───────────┘
                        │
                 ┌──────▼──────┐
                 │   Results   │
                 │  Processor  │
                 └──────┬──────┘
                        │
                 ┌──────▼──────┐
                 │  Database   │
                 │  (SQLite)   │
                 └─────────────┘
```

### Workflow

1. **Initialization**
   - Load configuration
   - Initialize database
   - Create rate limiter
   - Set up worker pool

2. **Job Distribution**
   - URLs are sent to jobs channel
   - Workers pick up jobs from channel
   - Rate limiter controls request rate

3. **Scraping**
   - Worker makes HTTP request
   - Parses HTML content
   - Extracts title, description, links
   - Sends result to results channel

4. **Result Processing**
   - Results processor saves to database
   - Handles retry logic for failures
   - Updates statistics

5. **Graceful Shutdown**
   - Ctrl+C triggers context cancellation
   - Workers finish current jobs
   - Clean resource cleanup

## 🧩 Key Components

### 1. Worker Pool

Manages concurrent workers using channels:

```go
// Create worker pool
pool := scraper.NewWorkerPool(ctx, workerCount, scraperInstance, db, rateLimiter)

// Start workers
pool.Start()

// Add jobs
pool.AddJobs(urls)

// Wait for completion
pool.Wait()
```

**Key Features:**
- Buffered job channel for efficient queuing
- Configurable number of workers
- Graceful shutdown with context
- Real-time statistics tracking

### 2. Rate Limiter

Token bucket algorithm to control request rate:

```go
// Create rate limiter (2 requests per second)
limiter := ratelimiter.NewRateLimiter(ctx, 2.0)

// Wait for token before making request
if err := limiter.Wait(ctx); err != nil {
    return err
}
```

**How it works:**
- Tokens are generated at specified rate
- Each request consumes one token
- Request blocks if no tokens available
- Prevents server overload

### 3. HTML Parser

Extracts structured data from HTML:

```go
parser := scraper.NewParser()
title, description, links, err := parser.ParseHTML(html, baseURL)
```

**Extracts:**
- Page title (from `<title>`, `og:title`, or `<h1>`)
- Meta description (from meta tags or first paragraph)
- All links with text

### 4. Database Storage

Concurrent-safe database operations:

```go
db, err := storage.NewDatabase(dbPath)
err = db.SavePage(result)
stats, err := db.GetStatistics()
```

**Features:**
- Auto-migration of schema
- Mutex-protected operations
- Duplicate URL handling
- Statistics aggregation

## 🔄 Concurrency Patterns

### Worker Pool Pattern

```go
// Create job channel
jobs := make(chan ScrapeJob, bufferSize)

// Start workers
for i := 0; i < workerCount; i++ {
    wg.Add(1)
    go worker(jobs)
}

// Send jobs
for _, url := range urls {
    jobs <- ScrapeJob{URL: url}
}

// Close channel and wait
close(jobs)
wg.Wait()
```

### Context for Cancellation

```go
// Create cancellable context
ctx, cancel := context.WithCancel(context.Background())

// Handle signals
go func() {
    <-sigChan
    cancel() // Cancels all operations
}()

// Use in operations
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
```

### Synchronization with WaitGroup

```go
var wg sync.WaitGroup

// Start workers
for i := 0; i < count; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // Do work
    }()
}

// Wait for all to complete
wg.Wait()
```

### Thread-Safe Statistics

```go
type Statistics struct {
    mu            sync.Mutex
    CompletedJobs int
}

func (s *Statistics) Increment() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.CompletedJobs++
}
```

## 📊 Database Schema

### ScrapedPage Table

| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER | Primary key |
| url | TEXT | Page URL (unique) |
| title | TEXT | Page title |
| description | TEXT | Meta description |
| link_count | INTEGER | Number of links found |
| status_code | INTEGER | HTTP status code |
| error | TEXT | Error message if failed |
| retry_count | INTEGER | Number of retry attempts |
| scraped_at | DATETIME | When scraped |
| duration | INTEGER | Request duration (ms) |
| created_at | DATETIME | Record creation time |
| updated_at | DATETIME | Last update time |

### Link Table

| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER | Primary key |
| page_id | INTEGER | Foreign key to scraped_page |
| url | TEXT | Link URL |
| text | TEXT | Link text |
| created_at | DATETIME | Record creation time |

## 🧪 Example Output

```
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║          🌐  Concurrent Web Scraper  🌐                      ║
║                                                              ║
║          Week 3 - Go Learning Path Project                  ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

✓ Database initialized at: ./data/scraper.db

📋 Loaded 10 URLs to scrape
⚙️  Configuration:
   Workers: 5
   Rate Limit: 2.0 req/s
   Max Retries: 3
   Timeout: 30s
   Database: ./data/scraper.db

🚀 Starting scraping...

Worker 1 started
Worker 2 started
Worker 3 started
Worker 4 started
Worker 5 started
Result processor started

Worker 1: scraping https://go.dev/ (attempt 1)
Worker 2: scraping https://go.dev/doc/ (attempt 1)
Worker 3: scraping https://go.dev/blog/ (attempt 1)

✓ Success: https://go.dev/ (status: 200, title: The Go Programming Language, links: 42, duration: 523ms)
✓ Success: https://go.dev/doc/ (status: 200, title: Documentation - The Go Programming Lang..., links: 38, duration: 445ms)
✓ Success: https://go.dev/blog/ (status: 200, title: Go Blog - The Go Programming Language, links: 56, duration: 612ms)

...

✓ Scraping completed!
⏱️  Total time: 5.2s

============================================================
Scraping Statistics
============================================================
Total Jobs:       10
Completed:        10
Successful:       8
Failed:           2
In Progress:      0
Total Retries:    3
============================================================

📊 Database Statistics:
============================================================
Total Pages:      10
Successful:       8
Failed:           2
Total Links:      247
Average Duration: 523ms
Total Duration:   4.8s
First Scraped:    2024-01-15 10:30:45
Last Scraped:     2024-01-15 10:31:12
============================================================
```

## 🏛️ Architecture & Design Patterns

### Patterns Demonstrated

1. **Worker Pool Pattern**
   - Fixed number of workers
   - Job queue with channels
   - Efficient resource utilization

2. **Producer-Consumer Pattern**
   - Main goroutine produces jobs
   - Workers consume and process
   - Results processor handles output

3. **Token Bucket Rate Limiting**
   - Tokens generated at fixed rate
   - Requests consume tokens
   - Smooth traffic distribution

4. **Graceful Degradation**
   - Retry failed requests
   - Continue on errors
   - Comprehensive error logging

5. **Repository Pattern**
   - Database abstraction
   - Clean separation of concerns
   - Easy to test and modify

## 🔒 Best Practices

### Concurrency
- ✅ Use channels for communication
- ✅ Use sync.WaitGroup for coordination
- ✅ Protect shared state with mutexes
- ✅ Use context for cancellation
- ✅ Avoid goroutine leaks

### Error Handling
- ✅ Always check errors
- ✅ Provide context in error messages
- ✅ Log errors appropriately
- ✅ Implement retry logic
- ✅ Fail gracefully

### Resource Management
- ✅ Close resources properly (defer)
- ✅ Use buffered channels appropriately
- ✅ Limit concurrent connections
- ✅ Implement timeouts
- ✅ Handle signals for cleanup

## 🎓 Learning Takeaways

After completing this project, you'll understand:

### Concurrency
- ✅ How to create and manage goroutines
- ✅ Channel patterns and best practices
- ✅ Worker pool implementation
- ✅ Synchronization with WaitGroup and Mutex
- ✅ Context usage for cancellation

### System Design
- ✅ Rate limiting algorithms
- ✅ Retry strategies
- ✅ Graceful shutdown patterns
- ✅ Resource pooling
- ✅ Error propagation

### Go Specifics
- ✅ HTTP client usage
- ✅ HTML parsing with goquery
- ✅ GORM database operations
- ✅ CLI with flag package
- ✅ Signal handling

## 🚀 Enhancement Ideas

1. **Add More Features:**
   - Sitemap generation
   - Robots.txt parsing
   - JavaScript rendering
   - Image downloading
   - CSV/JSON export

2. **Improve Performance:**
   - Connection pooling
   - HTTP/2 support
   - Compression handling
   - Caching responses
   - Batch database inserts

3. **Add Monitoring:**
   - Prometheus metrics
   - Request/response logging
   - Performance profiling
   - Memory usage tracking
   - Real-time dashboard

4. **Production Features:**
   - Distributed scraping
   - PostgreSQL support
   - Redis queue
   - Webhook notifications
   - API endpoints

## 🐛 Common Issues

### Too Many Open Files
**Problem:** Running out of file descriptors
**Solution:** Reduce worker count or increase system limits
```bash
ulimit -n 4096  # Increase file descriptor limit
```

### Rate Limited by Server
**Problem:** Server blocks requests
**Solution:** Decrease rate limit, add delays
```bash
go run main.go -rate 0.5  # Slower rate
```

### Database Locked
**Problem:** SQLite database lock errors
**Solution:** Ensure proper mutex usage, reduce concurrency

### Context Cancelled
**Problem:** Operations cancelled unexpectedly
**Solution:** Check for interrupt signals, implement proper cleanup

## 📝 License

This is a learning project. Feel free to use it for educational purposes.

## 🙏 Acknowledgments

- Built as part of the 1-Month Go Learning Path
- goquery: https://github.com/PuerkitoBio/goquery
- GORM: https://gorm.io/
- Go Standard Library: https://pkg.go.dev/std

---

**Happy Scraping!** 🚀

For a quick start, see [QUICK_START.md](QUICK_START.md)
For learning guidance, see [START_HERE.md](START_HERE.md)
