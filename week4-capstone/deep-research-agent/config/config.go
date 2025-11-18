package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Config represents the complete application configuration
type Config struct {
	Agent      AgentConfig      `json:"agent"`
	Tools      ToolsConfig      `json:"tools"`
	Storage    StorageConfig    `json:"storage"`
	API        APIConfig        `json:"api"`
	Auth       AuthConfig       `json:"auth"`
	RateLimit  RateLimitConfig  `json:"rate_limit"`
}

// AgentConfig contains configuration for the research agent
type AgentConfig struct {
	MaxSteps        int           `json:"max_steps"`
	Timeout         time.Duration `json:"timeout"`
	StepTimeout     time.Duration `json:"step_timeout"`
	MaxSources      int           `json:"max_sources"`
	ConcurrentTools int           `json:"concurrent_tools"`
	RetryAttempts   int           `json:"retry_attempts"`
	RetryDelay      time.Duration `json:"retry_delay"`
	EnableCaching   bool          `json:"enable_caching"`
	CacheTTL        time.Duration `json:"cache_ttl"`
	MinConfidence   float64       `json:"min_confidence"`
	WorkerPoolSize  int           `json:"worker_pool_size"`
	MaxConcurrent   int           `json:"max_concurrent"`
}

// ToolsConfig contains configuration for research tools
type ToolsConfig struct {
	EnableWebSearch bool   `json:"enable_web_search"`
	EnableWikipedia bool   `json:"enable_wikipedia"`
	EnablePDF       bool   `json:"enable_pdf"`
	EnableDOCX      bool   `json:"enable_docx"`
	WebSearchAPI    string `json:"web_search_api"`
	WebSearchKey    string `json:"web_search_key"`
	MaxResultsPerTool int  `json:"max_results_per_tool"`
	ToolTimeout     time.Duration `json:"tool_timeout"`
}

// StorageConfig contains configuration for data persistence
type StorageConfig struct {
	DatabasePath    string `json:"database_path"`
	EnableAutoSave  bool   `json:"enable_auto_save"`
	MaxSessionAge   time.Duration `json:"max_session_age"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
	DocumentsPath   string `json:"documents_path"`
}

// APIConfig contains configuration for the API server
type APIConfig struct {
	Enabled      bool   `json:"enabled"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	EnableCORS   bool   `json:"enable_cors"`
	EnableSwagger bool  `json:"enable_swagger"`
	TrustedProxies []string `json:"trusted_proxies"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret      string        `json:"-"` // Load from env, never save to file
	TokenExpiry    time.Duration `json:"token_expiry"`
	RefreshExpiry  time.Duration `json:"refresh_expiry"`
	EnableAuth     bool          `json:"enable_auth"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled            bool          `json:"enabled"`
	MaxRequestsPerHour int           `json:"max_requests_per_hour"`
	MaxConcurrentJobs  int           `json:"max_concurrent_jobs"`
	WindowDuration     time.Duration `json:"window_duration"`
	CleanupInterval    time.Duration `json:"cleanup_interval"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".deep-research-agent")

	return &Config{
		Agent: AgentConfig{
			MaxSteps:        10,
			Timeout:         5 * time.Minute,
			StepTimeout:     30 * time.Second,
			MaxSources:      10,
			ConcurrentTools: 3,
			RetryAttempts:   2,
			RetryDelay:      2 * time.Second,
			EnableCaching:   true,
			CacheTTL:        30 * time.Minute,
			MinConfidence:   0.3,
			WorkerPoolSize:  5,
			MaxConcurrent:   3,
		},
		Tools: ToolsConfig{
			EnableWebSearch:   true,
			EnableWikipedia:   true,
			EnablePDF:         true,
			EnableDOCX:        true,
			WebSearchAPI:      "duckduckgo", // Default to DuckDuckGo
			MaxResultsPerTool: 5,
			ToolTimeout:       30 * time.Second,
		},
		Storage: StorageConfig{
			DatabasePath:    filepath.Join(dataDir, "data", "research.db"),
			EnableAutoSave:  true,
			MaxSessionAge:   30 * 24 * time.Hour, // 30 days
			CleanupInterval: 24 * time.Hour,       // Daily cleanup
			DocumentsPath:   filepath.Join(dataDir, "documents"),
		},
		API: APIConfig{
			Enabled:        false, // API disabled by default
			Host:           "0.0.0.0",
			Port:           8080,
			EnableCORS:     true,
			EnableSwagger:  true,
			TrustedProxies: []string{},
		},
		Auth: AuthConfig{
			JWTSecret:     os.Getenv("JWT_SECRET"), // Load from environment
			TokenExpiry:   24 * time.Hour,          // 24 hours
			RefreshExpiry: 7 * 24 * time.Hour,      // 7 days
			EnableAuth:    true,
		},
		RateLimit: RateLimitConfig{
			Enabled:            true,
			MaxRequestsPerHour: 100,
			MaxConcurrentJobs:  10,
			WindowDuration:     1 * time.Hour,
			CleanupInterval:    10 * time.Minute,
		},
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves configuration to a JSON file
func (c *Config) SaveConfig(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Agent validation
	if c.Agent.MaxSteps <= 0 {
		c.Agent.MaxSteps = 10
	}
	if c.Agent.ConcurrentTools <= 0 {
		c.Agent.ConcurrentTools = 3
	}
	if c.Agent.MaxSources <= 0 {
		c.Agent.MaxSources = 10
	}
	if c.Agent.WorkerPoolSize <= 0 {
		c.Agent.WorkerPoolSize = 5
	}
	if c.Agent.Timeout <= 0 {
		c.Agent.Timeout = 5 * time.Minute
	}

	// Storage validation
	if c.Storage.DatabasePath == "" {
		homeDir, _ := os.UserHomeDir()
		c.Storage.DatabasePath = filepath.Join(homeDir, ".deep-research-agent", "data", "research.db")
	}
	if c.Storage.DocumentsPath == "" {
		homeDir, _ := os.UserHomeDir()
		c.Storage.DocumentsPath = filepath.Join(homeDir, ".deep-research-agent", "documents")
	}

	// Ensure directories exist
	dbDir := filepath.Dir(c.Storage.DatabasePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(c.Storage.DocumentsPath, 0755); err != nil {
		return err
	}

	return nil
}

// GetConfigPath returns the default config file path
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".deep-research-agent", "config.json")
}

// LoadOrCreateConfig loads config from file or creates default
func LoadOrCreateConfig() (*Config, error) {
	configPath := GetConfigPath()

	// Try to load existing config
	config, err := LoadConfig(configPath)
	if err != nil {
		// Create default config
		config = DefaultConfig()
		if err := config.Validate(); err != nil {
			return nil, err
		}
		// Save default config
		if err := config.SaveConfig(configPath); err != nil {
			// Log warning but continue with default config
			return config, nil
		}
	}

	// Validate loaded config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Global config instance for backward compatibility
var globalConfig *Config

// InitConfig initializes configuration (for backward compatibility with viper-based code)
func InitConfig(cfgFile string) error {
	var err error
	if cfgFile != "" {
		globalConfig, err = LoadConfig(cfgFile)
	} else {
		globalConfig, err = LoadOrCreateConfig()
	}
	return err
}

// Helper functions for backward compatibility

// GetDatabasePath returns the database path
func GetDatabasePath() string {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	return globalConfig.Storage.DatabasePath
}

// GetDocumentsDir returns the documents directory
func GetDocumentsDir() string {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	return globalConfig.Storage.DocumentsPath
}

// GetDefaultDepth returns the default research depth
func GetDefaultDepth() string {
	return "medium" // Default depth
}

// GetMaxSources returns the maximum number of sources
func GetMaxSources() int {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	return globalConfig.Agent.MaxSources
}

// GetConcurrentTools returns the number of concurrent tool executions
func GetConcurrentTools() int {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	return globalConfig.Agent.ConcurrentTools
}

// GetCitationStyle returns the default citation style
func GetCitationStyle() string {
	return "APA" // Default citation style
}

// GetDefaultFormat returns the default output format
func GetDefaultFormat() string {
	return "markdown" // Default output format
}

// GetConfig returns the global config instance
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	return globalConfig
}
