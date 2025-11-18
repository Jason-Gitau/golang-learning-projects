package tools

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
	"time"

	"deep-research-agent/models"
)

// CitationManager manages citations and bibliography
type CitationManager struct {
	sources map[string]models.Source // sourceID -> Source
}

// NewCitationManager creates a new citation manager tool
func NewCitationManager() *CitationManager {
	return &CitationManager{
		sources: make(map[string]models.Source),
	}
}

// Name returns the tool name
func (c *CitationManager) Name() string {
	return "citation_manager"
}

// Description returns the tool description
func (c *CitationManager) Description() string {
	return "Manage citations and generate bibliographies in various formats (APA, MLA, Chicago). Track sources and detect duplicates."
}

// Parameters returns the tool parameters
func (c *CitationManager) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "action",
			Type:        "string",
			Required:    true,
			Description: "Action: add_source, generate_citation, get_bibliography, deduplicate",
		},
		{
			Name:        "source",
			Type:        "object",
			Required:    false,
			Description: "Source object for add_source action",
		},
		{
			Name:        "sources",
			Type:        "array",
			Required:    false,
			Description: "Array of sources for deduplicate action",
		},
		{
			Name:        "style",
			Type:        "string",
			Required:    false,
			Description: "Citation style: APA, MLA, Chicago (default: APA)",
			Default:     "APA",
		},
		{
			Name:        "source_id",
			Type:        "string",
			Required:    false,
			Description: "Source ID for generate_citation action",
		},
	}
}

// Execute runs the citation manager
func (c *CitationManager) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	// Extract parameters
	action, ok := params["action"].(string)
	if !ok || action == "" {
		return nil, fmt.Errorf("action parameter is required")
	}

	style := "APA"
	if s, ok := params["style"].(string); ok {
		style = strings.ToUpper(s)
	}

	// Execute action
	switch action {
	case "add_source":
		sourceData, ok := params["source"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("source parameter is required for add_source action")
		}
		return c.addSource(sourceData)

	case "generate_citation":
		sourceID, ok := params["source_id"].(string)
		if !ok || sourceID == "" {
			return nil, fmt.Errorf("source_id parameter is required for generate_citation action")
		}
		return c.generateCitation(sourceID, style)

	case "get_bibliography":
		return c.getBibliography(style)

	case "deduplicate":
		sourcesData, ok := params["sources"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("sources parameter is required for deduplicate action")
		}
		return c.deduplicateSources(sourcesData)

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// addSource adds a source to the manager
func (c *CitationManager) addSource(sourceData map[string]interface{}) (*ToolResult, error) {
	// Convert map to Source
	source := c.mapToSource(sourceData)

	// Generate ID if not provided
	if source.ID == "" {
		source.ID = c.generateSourceID(source)
	}

	// Check for duplicates
	if existing, exists := c.sources[source.ID]; exists {
		return &ToolResult{
			Success: true,
			Data: map[string]interface{}{
				"added":     false,
				"duplicate": true,
				"existing":  existing,
			},
			Metadata: map[string]interface{}{
				"message": "Source already exists",
			},
		}, nil
	}

	// Add source
	c.sources[source.ID] = source

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"added":     true,
			"duplicate": false,
			"source":    source,
		},
		Metadata: map[string]interface{}{
			"source_id":    source.ID,
			"total_sources": len(c.sources),
		},
	}, nil
}

// generateCitation generates a citation for a source
func (c *CitationManager) generateCitation(sourceID, style string) (*ToolResult, error) {
	source, exists := c.sources[sourceID]
	if !exists {
		return &ToolResult{
			Success: false,
			Error:   fmt.Errorf("source not found: %s", sourceID),
		}, nil
	}

	var citation string
	switch style {
	case "APA":
		citation = c.formatAPA(source)
	case "MLA":
		citation = c.formatMLA(source)
	case "CHICAGO":
		citation = c.formatChicago(source)
	default:
		return nil, fmt.Errorf("unknown citation style: %s", style)
	}

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"citation":  citation,
			"style":     style,
			"source_id": sourceID,
		},
		Metadata: map[string]interface{}{
			"style": style,
		},
	}, nil
}

// getBibliography generates a complete bibliography
func (c *CitationManager) getBibliography(style string) (*ToolResult, error) {
	if len(c.sources) == 0 {
		return &ToolResult{
			Success: true,
			Data: &CitationData{
				Sources:      []models.Source{},
				Citations:    make(map[string]string),
				Bibliography: "",
				TotalSources: 0,
			},
		}, nil
	}

	// Generate citations for all sources
	citations := make(map[string]string)
	bibliography := make([]string, 0, len(c.sources))

	// Get sorted sources (by title)
	sources := c.getSortedSources()

	for _, source := range sources {
		var citation string
		switch style {
		case "APA":
			citation = c.formatAPA(source)
		case "MLA":
			citation = c.formatMLA(source)
		case "CHICAGO":
			citation = c.formatChicago(source)
		default:
			citation = c.formatAPA(source)
		}

		citations[source.ID] = citation
		bibliography = append(bibliography, citation)
	}

	bibliographyText := strings.Join(bibliography, "\n\n")

	citationData := &CitationData{
		Sources:      sources,
		Citations:    citations,
		Bibliography: bibliographyText,
		TotalSources: len(sources),
	}

	return &ToolResult{
		Success: true,
		Data:    citationData,
		Metadata: map[string]interface{}{
			"style":        style,
			"source_count": len(sources),
		},
	}, nil
}

// deduplicateSources removes duplicate sources
func (c *CitationManager) deduplicateSources(sourcesData []interface{}) (*ToolResult, error) {
	sources := make([]models.Source, 0)
	for _, sourceData := range sourcesData {
		if sourceMap, ok := sourceData.(map[string]interface{}); ok {
			sources = append(sources, c.mapToSource(sourceMap))
		}
	}

	// Use map to track unique sources
	uniqueMap := make(map[string]models.Source)
	duplicates := make([]models.Source, 0)

	for _, source := range sources {
		id := c.generateSourceID(source)

		if existing, exists := uniqueMap[id]; exists {
			duplicates = append(duplicates, existing)
		} else {
			source.ID = id
			uniqueMap[id] = source
		}
	}

	// Convert map to slice
	uniqueSources := make([]models.Source, 0, len(uniqueMap))
	for _, source := range uniqueMap {
		uniqueSources = append(uniqueSources, source)
	}

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"unique_sources":   uniqueSources,
			"duplicates":       duplicates,
			"original_count":   len(sources),
			"unique_count":     len(uniqueSources),
			"duplicate_count":  len(duplicates),
		},
		Metadata: map[string]interface{}{
			"removed": len(duplicates),
		},
	}, nil
}

// formatAPA formats a source in APA style
func (c *CitationManager) formatAPA(source models.Source) string {
	var parts []string

	// Author
	if source.Author != "" {
		parts = append(parts, source.Author+".")
	}

	// Year
	year := ""
	if !source.PublishDate.IsZero() {
		year = fmt.Sprintf("(%d)", source.PublishDate.Year())
	} else {
		year = fmt.Sprintf("(%d)", source.AccessDate.Year())
	}
	parts = append(parts, year)

	// Title
	if source.Title != "" {
		parts = append(parts, source.Title+".")
	}

	// URL and access date
	if source.URL != "" {
		parts = append(parts, fmt.Sprintf("Retrieved from %s", source.URL))
	} else if source.FilePath != "" {
		parts = append(parts, fmt.Sprintf("[%s]", source.FilePath))
	}

	return strings.Join(parts, " ")
}

// formatMLA formats a source in MLA style
func (c *CitationManager) formatMLA(source models.Source) string {
	var parts []string

	// Author (Last, First)
	if source.Author != "" {
		parts = append(parts, source.Author+".")
	}

	// Title (in quotes for articles, italics for books)
	if source.Title != "" {
		if source.Type == "web" || source.Type == "wikipedia" {
			parts = append(parts, fmt.Sprintf("\"%s.\"", source.Title))
		} else {
			parts = append(parts, source.Title+".")
		}
	}

	// Publisher
	if source.Publisher != "" {
		parts = append(parts, source.Publisher+",")
	}

	// Date
	if !source.PublishDate.IsZero() {
		parts = append(parts, source.PublishDate.Format("2 Jan. 2006")+".")
	}

	// URL
	if source.URL != "" {
		parts = append(parts, source.URL+".")
	}

	// Access date
	if source.Type == "web" || source.Type == "wikipedia" {
		parts = append(parts, fmt.Sprintf("Accessed %s.", source.AccessDate.Format("2 Jan. 2006")))
	}

	return strings.Join(parts, " ")
}

// formatChicago formats a source in Chicago style
func (c *CitationManager) formatChicago(source models.Source) string {
	var parts []string

	// Author
	if source.Author != "" {
		parts = append(parts, source.Author+".")
	}

	// Title
	if source.Title != "" {
		parts = append(parts, source.Title+".")
	}

	// Publisher and date
	if source.Publisher != "" && !source.PublishDate.IsZero() {
		parts = append(parts, fmt.Sprintf("%s, %d.", source.Publisher, source.PublishDate.Year()))
	} else if !source.PublishDate.IsZero() {
		parts = append(parts, fmt.Sprintf("%d.", source.PublishDate.Year()))
	}

	// URL
	if source.URL != "" {
		parts = append(parts, source.URL+".")
	}

	return strings.Join(parts, " ")
}

// generateSourceID generates a unique ID for a source based on its content
func (c *CitationManager) generateSourceID(source models.Source) string {
	// Create a hash based on key fields
	key := fmt.Sprintf("%s|%s|%s|%s",
		strings.ToLower(source.Title),
		strings.ToLower(source.Author),
		source.URL,
		source.FilePath,
	)

	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)[:16]
}

// getSortedSources returns sources sorted by title
func (c *CitationManager) getSortedSources() []models.Source {
	sources := make([]models.Source, 0, len(c.sources))
	for _, source := range c.sources {
		sources = append(sources, source)
	}

	sort.Slice(sources, func(i, j int) bool {
		// Sort by author first, then title
		if sources[i].Author != sources[j].Author {
			return sources[i].Author < sources[j].Author
		}
		return sources[i].Title < sources[j].Title
	})

	return sources
}

// mapToSource converts a map to a Source struct
func (c *CitationManager) mapToSource(data map[string]interface{}) models.Source {
	source := models.Source{
		AccessDate: time.Now(),
		Metadata:   make(map[string]interface{}),
	}

	if id, ok := data["id"].(string); ok {
		source.ID = id
	}
	if sourceType, ok := data["type"].(string); ok {
		source.Type = sourceType
	}
	if title, ok := data["title"].(string); ok {
		source.Title = title
	}
	if url, ok := data["url"].(string); ok {
		source.URL = url
	}
	if filePath, ok := data["file_path"].(string); ok {
		source.FilePath = filePath
	}
	if author, ok := data["author"].(string); ok {
		source.Author = author
	}
	if publisher, ok := data["publisher"].(string); ok {
		source.Publisher = publisher
	}
	if content, ok := data["content"].(string); ok {
		source.Content = content
	}
	if excerpt, ok := data["excerpt"].(string); ok {
		source.Excerpt = excerpt
	}
	if pageNum, ok := data["page_number"].(float64); ok {
		source.PageNumber = int(pageNum)
	}
	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		source.Metadata = metadata
	}

	return source
}

// CitationData represents citation information
type CitationData struct {
	Sources      []models.Source   `json:"sources"`
	Citations    map[string]string `json:"citations"`
	Bibliography string            `json:"bibliography"`
	TotalSources int               `json:"total_sources"`
}
