#!/bin/bash

# Document Analysis Example
# This script demonstrates document analysis capabilities

echo "====================================="
echo "Document Analysis Example"
echo "====================================="
echo ""

# Example 1: Analyze a single PDF
echo "Example 1: Analyzing a single PDF document"
./research-agent research "Summary of research findings" \
  --pdf papers/research_paper.pdf \
  --no-web \
  --output reports/paper-summary.md
echo ""

# Example 2: Analyze multiple PDFs
echo "Example 2: Analyzing multiple PDF documents"
./research-agent research "Comparative analysis of papers" \
  --pdf papers/paper1.pdf \
  --pdf papers/paper2.pdf \
  --pdf papers/paper3.pdf \
  --depth deep \
  --output reports/comparative-analysis.md
echo ""

# Example 3: Analyze DOCX files
echo "Example 3: Analyzing DOCX documents"
./research-agent research "Summary of meeting notes" \
  --docx notes/meeting_notes.docx \
  --docx notes/project_updates.docx \
  --no-web \
  --output reports/meeting-summary.md
echo ""

# Example 4: Mixed document types
echo "Example 4: Analyzing mixed document types"
./research-agent research "Comprehensive project overview" \
  --pdf reports/technical_spec.pdf \
  --docx notes/requirements.docx \
  --depth medium \
  --output reports/project-overview.md
echo ""

echo "Document analysis complete! Check the reports/ directory."
