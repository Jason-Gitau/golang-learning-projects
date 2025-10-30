package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jason/file-organizer/organizer"
	"github.com/jason/file-organizer/utils"
)

// This is the entry point for our File Organizer CLI application
// It demonstrates:
// 1. Command-line flag parsing (flag package)
// 2. Error handling and logging
// 3. Calling functions from other packages (organizer and utils)
// 4. User interaction and validation

func main() {
	// Define command-line FLAGS
	// Flags are boolean, string, or other types that can be passed to the program
	// Example: go run main.go -source /path/to/files -dry-run

	sourceDir := flag.String("source", "", "Source directory containing files to organize (required)")
	outputDir := flag.String("output", "", "Output directory where organized files will be placed (default: same as source)")
	configFile := flag.String("config", "", "Path to JSON config file with extension mappings")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without actually moving files")
	help := flag.Bool("help", false, "Show help message")
	createConfig := flag.Bool("create-config", false, "Create a default config file and exit")
	listMappings := flag.Bool("list", false, "List all file extension mappings")

	// Parse the command-line arguments
	// This reads os.Args[1:] and sets the flag variables
	flag.Parse()

	// Show help if requested
	if *help {
		printUsage()
		os.Exit(0)
	}

	// Handle create-config flag
	if *createConfig {
		handleCreateConfig(*configFile)
		os.Exit(0)
	}

	// Validate that source directory is provided (unless only listing)
	if *sourceDir == "" && !*listMappings {
		fmt.Println("Error: source directory is required")
		printUsage()
		os.Exit(1)
	}

	// If output directory is not specified, use the source directory
	if *outputDir == "" {
		*outputDir = *sourceDir
	}

	// Create the FileOrganizer instance
	// This is an OBJECT in Go (technically a pointer to a struct)
	fo := organizer.NewFileOrganizer(*sourceDir, *outputDir)
	fo.DryRun = *dryRun

	// Load custom configuration if provided
	if *configFile != "" {
		config, err := organizer.LoadConfig(*configFile)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		// Apply the config to the organizer
		config.ApplyToOrganizer(fo)
		fmt.Printf("Loaded config from: %s\n", *configFile)
	}

	// Show mappings if requested
	if *listMappings {
		fmt.Println("\n=== Current Extension Mappings ===")
		for _, mapping := range fo.ListMappings() {
			fmt.Println(mapping)
		}
		fmt.Println()
		os.Exit(0)
	}

	// Validate source directory exists
	if !utils.DirectoryExists(*sourceDir) {
		fmt.Printf("Error: source directory does not exist: %s\n", *sourceDir)
		os.Exit(1)
	}

	// Check if source is actually a directory
	isDir, err := utils.IsDirectory(*sourceDir)
	if err != nil {
		fmt.Printf("Error: failed to check if source is a directory: %v\n", err)
		os.Exit(1)
	}
	if !isDir {
		fmt.Printf("Error: source is not a directory: %s\n", *sourceDir)
		os.Exit(1)
	}

	// Get directory statistics before organizing
	fmt.Printf("Source directory: %s\n", *sourceDir)
	fmt.Printf("Output directory: %s\n", *outputDir)

	files, err := utils.ListFiles(*sourceDir)
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	fmt.Printf("Files found: %d\n", len(files))

	extCounts, err := utils.CountFilesByExtension(*sourceDir)
	if err != nil {
		log.Fatalf("Failed to count extensions: %v", err)
	}

	if len(extCounts) > 0 {
		fmt.Println("\nFile types found:")
		for ext, count := range extCounts {
			fmt.Printf("  %s: %d file(s)\n", ext, count)
		}
	}

	if *dryRun {
		fmt.Println("\n[DRY RUN MODE] - No files will actually be moved")
	}

	fmt.Println("\n=== Starting Organization ===")

	// Run the organizer
	if err := fo.Organize(); err != nil {
		log.Fatalf("Organization failed: %v", err)
	}

	fmt.Println("\n=== Organization Complete ===")

	// Show what was organized
	if !*dryRun {
		fmt.Printf("\nOrganized files in: %s\n", *outputDir)
		showDirectoryStructure(*outputDir)
	}
}

// printUsage displays help information for the program
// This demonstrates string formatting with fmt.Println
func printUsage() {
	fmt.Println(`
File Organizer - Command-line tool to organize files by extension

USAGE:
  file-organizer -source <directory> [options]

REQUIRED FLAGS:
  -source <directory>   : Directory containing files to organize

OPTIONAL FLAGS:
  -output <directory>   : Directory where organized files go (default: source directory)
  -config <filepath>    : JSON config file with extension mappings
  -dry-run             : Show what would happen without moving files
  -list                : Show all extension mappings
  -create-config       : Create a default config file
  -help                : Show this help message

EXAMPLES:
  # Organize files in current directory
  file-organizer -source ./Downloads

  # Dry run to see what would happen
  file-organizer -source ./Downloads -dry-run

  # Use custom config
  file-organizer -source ./Downloads -config config.json

  # Create default config file
  file-organizer -create-config -config config.json

  # Organize into different output directory
  file-organizer -source ./Downloads -output ./Organized
`)
}

// handleCreateConfig creates a default configuration file
// This demonstrates error handling and user feedback
func handleCreateConfig(configPath string) {
	if configPath == "" {
		configPath = "config.json"
	}

	// Check if file already exists
	if utils.DirectoryExists(configPath) {
		fmt.Printf("Config file already exists: %s\n", configPath)
		return
	}

	if err := organizer.CreateDefaultConfig(configPath); err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	fmt.Printf("Created default config at: %s\n", configPath)
	fmt.Println("Edit this file to customize your extension mappings")
}

// showDirectoryStructure displays a tree view of the organized directory
// This demonstrates recursion and using utility functions
func showDirectoryStructure(dirPath string) {
	fmt.Println("\nDirectory structure:")
	if err := utils.PrintDirectoryTree(dirPath, "", 3); err != nil {
		fmt.Printf("Failed to show directory structure: %v\n", err)
	}
}
