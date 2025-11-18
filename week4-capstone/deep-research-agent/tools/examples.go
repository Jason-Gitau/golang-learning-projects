package tools

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example: How to use the tools system

func ExamplePDFProcessor() {
	ctx := context.Background()

	// Create PDF processor
	pdfTool := NewPDFProcessor()

	// Extract all text from PDF
	result, err := pdfTool.Execute(ctx, map[string]interface{}{
		"file_path": "/path/to/document.pdf",
		"action":    "extract_text",
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		content := result.Data.(*PDFContent)
		fmt.Printf("Extracted %d pages\n", content.PageCount)
		fmt.Printf("Title: %s\n", content.Metadata.Title)
	}

	// Search in PDF
	result, err = pdfTool.Execute(ctx, map[string]interface{}{
		"file_path": "/path/to/document.pdf",
		"action":    "search",
		"query":     "machine learning",
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		searchResults := result.Data.(*PDFSearchResults)
		fmt.Printf("Found %d matches\n", searchResults.MatchCount)
	}
}

func ExampleWebSearch() {
	ctx := context.Background()

	// Create web search tool
	searchTool := NewWebSearch()

	// Perform search
	result, err := searchTool.Execute(ctx, map[string]interface{}{
		"query":       "Go concurrency patterns",
		"max_results": 5,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		searchResults := result.Data.(*SearchResults)
		for i, res := range searchResults.Results {
			fmt.Printf("%d. %s\n", i+1, res.Title)
			fmt.Printf("   %s\n", res.URL)
			fmt.Printf("   %s\n\n", res.Snippet)
		}
	}
}

func ExampleWikipedia() {
	ctx := context.Background()

	// Create Wikipedia tool
	wikiTool := NewWikipedia()

	// Get article summary
	result, err := wikiTool.Execute(ctx, map[string]interface{}{
		"action": "summary",
		"query":  "Artificial Intelligence",
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		article := result.Data.(*WikiArticle)
		fmt.Printf("Title: %s\n", article.Title)
		fmt.Printf("Summary: %s\n", article.Summary)
		fmt.Printf("URL: %s\n", article.URL)
	}
}

func ExampleURLFetcher() {
	ctx := context.Background()

	// Create URL fetcher
	fetchTool := NewURLFetcher()

	// Fetch page content
	result, err := fetchTool.Execute(ctx, map[string]interface{}{
		"url":    "https://golang.org/",
		"action": "extract_text",
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		content := result.Data.(*WebContent)
		fmt.Printf("Title: %s\n", content.Title)
		fmt.Printf("Content length: %d bytes\n", len(content.CleanText))
	}
}

func ExampleSummarizer() {
	ctx := context.Background()

	// Create summarizer
	summTool := NewSummarizer()

	longText := `Your long text here...`

	// Generate summary
	result, err := summTool.Execute(ctx, map[string]interface{}{
		"text":              longText,
		"action":            "summarize",
		"max_sentences":     5,
		"compression_ratio": 0.3,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		summary := result.Data.(*Summary)
		fmt.Printf("Original length: %d\n", summary.OriginalLength)
		fmt.Printf("Summary length: %d\n", summary.SummaryLength)
		fmt.Printf("Compression: %.2f%%\n", summary.CompressionRatio*100)
		fmt.Printf("Summary:\n%s\n", summary.Summary)
	}
}

func ExampleCitationManager() {
	ctx := context.Background()

	// Create citation manager
	citeTool := NewCitationManager()

	// Add a source
	sourceData := map[string]interface{}{
		"type":      "web",
		"title":     "The Go Programming Language",
		"url":       "https://golang.org/",
		"author":    "Google",
		"publisher": "Google Inc.",
	}

	result, err := citeTool.Execute(ctx, map[string]interface{}{
		"action": "add_source",
		"source": sourceData,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Generate bibliography
	result, err = citeTool.Execute(ctx, map[string]interface{}{
		"action": "get_bibliography",
		"style":  "APA",
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		bibData := result.Data.(*CitationData)
		fmt.Printf("Bibliography (%d sources):\n", bibData.TotalSources)
		fmt.Println(bibData.Bibliography)
	}
}

func ExampleFactChecker() {
	ctx := context.Background()

	// Create fact checker
	factTool := NewFactChecker()

	sources := []interface{}{
		map[string]interface{}{
			"title":   "Go Concurrency",
			"content": "Go supports concurrency through goroutines and channels.",
		},
		map[string]interface{}{
			"title":   "Go Documentation",
			"content": "Goroutines are lightweight threads managed by Go runtime.",
		},
	}

	// Check a fact
	result, err := factTool.Execute(ctx, map[string]interface{}{
		"action":  "check_fact",
		"claim":   "Go supports concurrency through goroutines",
		"sources": sources,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		checkResult := result.Data.(*FactCheckResult)
		fmt.Printf("Claim: %s\n", checkResult.Claim)
		fmt.Printf("Verified: %v\n", checkResult.Verified)
		fmt.Printf("Confidence: %.2f\n", checkResult.Confidence)
		fmt.Printf("Supporting sources: %d\n", checkResult.SupportCount)
	}
}

func ExampleToolRegistry() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create registry and register all tools
	registry := NewToolRegistry()
	if err := RegisterAllTools(registry); err != nil {
		log.Fatal(err)
	}

	// List all available tools
	tools := registry.List()
	fmt.Printf("Available tools: %d\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("- %s: %s\n", tool.Name(), tool.Description())
	}

	// Get tool info
	toolInfo := registry.GetToolInfo()
	for name, info := range toolInfo {
		fmt.Printf("\nTool: %s\n", name)
		fmt.Printf("Description: %s\n", info.Description)
		fmt.Printf("Parameters:\n")
		for _, param := range info.Parameters {
			required := ""
			if param.Required {
				required = " (required)"
			}
			fmt.Printf("  - %s (%s)%s: %s\n",
				param.Name,
				param.Type,
				required,
				param.Description,
			)
		}
	}

	// Execute tool through registry
	result, err := registry.Execute(ctx, "web_search", map[string]interface{}{
		"query":       "Go programming",
		"max_results": 5,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		fmt.Printf("Search completed successfully\n")
		fmt.Printf("Found %d sources\n", len(result.Sources))
	}
}
