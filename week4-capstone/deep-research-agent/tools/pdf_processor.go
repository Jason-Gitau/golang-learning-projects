package tools

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
	"deep-research-agent/models"
)

// PDFProcessor handles PDF document processing
type PDFProcessor struct{}

// NewPDFProcessor creates a new PDF processor tool
func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{}
}

// Name returns the tool name
func (p *PDFProcessor) Name() string {
	return "pdf_processor"
}

// Description returns the tool description
func (p *PDFProcessor) Description() string {
	return "Extract text, metadata, and content from PDF files. Supports page-by-page extraction, search, and metadata retrieval."
}

// Parameters returns the tool parameters
func (p *PDFProcessor) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "file_path",
			Type:        "string",
			Required:    true,
			Description: "Path to the PDF file",
		},
		{
			Name:        "action",
			Type:        "string",
			Required:    true,
			Description: "Action to perform: extract_text, extract_page, search, get_metadata",
		},
		{
			Name:        "page_number",
			Type:        "int",
			Required:    false,
			Description: "Page number for extract_page action (1-indexed)",
		},
		{
			Name:        "query",
			Type:        "string",
			Required:    false,
			Description: "Search query for search action",
		},
	}
}

// Execute runs the PDF processor
func (p *PDFProcessor) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	// Extract parameters
	filePath, ok := params["file_path"].(string)
	if !ok || filePath == "" {
		return nil, fmt.Errorf("file_path parameter is required")
	}

	action, ok := params["action"].(string)
	if !ok || action == "" {
		return nil, fmt.Errorf("action parameter is required")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("PDF file not found: %s", filePath),
		}, nil
	}

	// Open PDF file
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("failed to open PDF: %w", err),
		}, nil
	}
	defer f.Close()

	// Execute action
	switch action {
	case "extract_text":
		return p.extractText(r, filePath)
	case "extract_page":
		pageNum, ok := params["page_number"].(float64)
		if !ok {
			return nil, fmt.Errorf("page_number parameter is required for extract_page action")
		}
		return p.extractPage(r, filePath, int(pageNum))
	case "search":
		query, ok := params["query"].(string)
		if !ok || query == "" {
			return nil, fmt.Errorf("query parameter is required for search action")
		}
		return p.searchPDF(r, filePath, query)
	case "get_metadata":
		return p.getMetadata(r, filePath)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// extractText extracts all text from the PDF
func (p *PDFProcessor) extractText(r *pdf.Reader, filePath string) (*ToolResult, error) {
	totalPages := r.NumPage()
	pages := make([]PageContent, 0, totalPages)
	var fullText strings.Builder

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		pages = append(pages, PageContent{
			Number: pageNum,
			Text:   text,
		})

		fullText.WriteString(text)
		fullText.WriteString("\n")
	}

	content := &PDFContent{
		FullText:  fullText.String(),
		Pages:     pages,
		PageCount: totalPages,
		Metadata:  p.extractMetadata(r),
	}

	source := models.Source{
		ID:         fmt.Sprintf("pdf_%d", time.Now().Unix()),
		Type:       "pdf",
		Title:      content.Metadata.Title,
		FilePath:   filePath,
		Author:     content.Metadata.Author,
		AccessDate: time.Now(),
		Content:    content.FullText,
		Metadata: map[string]interface{}{
			"page_count": totalPages,
			"subject":    content.Metadata.Subject,
			"keywords":   content.Metadata.Keywords,
		},
	}

	return &ToolResult{
		Success: true,
		Data:    content,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"page_count": totalPages,
			"file_path":  filePath,
		},
	}, nil
}

// extractPage extracts text from a specific page
func (p *PDFProcessor) extractPage(r *pdf.Reader, filePath string, pageNum int) (*ToolResult, error) {
	totalPages := r.NumPage()
	if pageNum < 1 || pageNum > totalPages {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("page number %d out of range (1-%d)", pageNum, totalPages),
		}, nil
	}

	page := r.Page(pageNum)
	if page.V.IsNull() {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("page %d is empty or invalid", pageNum),
		}, nil
	}

	text, err := page.GetPlainText(nil)
	if err != nil {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("failed to extract text from page %d: %w", pageNum, err),
		}, nil
	}

	pageContent := &PageContent{
		Number: pageNum,
		Text:   text,
	}

	source := models.Source{
		ID:         fmt.Sprintf("pdf_%d_page_%d", time.Now().Unix(), pageNum),
		Type:       "pdf",
		FilePath:   filePath,
		AccessDate: time.Now(),
		Content:    text,
		PageNumber: pageNum,
	}

	return &ToolResult{
		Success: true,
		Data:    pageContent,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"page_number": pageNum,
			"total_pages": totalPages,
			"file_path":   filePath,
		},
	}, nil
}

// searchPDF searches for a query in the PDF
func (p *PDFProcessor) searchPDF(r *pdf.Reader, filePath string, query string) (*ToolResult, error) {
	totalPages := r.NumPage()
	matches := make([]SearchMatch, 0)
	queryLower := strings.ToLower(query)

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		textLower := strings.ToLower(text)
		if strings.Contains(textLower, queryLower) {
			// Extract context around the match
			context := p.extractContext(text, query, 100)
			matches = append(matches, SearchMatch{
				PageNumber: pageNum,
				Context:    context,
				Text:       text,
			})
		}
	}

	searchResults := &PDFSearchResults{
		Query:      query,
		Matches:    matches,
		MatchCount: len(matches),
		TotalPages: totalPages,
	}

	sources := make([]models.Source, 0, len(matches))
	for _, match := range matches {
		sources = append(sources, models.Source{
			ID:         fmt.Sprintf("pdf_%d_page_%d", time.Now().Unix(), match.PageNumber),
			Type:       "pdf",
			FilePath:   filePath,
			AccessDate: time.Now(),
			Content:    match.Text,
			Excerpt:    match.Context,
			PageNumber: match.PageNumber,
		})
	}

	return &ToolResult{
		Success: true,
		Data:    searchResults,
		Sources: sources,
		Metadata: map[string]interface{}{
			"match_count": len(matches),
			"query":       query,
			"file_path":   filePath,
		},
	}, nil
}

// getMetadata extracts PDF metadata
func (p *PDFProcessor) getMetadata(r *pdf.Reader, filePath string) (*ToolResult, error) {
	metadata := p.extractMetadata(r)
	totalPages := r.NumPage()

	metadata.PageCount = totalPages

	return &ToolResult{
		Success: true,
		Data:    metadata,
		Metadata: map[string]interface{}{
			"page_count": totalPages,
			"file_path":  filePath,
		},
	}, nil
}

// extractMetadata extracts metadata from PDF reader
func (p *PDFProcessor) extractMetadata(r *pdf.Reader) PDFMetadata {
	metadata := PDFMetadata{
		PageCount: r.NumPage(),
	}

	// Try to extract metadata from Info dictionary
	if info := r.Trailer().Key("Info"); !info.IsNull() {
		if title := info.Key("Title"); !title.IsNull() {
			metadata.Title = title.String()
		}
		if author := info.Key("Author"); !author.IsNull() {
			metadata.Author = author.String()
		}
		if subject := info.Key("Subject"); !subject.IsNull() {
			metadata.Subject = subject.String()
		}
		if keywords := info.Key("Keywords"); !keywords.IsNull() {
			metadata.Keywords = keywords.String()
		}
	}

	return metadata
}

// extractContext extracts context around a search match
func (p *PDFProcessor) extractContext(text, query string, contextLen int) string {
	queryLower := strings.ToLower(query)
	textLower := strings.ToLower(text)

	index := strings.Index(textLower, queryLower)
	if index == -1 {
		return ""
	}

	start := index - contextLen
	if start < 0 {
		start = 0
	}

	end := index + len(query) + contextLen
	if end > len(text) {
		end = len(text)
	}

	context := text[start:end]
	if start > 0 {
		context = "..." + context
	}
	if end < len(text) {
		context = context + "..."
	}

	return context
}

// PDFContent represents extracted PDF content
type PDFContent struct {
	FullText  string        `json:"full_text"`
	Pages     []PageContent `json:"pages"`
	Metadata  PDFMetadata   `json:"metadata"`
	PageCount int           `json:"page_count"`
}

// PageContent represents a single page's content
type PageContent struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

// PDFMetadata represents PDF metadata
type PDFMetadata struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Subject   string `json:"subject"`
	Keywords  string `json:"keywords"`
	PageCount int    `json:"page_count"`
}

// PDFSearchResults represents search results in a PDF
type PDFSearchResults struct {
	Query      string        `json:"query"`
	Matches    []SearchMatch `json:"matches"`
	MatchCount int           `json:"match_count"`
	TotalPages int           `json:"total_pages"`
}

// SearchMatch represents a search match in the PDF
type SearchMatch struct {
	PageNumber int    `json:"page_number"`
	Context    string `json:"context"`
	Text       string `json:"text"`
}
