package reporting

import (
	"fmt"
	"strings"
	"time"
)

// Report represents a research report
type Report struct {
	Query       string
	Summary     string
	Findings    []Finding
	Sources     []Source
	Methodology Methodology
	GeneratedAt time.Time
	Depth       string
	CiteStyle   string
}

// Finding represents a key finding
type Finding struct {
	Title       string
	Description string
	Sources     []int // Indices into Sources array
	Confidence  float64
}

// Source represents a research source
type Source struct {
	ID          int
	Type        string
	Title       string
	Author      string
	URL         string
	AccessDate  time.Time
	PublishDate string
	Citation    string
}

// Methodology describes the research process
type Methodology struct {
	ToolsUsed        []string
	DocumentsAnalyzed []string
	WebSources       int
	ResearchSteps    int
	Duration         time.Duration
}

// GenerateMockReport creates a sample research report
func GenerateMockReport(query, citeStyle string) string {
	report := &Report{
		Query:       query,
		Summary:     generateSummary(query),
		Findings:    generateMockFindings(),
		Sources:     generateMockSources(citeStyle),
		Methodology: generateMockMethodology(),
		GeneratedAt: time.Now(),
		Depth:       "medium",
		CiteStyle:   citeStyle,
	}

	return FormatMarkdownReport(report)
}

// FormatMarkdownReport generates a markdown-formatted report
func FormatMarkdownReport(report *Report) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# Research Report: %s\n\n", report.Query))
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n", report.GeneratedAt.Format("January 2, 2006 at 3:04 PM")))
	sb.WriteString(fmt.Sprintf("**Research Depth:** %s\n", strings.Title(report.Depth)))
	sb.WriteString(fmt.Sprintf("**Sources:** %d\n", len(report.Sources)))
	sb.WriteString(fmt.Sprintf("**Citation Style:** %s\n\n", report.CiteStyle))

	// Executive Summary
	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(report.Summary + "\n\n")

	// Key Findings
	sb.WriteString("## Key Findings\n\n")
	for i, finding := range report.Findings {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, finding.Title))
		sb.WriteString(fmt.Sprintf("   - %s\n", finding.Description))
		if len(finding.Sources) > 0 {
			sb.WriteString("   - Sources: ")
			for j, srcIdx := range finding.Sources {
				if j > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(fmt.Sprintf("[%d]", srcIdx+1))
			}
			sb.WriteString("\n")
		}
		if finding.Confidence > 0 {
			sb.WriteString(fmt.Sprintf("   - Confidence: %.0f%%\n", finding.Confidence*100))
		}
		sb.WriteString("\n")
	}

	// Detailed Analysis
	sb.WriteString("## Detailed Analysis\n\n")
	sb.WriteString("### Overview\n\n")
	sb.WriteString(fmt.Sprintf("This research on \"%s\" was conducted using multiple authoritative sources ", report.Query))
	sb.WriteString("and advanced AI-powered analysis. The findings represent a comprehensive synthesis ")
	sb.WriteString("of information from academic papers, web sources, and expert documentation.\n\n")

	sb.WriteString("### Main Themes\n\n")
	sb.WriteString("The research identified several key themes:\n\n")
	sb.WriteString("1. **Current State**: ")
	sb.WriteString(fmt.Sprintf("Research into %s shows significant developments in recent years [1][2].\n\n", report.Query))
	sb.WriteString("2. **Challenges**: ")
	sb.WriteString("Several challenges remain in this area, including technical limitations and ethical considerations [3].\n\n")
	sb.WriteString("3. **Future Directions**: ")
	sb.WriteString("Emerging trends suggest promising developments in the near future [4].\n\n")

	// Sources
	sb.WriteString("## Sources\n\n")
	for i, source := range report.Sources {
		sb.WriteString(fmt.Sprintf("[%d] %s\n\n", i+1, source.Citation))
	}

	// Research Methodology
	sb.WriteString("## Research Methodology\n\n")
	sb.WriteString(fmt.Sprintf("- **Tools used:** %s\n", strings.Join(report.Methodology.ToolsUsed, ", ")))
	if len(report.Methodology.DocumentsAnalyzed) > 0 {
		sb.WriteString(fmt.Sprintf("- **Documents analyzed:** %s\n", strings.Join(report.Methodology.DocumentsAnalyzed, ", ")))
	}
	sb.WriteString(fmt.Sprintf("- **Web sources:** %d\n", report.Methodology.WebSources))
	sb.WriteString(fmt.Sprintf("- **Research steps:** %d\n", report.Methodology.ResearchSteps))
	sb.WriteString(fmt.Sprintf("- **Duration:** %s\n\n", formatDuration(report.Methodology.Duration)))

	// Appendix
	sb.WriteString("## Appendix\n\n")
	sb.WriteString("### Research Steps\n\n")
	sb.WriteString("1. Query analysis and research planning\n")
	sb.WriteString("2. Web search for relevant sources\n")
	sb.WriteString("3. Wikipedia lookup for background information\n")
	sb.WriteString("4. Document analysis and text extraction\n")
	sb.WriteString("5. Information synthesis and fact checking\n")
	sb.WriteString("6. Report generation with citations\n\n")

	sb.WriteString("### Confidence Levels\n\n")
	sb.WriteString("- **High (>80%)**: Information verified by multiple authoritative sources\n")
	sb.WriteString("- **Medium (50-80%)**: Information from credible sources but limited verification\n")
	sb.WriteString("- **Low (<50%)**: Preliminary findings requiring additional verification\n\n")

	// Footer
	sb.WriteString("---\n\n")
	sb.WriteString("*This report was generated by Deep Research Agent, an AI-powered research assistant.*\n")

	return sb.String()
}

func generateSummary(query string) string {
	return fmt.Sprintf("This comprehensive research report on \"%s\" synthesizes information from multiple "+
		"authoritative sources including academic papers, web articles, and expert documentation. "+
		"The analysis reveals key insights, current trends, and future directions in this field. "+
		"All findings are supported by cited sources and presented with confidence levels.", query)
}

func generateMockFindings() []Finding {
	return []Finding{
		{
			Title:       "Current State of Research",
			Description: "Significant progress has been made in recent years with multiple breakthroughs",
			Sources:     []int{0, 1},
			Confidence:  0.85,
		},
		{
			Title:       "Key Challenges Identified",
			Description: "Technical and ethical challenges remain as primary obstacles to advancement",
			Sources:     []int{2, 3},
			Confidence:  0.78,
		},
		{
			Title:       "Emerging Trends",
			Description: "New approaches and methodologies show promising results",
			Sources:     []int{1, 4},
			Confidence:  0.72,
		},
		{
			Title:       "Practical Applications",
			Description: "Real-world implementations demonstrate feasibility and effectiveness",
			Sources:     []int{3, 4},
			Confidence:  0.80,
		},
	}
}

func generateMockSources(citeStyle string) []Source {
	sources := []Source{
		{
			ID:          1,
			Type:        "web",
			Title:       "Recent Advances in the Field",
			Author:      "Smith, J.",
			URL:         "https://example.com/article1",
			AccessDate:  time.Now(),
			PublishDate: "2024",
		},
		{
			ID:          2,
			Type:        "wikipedia",
			Title:       "Comprehensive Overview",
			Author:      "Wikipedia Contributors",
			URL:         "https://en.wikipedia.org/wiki/Topic",
			AccessDate:  time.Now(),
			PublishDate: "2024",
		},
		{
			ID:          3,
			Type:        "pdf",
			Title:       "Technical Analysis and Review",
			Author:      "Johnson, M. & Brown, K.",
			URL:         "",
			AccessDate:  time.Now(),
			PublishDate: "2023",
		},
		{
			ID:          4,
			Type:        "web",
			Title:       "Practical Applications Guide",
			Author:      "Davis, R.",
			URL:         "https://example.com/guide",
			AccessDate:  time.Now(),
			PublishDate: "2024",
		},
		{
			ID:          5,
			Type:        "web",
			Title:       "Future Trends and Predictions",
			Author:      "Wilson, A.",
			URL:         "https://example.com/trends",
			AccessDate:  time.Now(),
			PublishDate: "2024",
		},
	}

	// Format citations based on style
	for i := range sources {
		sources[i].Citation = formatCitation(&sources[i], citeStyle)
	}

	return sources
}

func formatCitation(source *Source, style string) string {
	switch strings.ToUpper(style) {
	case "APA":
		return formatAPACitation(source)
	case "MLA":
		return formatMLACitation(source)
	case "CHICAGO":
		return formatChicagoCitation(source)
	default:
		return formatAPACitation(source)
	}
}

func formatAPACitation(source *Source) string {
	citation := source.Author + ". (" + source.PublishDate + "). "
	citation += "*" + source.Title + "*. "
	if source.URL != "" {
		citation += "Retrieved " + source.AccessDate.Format("January 2, 2006") + ", from " + source.URL
	}
	return citation
}

func formatMLACitation(source *Source) string {
	citation := source.Author + ". \"" + source.Title + ".\" "
	if source.URL != "" {
		citation += "*Web*. " + source.AccessDate.Format("2 Jan. 2006") + ". <" + source.URL + ">"
	} else {
		citation += source.PublishDate + "."
	}
	return citation
}

func formatChicagoCitation(source *Source) string {
	citation := source.Author + ". \"" + source.Title + ".\" "
	if source.URL != "" {
		citation += "Accessed " + source.AccessDate.Format("January 2, 2006") + ". " + source.URL + "."
	} else {
		citation += source.PublishDate + "."
	}
	return citation
}

func generateMockMethodology() Methodology {
	return Methodology{
		ToolsUsed:         []string{"Web Search", "Wikipedia", "PDF Analysis", "Text Extraction", "Summarization"},
		DocumentsAnalyzed: []string{"research_paper.pdf", "technical_report.pdf"},
		WebSources:        10,
		ResearchSteps:     8,
		Duration:          3*time.Minute + 24*time.Second,
	}
}

func formatDuration(d time.Duration) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	if mins > 0 {
		return fmt.Sprintf("%dm %ds", mins, secs)
	}
	return fmt.Sprintf("%ds", secs)
}
