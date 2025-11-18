package tools

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fumiama/go-docx"
	"deep-research-agent/models"
)

// DOCXProcessor handles Word document processing
type DOCXProcessor struct{}

// NewDOCXProcessor creates a new DOCX processor tool
func NewDOCXProcessor() *DOCXProcessor {
	return &DOCXProcessor{}
}

// Name returns the tool name
func (d *DOCXProcessor) Name() string {
	return "docx_processor"
}

// Description returns the tool description
func (d *DOCXProcessor) Description() string {
	return "Extract text, structure, and content from Word (.docx) documents. Supports text extraction, heading extraction, and search."
}

// Parameters returns the tool parameters
func (d *DOCXProcessor) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "file_path",
			Type:        "string",
			Required:    true,
			Description: "Path to the DOCX file",
		},
		{
			Name:        "action",
			Type:        "string",
			Required:    true,
			Description: "Action to perform: extract_text, extract_structured, search",
		},
		{
			Name:        "query",
			Type:        "string",
			Required:    false,
			Description: "Search query for search action",
		},
	}
}

// Execute runs the DOCX processor
func (d *DOCXProcessor) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
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
			Error:   fmt.Errorf("DOCX file not found: %s", filePath),
		}, nil
	}

	// Open DOCX file
	doc, err := docx.Open(filePath)
	if err != nil {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("failed to open DOCX: %w", err),
		}, nil
	}
	defer doc.Close()

	// Execute action
	switch action {
	case "extract_text":
		return d.extractText(doc, filePath)
	case "extract_structured":
		return d.extractStructured(doc, filePath)
	case "search":
		query, ok := params["query"].(string)
		if !ok || query == "" {
			return nil, fmt.Errorf("query parameter is required for search action")
		}
		return d.searchDOCX(doc, filePath, query)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// extractText extracts all text from the DOCX
func (d *DOCXProcessor) extractText(doc *docx.Docx, filePath string) (*ToolResult, error) {
	var fullText strings.Builder
	paragraphs := make([]string, 0)

	// Extract text from all paragraphs
	for _, para := range doc.Paragraphs {
		text := para.String()
		if text != "" {
			paragraphs = append(paragraphs, text)
			fullText.WriteString(text)
			fullText.WriteString("\n")
		}
	}

	content := &DOCXContent{
		FullText:   fullText.String(),
		Paragraphs: paragraphs,
	}

	source := models.Source{
		ID:         fmt.Sprintf("docx_%d", time.Now().Unix()),
		Type:       "docx",
		FilePath:   filePath,
		AccessDate: time.Now(),
		Content:    content.FullText,
		Metadata: map[string]interface{}{
			"paragraph_count": len(paragraphs),
		},
	}

	return &ToolResult{
		Success: true,
		Data:    content,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"paragraph_count": len(paragraphs),
			"file_path":       filePath,
		},
	}, nil
}

// extractStructured extracts text with structure (headings, paragraphs)
func (d *DOCXProcessor) extractStructured(doc *docx.Docx, filePath string) (*ToolResult, error) {
	var fullText strings.Builder
	paragraphs := make([]string, 0)
	headings := make([]Heading, 0)
	tables := make([]Table, 0)

	// Extract paragraphs and detect headings
	for _, para := range doc.Paragraphs {
		text := para.String()
		if text == "" {
			continue
		}

		paragraphs = append(paragraphs, text)
		fullText.WriteString(text)
		fullText.WriteString("\n")

		// Check if this is a heading (simple heuristic)
		if d.isHeading(para) {
			level := d.getHeadingLevel(para)
			headings = append(headings, Heading{
				Level: level,
				Text:  text,
			})
		}
	}

	// Extract tables
	for i, table := range doc.Tables {
		rows := make([][]string, 0)
		for _, row := range table.TableRows {
			cells := make([]string, 0)
			for _, cell := range row.TableCells {
				cellText := strings.TrimSpace(cell.String())
				cells = append(cells, cellText)
			}
			if len(cells) > 0 {
				rows = append(rows, cells)
			}
		}

		if len(rows) > 0 {
			tables = append(tables, Table{
				Index: i,
				Rows:  rows,
			})
		}
	}

	content := &DOCXContent{
		FullText:   fullText.String(),
		Paragraphs: paragraphs,
		Headings:   headings,
		Tables:     tables,
	}

	source := models.Source{
		ID:         fmt.Sprintf("docx_%d", time.Now().Unix()),
		Type:       "docx",
		FilePath:   filePath,
		AccessDate: time.Now(),
		Content:    content.FullText,
		Metadata: map[string]interface{}{
			"paragraph_count": len(paragraphs),
			"heading_count":   len(headings),
			"table_count":     len(tables),
		},
	}

	return &ToolResult{
		Success: true,
		Data:    content,
		Sources: []models.Source{source},
		Metadata: map[string]interface{}{
			"paragraph_count": len(paragraphs),
			"heading_count":   len(headings),
			"table_count":     len(tables),
			"file_path":       filePath,
		},
	}, nil
}

// searchDOCX searches for a query in the DOCX
func (d *DOCXProcessor) searchDOCX(doc *docx.Docx, filePath string, query string) (*ToolResult, error) {
	matches := make([]DOCXSearchMatch, 0)
	queryLower := strings.ToLower(query)

	// Search in paragraphs
	for i, para := range doc.Paragraphs {
		text := para.String()
		textLower := strings.ToLower(text)

		if strings.Contains(textLower, queryLower) {
			matches = append(matches, DOCXSearchMatch{
				Type:      "paragraph",
				Index:     i,
				Text:      text,
				Context:   d.extractContext(text, query, 100),
				IsHeading: d.isHeading(para),
			})
		}
	}

	searchResults := &DOCXSearchResults{
		Query:      query,
		Matches:    matches,
		MatchCount: len(matches),
	}

	sources := make([]models.Source, 0, len(matches))
	for _, match := range matches {
		sources = append(sources, models.Source{
			ID:         fmt.Sprintf("docx_%d_match_%d", time.Now().Unix(), match.Index),
			Type:       "docx",
			FilePath:   filePath,
			AccessDate: time.Now(),
			Content:    match.Text,
			Excerpt:    match.Context,
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

// isHeading checks if a paragraph is a heading (simple heuristic)
func (d *DOCXProcessor) isHeading(para docx.Paragraph) bool {
	// Check if paragraph has heading style or is short and bold
	text := para.String()
	if len(text) == 0 {
		return false
	}

	// Simple heuristics:
	// 1. Short text (< 100 chars) and ends without period
	// 2. All caps
	// 3. Contains "Chapter", "Section", etc.
	if len(text) < 100 && !strings.HasSuffix(text, ".") {
		if text == strings.ToUpper(text) {
			return true
		}
		headingKeywords := []string{"Chapter", "Section", "Part", "Appendix"}
		for _, keyword := range headingKeywords {
			if strings.HasPrefix(text, keyword) {
				return true
			}
		}
	}

	return false
}

// getHeadingLevel determines heading level (1-6)
func (d *DOCXProcessor) getHeadingLevel(para docx.Paragraph) int {
	text := para.String()

	// Simple heuristic based on text patterns
	if strings.HasPrefix(text, "Chapter") || strings.HasPrefix(text, "Part") {
		return 1
	}
	if strings.HasPrefix(text, "Section") {
		return 2
	}
	if text == strings.ToUpper(text) {
		return 2
	}

	return 3 // Default level
}

// extractContext extracts context around a search match
func (d *DOCXProcessor) extractContext(text, query string, contextLen int) string {
	queryLower := strings.ToLower(query)
	textLower := strings.ToLower(text)

	index := strings.Index(textLower, queryLower)
	if index == -1 {
		return text
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

// DOCXContent represents extracted DOCX content
type DOCXContent struct {
	FullText   string    `json:"full_text"`
	Paragraphs []string  `json:"paragraphs"`
	Headings   []Heading `json:"headings,omitempty"`
	Tables     []Table   `json:"tables,omitempty"`
}

// Heading represents a document heading
type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

// Table represents a document table
type Table struct {
	Index int        `json:"index"`
	Rows  [][]string `json:"rows"`
}

// DOCXSearchResults represents search results in a DOCX
type DOCXSearchResults struct {
	Query      string            `json:"query"`
	Matches    []DOCXSearchMatch `json:"matches"`
	MatchCount int               `json:"match_count"`
}

// DOCXSearchMatch represents a search match in the DOCX
type DOCXSearchMatch struct {
	Type      string `json:"type"` // paragraph, heading, table
	Index     int    `json:"index"`
	Text      string `json:"text"`
	Context   string `json:"context"`
	IsHeading bool   `json:"is_heading"`
}
