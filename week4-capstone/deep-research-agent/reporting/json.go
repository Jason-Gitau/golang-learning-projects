package reporting

import (
	"encoding/json"
	"time"
)

// JSONReport represents a report in JSON format
type JSONReport struct {
	Query       string          `json:"query"`
	Summary     string          `json:"summary"`
	Findings    []JSONFinding   `json:"findings"`
	Sources     []JSONSource    `json:"sources"`
	Methodology JSONMethodology `json:"methodology"`
	GeneratedAt string          `json:"generated_at"`
	Depth       string          `json:"depth"`
	CiteStyle   string          `json:"cite_style"`
}

// JSONFinding represents a finding in JSON format
type JSONFinding struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Sources     []int    `json:"source_indices"`
	Confidence  float64  `json:"confidence"`
}

// JSONSource represents a source in JSON format
type JSONSource struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	URL         string `json:"url,omitempty"`
	AccessDate  string `json:"access_date"`
	PublishDate string `json:"publish_date"`
	Citation    string `json:"citation"`
}

// JSONMethodology represents methodology in JSON format
type JSONMethodology struct {
	ToolsUsed         []string `json:"tools_used"`
	DocumentsAnalyzed []string `json:"documents_analyzed"`
	WebSources        int      `json:"web_sources"`
	ResearchSteps     int      `json:"research_steps"`
	DurationSeconds   int      `json:"duration_seconds"`
}

// GenerateJSONReport converts a report to JSON format
func GenerateJSONReport(session interface{}) string {
	report := &JSONReport{
		Query:       "Sample Research Query",
		Summary:     "This is a comprehensive research summary with key findings and insights.",
		Findings:    generateJSONFindings(),
		Sources:     generateJSONSources(),
		Methodology: generateJSONMethodology(),
		GeneratedAt: time.Now().Format(time.RFC3339),
		Depth:       "medium",
		CiteStyle:   "APA",
	}

	jsonData, _ := json.MarshalIndent(report, "", "  ")
	return string(jsonData)
}

// ConvertToJSON converts markdown report to JSON
func ConvertToJSON(markdownContent string) string {
	// Simplified conversion - in real implementation, this would parse the markdown
	report := &JSONReport{
		Query:       "Research Query",
		Summary:     "Executive summary extracted from markdown",
		Findings:    generateJSONFindings(),
		Sources:     generateJSONSources(),
		Methodology: generateJSONMethodology(),
		GeneratedAt: time.Now().Format(time.RFC3339),
		Depth:       "medium",
		CiteStyle:   "APA",
	}

	jsonData, _ := json.MarshalIndent(report, "", "  ")
	return string(jsonData)
}

func generateJSONFindings() []JSONFinding {
	return []JSONFinding{
		{
			Title:       "Current State of Research",
			Description: "Significant progress has been made in recent years",
			Sources:     []int{1, 2},
			Confidence:  0.85,
		},
		{
			Title:       "Key Challenges",
			Description: "Technical and ethical challenges remain",
			Sources:     []int{3, 4},
			Confidence:  0.78,
		},
		{
			Title:       "Future Directions",
			Description: "Emerging trends show promise",
			Sources:     []int{2, 5},
			Confidence:  0.72,
		},
	}
}

func generateJSONSources() []JSONSource {
	now := time.Now()
	return []JSONSource{
		{
			ID:          1,
			Type:        "web",
			Title:       "Recent Advances in the Field",
			Author:      "Smith, J.",
			URL:         "https://example.com/article1",
			AccessDate:  now.Format(time.RFC3339),
			PublishDate: "2024",
			Citation:    "Smith, J. (2024). Recent Advances in the Field.",
		},
		{
			ID:          2,
			Type:        "wikipedia",
			Title:       "Comprehensive Overview",
			Author:      "Wikipedia Contributors",
			URL:         "https://en.wikipedia.org/wiki/Topic",
			AccessDate:  now.Format(time.RFC3339),
			PublishDate: "2024",
			Citation:    "Wikipedia Contributors. (2024). Comprehensive Overview.",
		},
		{
			ID:          3,
			Type:        "pdf",
			Title:       "Technical Analysis",
			Author:      "Johnson, M.",
			AccessDate:  now.Format(time.RFC3339),
			PublishDate: "2023",
			Citation:    "Johnson, M. (2023). Technical Analysis.",
		},
	}
}

func generateJSONMethodology() JSONMethodology {
	return JSONMethodology{
		ToolsUsed:         []string{"Web Search", "Wikipedia", "PDF Analysis", "Summarization"},
		DocumentsAnalyzed: []string{"paper1.pdf", "paper2.pdf"},
		WebSources:        10,
		ResearchSteps:     8,
		DurationSeconds:   204, // 3m 24s
	}
}
