package reporting

import (
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDF creates a PDF report from markdown content
func GeneratePDF(markdownContent, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set up fonts
	pdf.SetFont("Arial", "B", 16)

	// Parse and render markdown (simplified version)
	lines := strings.Split(markdownContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			pdf.Ln(5)
			continue
		}

		// Headers
		if strings.HasPrefix(line, "# ") {
			pdf.SetFont("Arial", "B", 18)
			pdf.Cell(0, 10, strings.TrimPrefix(line, "# "))
			pdf.Ln(10)
			pdf.SetFont("Arial", "", 11)
		} else if strings.HasPrefix(line, "## ") {
			pdf.SetFont("Arial", "B", 14)
			pdf.Cell(0, 8, strings.TrimPrefix(line, "## "))
			pdf.Ln(8)
			pdf.SetFont("Arial", "", 11)
		} else if strings.HasPrefix(line, "### ") {
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(0, 7, strings.TrimPrefix(line, "### "))
			pdf.Ln(7)
			pdf.SetFont("Arial", "", 11)
		} else if strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**") {
			// Bold text
			pdf.SetFont("Arial", "B", 11)
			text := strings.Trim(line, "**")
			pdf.MultiCell(0, 6, text, "", "", false)
			pdf.SetFont("Arial", "", 11)
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			// List items
			pdf.Cell(10, 6, "")
			pdf.MultiCell(0, 6, strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "), "", "", false)
		} else if strings.HasPrefix(line, "---") {
			// Horizontal rule
			pdf.Ln(3)
			pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
			pdf.Ln(3)
		} else {
			// Regular text
			pdf.MultiCell(0, 6, line, "", "", false)
		}
	}

	// Save PDF
	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return fmt.Errorf("error saving PDF: %w", err)
	}

	return nil
}

// GeneratePDFFromReport creates a PDF from a Report struct
func GeneratePDFFromReport(report *Report, outputPath string) error {
	markdown := FormatMarkdownReport(report)
	return GeneratePDF(markdown, outputPath)
}
