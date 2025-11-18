# 🚀 Quick Start Guide

Get the Concurrent Web Scraper running in **2 minutes**!

## Prerequisites

- Go 1.21 or higher installed
- Terminal/Command prompt
- Internet connection

## Steps

### 1. Navigate to Project Directory

```bash
cd week3-projects/01-concurrent-web-scraper
```

### 2. Install Dependencies

```bash
go mod download
```

This will download:
- `goquery` - HTML parsing
- `gorm` - Database ORM
- `sqlite` - Database driver

### 3. Run the Scraper

```bash
go run main.go
```

You should see:

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

Worker 1: scraping https://go.dev/ (attempt 1)
✓ Success: https://go.dev/ (status: 200, title: The Go Programming Language, links: 42, duration: 523ms)
...

✓ Scraping completed!
```

## Common Commands

### Basic Scraping

```bash
# Use default settings (5 workers, 2 req/s)
go run main.go
```

### Custom Configuration

```bash
# 10 workers, 5 requests per second
go run main.go -workers 10 -rate 5.0

# Aggressive scraping (careful!)
go run main.go -workers 20 -rate 10.0 -retries 5

# Slow and steady
go run main.go -workers 2 -rate 0.5
```

### View Statistics

```bash
# Show database statistics
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
# Delete all scraped data
go run main.go -clear
```

### Use Custom URLs

Create your own `my-urls.json`:
```json
{
  "urls": [
    "https://example.com",
    "https://your-site.com"
  ]
}
```

Run with custom file:
```bash
go run main.go -urls my-urls.json
```

## Build Executable

### Linux/Mac

```bash
go build -o scraper main.go
./scraper -workers 5
```

### Windows

```bash
go build -o scraper.exe main.go
scraper.exe -workers 5
```

### Cross-Platform Build

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o scraper-linux main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o scraper.exe main.go

# Build for Mac
GOOS=darwin GOARCH=amd64 go build -o scraper-mac main.go
```

## Quick Examples

### Example 1: Fast Scraping

```bash
# 20 workers, 10 req/s - scrapes quickly
go run main.go -workers 20 -rate 10.0
```

**Use when:** You control the target servers or they can handle high load

### Example 2: Respectful Scraping

```bash
# 2 workers, 0.5 req/s - slow and respectful
go run main.go -workers 2 -rate 0.5
```

**Use when:** Scraping third-party sites, being a good netizen

### Example 3: With Custom Settings

```bash
# Full customization
go run main.go \
  -workers 10 \
  -rate 3.0 \
  -retries 5 \
  -timeout 60 \
  -db ./my-data/scraper.db \
  -urls ./my-urls.json \
  -user-agent "MyBot/2.0"
```

## All Available Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-workers` | 5 | Number of concurrent workers |
| `-rate` | 2.0 | Requests per second |
| `-retries` | 3 | Max retry attempts |
| `-timeout` | 30 | Request timeout (seconds) |
| `-db` | ./data/scraper.db | Database path |
| `-urls` | ./urls.json | URLs file path |
| `-user-agent` | GoWebScraper/1.0 | User agent string |
| `-stats` | false | Show stats only |
| `-clear` | false | Clear database |

## Understanding the Output

### Success Message
```
✓ Success: https://go.dev/ (status: 200, title: The Go Programming Language, links: 42, duration: 523ms)
```
- ✓ = Successfully scraped
- status: HTTP status code
- title: Page title (truncated to 50 chars)
- links: Number of links found
- duration: How long the request took

### Failure Message
```
✗ Failed: https://example.com - Get "https://example.com": dial tcp: lookup example.com: no such host
```
- ✗ = Failed to scrape
- Shows the error message

### Retry Message
```
↻ Retrying: https://example.com (attempt 2)
```
- ↻ = Retrying after failure
- Shows which attempt number

## Tips

### 1. Start Small
Test with a few URLs first:
```json
{
  "urls": [
    "https://go.dev/",
    "https://go.dev/doc/"
  ]
}
```

### 2. Monitor the Output
Watch for:
- Success/failure ratio
- Response times
- Error messages

### 3. Adjust Rate Based on Results
- Too many failures? Slow down: `-rate 1.0`
- All successful? Speed up: `-rate 5.0`

### 4. Use Stats to Track Progress
```bash
# After scraping
go run main.go -stats
```

### 5. Stop Anytime
Press `Ctrl+C` to stop gracefully. Current jobs will complete.

## Troubleshooting

### "No URLs found"
**Problem:** urls.json is empty or malformed
**Solution:** Check urls.json format, ensure it has valid JSON

### "Failed to connect to database"
**Problem:** Cannot create database file
**Solution:** Check write permissions, ensure parent directory exists

### "Too many open files"
**Problem:** System file descriptor limit reached
**Solution:** Reduce workers: `-workers 3`

### "All requests failing"
**Problem:** Network issues or blocked
**Solution:**
- Check internet connection
- Try lower rate: `-rate 0.5`
- Check if sites are accessible in browser

### "Context cancelled"
**Problem:** You pressed Ctrl+C
**Solution:** This is normal - scraper is shutting down gracefully

## What's Next?

Now that you have the scraper running:

1. **Read the full README** - [README.md](README.md)
   - Understand the architecture
   - Learn about design patterns
   - See enhancement ideas

2. **Study the code** - [START_HERE.md](START_HERE.md)
   - Follow the learning guide
   - Understand concurrency patterns
   - Complete exercises

3. **Experiment:**
   - Try different worker counts
   - Adjust rate limits
   - Add your own URLs
   - Monitor database growth

4. **Customize:**
   - Modify parsing logic
   - Add new features
   - Improve error handling
   - Export data to CSV

## Quick Reference

### Standard Run
```bash
go run main.go
```

### Show Help
```bash
go run main.go -h
```

### View Stats
```bash
go run main.go -stats
```

### Clear Data
```bash
go run main.go -clear
```

### Custom Config
```bash
go run main.go -workers 10 -rate 5.0 -urls custom.json
```

---

**You're ready to scrape!** 🚀

For complete documentation, see [README.md](README.md)
For learning guidance, see [START_HERE.md](START_HERE.md)
