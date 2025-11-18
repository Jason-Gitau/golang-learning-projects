package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"deep-research-agent/config"
	"deep-research-agent/reporting"
	"deep-research-agent/ui"
)

var (
	// Research flags
	depth         string
	maxSources    int
	pdfFiles      []string
	docxFiles     []string
	documents     string
	useWeb        bool
	noWeb         bool
	citeStyle     string
	outputFile    string
	format        string
	concurrent    int
)

// researchCmd represents the research command
var researchCmd = &cobra.Command{
	Use:   "research [query]",
	Short: "Conduct AI-powered research on a topic",
	Long: `Research a topic using multiple sources and AI analysis.

The research agent will:
1. Analyze your query and plan research steps
2. Search multiple sources (web, Wikipedia, documents)
3. Extract and synthesize information
4. Generate a comprehensive report with citations

Examples:
  # Basic research with web sources
  research-agent research "AI ethics in healthcare"

  # Deep research with specific sources
  research-agent research "quantum computing" --depth deep --max-sources 20

  # Analyze documents without web search
  research-agent research "summary of papers" --pdf paper1.pdf --pdf paper2.pdf --no-web

  # Research with custom output
  research-agent research "climate change" --output report.md --cite-style APA`,
	Args: cobra.MinimumNArgs(1),
	RunE: runResearch,
}

func init() {
	rootCmd.AddCommand(researchCmd)

	// Research depth
	researchCmd.Flags().StringVar(&depth, "depth", "", "research depth: shallow, medium, or deep (default from config)")
	researchCmd.Flags().IntVar(&maxSources, "max-sources", 0, "maximum number of sources to use (default from config)")

	// Document sources
	researchCmd.Flags().StringArrayVar(&pdfFiles, "pdf", []string{}, "PDF files to analyze (can specify multiple)")
	researchCmd.Flags().StringArrayVar(&docxFiles, "docx", []string{}, "DOCX files to analyze (can specify multiple)")
	researchCmd.Flags().StringVar(&documents, "documents", "", "document pattern (e.g., 'papers/*.pdf')")

	// Web search
	researchCmd.Flags().BoolVar(&useWeb, "use-web", false, "include web search (default: true unless --no-web)")
	researchCmd.Flags().BoolVar(&noWeb, "no-web", false, "disable web search (only use provided documents)")

	// Output options
	researchCmd.Flags().StringVar(&citeStyle, "cite-style", "", "citation style: APA, MLA, or Chicago (default from config)")
	researchCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path")
	researchCmd.Flags().StringVar(&format, "format", "", "output format: markdown, json, or pdf (default from config)")

	// Performance
	researchCmd.Flags().IntVar(&concurrent, "concurrent", 0, "concurrent tool executions (default from config)")
}

func runResearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")

	// Use config defaults if not specified
	if depth == "" {
		depth = config.GetDefaultDepth()
	}
	if maxSources == 0 {
		maxSources = config.GetMaxSources()
	}
	if citeStyle == "" {
		citeStyle = config.GetCitationStyle()
	}
	if format == "" {
		format = config.GetDefaultFormat()
	}
	if concurrent == 0 {
		concurrent = config.GetConcurrentTools()
	}

	// Validate depth
	if !isValidDepth(depth) {
		return fmt.Errorf("invalid depth: %s (must be shallow, medium, or deep)", depth)
	}

	// Validate citation style
	if !isValidCitationStyle(citeStyle) {
		return fmt.Errorf("invalid citation style: %s (must be APA, MLA, or Chicago)", citeStyle)
	}

	// Validate format
	if !isValidFormat(format) {
		return fmt.Errorf("invalid format: %s (must be markdown, json, or pdf)", format)
	}

	// Check for document files
	allDocs := append(pdfFiles, docxFiles...)

	// Web search logic
	includeWeb := true
	if noWeb {
		includeWeb = false
	} else if useWeb {
		includeWeb = true
	} else if len(allDocs) == 0 {
		// Default to web if no documents provided
		includeWeb = true
	}

	fmt.Printf("Starting research: %s\n", query)
	fmt.Printf("Depth: %s | Max Sources: %d | Citation: %s\n", depth, maxSources, citeStyle)
	if len(allDocs) > 0 {
		fmt.Printf("Documents: %d files\n", len(allDocs))
	}
	if includeWeb {
		fmt.Println("Web search: Enabled")
	} else {
		fmt.Println("Web search: Disabled")
	}
	fmt.Println()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create progress tracker
	progress := ui.NewProgressTracker(query)
	progress.Start()

	// Simulate research process (in real implementation, this would call the agent)
	err := conductResearch(ctx, query, ResearchOptions{
		Depth:       depth,
		MaxSources:  maxSources,
		PDFFiles:    pdfFiles,
		DOCXFiles:   docxFiles,
		Documents:   documents,
		UseWeb:      includeWeb,
		CiteStyle:   citeStyle,
		Concurrent:  concurrent,
	}, progress)

	progress.Stop()

	if err != nil {
		return fmt.Errorf("research failed: %w", err)
	}

	// Generate report
	report := generateReport(query, citeStyle)

	// Output report
	if outputFile != "" {
		if err := saveReport(report, outputFile, format); err != nil {
			return fmt.Errorf("failed to save report: %w", err)
		}
		fmt.Printf("\nReport saved to: %s\n", outputFile)
	} else {
		// Print to stdout
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println(report)
		fmt.Println(strings.Repeat("=", 80))
	}

	return nil
}

// ResearchOptions contains all research configuration
type ResearchOptions struct {
	Depth      string
	MaxSources int
	PDFFiles   []string
	DOCXFiles  []string
	Documents  string
	UseWeb     bool
	CiteStyle  string
	Concurrent int
}

func conductResearch(ctx context.Context, query string, opts ResearchOptions, progress *ui.ProgressTracker) error {
	// Step 1: Initialize research
	progress.UpdateStep("Initializing research plan", 0, 10)
	time.Sleep(500 * time.Millisecond)

	// Step 2: Web search (if enabled)
	if opts.UseWeb {
		progress.UpdateStep("Searching the web", 1, 10)
		progress.AddSource("web")
		time.Sleep(1 * time.Second)

		progress.UpdateStep("Searching Wikipedia", 2, 10)
		progress.AddSource("wikipedia")
		time.Sleep(800 * time.Millisecond)
	}

	// Step 3: Document analysis
	if len(opts.PDFFiles) > 0 {
		for i, pdf := range opts.PDFFiles {
			progress.UpdateStep(fmt.Sprintf("Analyzing %s", pdf), 3+i, 10)
			progress.AddSource("pdf")
			time.Sleep(1 * time.Second)
		}
	}

	// Step 4: Information extraction
	progress.UpdateStep("Extracting key information", 7, 10)
	time.Sleep(1 * time.Second)

	// Step 5: Fact checking
	progress.UpdateStep("Fact checking findings", 8, 10)
	time.Sleep(1 * time.Second)

	// Step 6: Synthesis
	progress.UpdateStep("Synthesizing research", 9, 10)
	time.Sleep(1 * time.Second)

	// Step 7: Report generation
	progress.UpdateStep("Generating report", 10, 10)
	time.Sleep(500 * time.Millisecond)

	return nil
}

func generateReport(query, citeStyle string) string {
	// This is a placeholder - in real implementation, this would use the reporting package
	report := reporting.GenerateMockReport(query, citeStyle)
	return report
}

func saveReport(content, filename, format string) error {
	switch format {
	case "markdown", "md":
		return os.WriteFile(filename, []byte(content), 0644)
	case "json":
		// Convert to JSON format
		jsonReport := reporting.ConvertToJSON(content)
		return os.WriteFile(filename, []byte(jsonReport), 0644)
	case "pdf":
		return reporting.GeneratePDF(content, filename)
	default:
		return os.WriteFile(filename, []byte(content), 0644)
	}
}

func isValidDepth(d string) bool {
	return d == "shallow" || d == "medium" || d == "deep"
}

func isValidCitationStyle(s string) bool {
	s = strings.ToUpper(s)
	return s == "APA" || s == "MLA" || s == "CHICAGO"
}

func isValidFormat(f string) bool {
	f = strings.ToLower(f)
	return f == "markdown" || f == "md" || f == "json" || f == "pdf"
}
