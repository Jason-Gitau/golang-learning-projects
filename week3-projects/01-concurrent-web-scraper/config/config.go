package config

import (
	"encoding/json"
	"os"
)

// Config holds all configuration for the scraper
type Config struct {
	// Worker pool configuration
	WorkerCount int `json:"worker_count"`

	// Rate limiting configuration (requests per second)
	RateLimit float64 `json:"rate_limit"`

	// Retry configuration
	MaxRetries int `json:"max_retries"`

	// Timeout configuration (seconds)
	RequestTimeout int `json:"request_timeout"`

	// Database configuration
	DatabasePath string `json:"database_path"`

	// URLs configuration
	URLsFile string `json:"urls_file"`

	// User agent for requests
	UserAgent string `json:"user_agent"`

	// Follow redirects
	FollowRedirects bool `json:"follow_redirects"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		WorkerCount:     5,
		RateLimit:       2.0, // 2 requests per second
		MaxRetries:      3,
		RequestTimeout:  30,
		DatabasePath:    "./data/scraper.db",
		URLsFile:        "./urls.json",
		UserAgent:       "GoWebScraper/1.0",
		FollowRedirects: true,
	}
}

// LoadFromFile loads configuration from a JSON file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveToFile saves configuration to a JSON file
func (c *Config) SaveToFile(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.WorkerCount < 1 {
		c.WorkerCount = 1
	}
	if c.WorkerCount > 100 {
		c.WorkerCount = 100
	}

	if c.RateLimit < 0.1 {
		c.RateLimit = 0.1
	}
	if c.RateLimit > 100 {
		c.RateLimit = 100
	}

	if c.MaxRetries < 0 {
		c.MaxRetries = 0
	}
	if c.MaxRetries > 10 {
		c.MaxRetries = 10
	}

	if c.RequestTimeout < 1 {
		c.RequestTimeout = 10
	}
	if c.RequestTimeout > 300 {
		c.RequestTimeout = 300
	}

	if c.DatabasePath == "" {
		c.DatabasePath = "./data/scraper.db"
	}

	if c.URLsFile == "" {
		c.URLsFile = "./urls.json"
	}

	if c.UserAgent == "" {
		c.UserAgent = "GoWebScraper/1.0"
	}

	return nil
}

// URLList represents a list of URLs to scrape
type URLList struct {
	URLs []string `json:"urls"`
}

// LoadURLs loads URLs from a JSON file
func LoadURLs(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var urlList URLList
	if err := json.Unmarshal(data, &urlList); err != nil {
		return nil, err
	}

	return urlList.URLs, nil
}
