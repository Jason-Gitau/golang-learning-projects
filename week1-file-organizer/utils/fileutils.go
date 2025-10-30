package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// DirectoryExists checks if a directory exists and is accessible
// This demonstrates ERROR HANDLING and working with os.Stat
func DirectoryExists(path string) bool {
	// os.Stat returns FileInfo and error
	// If error is nil, the path exists
	// We use the blank identifier _ to ignore the FileInfo since we only care about existence
	_, err := os.Stat(path)
	// err == nil means the path exists
	return err == nil
}

// IsDirectory returns true if the path is a directory, false if it's a file
// This demonstrates working with FileInfo.IsDir()
func IsDirectory(path string) (bool, error) {
	// os.Stat returns information about a file/directory
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat path: %w", err)
	}

	// IsDir() returns true if it's a directory, false if it's a regular file
	return info.IsDir(), nil
}

// GetFileExtension returns the extension of a file (e.g., ".txt")
// This demonstrates strings and filepath manipulation
func GetFileExtension(filename string) string {
	// filepath.Ext returns the extension including the dot
	// For "document.txt" it returns ".txt"
	// For "archive.tar.gz" it returns ".gz" (only the last extension)
	return filepath.Ext(filename)
}

// GetDirectorySize calculates the total size of all files in a directory (non-recursive)
// This demonstrates:
// 1. Iterating over directory contents with os.ReadDir
// 2. Getting file information with FileInfo
// 3. SLICES and array indexing
// 4. Error handling
func GetDirectorySize(dirPath string) (int64, error) {
	// os.ReadDir returns a slice of DirEntry objects
	// DirEntry is more efficient than os.Stat for just listing directories
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read directory: %w", err)
	}

	// Initialize a variable to accumulate total size
	var totalSize int64

	// Iterate over the slice of DirEntry objects
	for _, entry := range entries {
		// Skip directories - we only want file sizes
		if entry.IsDir() {
			continue
		}

		// Get file information to access the Size() method
		info, err := entry.Info()
		if err != nil {
			// Log the error but continue processing other files
			fmt.Printf("Warning: failed to get info for %s: %v\n", entry.Name(), err)
			continue
		}

		// Add this file's size to the total
		// Size() returns int64 which is important for large file sizes
		totalSize += info.Size()
	}

	return totalSize, nil
}

// ListFiles returns a slice of all filenames in a directory
// This demonstrates SLICES and the make() built-in function
func ListFiles(dirPath string) ([]string, error) {
	// os.ReadDir returns a slice of DirEntry
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// make([]string, 0) creates an empty slice with capacity 0
	// It will grow as we append to it
	// We could also use make([]string, 0, len(entries)) to pre-allocate capacity
	var files []string

	// Iterate and collect only files (not directories)
	for _, entry := range entries {
		if !entry.IsDir() {
			// append() adds an element to the slice
			// If capacity is exceeded, Go allocates a larger backing array
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// CountFilesByExtension counts how many files have each extension in a directory
// This demonstrates MAPS for collecting statistics
// Returns a map where keys are extensions and values are counts
func CountFilesByExtension(dirPath string) (map[string]int, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Create a map to store extension counts
	// map[string]int: keys are strings (extensions), values are integers (counts)
	extensionCounts := make(map[string]int)

	// Iterate over files
	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Get the file extension
		ext := filepath.Ext(entry.Name())

		// If extension is empty, use a special key
		if ext == "" {
			ext = "[no extension]"
		}

		// Increment the count for this extension
		// If the key doesn't exist, Go initializes the value to 0
		// Then we increment it
		extensionCounts[ext]++
	}

	return extensionCounts, nil
}

// CreateDirectoryIfNotExists creates a directory if it doesn't already exist
// This demonstrates working with os.MkdirAll
// The 0755 is the Unix permission (rwxr-xr-x)
func CreateDirectoryIfNotExists(path string) error {
	// Check if directory already exists
	if info, err := os.Stat(path); err == nil {
		// Path exists - check if it's a directory
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", path)
		}
		// Directory exists, nothing to do
		return nil
	}

	// os.MkdirAll creates the directory and any parent directories needed
	// The second parameter (0755) sets file permissions
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// SafePath joins path components and validates the result is within a base directory
// This is useful for preventing directory traversal attacks
// This demonstrates filepath.Join and path validation
func SafePath(basePath string, components ...string) (string, error) {
	// filepath.Join joins path elements with the OS separator
	// It cleans up . and .. in the path
	fullPath := filepath.Join(basePath, filepath.Join(components...))

	// Convert to absolute path for comparison
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for base: %w", err)
	}

	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify that absPath starts with absBase to prevent directory traversal
	// This prevents paths like "../../etc/passwd" from escaping the base directory
	if !isPathWithin(absPath, absBase) {
		return "", fmt.Errorf("path escapes base directory: %s", fullPath)
	}

	return fullPath, nil
}

// isPathWithin checks if a path is within a base directory
// Helper function for SafePath
func isPathWithin(path, basePath string) bool {
	// Ensure the base path ends with a separator for proper comparison
	if basePath[len(basePath)-1] != filepath.Separator {
		basePath += string(filepath.Separator)
	}

	// Check if path starts with basePath
	// Using filepath.HasPrefix would be ideal but it doesn't exist
	// So we use len and string comparison
	return len(path) >= len(basePath)-1 && path[:len(basePath)-1] == basePath[:len(basePath)-1]
}

// PrintDirectoryTree prints a simple text-based view of directory structure
// This demonstrates working with os.ReadDir recursively
func PrintDirectoryTree(dirPath string, prefix string, maxDepth int) error {
	if maxDepth <= 0 {
		return nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for i, entry := range entries {
		// Determine the prefix characters for tree display
		isLast := i == len(entries)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		fmt.Printf("%s%s%s\n", prefix, connector, entry.Name())

		// Recursively print subdirectories
		if entry.IsDir() {
			extension := "│   "
			if isLast {
				extension = "    "
			}

			subPath := filepath.Join(dirPath, entry.Name())
			PrintDirectoryTree(subPath, prefix+extension, maxDepth-1)
		}
	}

	return nil
}
