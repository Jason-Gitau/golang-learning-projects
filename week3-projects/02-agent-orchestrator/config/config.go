package config

import "time"

// Config holds the application configuration
type Config struct {
	// Agent pool configuration
	NumAgents       int           `json:"num_agents"`
	RequestQueueSize int          `json:"request_queue_size"`
	DefaultTimeout  time.Duration `json:"default_timeout"`

	// API configuration
	APIPort         string `json:"api_port"`
	EnableCORS      bool   `json:"enable_cors"`

	// Tool configuration
	WeatherAPIKey   string `json:"weather_api_key,omitempty"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		NumAgents:        5,
		RequestQueueSize: 100,
		DefaultTimeout:   30 * time.Second,
		APIPort:          "8080",
		EnableCORS:       true,
		WeatherAPIKey:    "", // Optional: set via environment variable
	}
}

// WithAgents sets the number of agents
func (c *Config) WithAgents(n int) *Config {
	c.NumAgents = n
	return c
}

// WithQueueSize sets the request queue size
func (c *Config) WithQueueSize(size int) *Config {
	c.RequestQueueSize = size
	return c
}

// WithPort sets the API port
func (c *Config) WithPort(port string) *Config {
	c.APIPort = port
	return c
}

// WithWeatherAPIKey sets the weather API key
func (c *Config) WithWeatherAPIKey(key string) *Config {
	c.WeatherAPIKey = key
	return c
}
