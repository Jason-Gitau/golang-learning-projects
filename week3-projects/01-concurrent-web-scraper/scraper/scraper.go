package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/user/web-scraper/models"
)

// Scraper handles web scraping operations
type Scraper struct {
	client      *http.Client
	parser      *Parser
	userAgent   string
	maxRetries  int
}

// NewScraper creates a new scraper instance
func NewScraper(timeout time.Duration, userAgent string, maxRetries int, followRedirects bool) *Scraper {
	client := &http.Client{
		Timeout: timeout,
	}

	// Configure redirect policy
	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &Scraper{
		client:     client,
		parser:     NewParser(),
		userAgent:  userAgent,
		maxRetries: maxRetries,
	}
}

// ScrapeURL scrapes a single URL and returns the result
func (s *Scraper) ScrapeURL(ctx context.Context, url string) *models.ScrapeResult {
	startTime := time.Now()

	result := &models.ScrapeResult{
		URL:       url,
		ScrapedAt: startTime,
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result.Error = fmt.Errorf("failed to create request: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	// Set headers
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Perform request
	resp, err := s.client.Do(req)
	if err != nil {
		result.Error = fmt.Errorf("request failed: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// Check status code
	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Errorf("non-OK status code: %d", resp.StatusCode)
		result.Duration = time.Since(startTime)
		return result
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Errorf("failed to read response: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	// Parse HTML
	title, description, links, err := s.parser.ParseHTML(string(body), url)
	if err != nil {
		result.Error = fmt.Errorf("failed to parse HTML: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Title = title
	result.Description = description
	result.Links = links
	result.Duration = time.Since(startTime)

	return result
}

// Close closes the scraper and releases resources
func (s *Scraper) Close() error {
	// Close idle connections
	s.client.CloseIdleConnections()
	return nil
}
