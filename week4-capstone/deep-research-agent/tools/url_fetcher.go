package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"deep-research-agent/models"
)

// URLFetcher fetches and processes web pages
type URLFetcher struct {
	client *http.Client
}

// NewURLFetcher creates a new URL fetcher tool
func NewURLFetcher() *URLFetcher {
	return &URLFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the tool name
func (u *URLFetcher) Name() string {
	return "url_fetcher"
}

// Description returns the tool description
func (u *URLFetcher) Description() string {
	return "Fetch web page content, extract main text, and retrieve metadata. Handles HTML parsing and content extraction."
}

// Parameters returns the tool parameters
func (u *URLFetcher) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "url",
			Type:        "string",
			Required:    true,
			Description: "URL to fetch",
		},
		{
			Name:        "action",
			Type:        "string",
			Required:    false,
			Description: "Action to perform: fetch, extract_text, extract_metadata (default: fetch)",
			Default:     "fetch",
		},
	}
}

// Execute runs the URL fetcher
func (u *URLFetcher) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	// Extract parameters
	urlStr, ok := params["url"].(string)
	if !ok || urlStr == "" {
		return nil, fmt.Errorf("url parameter is required")
	}

	action := "fetch"
	if a, ok := params["action"].(string); ok {
		action = a
	}

	// Fetch the URL
	doc, resp, err := u.fetchURL(ctx, urlStr)
	if err != nil {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("failed to fetch URL: %w", err),
		}, nil
	}

	// Execute action
	switch action {
	case "fetch":
		return u.fetchPage(doc, urlStr, resp)
	case "extract_text":
		return u.extractMainContent(doc, urlStr)
	case "extract_metadata":
		return u.extractMetadata(doc, urlStr)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// fetchURL fetches a URL and parses it with goquery
func (u *URLFetcher) fetchURL(ctx context.Context, urlStr string) (*goquery.Document, *http.Response, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ResearchAgent/1.0)")

	// Make request
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, nil, err
	}

	return doc, resp, nil
}

// fetchPage fetches complete page information
func (u *URLFetcher) fetchPage(doc *goquery.Document, urlStr string, resp *http.Response) (*ToolResult, error) {
	// Extract title
	title := doc.Find("title").First().Text()
	if title == "" {
		title = urlStr
	}

	// Extract metadata
	metadata := u.getMetadataMap(doc)

	// Extract main content
	cleanText := u.extractCleanText(doc)

	// Get raw HTML
	html, _ := doc.Html()

	page := &WebPage{
		URL:       urlStr,
		Title:     title,
		Content:   html,
		CleanText: cleanText,
		Metadata:  metadata,
	}

	source := models.Source{
		ID:         fmt.Sprintf("url_%d", time.Now().Unix()),
		Type:       "web",
		Title:      title,
		URL:        urlStr,
		Content:    cleanText,
		Excerpt:    u.truncate(cleanText, 300),
		AccessDate: time.Now(),
		Metadata:   metadata,
	}

	return &ToolResult{
		Success: true,
		Data:    page,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"url":          urlStr,
			"title":        title,
			"content_size": len(cleanText),
		},
	}, nil
}

// extractMainContent extracts main text content
func (u *URLFetcher) extractMainContent(doc *goquery.Document, urlStr string) (*ToolResult, error) {
	cleanText := u.extractCleanText(doc)
	title := doc.Find("title").First().Text()

	content := &WebContent{
		URL:       urlStr,
		Title:     title,
		CleanText: cleanText,
	}

	source := models.Source{
		ID:         fmt.Sprintf("url_%d", time.Now().Unix()),
		Type:       "web",
		Title:      title,
		URL:        urlStr,
		Content:    cleanText,
		Excerpt:    u.truncate(cleanText, 300),
		AccessDate: time.Now(),
	}

	return &ToolResult{
		Success: true,
		Data:    content,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"url":          urlStr,
			"content_size": len(cleanText),
		},
	}, nil
}

// extractMetadata extracts page metadata
func (u *URLFetcher) extractMetadata(doc *goquery.Document, urlStr string) (*ToolResult, error) {
	metadata := u.getMetadataMap(doc)
	title := doc.Find("title").First().Text()

	metaData := &WebMetadata{
		URL:      urlStr,
		Title:    title,
		Metadata: metadata,
	}

	return &ToolResult{
		Success: true,
		Data:    metaData,
		Metadata: map[string]interface{}{
			"url":   urlStr,
			"title": title,
		},
	}, nil
}

// extractCleanText extracts clean text from HTML, removing scripts and styles
func (u *URLFetcher) extractCleanText(doc *goquery.Document) string {
	// Remove script and style elements
	doc.Find("script, style, nav, footer, header, aside").Remove()

	// Try to find main content area
	var text string

	// Try common content selectors
	contentSelectors := []string{
		"article",
		"main",
		".content",
		".main-content",
		"#content",
		"#main",
		".post-content",
		".entry-content",
	}

	for _, selector := range contentSelectors {
		content := doc.Find(selector).First()
		if content.Length() > 0 {
			text = content.Text()
			break
		}
	}

	// Fallback to body if no content area found
	if text == "" {
		text = doc.Find("body").Text()
	}

	// Clean up whitespace
	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")
	cleanLines := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// getMetadataMap extracts metadata from HTML meta tags
func (u *URLFetcher) getMetadataMap(doc *goquery.Document) map[string]string {
	metadata := make(map[string]string)

	// Extract meta tags
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			if content, exists := s.Attr("content"); exists {
				metadata[name] = content
			}
		}

		if property, exists := s.Attr("property"); exists {
			if content, exists := s.Attr("content"); exists {
				metadata[property] = content
			}
		}
	})

	// Extract common metadata
	if desc := doc.Find("meta[name='description']").First(); desc.Length() > 0 {
		if content, exists := desc.Attr("content"); exists {
			metadata["description"] = content
		}
	}

	if keywords := doc.Find("meta[name='keywords']").First(); keywords.Length() > 0 {
		if content, exists := keywords.Attr("content"); exists {
			metadata["keywords"] = content
		}
	}

	if author := doc.Find("meta[name='author']").First(); author.Length() > 0 {
		if content, exists := author.Attr("content"); exists {
			metadata["author"] = content
		}
	}

	// Open Graph tags
	if ogTitle := doc.Find("meta[property='og:title']").First(); ogTitle.Length() > 0 {
		if content, exists := ogTitle.Attr("content"); exists {
			metadata["og:title"] = content
		}
	}

	if ogDesc := doc.Find("meta[property='og:description']").First(); ogDesc.Length() > 0 {
		if content, exists := ogDesc.Attr("content"); exists {
			metadata["og:description"] = content
		}
	}

	return metadata
}

// truncate truncates text to specified length
func (u *URLFetcher) truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	truncated := text[:maxLen]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}

// WebPage represents a fetched web page
type WebPage struct {
	URL       string            `json:"url"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	CleanText string            `json:"clean_text"`
	Metadata  map[string]string `json:"metadata"`
}

// WebContent represents extracted web content
type WebContent struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	CleanText string `json:"clean_text"`
}

// WebMetadata represents web page metadata
type WebMetadata struct {
	URL      string            `json:"url"`
	Title    string            `json:"title"`
	Metadata map[string]string `json:"metadata"`
}
