package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileOrganizer handles the organization of files into directories by extension
// This demonstrates the concept of STRUCTS in Go - a custom data type that groups
// related data (fields) together
type FileOrganizer struct {
	// SourceDir is the directory to organize files from
	SourceDir string
	// OutputDir is the base directory where organized files will be placed
	OutputDir string
	// ExtensionMap maps file extensions to folder names (e.g., ".txt" -> "Documents")
	// This demonstrates MAPS in Go - key-value pairs with O(1) lookup time
	ExtensionMap map[string]string
	// DryRun if true, shows what would be moved without actually moving files
	DryRun bool
}

// NewFileOrganizer creates and returns a new FileOrganizer instance
// This demonstrates FUNCTIONS and POINTERS in Go
// The asterisk (*) means we're returning a pointer to a FileOrganizer, not a copy
// This is more efficient for larger structs
func NewFileOrganizer(sourceDir, outputDir string) *FileOrganizer {
	return &FileOrganizer{
		SourceDir: sourceDir,
		OutputDir: outputDir,
		// Initialize map with default extension mappings
		// This demonstrates MAP INITIALIZATION with values
		ExtensionMap: map[string]string{
			".txt":   "Documents",
			".pdf":   "Documents",
			".doc":   "Documents",
			".docx":  "Documents",
			".jpg":   "Images",
			".jpeg":  "Images",
			".png":   "Images",
			".gif":   "Images",
			".mp4":   "Videos",
			".mkv":   "Videos",
			".mov":   "Videos",
			".mp3":   "Music",
			".wav":   "Music",
			".flac":  "Music",
			".zip":   "Archives",
			".rar":   "Archives",
			".7z":    "Archives",
			".exe":   "Executables",
			".msi":   "Executables",
		},
		DryRun: false,
	}
}

// Organize is a METHOD on the FileOrganizer struct
// Methods are functions that have a receiver (the struct instance before the function name)
// In Go, we don't have classes, but we use structs with methods to achieve similar behavior
func (fo *FileOrganizer) Organize() error {
	// Validate that the source directory exists
	if err := fo.validateSourceDir(); err != nil {
		return err
	}

	// Read all files in the source directory
	// This demonstrates working with os.ReadDir which returns a slice of DirEntry
	files, err := os.ReadDir(fo.SourceDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Initialize a map to track statistics
	// This demonstrates how maps can be used to collect data
	stats := map[string]int{
		"processed": 0,
		"skipped":   0,
		"moved":     0,
		"errors":    0,
	}

	// Iterate over files using a for loop with range
	// SLICES and ARRAYS are fundamental Go data structures
	// This demonstrates iterating over a slice with for...range
	for _, file := range files {
		// Skip directories - we only want to organize files
		if file.IsDir() {
			stats["skipped"]++
			continue
		}

		// Process each file
		if err := fo.processFile(file.Name()); err != nil {
			fmt.Printf("Error processing %s: %v\n", file.Name(), err)
			stats["errors"]++
		} else {
			stats["processed"]++
			if !fo.DryRun {
				stats["moved"]++
			}
		}
	}

	// Print summary statistics
	fmt.Println("\n=== Organization Summary ===")
	fmt.Printf("Files processed: %d\n", stats["processed"])
	fmt.Printf("Files moved: %d\n", stats["moved"])
	fmt.Printf("Errors: %d\n", stats["errors"])

	return nil
}

// processFile handles the organization of a single file
// This demonstrates ERROR HANDLING - returning error as the last return value
// Go's convention is to return (result, error) not throwing exceptions
func (fo *FileOrganizer) processFile(filename string) error {
	// Get the file extension
	ext := strings.ToLower(filepath.Ext(filename))

	// Look up the folder name for this extension in the ExtensionMap
	// This demonstrates MAP LOOKUP with the comma-ok idiom
	// folderName is the value, exists tells us if the key was found
	folderName, exists := fo.ExtensionMap[ext]
	if !exists {
		folderName = "Other" // Default folder for unknown extensions
	}

	// Construct the full source and destination paths
	sourcePath := filepath.Join(fo.SourceDir, filename)
	destDir := filepath.Join(fo.OutputDir, folderName)
	destPath := filepath.Join(destDir, filename)

	// If file already exists at destination, skip it
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("destination file already exists")
	}

	if fo.DryRun {
		fmt.Printf("[DRY RUN] Would move: %s -> %s\n", filename, folderName)
		return nil
	}

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Move the file (rename is the way to move files in Go)
	if err := os.Rename(sourcePath, destPath); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	fmt.Printf("Moved: %s -> %s/\n", filename, folderName)
	return nil
}

// validateSourceDir checks if the source directory exists and is a directory
// This demonstrates ERROR HANDLING and working with os.Stat
func (fo *FileOrganizer) validateSourceDir() error {
	// os.Stat returns file information and an error
	// If the error is not nil, the file/directory doesn't exist or can't be accessed
	info, err := os.Stat(fo.SourceDir)
	if err != nil {
		return fmt.Errorf("source directory error: %w", err)
	}

	// Check if the path is actually a directory
	if !info.IsDir() {
		return fmt.Errorf("source path is not a directory")
	}

	return nil
}

// UpdateExtensionMap allows users to add or modify extension mappings
// This demonstrates METHODS with parameters and how to modify struct state
func (fo *FileOrganizer) UpdateExtensionMap(ext, folder string) {
	// Ensure extension has a leading dot
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	// Normalize to lowercase for consistency
	ext = strings.ToLower(ext)
	// Add to map
	fo.ExtensionMap[ext] = folder
}

// ListMappings returns all current extension mappings as a slice of strings
// This demonstrates SLICES and how to iterate over maps
func (fo *FileOrganizer) ListMappings() []string {
	// SLICES are dynamic arrays - their length can change
	// We initialize with make([]string, 0) which creates an empty slice
	var mappings []string

	// Iterate over the map using for...range
	// When iterating over a map, the order is random (by design in Go)
	for ext, folder := range fo.ExtensionMap {
		// Append adds an element to the slice and returns the updated slice
		mapping := fmt.Sprintf("%s -> %s", ext, folder)
		mappings = append(mappings, mapping)
	}

	return mappings
}
