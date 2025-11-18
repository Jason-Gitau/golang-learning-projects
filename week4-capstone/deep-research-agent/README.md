# Deep Research Agent

A production-ready AI-powered research assistant that helps you conduct comprehensive research, analyze documents, and generate professional reports with citations.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage Guide](#usage-guide)
  - [Basic Research](#basic-research)
  - [Document Analysis](#document-analysis)
  - [Session Management](#session-management)
  - [Report Export](#report-export)
- [Commands Reference](#commands-reference)
- [Configuration](#configuration)
- [Research Tools](#research-tools)
- [Report Formats](#report-formats)
- [Examples](#examples)
- [Architecture](#architecture)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Features

### Core Capabilities

- **Multi-Source Research**: Gather information from web searches, Wikipedia, PDF documents, and DOCX files
- **AI-Powered Analysis**: Intelligent query planning, information extraction, and synthesis
- **Professional Reports**: Generate well-structured reports with proper citations in APA, MLA, or Chicago style
- **Document Processing**: Extract and analyze text from PDF and DOCX documents
- **Session Management**: Save, resume, and manage research sessions
- **Interactive Mode**: Guided research through an interactive CLI interface
- **Real-Time Progress**: Visual progress tracking with status updates
- **Multiple Export Formats**: Export reports as Markdown, PDF, or JSON

### Advanced Features

- **Concurrent Processing**: Execute multiple research tools in parallel for faster results
- **Smart Caching**: Cache results to speed up repeated queries
- **Fact Checking**: Verify information across multiple sources
- **Citation Management**: Automatic source citation in multiple styles
- **Document Indexing**: Build a searchable index of your research documents
- **Research Statistics**: Track usage and research performance
- **Customizable Configuration**: Configure defaults, paths, and behavior via YAML

## Quick Start

```bash
# Build the agent
go build -o research-agent

# Conduct basic research
./research-agent research "AI ethics in healthcare"

# Research with documents
./research-agent research "Summary of papers" --pdf paper1.pdf --pdf paper2.pdf

# Export a report
./research-agent export latest --format markdown --output report.md

# Start interactive mode
./research-agent interactive
```

See [QUICK_START.md](QUICK_START.md) for a 5-minute getting started guide.

## Installation

### Prerequisites

- Go 1.24 or later
- (Optional) pandoc or pdflatex for PDF export

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/deep-research-agent.git
cd deep-research-agent

# Install dependencies
go mod download

# Build the binary
go build -o research-agent

# (Optional) Install globally
sudo mv research-agent /usr/local/bin/
```

### Configuration

Create a configuration file at `~/.research-agent.yaml`:

```yaml
agent:
  default_depth: medium
  max_sources: 15
  concurrent_tools: 3
  timeout: 300

citation:
  default_style: APA

output:
  default_format: markdown

storage:
  database_path: ~/.research-agent/data/research.db
  documents_dir: ~/.research-agent/documents
```

## Usage Guide

### Basic Research

Conduct research on any topic with web sources:

```bash
# Simple research
research-agent research "Go programming best practices"

# Specify research depth
research-agent research "quantum computing" --depth deep

# Limit number of sources
research-agent research "climate change" --max-sources 20

# Choose citation style
research-agent research "AI safety" --cite-style APA

# Save to file
research-agent research "blockchain technology" --output blockchain.md
```

### Research Depth Levels

- **Shallow**: Quick overview with 5-10 sources (1-2 minutes)
- **Medium**: Balanced research with 10-15 sources (3-5 minutes)
- **Deep**: Comprehensive research with 15-25+ sources (5-10 minutes)

### Document Analysis

Analyze PDF and DOCX documents:

```bash
# Single PDF
research-agent research "Summary of paper" --pdf research_paper.pdf --no-web

# Multiple PDFs
research-agent research "Comparative analysis" \
  --pdf paper1.pdf \
  --pdf paper2.pdf \
  --pdf paper3.pdf

# DOCX files
research-agent research "Meeting summary" --docx notes.docx

# Mix documents with web search
research-agent research "Comprehensive review" \
  --pdf technical_spec.pdf \
  --docx requirements.docx \
  --use-web

# Document patterns
research-agent research "All papers analysis" --documents "papers/*.pdf"
```

### Session Management

Manage your research sessions:

```bash
# List all sessions
research-agent session list

# View session details
research-agent session show abc123

# Resume a session
research-agent session resume abc123

# Delete a session
research-agent session delete abc123
```

### Document Management

Build and search a document index:

```bash
# Add document to index
research-agent document add paper.pdf --index

# List indexed documents
research-agent document list

# Search documents
research-agent document search "machine learning"

# Analyze a document
research-agent document analyze report.pdf
```

### Report Export

Export research in various formats:

```bash
# Export as Markdown (default)
research-agent export abc123 --format markdown --output report.md

# Export as PDF
research-agent export abc123 --format pdf --output report.pdf

# Export as JSON
research-agent export abc123 --format json --output data.json

# Export latest session
research-agent export latest --format markdown
```

### Interactive Mode

Start an interactive research session:

```bash
research-agent interactive
```

The interactive mode guides you through:
1. Defining your research query
2. Selecting depth and sources
3. Viewing results
4. Exporting reports
5. Adding more sources

### Statistics

View usage statistics:

```bash
research-agent stats
```

Shows:
- Total research sessions
- Sources analyzed
- Document library stats
- Tool usage metrics
- Performance statistics

## Commands Reference

### Research Commands

```bash
research-agent research [query] [flags]
```

**Flags:**
- `--depth` - Research depth: shallow, medium, or deep
- `--max-sources` - Maximum number of sources to use
- `--pdf` - PDF file(s) to analyze (can specify multiple)
- `--docx` - DOCX file(s) to analyze (can specify multiple)
- `--documents` - Document pattern (e.g., `papers/*.pdf`)
- `--use-web` - Include web search
- `--no-web` - Disable web search (only use documents)
- `--cite-style` - Citation style: APA, MLA, or Chicago
- `--output, -o` - Output file path
- `--format` - Output format: markdown, json, or pdf
- `--concurrent` - Concurrent tool executions

### Session Commands

```bash
research-agent session list [--format table|json]
research-agent session show [session-id]
research-agent session resume [session-id]
research-agent session delete [session-id]
```

### Document Commands

```bash
research-agent document add [file] [--index]
research-agent document list
research-agent document search [query]
research-agent document analyze [file]
```

### Export Commands

```bash
research-agent export [session-id] [flags]
```

**Flags:**
- `--format` - Export format: markdown, pdf, or json
- `--output, -o` - Output file (default: `<session-id>.<format>`)

### Utility Commands

```bash
research-agent interactive      # Start interactive mode
research-agent stats            # Show statistics
research-agent help             # Show help
research-agent version          # Show version
```

### Global Flags

- `--config` - Config file path (default: `~/.research-agent.yaml`)
- `--verbose, -v` - Verbose output
- `--help, -h` - Help for any command

## Configuration

### Configuration File

The agent uses a YAML configuration file located at:
- `~/.research-agent.yaml` (home directory)
- `./.research-agent.yaml` (current directory)
- Custom path via `--config` flag

### Configuration Options

```yaml
agent:
  default_depth: medium          # Default research depth
  max_sources: 15                # Max sources to gather
  concurrent_tools: 3            # Concurrent tool executions
  timeout: 300                   # Timeout in seconds

citation:
  default_style: APA             # Default citation style

output:
  default_format: markdown       # Default output format
  template: ""                   # Custom template path

storage:
  database_path: ~/.research-agent/data/research.db
  documents_dir: ~/.research-agent/documents
```

### Environment Variables

All config options can be set via environment variables with the `RESEARCH_AGENT_` prefix:

```bash
export RESEARCH_AGENT_AGENT_MAX_SOURCES=20
export RESEARCH_AGENT_CITATION_DEFAULT_STYLE=MLA
export RESEARCH_AGENT_OUTPUT_DEFAULT_FORMAT=pdf
```

## Research Tools

The agent uses 8 specialized tools for research:

### 1. Web Search
Searches the internet for relevant information using DuckDuckGo or other search engines.

**Use Cases:**
- Finding current information
- Discovering multiple perspectives
- Gathering recent sources

### 2. Wikipedia Search
Queries Wikipedia for encyclopedic knowledge and background information.

**Use Cases:**
- Understanding fundamental concepts
- Getting historical context
- Finding authoritative overviews

### 3. PDF Text Extraction
Extracts text content from PDF documents.

**Use Cases:**
- Analyzing research papers
- Extracting information from reports
- Processing academic literature

### 4. DOCX Text Extraction
Extracts text content from Microsoft Word documents.

**Use Cases:**
- Analyzing meeting notes
- Processing reports
- Extracting requirements documents

### 5. Text Summarization
Generates concise summaries of long text content.

**Use Cases:**
- Creating executive summaries
- Condensing long documents
- Extracting key points

### 6. Information Extraction
Identifies and extracts key facts, entities, and relationships.

**Use Cases:**
- Finding specific data points
- Extracting structured information
- Identifying key concepts

### 7. Fact Checking
Verifies information across multiple sources.

**Use Cases:**
- Validating claims
- Ensuring accuracy
- Building confidence in findings

### 8. Research Synthesis
Combines information from multiple sources into coherent findings.

**Use Cases:**
- Creating comprehensive reports
- Identifying patterns across sources
- Generating insights

See [TOOLS_GUIDE.md](TOOLS_GUIDE.md) for detailed documentation on each tool.

## Report Formats

### Markdown Report

Professional markdown reports include:

```markdown
# Research Report: [Query]

**Generated:** [Date]
**Research Depth:** [Depth]
**Sources:** [N sources]
**Citation Style:** [Style]

## Executive Summary
[Comprehensive summary of findings]

## Key Findings
1. **[Finding Title]**
   - [Description]
   - Sources: [1], [2]
   - Confidence: 85%

## Detailed Analysis
[In-depth analysis with citations]

## Sources
[1] Smith, J. (2024). Title. URL
[2] Johnson, M. (2023). Title. URL

## Research Methodology
- Tools used: [List]
- Documents analyzed: [List]
- Web sources: N
- Research steps: N
- Duration: Xm Ys

## Appendix
[Additional information]
```

### JSON Report

Structured JSON export for programmatic use:

```json
{
  "query": "Research query",
  "summary": "Executive summary",
  "findings": [
    {
      "title": "Finding title",
      "description": "Details",
      "sources": [1, 2],
      "confidence": 0.85
    }
  ],
  "sources": [
    {
      "id": 1,
      "type": "web",
      "title": "Source title",
      "author": "Author name",
      "citation": "Formatted citation"
    }
  ],
  "methodology": {
    "tools_used": ["Web Search", "PDF Analysis"],
    "web_sources": 10,
    "research_steps": 8
  }
}
```

### PDF Report

Professional PDF documents with:
- Formatted headings and sections
- Proper typography
- Citations and references
- Page numbers

## Examples

### Example 1: Basic Web Research

```bash
./research-agent research "Go programming concurrency patterns" \
  --depth medium \
  --max-sources 15 \
  --cite-style APA \
  --output go-concurrency-report.md
```

**Output:** A comprehensive markdown report on Go concurrency patterns with 15 web sources cited in APA style.

### Example 2: Document Analysis

```bash
./research-agent research "Summary of climate change research" \
  --pdf ipcc-report.pdf \
  --pdf climate-study-2024.pdf \
  --no-web \
  --depth deep \
  --output climate-summary.md
```

**Output:** Deep analysis of the provided PDF documents without web sources.

### Example 3: Mixed Sources

```bash
./research-agent research "AI safety comprehensive review" \
  --pdf ai-safety-paper.pdf \
  --docx ai-ethics-notes.docx \
  --use-web \
  --depth deep \
  --max-sources 30 \
  --cite-style Chicago \
  --output ai-safety-review.md
```

**Output:** Comprehensive research combining local documents and web sources.

### Example 4: Batch Processing

```bash
# Research multiple topics
for topic in "quantum computing" "blockchain" "machine learning"; do
  ./research-agent research "$topic" \
    --depth medium \
    --output "reports/${topic// /-}.md"
done
```

### Example 5: Session Workflow

```bash
# Start research
./research-agent research "Neural networks overview" --output nn-report.md

# List sessions to find ID
./research-agent session list

# Show session details
./research-agent session show abc123

# Export in multiple formats
./research-agent export abc123 --format markdown --output nn-report.md
./research-agent export abc123 --format pdf --output nn-report.pdf
./research-agent export abc123 --format json --output nn-data.json
```

See the `examples/` directory for more complete scripts.

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────┐
│                     CLI Interface                        │
│  (Cobra Commands + Interactive Mode + Progress UI)      │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                  Research Agent                          │
│    (Query Planning + Tool Orchestration + Synthesis)    │
└────────────────────┬────────────────────────────────────┘
                     │
        ┌────────────┴────────────┐
        │                         │
┌───────▼────────┐      ┌────────▼─────────┐
│  Research Tools │      │  Storage Layer   │
│  (8 Tools)     │      │  (SQLite + Docs) │
└───────┬────────┘      └────────┬─────────┘
        │                         │
┌───────▼─────────────────────────▼─────────┐
│         Report Generation                  │
│  (Markdown + JSON + PDF)                  │
└───────────────────────────────────────────┘
```

### Components

- **CLI Layer**: Command parsing, user interaction, progress display
- **Agent Layer**: Research planning, tool orchestration, result synthesis
- **Tools Layer**: Specialized research capabilities (web search, PDF analysis, etc.)
- **Storage Layer**: Session persistence, document indexing, caching
- **Reporting Layer**: Report generation in multiple formats

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed system design.

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/yourusername/deep-research-agent.git
cd deep-research-agent

# Install dependencies
go mod download

# Build
go build -o research-agent

# Run tests
go test ./...

# Run with verbose logging
./research-agent research "test query" --verbose
```

### Project Structure

```
deep-research-agent/
├── main.go                     # Entry point
├── cmd/                        # CLI commands
│   ├── root.go                # Root command
│   ├── research.go            # Research command
│   ├── session.go             # Session management
│   ├── document.go            # Document management
│   ├── export.go              # Export command
│   ├── stats.go               # Statistics command
│   └── interactive.go         # Interactive mode
├── config/                     # Configuration management
│   └── config.go
├── reporting/                  # Report generation
│   ├── generator.go           # Main generator
│   ├── markdown.go            # Markdown formatting
│   ├── json.go                # JSON export
│   └── pdf.go                 # PDF generation
├── ui/                         # User interface
│   ├── progress.go            # Progress tracking
│   └── interactive.go         # Interactive session
├── tools/                      # Research tools
│   └── interface.go           # Tool interface
├── storage/                    # Data persistence
├── models/                     # Data models
│   └── source.go
└── examples/                   # Example scripts
    ├── basic_research.sh
    ├── document_analysis.sh
    └── advanced_research.sh
```

### Adding New Tools

Implement the `Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() []Parameter
    Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}
```

See existing tools in `tools/` for examples.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new features
- Update documentation
- Use meaningful commit messages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [Viper](https://github.com/spf13/viper) for configuration
- Progress bars by [progressbar](https://github.com/schollz/progressbar)
- PDF generation with [gofpdf](https://github.com/jung-kurt/gofpdf)

## Support

- Documentation: See [START_HERE.md](START_HERE.md) for learning guide
- Quick Start: See [QUICK_START.md](QUICK_START.md)
- Tools Guide: See [TOOLS_GUIDE.md](TOOLS_GUIDE.md)
- Architecture: See [ARCHITECTURE.md](ARCHITECTURE.md)
- Issues: [GitHub Issues](https://github.com/yourusername/deep-research-agent/issues)

---

**Deep Research Agent** - Conduct comprehensive research with AI assistance.
