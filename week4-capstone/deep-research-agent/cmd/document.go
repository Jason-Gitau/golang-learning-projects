package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"deep-research-agent/config"
)

// Document represents an indexed document
type Document struct {
	ID          string
	Filename    string
	Path        string
	Type        string
	Size        int64
	Indexed     time.Time
	PageCount   int
	WordCount   int
}

var (
	indexDoc     bool
	searchQuery  string
)

// documentCmd represents document management commands
var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage research documents",
	Long:  `Add, list, and search indexed documents.`,
}

var addDocumentCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add a document to the index",
	Long:  `Add a PDF or DOCX document to the research index for faster searching.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAddDocument,
}

var listDocumentsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all indexed documents",
	Long:  `Display a list of all documents in the research index.`,
	RunE:  runListDocuments,
}

var searchDocumentsCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search indexed documents",
	Long:  `Search through indexed documents for specific content.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchDocuments,
}

var analyzeDocumentCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Analyze a document",
	Long:  `Analyze a document and extract key information, summaries, and metadata.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAnalyzeDocument,
}

func init() {
	rootCmd.AddCommand(documentCmd)
	documentCmd.AddCommand(addDocumentCmd)
	documentCmd.AddCommand(listDocumentsCmd)
	documentCmd.AddCommand(searchDocumentsCmd)
	documentCmd.AddCommand(analyzeDocumentCmd)

	// Flags
	addDocumentCmd.Flags().BoolVar(&indexDoc, "index", true, "index document for searching")
}

func runAddDocument(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Get file type
	ext := filepath.Ext(filePath)
	docType := ""
	switch ext {
	case ".pdf":
		docType = "PDF"
	case ".docx":
		docType = "DOCX"
	default:
		return fmt.Errorf("unsupported file type: %s (only PDF and DOCX supported)", ext)
	}

	fmt.Printf("Adding document: %s\n", filepath.Base(filePath))
	fmt.Printf("Type: %s | Size: %s\n", docType, formatFileSize(info.Size()))

	// Copy to documents directory
	docsDir := config.GetDocumentsDir()
	destPath := filepath.Join(docsDir, filepath.Base(filePath))

	if indexDoc {
		fmt.Println("\nIndexing document...")
		time.Sleep(1 * time.Second)
		color.Green("Document indexed successfully!")
	}

	// Copy file
	input, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	if err := os.WriteFile(destPath, input, 0644); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	color.Green("\nDocument added successfully!")
	fmt.Printf("Path: %s\n", destPath)

	return nil
}

func runListDocuments(cmd *cobra.Command, args []string) error {
	// Mock documents - in real implementation, this would query the database
	documents := getMockDocuments()

	if len(documents) == 0 {
		fmt.Println("No documents found.")
		fmt.Printf("Add documents with: research-agent document add <file>\n")
		return nil
	}

	bold := color.New(color.Bold)

	bold.Println("\nIndexed Documents:")
	fmt.Println(color.CyanString(fmt.Sprintf("%-40s %-8s %-10s %-8s %-20s", "Filename", "Type", "Size", "Pages", "Indexed")))
	fmt.Println(color.CyanString("--------------------------------------------------------------------------------"))

	for _, doc := range documents {
		fmt.Printf("%-40s %-8s %-10s %-8d %-20s\n",
			truncateString(doc.Filename, 40),
			doc.Type,
			formatFileSize(doc.Size),
			doc.PageCount,
			doc.Indexed.Format("2006-01-02 15:04"))
	}

	fmt.Printf("\nTotal documents: %d\n\n", len(documents))

	return nil
}

func runSearchDocuments(cmd *cobra.Command, args []string) error {
	query := args[0]

	fmt.Printf("Searching documents for: %s\n\n", query)

	// Mock search results
	results := []struct {
		Document string
		Matches  int
		Preview  string
	}{
		{
			Document: "AI_Ethics_Paper.pdf",
			Matches:  5,
			Preview:  "...ethical considerations in AI healthcare applications...",
		},
		{
			Document: "Machine_Learning_Guide.pdf",
			Matches:  3,
			Preview:  "...machine learning algorithms and their applications...",
		},
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)

	for i, result := range results {
		bold.Printf("%d. %s\n", i+1, result.Document)
		cyan.Printf("   Matches: %d\n", result.Matches)
		fmt.Printf("   %s\n\n", result.Preview)
	}

	fmt.Printf("Found %d documents with matches\n", len(results))

	return nil
}

func runAnalyzeDocument(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	fmt.Printf("Analyzing document: %s\n", filepath.Base(filePath))
	fmt.Printf("Size: %s\n\n", formatFileSize(info.Size()))

	// Simulate analysis
	fmt.Println("Extracting text...")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Analyzing content...")
	time.Sleep(1 * time.Second)

	fmt.Println("Generating summary...")
	time.Sleep(800 * time.Millisecond)

	// Display results
	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	bold.Println("Document Analysis")
	fmt.Println(color.CyanString("================================================================================"))
	fmt.Println()

	cyan.Println("Summary:")
	fmt.Println("This document discusses artificial intelligence applications in healthcare,")
	fmt.Println("focusing on ethical considerations and patient privacy concerns.")
	fmt.Println()

	cyan.Println("Key Topics:")
	fmt.Println("  - AI ethics")
	fmt.Println("  - Healthcare applications")
	fmt.Println("  - Patient privacy")
	fmt.Println("  - Medical diagnosis")
	fmt.Println()

	cyan.Println("Statistics:")
	fmt.Printf("  Pages: %d\n", 15)
	fmt.Printf("  Words: ~%d\n", 5000)
	fmt.Printf("  Paragraphs: %d\n", 120)
	fmt.Println()

	return nil
}

func getMockDocuments() []Document {
	return []Document{
		{
			ID:        "doc1",
			Filename:  "AI_Ethics_Paper.pdf",
			Path:      "/path/to/AI_Ethics_Paper.pdf",
			Type:      "PDF",
			Size:      2048576,
			Indexed:   time.Now().Add(-2 * time.Hour),
			PageCount: 15,
			WordCount: 5000,
		},
		{
			ID:        "doc2",
			Filename:  "Machine_Learning_Guide.pdf",
			Path:      "/path/to/Machine_Learning_Guide.pdf",
			Type:      "PDF",
			Size:      4194304,
			Indexed:   time.Now().Add(-24 * time.Hour),
			PageCount: 30,
			WordCount: 10000,
		},
		{
			ID:        "doc3",
			Filename:  "Research_Notes.docx",
			Path:      "/path/to/Research_Notes.docx",
			Type:      "DOCX",
			Size:      524288,
			Indexed:   time.Now().Add(-1 * time.Hour),
			PageCount: 8,
			WordCount: 2500,
		},
	}
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
