package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/user/web-scraper/config"
	"github.com/user/web-scraper/ratelimiter"
	"github.com/user/web-scraper/scraper"
	"github.com/user/web-scraper/storage"
)

var (
	// Command line flags
	workersFlag   = flag.Int("workers", 5, "Number of concurrent workers")
	rateFlag      = flag.Float64("rate", 2.0, "Requests per second")
	retriesFlag   = flag.Int("retries", 3, "Maximum retry attempts")
	timeoutFlag   = flag.Int("timeout", 30, "Request timeout in seconds")
	dbPathFlag    = flag.String("db", "./data/scraper.db", "Database file path")
	urlsFileFlag  = flag.String("urls", "./urls.json", "URLs file path")
	configFlag    = flag.String("config", "", "Config file path (overrides other flags)")
	statsFlag     = flag.Bool("stats", false, "Show database statistics and exit")
	clearFlag     = flag.Bool("clear", false, "Clear all data from database")
	userAgentFlag = flag.String("user-agent", "GoWebScraper/1.0", "User agent string")
)

func main() {
	flag.Parse()

	// Print banner
	printBanner()

	// Load or create configuration
	cfg := loadConfiguration()

	// Initialize database
	db, err := storage.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Handle stats flag
	if *statsFlag {
		showStatistics(db)
		return
	}

	// Handle clear flag
	if *clearFlag {
		if err := db.DeleteAllPages(); err != nil {
			log.Fatalf("Failed to clear database: %v", err)
		}
		fmt.Println("✓ Database cleared successfully")
		return
	}

	// Load URLs
	urls, err := config.LoadURLs(cfg.URLsFile)
	if err != nil {
		log.Fatalf("Failed to load URLs from %s: %v", cfg.URLsFile, err)
	}

	if len(urls) == 0 {
		log.Fatalf("No URLs found in %s", cfg.URLsFile)
	}

	fmt.Printf("\n📋 Loaded %d URLs to scrape\n", len(urls))
	fmt.Printf("⚙️  Configuration:\n")
	fmt.Printf("   Workers: %d\n", cfg.WorkerCount)
	fmt.Printf("   Rate Limit: %.1f req/s\n", cfg.RateLimit)
	fmt.Printf("   Max Retries: %d\n", cfg.MaxRetries)
	fmt.Printf("   Timeout: %ds\n", cfg.RequestTimeout)
	fmt.Printf("   Database: %s\n\n", cfg.DatabasePath)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\n🛑 Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	// Create rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(ctx, cfg.RateLimit)
	defer rateLimiter.Stop()

	// Create scraper
	scraperInstance := scraper.NewScraper(
		time.Duration(cfg.RequestTimeout)*time.Second,
		cfg.UserAgent,
		cfg.MaxRetries,
		cfg.FollowRedirects,
	)
	defer scraperInstance.Close()

	// Create worker pool
	pool := scraper.NewWorkerPool(ctx, cfg.WorkerCount, scraperInstance, db, rateLimiter)

	// Start worker pool
	pool.Start()

	// Add jobs
	fmt.Println("🚀 Starting scraping...\n")
	startTime := time.Now()
	pool.AddJobs(urls)

	// Wait for completion or cancellation
	pool.Wait()

	// Print final statistics
	fmt.Println("\n✓ Scraping completed!")
	duration := time.Since(startTime)
	fmt.Printf("⏱️  Total time: %s\n\n", duration.Round(time.Millisecond))

	pool.PrintStatistics()

	// Show database statistics
	fmt.Println("\n📊 Database Statistics:")
	showStatistics(db)
}

// loadConfiguration loads configuration from file or flags
func loadConfiguration() *config.Config {
	var cfg *config.Config

	if *configFlag != "" {
		// Load from config file
		var err error
		cfg, err = config.LoadFromFile(*configFlag)
		if err != nil {
			log.Fatalf("Failed to load config from %s: %v", *configFlag, err)
		}
		fmt.Printf("✓ Loaded configuration from %s\n", *configFlag)
	} else {
		// Use flags
		cfg = &config.Config{
			WorkerCount:     *workersFlag,
			RateLimit:       *rateFlag,
			MaxRetries:      *retriesFlag,
			RequestTimeout:  *timeoutFlag,
			DatabasePath:    *dbPathFlag,
			URLsFile:        *urlsFileFlag,
			UserAgent:       *userAgentFlag,
			FollowRedirects: true,
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	return cfg
}

// showStatistics displays database statistics
func showStatistics(db *storage.Database) {
	stats, err := db.GetStatistics()
	if err != nil {
		log.Printf("Failed to get statistics: %v", err)
		return
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Pages:      %d\n", stats.TotalPages)
	fmt.Printf("Successful:       %d\n", stats.SuccessfulPages)
	fmt.Printf("Failed:           %d\n", stats.FailedPages)
	fmt.Printf("Total Links:      %d\n", stats.TotalLinks)
	if stats.TotalPages > 0 {
		fmt.Printf("Average Duration: %s\n", stats.AverageDuration)
		fmt.Printf("Total Duration:   %s\n", stats.TotalDuration)
	}
	if !stats.StartTime.IsZero() {
		fmt.Printf("First Scraped:    %s\n", stats.StartTime.Format("2006-01-02 15:04:05"))
	}
	if !stats.EndTime.IsZero() {
		fmt.Printf("Last Scraped:     %s\n", stats.EndTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Println(strings.Repeat("=", 60))
}

// printBanner prints the application banner
func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║          🌐  Concurrent Web Scraper  🌐                      ║
║                                                              ║
║          Week 3 - Go Learning Path Project                  ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
