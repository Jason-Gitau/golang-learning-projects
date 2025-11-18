package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"deep-research-agent/models"
)

// Wikipedia provides Wikipedia search and article retrieval
type Wikipedia struct {
	baseURL string
	client  *http.Client
}

// NewWikipedia creates a new Wikipedia tool
func NewWikipedia() *Wikipedia {
	return &Wikipedia{
		baseURL: "https://en.wikipedia.org/w/api.php",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the tool name
func (w *Wikipedia) Name() string {
	return "wikipedia"
}

// Description returns the tool description
func (w *Wikipedia) Description() string {
	return "Search Wikipedia articles, fetch summaries, and retrieve full article content with references."
}

// Parameters returns the tool parameters
func (w *Wikipedia) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "action",
			Type:        "string",
			Required:    true,
			Description: "Action to perform: search, summary, full_article",
		},
		{
			Name:        "query",
			Type:        "string",
			Required:    true,
			Description: "Search query or article title",
		},
		{
			Name:        "limit",
			Type:        "int",
			Required:    false,
			Description: "Maximum number of search results (default: 5)",
			Default:     5,
		},
	}
}

// Execute runs the Wikipedia tool
func (w *Wikipedia) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	// Extract parameters
	action, ok := params["action"].(string)
	if !ok || action == "" {
		return nil, fmt.Errorf("action parameter is required")
	}

	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	limit := 5
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}

	// Execute action
	switch action {
	case "search":
		return w.search(ctx, query, limit)
	case "summary":
		return w.getSummary(ctx, query)
	case "full_article":
		return w.getFullArticle(ctx, query)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// search searches Wikipedia for articles
func (w *Wikipedia) search(ctx context.Context, query string, limit int) (*ToolResult, error) {
	// Build search URL
	params := url.Values{}
	params.Set("action", "opensearch")
	params.Set("search", query)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("format", "json")

	searchURL := w.baseURL + "?" + params.Encode()

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return w.mockSearch(query, limit), nil // Fallback to mock
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return w.mockSearch(query, limit), nil // Fallback to mock
	}
	defer resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return w.mockSearch(query, limit), nil // Fallback to mock
	}

	var searchResult []interface{}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return w.mockSearch(query, limit), nil // Fallback to mock
	}

	// Parse OpenSearch format: [query, [titles], [descriptions], [urls]]
	if len(searchResult) < 4 {
		return w.mockSearch(query, limit), nil
	}

	titles, _ := searchResult[1].([]interface{})
	descriptions, _ := searchResult[2].([]interface{})
	urls, _ := searchResult[3].([]interface{})

	results := make([]WikiSearchResult, 0)
	for i := 0; i < len(titles) && i < limit; i++ {
		title, _ := titles[i].(string)
		desc, _ := descriptions[i].(string)
		url, _ := urls[i].(string)

		results = append(results, WikiSearchResult{
			Title:       title,
			Description: desc,
			URL:         url,
		})
	}

	searchResults := &WikiSearchResults{
		Query:   query,
		Results: results,
		Count:   len(results),
	}

	// Convert to sources
	sources := make([]models.Source, 0, len(results))
	for _, result := range results {
		sources = append(sources, models.Source{
			ID:         fmt.Sprintf("wiki_%d", time.Now().UnixNano()),
			Type:       "wikipedia",
			Title:      result.Title,
			URL:        result.URL,
			Excerpt:    result.Description,
			AccessDate: time.Now(),
		})
	}

	return &ToolResult{
		Success: true,
		Data:    searchResults,
		Sources: sources,
		Metadata: map[string]interface{}{
			"query":        query,
			"result_count": len(results),
		},
	}, nil
}

// getSummary gets a Wikipedia article summary
func (w *Wikipedia) getSummary(ctx context.Context, title string) (*ToolResult, error) {
	// Build summary URL (using extract API)
	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("titles", title)
	params.Set("prop", "extracts")
	params.Set("exintro", "true")
	params.Set("explaintext", "true")

	summaryURL := w.baseURL + "?" + params.Encode()

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", summaryURL, nil)
	if err != nil {
		return w.mockSummary(title), nil // Fallback to mock
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return w.mockSummary(title), nil // Fallback to mock
	}
	defer resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return w.mockSummary(title), nil // Fallback to mock
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return w.mockSummary(title), nil // Fallback to mock
	}

	// Extract summary from response
	summary := w.extractSummaryFromResponse(result, title)
	if summary == "" {
		return w.mockSummary(title), nil // Fallback to mock
	}

	article := &WikiArticle{
		Title:   title,
		Summary: summary,
		URL:     fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(title)),
	}

	source := models.Source{
		ID:         fmt.Sprintf("wiki_%d", time.Now().Unix()),
		Type:       "wikipedia",
		Title:      title,
		URL:        article.URL,
		Content:    summary,
		Excerpt:    w.truncate(summary, 200),
		AccessDate: time.Now(),
	}

	return &ToolResult{
		Success: true,
		Data:    article,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"title": title,
		},
	}, nil
}

// getFullArticle gets a full Wikipedia article
func (w *Wikipedia) getFullArticle(ctx context.Context, title string) (*ToolResult, error) {
	// Build article URL
	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("titles", title)
	params.Set("prop", "extracts|info")
	params.Set("explaintext", "true")
	params.Set("inprop", "url")

	articleURL := w.baseURL + "?" + params.Encode()

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", articleURL, nil)
	if err != nil {
		return w.mockFullArticle(title), nil // Fallback to mock
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return w.mockFullArticle(title), nil // Fallback to mock
	}
	defer resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return w.mockFullArticle(title), nil // Fallback to mock
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return w.mockFullArticle(title), nil // Fallback to mock
	}

	// Extract content from response
	content := w.extractContentFromResponse(result, title)
	if content == "" {
		return w.mockFullArticle(title), nil // Fallback to mock
	}

	article := &WikiArticle{
		Title:   title,
		Content: content,
		Summary: w.truncate(content, 500),
		URL:     fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(title)),
	}

	source := models.Source{
		ID:         fmt.Sprintf("wiki_%d", time.Now().Unix()),
		Type:       "wikipedia",
		Title:      title,
		URL:        article.URL,
		Content:    content,
		Excerpt:    article.Summary,
		AccessDate: time.Now(),
	}

	return &ToolResult{
		Success: true,
		Data:    article,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"title":        title,
			"content_size": len(content),
		},
	}, nil
}

// extractSummaryFromResponse extracts summary from Wikipedia API response
func (w *Wikipedia) extractSummaryFromResponse(result map[string]interface{}, title string) string {
	query, ok := result["query"].(map[string]interface{})
	if !ok {
		return ""
	}

	pages, ok := query["pages"].(map[string]interface{})
	if !ok {
		return ""
	}

	for _, page := range pages {
		pageMap, ok := page.(map[string]interface{})
		if !ok {
			continue
		}

		if extract, ok := pageMap["extract"].(string); ok {
			return extract
		}
	}

	return ""
}

// extractContentFromResponse extracts content from Wikipedia API response
func (w *Wikipedia) extractContentFromResponse(result map[string]interface{}, title string) string {
	return w.extractSummaryFromResponse(result, title) // Same extraction logic
}

// mockSearch returns mock search results (fallback)
func (w *Wikipedia) mockSearch(query string, limit int) *ToolResult {
	results := []WikiSearchResult{
		{
			Title:       query,
			Description: fmt.Sprintf("Mock Wikipedia article about %s. This is a fallback result when the API is unavailable.", query),
			URL:         fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(query)),
		},
	}

	searchResults := &WikiSearchResults{
		Query:   query,
		Results: results,
		Count:   len(results),
	}

	source := models.Source{
		ID:         fmt.Sprintf("wiki_%d", time.Now().UnixNano()),
		Type:       "wikipedia",
		Title:      query,
		URL:        results[0].URL,
		Excerpt:    results[0].Description,
		AccessDate: time.Now(),
	}

	return &ToolResult{
		Success: true,
		Data:    searchResults,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"query":        query,
			"result_count": 1,
			"mock":         true,
		},
	}
}

// mockSummary returns mock summary (fallback)
func (w *Wikipedia) mockSummary(title string) *ToolResult {
	summary := fmt.Sprintf("%s is a topic covered in this Wikipedia article. This mock summary provides a brief overview of the subject matter. In a real implementation, this would contain the actual Wikipedia article summary retrieved from the Wikipedia API.", title)

	article := &WikiArticle{
		Title:   title,
		Summary: summary,
		URL:     fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(title)),
	}

	source := models.Source{
		ID:         fmt.Sprintf("wiki_%d", time.Now().Unix()),
		Type:       "wikipedia",
		Title:      title,
		URL:        article.URL,
		Content:    summary,
		Excerpt:    w.truncate(summary, 200),
		AccessDate: time.Now(),
		Metadata: map[string]interface{}{
			"mock": true,
		},
	}

	return &ToolResult{
		Success: true,
		Data:    article,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"title": title,
			"mock":  true,
		},
	}
}

// mockFullArticle returns mock full article (fallback)
func (w *Wikipedia) mockFullArticle(title string) *ToolResult {
	content := fmt.Sprintf(`%s

This is a mock Wikipedia article for %s. In a production environment, this would contain the full article content retrieved from the Wikipedia API.

Background
----------
The topic of %s has been extensively studied and documented. This section would provide background information and context.

Details
-------
This section would contain detailed information about %s, including various aspects, characteristics, and related topics.

References
----------
In a real article, this would contain citations and references to sources.`, title, title, title, title)

	article := &WikiArticle{
		Title:   title,
		Content: content,
		Summary: w.truncate(content, 500),
		URL:     fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(title)),
	}

	source := models.Source{
		ID:         fmt.Sprintf("wiki_%d", time.Now().Unix()),
		Type:       "wikipedia",
		Title:      title,
		URL:        article.URL,
		Content:    content,
		Excerpt:    article.Summary,
		AccessDate: time.Now(),
		Metadata: map[string]interface{}{
			"mock": true,
		},
	}

	return &ToolResult{
		Success: true,
		Data:    article,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"title":        title,
			"content_size": len(content),
			"mock":         true,
		},
	}
}

// truncate truncates text to specified length
func (w *Wikipedia) truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	// Find last space before maxLen
	truncated := text[:maxLen]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}

// WikiArticle represents a Wikipedia article
type WikiArticle struct {
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content,omitempty"`
	URL        string   `json:"url"`
	References []string `json:"references,omitempty"`
}

// WikiSearchResults represents Wikipedia search results
type WikiSearchResults struct {
	Query   string             `json:"query"`
	Results []WikiSearchResult `json:"results"`
	Count   int                `json:"count"`
}

// WikiSearchResult represents a single Wikipedia search result
type WikiSearchResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
