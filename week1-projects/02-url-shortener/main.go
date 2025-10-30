package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jason/url-shortener/shortener"
	"github.com/jason/url-shortener/storage"
)

// Global variables for the shortener and storage
var (
	us  *shortener.URLShortener
	st  *storage.Storage
	dbFile = "urls.json"
)

func main() {
	// Check if at least one argument is provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Initialize shortener and storage
	us = shortener.NewURLShortener(6) // 6-character short codes
	st = storage.NewStorage(dbFile)

	// Load existing mappings from file
	mappings, err := st.LoadMappings()
	if err != nil {
		log.Fatalf("Failed to load mappings: %v", err)
	}

	// Restore mappings into the shortener
	// This demonstrates: rebuilding state from persistent storage
	for _, mapping := range mappings {
		us.Mappings[mapping.ShortCode] = mapping
	}

	// Get the command (first argument)
	command := os.Args[1]

	// Handle different commands
	// This demonstrates: SWITCH statements and command routing
	switch command {
	case "shorten":
		handleShorten()
	case "get":
		handleGet()
	case "list":
		handleList()
	case "stats":
		handleStats()
	case "delete":
		handleDelete()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

// handleShorten creates a new shortened URL
func handleShorten() {
	// Parse flags for the shorten command
	// This allows: url-shortener shorten https://example.com
	flagSet := flag.NewFlagSet("shorten", flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Println("Usage: url-shortener shorten <URL>")
		fmt.Println("\nExample: url-shortener shorten https://www.google.com")
	}

	// Check if URL argument is provided
	if len(os.Args) < 3 {
		flagSet.Usage()
		os.Exit(1)
	}

	url := os.Args[2]

	// Shorten the URL
	shortCode, err := us.ShortenURL(url)
	if err != nil {
		log.Fatalf("Failed to shorten URL: %v", err)
	}

	// Get the mapping for saving
	mapping, _ := us.GetMapping(shortCode)

	// Save to persistent storage
	if err := st.AppendMapping(mapping); err != nil {
		log.Fatalf("Failed to save mapping: %v", err)
	}

	// Display results
	fmt.Printf("\n✅ URL Shortened Successfully!\n\n")
	fmt.Printf("Original URL: %s\n", url)
	fmt.Printf("Short Code:  %s\n", shortCode)
	fmt.Printf("Short URL:   http://short.url/%s\n\n", shortCode)
}

// handleGet retrieves the original URL for a short code
func handleGet() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: url-shortener get <SHORT_CODE>")
		fmt.Println("Example: url-shortener get abc123")
		os.Exit(1)
	}

	shortCode := os.Args[2]

	// Get the original URL
	url, err := us.GetURL(shortCode)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Update the storage with new visit count
	mapping, _ := us.GetMapping(shortCode)
	st.UpdateMapping(mapping)

	fmt.Printf("\n✅ Found URL!\n\n")
	fmt.Printf("Short Code:  %s\n", shortCode)
	fmt.Printf("Original URL: %s\n", url)
	fmt.Printf("Visits:      %d\n\n", mapping.Visits)
}

// handleList displays all shortened URLs
func handleList() {
	mappings := us.ListAllMappings()

	if len(mappings) == 0 {
		fmt.Println("No URLs shortened yet.")
		return
	}

	fmt.Printf("\n📋 All Shortened URLs (%d total)\n\n", len(mappings))

	// Create a table writer for neat output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CODE\tORIGINAL URL\tVISITS\tCREATED")

	// Display each mapping in the table
	for _, m := range mappings {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", m.ShortCode, truncateURL(m.OriginalURL, 40), m.Visits, m.CreatedAt)
	}

	w.Flush()
	fmt.Println()
}

// handleStats displays statistics for a shortened URL
func handleStats() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: url-shortener stats <SHORT_CODE>")
		fmt.Println("Example: url-shortener stats abc123")
		os.Exit(1)
	}

	shortCode := os.Args[2]

	// Get statistics
	stats, err := us.GetStats(shortCode)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n📊 Statistics for %s\n\n", shortCode)
	fmt.Printf("Original URL:     %s\n", stats["original_url"])
	fmt.Printf("Created:          %s\n", stats["created_at"])
	fmt.Printf("Visits:           %v\n", stats["visits"])
	fmt.Printf("Original Length:  %v characters\n", stats["url_length"])
	fmt.Printf("Code Length:      %v characters\n", stats["code_length"])
	fmt.Printf("Compression:      %.2fx\n\n", stats["compression"])

	// Save updated stats
	mapping, _ := us.GetMapping(shortCode)
	st.UpdateMapping(mapping)
}

// handleDelete removes a shortened URL
func handleDelete() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: url-shortener delete <SHORT_CODE>")
		fmt.Println("Example: url-shortener delete abc123")
		os.Exit(1)
	}

	shortCode := os.Args[2]

	// Delete from shortener
	if err := us.DeleteURL(shortCode); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Delete from storage
	if err := st.RemoveMapping(shortCode); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Deleted short code: %s\n\n", shortCode)
}

// printUsage displays the help text
func printUsage() {
	fmt.Println(`
╔════════════════════════════════════════════════════════════════╗
║            URL Shortener - Command Line Tool                  ║
╚════════════════════════════════════════════════════════════════╝

USAGE:
  url-shortener <command> [arguments]

COMMANDS:
  shorten <URL>       Create a shortened URL
                      Example: url-shortener shorten https://google.com

  get <CODE>          Get the original URL for a short code
                      Example: url-shortener get abc123

  list                List all shortened URLs
                      Example: url-shortener list

  stats <CODE>        Show statistics for a shortened URL
                      Example: url-shortener stats abc123

  delete <CODE>       Delete a shortened URL
                      Example: url-shortener delete abc123

  help                Show this help message

EXAMPLES:
  # Create a short URL
  $ url-shortener shorten https://www.example.com/very/long/url

  # Look up what the short code points to
  $ url-shortener get abc123

  # See all your shortened URLs
  $ url-shortener list

  # Get detailed stats about a short URL
  $ url-shortener stats abc123

  # Delete a short URL
  $ url-shortener delete abc123

DATA:
  All shortened URLs are stored in: urls.json

FEATURES:
  ✓ Create shortened URLs with 6-character codes
  ✓ Track visit counts for each URL
  ✓ Store all data in JSON format
  ✓ Persistent storage (survives app restart)
  ✓ View statistics for any shortened URL
  ✓ Delete shortened URLs when no longer needed

`)
}

// truncateURL shortens a long URL for display
// Demonstrates STRING MANIPULATION
func truncateURL(url string, maxLen int) string {
	if len(url) <= maxLen {
		return url
	}
	// Return first maxLen-3 characters plus "..."
	return url[:maxLen-3] + "..."
}
