package tools

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"deep-research-agent/models"
)

// WebSearch provides web search capabilities
type WebSearch struct {
	// Can be extended to use real APIs like DuckDuckGo, SerpAPI, etc.
}

// NewWebSearch creates a new web search tool
func NewWebSearch() *WebSearch {
	return &WebSearch{}
}

// Name returns the tool name
func (w *WebSearch) Name() string {
	return "web_search"
}

// Description returns the tool description
func (w *WebSearch) Description() string {
	return "Search the web for information. Returns top results with titles, URLs, snippets, and relevance scores."
}

// Parameters returns the tool parameters
func (w *WebSearch) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "query",
			Type:        "string",
			Required:    true,
			Description: "Search query",
		},
		{
			Name:        "max_results",
			Type:        "int",
			Required:    false,
			Description: "Maximum number of results to return (default: 10)",
			Default:     10,
		},
	}
}

// Execute runs the web search
func (w *WebSearch) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	// Extract parameters
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	maxResults := 10
	if mr, ok := params["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	// Perform search (mock implementation)
	results := w.mockSearch(query, maxResults)

	// Convert to sources
	sources := make([]models.Source, 0, len(results.Results))
	for _, result := range results.Results {
		sources = append(sources, models.Source{
			ID:         fmt.Sprintf("web_%d", time.Now().UnixNano()),
			Type:       "web",
			Title:      result.Title,
			URL:        result.URL,
			Excerpt:    result.Snippet,
			AccessDate: time.Now(),
			Metadata: map[string]interface{}{
				"source":    result.Source,
				"relevance": result.Relevance,
			},
		})
	}

	return &ToolResult{
		Success: true,
		Data:    results,
		Sources: sources,
		Metadata: map[string]interface{}{
			"query":        query,
			"result_count": len(results.Results),
		},
	}, nil
}

// mockSearch generates realistic mock search results
func (w *WebSearch) mockSearch(query string, maxResults int) *SearchResults {
	// Generate context-aware results based on the query
	results := w.generateContextualResults(query, maxResults)

	return &SearchResults{
		Query:      query,
		Results:    results,
		TotalCount: len(results),
	}
}

// generateContextualResults creates realistic results based on query
func (w *WebSearch) generateContextualResults(query string, maxResults int) []SearchResult {
	queryLower := strings.ToLower(query)
	results := make([]SearchResult, 0, maxResults)

	// Define result templates for common topics
	templates := w.getResultTemplates(queryLower)

	// If we have templates, use them; otherwise generate generic results
	if len(templates) > 0 {
		count := maxResults
		if len(templates) < maxResults {
			count = len(templates)
		}

		for i := 0; i < count; i++ {
			results = append(results, templates[i])
		}
	} else {
		// Generate generic results
		for i := 0; i < maxResults; i++ {
			results = append(results, SearchResult{
				Title:     fmt.Sprintf("Result %d for '%s'", i+1, query),
				URL:       fmt.Sprintf("https://example.com/%s/%d", strings.ReplaceAll(query, " ", "-"), i+1),
				Snippet:   fmt.Sprintf("This page contains information about %s. Learn more about %s and related topics.", query, query),
				Source:    "example.com",
				Relevance: 0.9 - float64(i)*0.05,
			})
		}
	}

	return results
}

// getResultTemplates returns predefined templates for common queries
func (w *WebSearch) getResultTemplates(query string) []SearchResult {
	// Check for programming-related queries
	if strings.Contains(query, "go") || strings.Contains(query, "golang") {
		return []SearchResult{
			{
				Title:     "The Go Programming Language",
				URL:       "https://golang.org/",
				Snippet:   "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
				Source:    "golang.org",
				Relevance: 0.98,
			},
			{
				Title:     "Go by Example",
				URL:       "https://gobyexample.com/",
				Snippet:   "Go by Example is a hands-on introduction to Go using annotated example programs. Check out the first example or browse the full list below.",
				Source:    "gobyexample.com",
				Relevance: 0.95,
			},
			{
				Title:     "Effective Go - The Go Programming Language",
				URL:       "https://golang.org/doc/effective_go",
				Snippet:   "Effective Go provides tips for writing clear, idiomatic Go code. It's a good reference for writing clean Go programs.",
				Source:    "golang.org",
				Relevance: 0.92,
			},
		}
	}

	// Check for concurrency-related queries
	if strings.Contains(query, "concurrency") || strings.Contains(query, "concurrent") {
		return []SearchResult{
			{
				Title:     "Concurrency in Go - Goroutines and Channels",
				URL:       "https://go.dev/tour/concurrency",
				Snippet:   "Go provides concurrency features as part of the core language. This section covers goroutines and channels, and how they are used to implement concurrent programming patterns.",
				Source:    "go.dev",
				Relevance: 0.97,
			},
			{
				Title:     "Concurrency Patterns in Go",
				URL:       "https://blog.golang.org/pipelines",
				Snippet:   "Go's concurrency primitives make it easy to construct streaming data pipelines that make efficient use of I/O and multiple CPUs.",
				Source:    "blog.golang.org",
				Relevance: 0.94,
			},
		}
	}

	// Check for AI/ML related queries
	if strings.Contains(query, "ai") || strings.Contains(query, "artificial intelligence") || strings.Contains(query, "machine learning") {
		return []SearchResult{
			{
				Title:     "What is Artificial Intelligence (AI)? | IBM",
				URL:       "https://www.ibm.com/topics/artificial-intelligence",
				Snippet:   "Artificial intelligence leverages computers and machines to mimic the problem-solving and decision-making capabilities of the human mind.",
				Source:    "ibm.com",
				Relevance: 0.96,
			},
			{
				Title:     "Machine Learning | Coursera",
				URL:       "https://www.coursera.org/learn/machine-learning",
				Snippet:   "Machine learning is the science of getting computers to act without being explicitly programmed. Learn the fundamentals from Stanford.",
				Source:    "coursera.org",
				Relevance: 0.93,
			},
		}
	}

	// Check for research-related queries
	if strings.Contains(query, "research") {
		return []SearchResult{
			{
				Title:     "Research Methods: Types, Examples and Guide",
				URL:       "https://research.com/research/research-methods",
				Snippet:   "Research methods are the strategies, processes or techniques utilized in the collection of data or evidence for analysis in order to uncover new information.",
				Source:    "research.com",
				Relevance: 0.95,
			},
			{
				Title:     "Google Scholar",
				URL:       "https://scholar.google.com/",
				Snippet:   "Google Scholar provides a simple way to broadly search for scholarly literature across many disciplines and sources.",
				Source:    "scholar.google.com",
				Relevance: 0.92,
			},
		}
	}

	return nil // No templates found
}

// SearchResults represents web search results
type SearchResults struct {
	Query      string         `json:"query"`
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Title     string  `json:"title"`
	URL       string  `json:"url"`
	Snippet   string  `json:"snippet"`
	Source    string  `json:"source"`
	Relevance float64 `json:"relevance"`
}

// init initializes random seed for mock data generation
func init() {
	rand.Seed(time.Now().UnixNano())
}
