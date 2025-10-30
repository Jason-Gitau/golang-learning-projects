package shortener

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// URLMapping represents a shortened URL and its metadata
// Demonstrates STRUCTS in Go - grouping related data
type URLMapping struct {
	// ShortCode is the unique identifier (e.g., "abc123")
	ShortCode string `json:"short_code"`
	// OriginalURL is the full URL that was shortened
	OriginalURL string `json:"original_url"`
	// CreatedAt is when the URL was shortened
	CreatedAt string `json:"created_at"`
	// Visits tracks how many times this short URL was accessed
	Visits int `json:"visits"`
}

// URLShortener manages the creation and retrieval of shortened URLs
// Demonstrates METHODS on STRUCTS
type URLShortener struct {
	// Mappings stores all URL mappings in memory
	// Key: short code, Value: URL mapping
	// This demonstrates MAPS in Go for fast lookups
	Mappings map[string]*URLMapping
	// CodeLength is how long the generated short codes should be
	CodeLength int
	// charset is the characters used to generate short codes
	charset string
}

// NewURLShortener creates a new URLShortener instance
// Demonstrates FUNCTIONS that return POINTERS
// Returns a pointer (*URLShortener) for efficiency
func NewURLShortener(codeLength int) *URLShortener {
	return &URLShortener{
		// make() initializes the map with capacity 0
		// Go maps grow dynamically as needed
		Mappings:   make(map[string]*URLMapping),
		CodeLength: codeLength,
		// alphanumeric characters for URL-safe short codes
		charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
}

// ShortenURL creates a shortened version of a URL
// Demonstrates ERROR HANDLING - returning (value, error)
// This is Go's idiomatic way to handle errors
func (us *URLShortener) ShortenURL(originalURL string) (string, error) {
	// Validate the input
	if err := us.validateURL(originalURL); err != nil {
		return "", err
	}

	// Check if this URL is already shortened
	// This demonstrates iterating over a map with for...range
	for code, mapping := range us.Mappings {
		if mapping.OriginalURL == originalURL {
			// URL already exists, return existing code
			return code, nil
		}
	}

	// Generate a new unique short code
	// Loop until we find a code that hasn't been used
	var shortCode string
	var attempts int
	maxAttempts := 100 // Prevent infinite loops

	for attempts < maxAttempts {
		shortCode = us.generateCode()
		// Check if this code is already taken using comma-ok idiom
		if _, exists := us.Mappings[shortCode]; !exists {
			break // Found an unused code
		}
		attempts++
	}

	if attempts >= maxAttempts {
		return "", fmt.Errorf("failed to generate unique short code after %d attempts", maxAttempts)
	}

	// Create the mapping entry
	mapping := &URLMapping{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Visits:      0,
	}

	// Store in the map
	us.Mappings[shortCode] = mapping

	return shortCode, nil
}

// GetURL retrieves the original URL for a short code
// Demonstrates MAP LOOKUPS and ERROR HANDLING
func (us *URLShortener) GetURL(shortCode string) (string, error) {
	// Validate input
	if shortCode == "" {
		return "", fmt.Errorf("short code cannot be empty")
	}

	// Look up the mapping using comma-ok idiom
	// This safely checks if the key exists and gets the value
	mapping, exists := us.Mappings[shortCode]
	if !exists {
		return "", fmt.Errorf("short code not found: %s", shortCode)
	}

	// Increment visit counter
	mapping.Visits++

	return mapping.OriginalURL, nil
}

// GetMapping returns the full mapping details for a short code
// Demonstrates returning POINTERS to structs
func (us *URLShortener) GetMapping(shortCode string) (*URLMapping, error) {
	mapping, exists := us.Mappings[shortCode]
	if !exists {
		return nil, fmt.Errorf("short code not found: %s", shortCode)
	}

	return mapping, nil
}

// ListAllMappings returns all URL mappings as a slice
// Demonstrates SLICES and iterating over MAPS
func (us *URLShortener) ListAllMappings() []*URLMapping {
	// Create a slice to hold all mappings
	// We use make with initial length 0 and capacity equal to number of mappings
	var mappings []*URLMapping

	// Iterate over map (order is random by design)
	for _, mapping := range us.Mappings {
		mappings = append(mappings, mapping)
	}

	return mappings
}

// GetStats returns statistics about a shortened URL
// Demonstrates creating new data from existing data
func (us *URLShortener) GetStats(shortCode string) (map[string]interface{}, error) {
	mapping, err := us.GetMapping(shortCode)
	if err != nil {
		return nil, err
	}

	// Create a map with statistics
	// map[string]interface{} allows values of any type
	stats := map[string]interface{}{
		"short_code":    mapping.ShortCode,
		"original_url":  mapping.OriginalURL,
		"created_at":    mapping.CreatedAt,
		"visits":        mapping.Visits,
		"url_length":    len(mapping.OriginalURL),
		"code_length":   len(mapping.ShortCode),
		"compression":   float64(len(mapping.OriginalURL)) / float64(len(mapping.ShortCode)),
	}

	return stats, nil
}

// DeleteURL removes a shortened URL mapping
// Demonstrates deleting from MAPS
func (us *URLShortener) DeleteURL(shortCode string) error {
	if _, exists := us.Mappings[shortCode]; !exists {
		return fmt.Errorf("short code not found: %s", shortCode)
	}

	// delete() is a built-in function for removing map entries
	delete(us.Mappings, shortCode)
	return nil
}

// GetAllURLs returns all mappings as a slice for JSON export
// Demonstrates returning a SLICE of POINTERS
func (us *URLShortener) GetAllURLs() []*URLMapping {
	var result []*URLMapping

	for _, mapping := range us.Mappings {
		result = append(result, mapping)
	}

	return result
}

// generateCode creates a random short code
// Demonstrates using rand and string manipulation
func (us *URLShortener) generateCode() string {
	// Create a byte slice to build the code
	// SLICES in Go are dynamic arrays
	code := make([]byte, us.CodeLength)

	// Generate random indices into the charset
	for i := range code {
		// rand.Intn returns a random number between 0 and n-1
		// This selects a random character from charset
		code[i] = us.charset[rand.Intn(len(us.charset))]
	}

	// Convert byte slice to string
	return string(code)
}

// validateURL checks if a URL is valid
// Demonstrates INPUT VALIDATION and ERROR HANDLING
func (us *URLShortener) validateURL(url string) error {
	// Check if empty
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Check if it starts with http:// or https://
	// strings.HasPrefix checks if a string starts with a prefix
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	// Check minimum length (at least protocol + domain)
	if len(url) < 10 {
		return fmt.Errorf("URL is too short: %s", url)
	}

	// Check for spaces
	if strings.Contains(url, " ") {
		return fmt.Errorf("URL cannot contain spaces")
	}

	return nil
}

// GetOverallStats returns overall statistics
// Demonstrates aggregating data
func (us *URLShortener) GetOverallStats() map[string]interface{} {
	totalVisits := 0
	longestURL := ""
	shortestURL := ""

	// Iterate over all mappings to calculate stats
	for _, mapping := range us.Mappings {
		totalVisits += mapping.Visits

		if longestURL == "" || len(mapping.OriginalURL) > len(longestURL) {
			longestURL = mapping.OriginalURL
		}

		if shortestURL == "" || len(mapping.OriginalURL) < len(shortestURL) {
			shortestURL = mapping.OriginalURL
		}
	}

	return map[string]interface{}{
		"total_urls":      len(us.Mappings),
		"total_visits":    totalVisits,
		"longest_url":     longestURL,
		"shortest_url":    shortestURL,
		"avg_url_length":  float64(totalVisits) / float64(len(us.Mappings)+1),
	}
}

// init() runs when the package is imported
// We use it to seed the random number generator
// This ensures we get different random numbers each time we run
func init() {
	rand.Seed(time.Now().UnixNano())
}
