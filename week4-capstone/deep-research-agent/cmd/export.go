package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"deep-research-agent/reporting"
)

var (
	exportFormat string
	exportOutput string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [session-id]",
	Short: "Export a research session report",
	Long: `Export a research session report in various formats.

Supported formats:
  - markdown: Markdown report with citations
  - pdf: PDF document (requires pandoc or pdflatex)
  - json: JSON data export

Examples:
  # Export as markdown
  research-agent export abc123 --format markdown --output report.md

  # Export as PDF
  research-agent export abc123 --format pdf --output report.pdf

  # Export as JSON
  research-agent export abc123 --format json --output data.json

  # Export latest session
  research-agent export latest --format markdown`,
	Args: cobra.ExactArgs(1),
	RunE: runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVar(&exportFormat, "format", "markdown", "export format: markdown, pdf, or json")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "output file (default: <session-id>.<format>)")
}

func runExport(cmd *cobra.Command, args []string) error {
	sessionID := args[0]

	// Validate format
	if !isValidFormat(exportFormat) {
		return fmt.Errorf("invalid format: %s (must be markdown, pdf, or json)", exportFormat)
	}

	// Find session
	var session *Session
	if sessionID == "latest" {
		sessions := getMockSessions()
		if len(sessions) > 0 {
			session = &sessions[0]
			sessionID = session.ID
		}
	} else {
		session = findSessionByID(sessionID)
	}

	if session == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Default output filename
	if exportOutput == "" {
		extension := exportFormat
		if extension == "markdown" {
			extension = "md"
		}
		exportOutput = fmt.Sprintf("%s-report.%s", sessionID, extension)
	}

	fmt.Printf("Exporting session: %s\n", sessionID)
	fmt.Printf("Query: %s\n", session.Query)
	fmt.Printf("Format: %s\n", exportFormat)
	fmt.Printf("Output: %s\n\n", exportOutput)

	// Generate report content
	var content string
	switch exportFormat {
	case "markdown", "md":
		content = reporting.GenerateMockReport(session.Query, "APA")
	case "json":
		content = reporting.GenerateJSONReport(session)
	case "pdf":
		// Generate markdown first, then convert to PDF
		mdContent := reporting.GenerateMockReport(session.Query, "APA")
		if err := reporting.GeneratePDF(mdContent, exportOutput); err != nil {
			return fmt.Errorf("failed to generate PDF: %w", err)
		}
		color.Green("\nReport exported successfully!")
		fmt.Printf("Location: %s\n", exportOutput)
		return nil
	}

	// Save content
	if err := saveReport(content, exportOutput, exportFormat); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	color.Green("\nReport exported successfully!")
	fmt.Printf("Location: %s\n", exportOutput)
	fmt.Printf("Size: %s\n", formatFileSize(int64(len(content))))

	return nil
}
