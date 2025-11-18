#!/bin/bash

# Advanced Research Example
# This script demonstrates advanced features and workflows

echo "====================================="
echo "Advanced Research Example"
echo "====================================="
echo ""

# Example 1: Multi-source research with web and documents
echo "Example 1: Combining web sources and documents"
./research-agent research "Machine learning in climate science" \
  --pdf papers/ml_climate_paper.pdf \
  --pdf papers/climate_data_analysis.pdf \
  --use-web \
  --depth deep \
  --max-sources 30 \
  --cite-style APA \
  --output reports/ml-climate-research.md
echo ""

# Example 2: Research with document patterns
echo "Example 2: Using document patterns"
./research-agent research "Analysis of all research papers" \
  --documents "papers/*.pdf" \
  --depth medium \
  --concurrent 5 \
  --output reports/all-papers-analysis.md
echo ""

# Example 3: Export in different formats
echo "Example 3: Exporting research in multiple formats"

# First, conduct the research
./research-agent research "Blockchain technology overview" \
  --depth medium \
  --output reports/blockchain.md

# Get the session ID (in real usage, you'd capture this from the previous command)
SESSION_ID="latest"

# Export as PDF
echo "Exporting as PDF..."
./research-agent export $SESSION_ID --format pdf --output reports/blockchain.pdf

# Export as JSON
echo "Exporting as JSON..."
./research-agent export $SESSION_ID --format json --output reports/blockchain.json

echo ""

# Example 4: Session management workflow
echo "Example 4: Session management"

# List all sessions
echo "Listing all research sessions:"
./research-agent session list

# Show specific session details
echo ""
echo "Showing session details:"
./research-agent session show abc123

# Resume a previous session
echo ""
echo "Resuming previous session:"
./research-agent session resume abc123

echo ""

# Example 5: Document management workflow
echo "Example 5: Document management"

# Add documents to index
echo "Adding documents to index:"
./research-agent document add papers/important_paper.pdf --index

# List all indexed documents
echo ""
echo "Listing indexed documents:"
./research-agent document list

# Search indexed documents
echo ""
echo "Searching documents:"
./research-agent document search "machine learning"

# Analyze a document
echo ""
echo "Analyzing document:"
./research-agent document analyze papers/technical_report.pdf

echo ""

# Example 6: High-concurrency research
echo "Example 6: High-performance research with concurrency"
./research-agent research "Comprehensive AI survey" \
  --depth deep \
  --max-sources 50 \
  --concurrent 10 \
  --use-web \
  --documents "papers/*.pdf" \
  --cite-style Chicago \
  --output reports/ai-survey.md
echo ""

# Example 7: Interactive mode
echo "Example 7: Starting interactive mode"
echo "(This will open an interactive session)"
# Uncomment to run:
# ./research-agent interactive

echo ""
echo "Advanced examples complete!"
echo ""
echo "Tip: Use 'research-agent stats' to view usage statistics"
./research-agent stats
