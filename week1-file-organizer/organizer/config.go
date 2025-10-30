package organizer

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the structure of a JSON configuration file
// This demonstrates STRUCTS and JSON MARSHALING/UNMARSHALING in Go
// The `json:""` tags tell the encoder/decoder how to map JSON fields to Go fields
type Config struct {
	// Extensions is a map of file extensions to folder names
	Extensions map[string]string `json:"extensions"`
	// DefaultFolder is used for files with unknown extensions
	DefaultFolder string `json:"defaultFolder"`
}

// LoadConfig reads and parses a JSON configuration file
// This demonstrates:
// 1. Working with FILES (os.ReadFile)
// 2. JSON UNMARSHALING (converting JSON bytes to Go struct)
// 3. ERROR HANDLING with meaningful error messages
func LoadConfig(filepath string) (*Config, error) {
	// os.ReadFile reads the entire file into memory as a byte slice
	// This is fine for config files which are typically small
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Create a new Config struct to unmarshal into
	config := &Config{}

	// json.Unmarshal converts JSON bytes into our Go struct
	// It uses the json tags on our struct fields to know how to map the data
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveConfig writes the configuration to a JSON file
// This demonstrates:
// 1. JSON MARSHALING (converting Go struct to JSON bytes)
// 2. Writing FILES (os.WriteFile)
// 3. File permissions
func SaveConfig(filepath string, config *Config) error {
	// json.MarshalIndent converts our Go struct to formatted JSON
	// indent strings add pretty-printing with 2-space indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// os.WriteFile writes the JSON bytes to a file
	// 0644 is the file permission (rw-r--r--)
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// CreateDefaultConfig generates a default configuration file
// This is useful for users who haven't created a config yet
func CreateDefaultConfig(filepath string) error {
	defaultConfig := &Config{
		Extensions: map[string]string{
			".txt":   "Documents",
			".pdf":   "Documents",
			".jpg":   "Images",
			".png":   "Images",
			".mp4":   "Videos",
			".mp3":   "Music",
			".zip":   "Archives",
		},
		DefaultFolder: "Other",
	}

	return SaveConfig(filepath, defaultConfig)
}

// ApplyConfigToOrganizer updates the FileOrganizer with settings from Config
// This demonstrates how to use a config to initialize behavior
func (c *Config) ApplyToOrganizer(fo *FileOrganizer) {
	// Clear the existing extension map and replace with config values
	fo.ExtensionMap = c.Extensions
}
