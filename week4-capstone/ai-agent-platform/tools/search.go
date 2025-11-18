package tools

import (
	"context"
	"fmt"
	"strings"
)

// SearchTool performs web searches (mock results)
type SearchTool struct{}

// NewSearchTool creates a new search tool
func NewSearchTool() *SearchTool {
	return &SearchTool{}
}

// Name returns the tool name
func (t *SearchTool) Name() string {
	return "search"
}

// Description returns the tool description
func (t *SearchTool) Description() string {
	return "Searches the web for information and returns relevant results"
}

// Parameters returns the tool parameters schema
func (t *SearchTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "The search query",
			},
			"max_results": map[string]interface{}{
				"type":        "number",
				"description": "Maximum number of results to return",
				"default":     5,
			},
		},
		"required": []string{"query"},
	}
}

// Execute performs a web search
func (t *SearchTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query must be a non-empty string")
	}

	maxResults := 5
	if mr, ok := params["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	// Generate mock search results
	results := t.generateMockResults(query, maxResults)

	return map[string]interface{}{
		"query":         query,
		"total_results": len(results),
		"results":       results,
	}, nil
}

// generateMockResults creates realistic mock search results
func (t *SearchTool) generateMockResults(query string, maxResults int) []map[string]interface{} {
	// Define mock result templates based on query keywords
	templates := t.getResultTemplates(query)

	results := make([]map[string]interface{}, 0, maxResults)
	for i := 0; i < maxResults && i < len(templates); i++ {
		results = append(results, templates[i])
	}

	// If we need more results, generate generic ones
	for len(results) < maxResults {
		idx := len(results) + 1
		results = append(results, map[string]interface{}{
			"title":       fmt.Sprintf("Result %d: %s", idx, query),
			"url":         fmt.Sprintf("https://example.com/search/%d", idx),
			"snippet":     fmt.Sprintf("This is a search result for '%s'. It contains relevant information about your query.", query),
			"relevance":   0.5,
		})
	}

	return results
}

// getResultTemplates returns templates based on query content
func (t *SearchTool) getResultTemplates(query string) []map[string]interface{} {
	lowerQuery := strings.ToLower(query)

	// Go-related queries
	if strings.Contains(lowerQuery, "go") || strings.Contains(lowerQuery, "golang") {
		return []map[string]interface{}{
			{
				"title":     "Go Programming Language Official Website",
				"url":       "https://go.dev/",
				"snippet":   "Go is an open source programming language that makes it simple to build secure, scalable systems.",
				"relevance": 0.98,
			},
			{
				"title":     "A Tour of Go",
				"url":       "https://go.dev/tour/",
				"snippet":   "Learn the basics of Go with interactive tutorials and examples.",
				"relevance": 0.95,
			},
			{
				"title":     "Effective Go - The Go Programming Language",
				"url":       "https://go.dev/doc/effective_go",
				"snippet":   "A comprehensive guide to writing clear, idiomatic Go code.",
				"relevance": 0.92,
			},
		}
	}

	// AI-related queries
	if strings.Contains(lowerQuery, "ai") || strings.Contains(lowerQuery, "artificial intelligence") {
		return []map[string]interface{}{
			{
				"title":     "What is Artificial Intelligence (AI)?",
				"url":       "https://example.com/ai-overview",
				"snippet":   "Artificial Intelligence refers to computer systems that can perform tasks that typically require human intelligence.",
				"relevance": 0.97,
			},
			{
				"title":     "AI and Machine Learning Fundamentals",
				"url":       "https://example.com/ai-ml-basics",
				"snippet":   "Learn about neural networks, deep learning, and modern AI architectures.",
				"relevance": 0.94,
			},
		}
	}

	// API-related queries
	if strings.Contains(lowerQuery, "api") || strings.Contains(lowerQuery, "rest") {
		return []map[string]interface{}{
			{
				"title":     "What is a REST API?",
				"url":       "https://example.com/rest-api-guide",
				"snippet":   "REST APIs are a way for applications to communicate over HTTP using standard methods like GET, POST, PUT, and DELETE.",
				"relevance": 0.96,
			},
			{
				"title":     "Building RESTful APIs with Go",
				"url":       "https://example.com/go-rest-api",
				"snippet":   "A comprehensive guide to building production-ready REST APIs using Go and popular frameworks.",
				"relevance": 0.93,
			},
		}
	}

	// Default results
	return []map[string]interface{}{
		{
			"title":     fmt.Sprintf("Everything about %s", query),
			"url":       "https://example.com/article1",
			"snippet":   fmt.Sprintf("A comprehensive guide covering all aspects of %s with examples and best practices.", query),
			"relevance": 0.90,
		},
		{
			"title":     fmt.Sprintf("%s: Getting Started", query),
			"url":       "https://example.com/getting-started",
			"snippet":   fmt.Sprintf("Learn the fundamentals of %s in this beginner-friendly tutorial.", query),
			"relevance": 0.85,
		},
		{
			"title":     fmt.Sprintf("Advanced %s Techniques", query),
			"url":       "https://example.com/advanced",
			"snippet":   fmt.Sprintf("Master advanced concepts and patterns for working with %s.", query),
			"relevance": 0.80,
		},
	}
}
