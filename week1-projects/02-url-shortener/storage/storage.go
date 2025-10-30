package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jason/url-shortener/shortener"
)

// Storage handles persistence of URL mappings to JSON files
// Demonstrates FILE I/O and JSON MARSHALING/UNMARSHALING
type Storage struct {
	// filePath is where the JSON data is stored
	filePath string
}

// NewStorage creates a new Storage instance
// Demonstrates FUNCTIONS that return POINTERS
func NewStorage(filePath string) *Storage {
	return &Storage{
		filePath: filePath,
	}
}

// SaveMappings writes all URL mappings to a JSON file
// Demonstrates:
// 1. JSON MARSHALING (Go struct to JSON)
// 2. FILE WRITING (os.WriteFile)
// 3. ERROR HANDLING
func (s *Storage) SaveMappings(mappings []*shortener.URLMapping) error {
	// json.MarshalIndent converts Go structs to formatted JSON
	// The struct tags (json:"...") in URLMapping tell the encoder
	// which fields to include and what names to use
	data, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		// Error wrapping with %w allows error inspection with errors.Is()
		return fmt.Errorf("failed to marshal mappings: %w", err)
	}

	// os.WriteFile writes data to a file
	// 0644 is Unix file permissions (rw-r--r--)
	// If the file doesn't exist, it's created
	// If it exists, it's overwritten
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LoadMappings reads all URL mappings from the JSON file
// Demonstrates:
// 1. FILE READING (os.ReadFile)
// 2. JSON UNMARSHALING (JSON to Go struct)
// 3. ERROR HANDLING
func (s *Storage) LoadMappings() ([]*shortener.URLMapping, error) {
	// Check if file exists
	// os.Stat returns file information or an error if the file doesn't exist
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// File doesn't exist yet - return empty slice
		// This is not an error, just no data to load
		return []*shortener.URLMapping{}, nil
	}

	// Read the entire file into memory
	// For small files (like our JSON), this is fine
	// For large files, you'd stream or use a database
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return []*shortener.URLMapping{}, nil
	}

	// Create a slice to hold the mappings
	// We use a slice of pointers for efficiency
	var mappings []*shortener.URLMapping

	// json.Unmarshal converts JSON bytes to Go structs
	// The struct tags tell the decoder which fields to populate
	if err := json.Unmarshal(data, &mappings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mappings: %w", err)
	}

	return mappings, nil
}

// AppendMapping adds a single mapping to the file
// Instead of rewriting the entire file, we append
// Demonstrates READING THEN WRITING
func (s *Storage) AppendMapping(mapping *shortener.URLMapping) error {
	// Load existing mappings
	mappings, err := s.LoadMappings()
	if err != nil {
		return fmt.Errorf("failed to load existing mappings: %w", err)
	}

	// Append the new mapping to the slice
	// append() returns a new slice with the element added
	mappings = append(mappings, mapping)

	// Save all mappings back to file
	return s.SaveMappings(mappings)
}

// RemoveMapping removes a mapping from the file
// Demonstrates FILTERING a SLICE
func (s *Storage) RemoveMapping(shortCode string) error {
	// Load existing mappings
	mappings, err := s.LoadMappings()
	if err != nil {
		return fmt.Errorf("failed to load mappings: %w", err)
	}

	// Create a new slice without the mapping to remove
	// This demonstrates SLICES and filtering logic
	var filtered []*shortener.URLMapping

	found := false
	for _, mapping := range mappings {
		if mapping.ShortCode != shortCode {
			// Keep this mapping
			filtered = append(filtered, mapping)
		} else {
			// This is the one to remove
			found = true
		}
	}

	if !found {
		return fmt.Errorf("short code not found: %s", shortCode)
	}

	// Save the filtered mappings
	return s.SaveMappings(filtered)
}

// UpdateMapping updates an existing mapping in the file
// Demonstrates MODIFYING STRUCT VALUES
func (s *Storage) UpdateMapping(mapping *shortener.URLMapping) error {
	// Load all mappings
	mappings, err := s.LoadMappings()
	if err != nil {
		return fmt.Errorf("failed to load mappings: %w", err)
	}

	// Find and update the mapping
	found := false
	for i, m := range mappings {
		if m.ShortCode == mapping.ShortCode {
			// Update the mapping at this index
			// Use indexing to modify the slice element
			mappings[i] = mapping
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("short code not found: %s", mapping.ShortCode)
	}

	// Save all mappings
	return s.SaveMappings(mappings)
}

// BackupFile creates a backup of the current mappings file
// Demonstrates FILE OPERATIONS
func (s *Storage) BackupFile() error {
	// Read the current file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		// If file doesn't exist, nothing to backup
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create a backup filename with timestamp
	backupPath := s.filePath + ".backup"

	// Write backup
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

// FileExists checks if the storage file exists
// Demonstrates CHECKING FILE EXISTENCE
func (s *Storage) FileExists() bool {
	_, err := os.Stat(s.filePath)
	// err == nil means the file exists
	return err == nil
}

// GetFileSize returns the size of the storage file in bytes
// Demonstrates GETTING FILE INFORMATION
func (s *Storage) GetFileSize() (int64, error) {
	info, err := os.Stat(s.filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return info.Size(), nil
}

// ClearAll deletes all mappings from the storage file
// Demonstrates DELETING FILES
func (s *Storage) ClearAll() error {
	// Save an empty slice to clear everything
	return s.SaveMappings([]*shortener.URLMapping{})
}
