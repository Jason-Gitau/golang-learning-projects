# START HERE: Complete Learning Guide

Welcome to the Deep Research Agent! This comprehensive guide will teach you everything you need to know about using and understanding this AI-powered research assistant.

## Table of Contents

1. [Introduction](#introduction)
2. [What is Deep Research Agent?](#what-is-deep-research-agent)
3. [How It Works](#how-it-works)
4. [Core Concepts](#core-concepts)
5. [Installation and Setup](#installation-and-setup)
6. [Basic Usage Tutorial](#basic-usage-tutorial)
7. [Advanced Features](#advanced-features)
8. [Research Tools Explained](#research-tools-explained)
9. [Understanding Reports](#understanding-reports)
10. [Best Practices](#best-practices)
11. [Common Patterns](#common-patterns)
12. [Troubleshooting](#troubleshooting)
13. [Learning Exercises](#learning-exercises)
14. [Code Walkthrough](#code-walkthrough)

## Introduction

The Deep Research Agent is a production-ready command-line tool that helps you conduct comprehensive research on any topic. It combines multiple research tools, document analysis, and AI-powered synthesis to create professional reports with proper citations.

### Who Should Use This?

- **Researchers**: Conduct literature reviews and gather sources
- **Students**: Research topics for papers and projects
- **Developers**: Investigate technical topics and best practices
- **Analysts**: Gather and synthesize information from multiple sources
- **Anyone**: Who needs to quickly understand a topic in depth

### What You'll Learn

By the end of this guide, you'll understand:
- How to conduct effective research with the agent
- How the multi-tool architecture works
- How to analyze documents and combine sources
- How to generate professional reports
- How to manage research sessions
- How to extend and customize the agent

## What is Deep Research Agent?

The Deep Research Agent is an AI-powered research assistant that automates the research process. Instead of manually searching, reading, and synthesizing information, you provide a query and the agent:

1. **Plans** the research strategy
2. **Searches** multiple sources (web, Wikipedia, documents)
3. **Extracts** relevant information
4. **Verifies** facts across sources
5. **Synthesizes** findings into coherent insights
6. **Generates** professional reports with citations

### Key Features

**Multi-Source Research**
- Web search for current information
- Wikipedia for encyclopedic knowledge
- PDF document analysis
- DOCX document processing

**Intelligent Processing**
- Query understanding and planning
- Information extraction
- Fact checking and verification
- Cross-source synthesis

**Professional Output**
- Markdown reports with structure
- Multiple citation styles (APA, MLA, Chicago)
- PDF and JSON export
- Source attribution

**Session Management**
- Save and resume research
- Track research history
- Reuse and extend previous work

## How It Works

### The Research Pipeline

```
1. Query Input
   ↓
2. Query Analysis & Planning
   ↓
3. Multi-Tool Execution (Parallel)
   ├─ Web Search
   ├─ Wikipedia Search
   ├─ PDF Analysis
   └─ DOCX Analysis
   ↓
4. Information Extraction
   ↓
5. Fact Checking
   ↓
6. Synthesis & Report Generation
   ↓
7. Output (Markdown/PDF/JSON)
```

### Example Flow

Let's trace a simple research query through the system:

**Query:** "Go concurrency patterns"

**Step 1: Planning**
- Agent analyzes the query
- Determines it's a technical topic
- Plans to use web search and Wikipedia

**Step 2: Tool Execution**
- Web Search: Finds 10 articles on Go concurrency
- Wikipedia: Retrieves Go programming language article
- Concurrent execution for speed

**Step 3: Information Extraction**
- Extracts key patterns: goroutines, channels, select
- Identifies code examples
- Finds best practices

**Step 4: Fact Checking**
- Verifies information across sources
- Checks consistency of recommendations
- Assigns confidence scores

**Step 5: Synthesis**
- Combines findings into themes
- Organizes by importance
- Creates structured findings

**Step 6: Report Generation**
- Formats as markdown report
- Adds proper citations
- Includes methodology section

## Core Concepts

### 1. Research Depth

Research depth determines how comprehensive the investigation is:

**Shallow (Quick Overview)**
- 5-10 sources
- 1-2 minutes
- Best for: Quick facts, basic understanding
- Example: "What is Docker?"

**Medium (Balanced Research)**
- 10-15 sources
- 3-5 minutes
- Best for: Topic understanding, general research
- Example: "Docker vs Kubernetes comparison"

**Deep (Comprehensive Investigation)**
- 15-25+ sources
- 5-10 minutes
- Best for: Academic research, thorough analysis
- Example: "Container orchestration architectures: comprehensive review"

### 2. Sources

The agent gathers information from multiple source types:

**Web Sources**
- Current information
- Multiple perspectives
- Blog posts, articles, documentation

**Wikipedia**
- Encyclopedic knowledge
- Historical context
- Fundamental concepts

**PDF Documents**
- Academic papers
- Technical reports
- Research publications

**DOCX Documents**
- Meeting notes
- Project documentation
- Requirements documents

### 3. Citations

Proper source attribution in academic styles:

**APA (American Psychological Association)**
- Used in: Social sciences, education, psychology
- Format: Author, A. (Year). Title. Source.
- Example: Smith, J. (2024). Go Concurrency Patterns. Go Blog.

**MLA (Modern Language Association)**
- Used in: Humanities, literature, arts
- Format: Author. "Title." Source, Date.
- Example: Smith, John. "Go Concurrency Patterns." Go Blog, 2024.

**Chicago**
- Used in: History, business
- Format: Author. "Title." Source. Date.
- Example: Smith, John. "Go Concurrency Patterns." Go Blog. 2024.

### 4. Research Session

A session represents one complete research investigation:

```
Session {
  ID: "abc123"
  Query: "Go concurrency patterns"
  Created: 2024-01-15 14:30:00
  Sources: 15
  Status: "completed"
  Report: "reports/go-concurrency.md"
}
```

Sessions can be:
- Listed (`session list`)
- Viewed (`session show`)
- Resumed (`session resume`)
- Deleted (`session delete`)

### 5. Tools

Tools are specialized capabilities the agent uses:

1. **Web Search** - Search the internet
2. **Wikipedia Search** - Query Wikipedia
3. **PDF Extraction** - Read PDF files
4. **DOCX Extraction** - Read Word documents
5. **Summarization** - Create summaries
6. **Information Extraction** - Extract facts
7. **Fact Checking** - Verify information
8. **Synthesis** - Combine findings

See [TOOLS_GUIDE.md](TOOLS_GUIDE.md) for detailed documentation.

## Installation and Setup

### Prerequisites

```bash
# Check Go version (need 1.24+)
go version

# Install if needed
# Visit: https://golang.org/dl/
```

### Installation Steps

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/deep-research-agent.git
cd deep-research-agent

# 2. Install dependencies
go mod download

# 3. Build the binary
go build -o research-agent

# 4. Test the installation
./research-agent --version
./research-agent --help
```

### Optional: Global Installation

```bash
# Move to PATH
sudo mv research-agent /usr/local/bin/

# Now you can run from anywhere
research-agent --version
```

### Configuration

Create `~/.research-agent.yaml`:

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

## Basic Usage Tutorial

### Tutorial 1: Your First Research

Let's start with simple web research:

```bash
./research-agent research "What is Docker?"
```

**What happens:**
1. Agent searches the web
2. Queries Wikipedia
3. Extracts information
4. Generates a report
5. Displays to console

**Output:**
```
Starting research: What is Docker?
Depth: medium | Max Sources: 15 | Citation: APA
Web search: Enabled

🔍 Researching: What is Docker?
...
✓ Research completed!
  Sources found: 12
  Time elapsed: 2m 15s

# Research Report: What is Docker?
...
```

### Tutorial 2: Save Report to File

```bash
./research-agent research "Go programming best practices" --output go-practices.md
```

**What's different:**
- Report saved to `go-practices.md`
- Can be viewed later with any text editor

### Tutorial 3: Specify Research Depth

```bash
# Quick overview
./research-agent research "Kubernetes basics" --depth shallow

# Balanced research
./research-agent research "Docker vs Kubernetes" --depth medium

# Comprehensive study
./research-agent research "Container orchestration comparison" --depth deep
```

**Compare outputs:**
- Shallow: ~5 sources, high-level overview
- Medium: ~12 sources, detailed explanation
- Deep: ~20 sources, comprehensive analysis

### Tutorial 4: Analyze Documents

```bash
# Create a sample PDF (or use your own)
echo "Sample content" > sample.txt
# (Convert to PDF using your preferred tool)

# Analyze the document
./research-agent research "Summary of document" --pdf sample.pdf --no-web
```

**What's different:**
- Uses only the PDF (no web search)
- Extracts text from document
- Summarizes content

### Tutorial 5: Multiple Documents

```bash
./research-agent research "Comparative analysis" \
  --pdf paper1.pdf \
  --pdf paper2.pdf \
  --pdf paper3.pdf \
  --depth deep \
  --output comparison.md
```

**Use case:** Compare multiple research papers

### Tutorial 6: Choose Citation Style

```bash
# APA style (default)
./research-agent research "AI ethics" --cite-style APA

# MLA style
./research-agent research "Shakespeare's influence" --cite-style MLA

# Chicago style
./research-agent research "American Revolution" --cite-style Chicago
```

### Tutorial 7: Session Management

```bash
# Conduct research
./research-agent research "Machine learning basics" --output ml-basics.md

# List all sessions
./research-agent session list

# Output:
# ID          Query                    Status      Sources  Created
# abc123      Machine learning basics  completed   14       2024-01-15 14:30

# View session details
./research-agent session show abc123

# Resume for more research
./research-agent session resume abc123
```

### Tutorial 8: Export Formats

```bash
# Research first
./research-agent research "Blockchain technology" --output blockchain.md

# Get session ID
./research-agent session list

# Export as PDF
./research-agent export abc123 --format pdf --output blockchain.pdf

# Export as JSON
./research-agent export abc123 --format json --output blockchain.json
```

### Tutorial 9: Interactive Mode

```bash
./research-agent interactive
```

Follow the prompts:
```
Welcome to Research Agent 🔍

What would you like to research?
> Climate change impacts

Analyzing query: climate change impacts

Should I include web sources? (y/n) y

How deep should I research?
  1. Shallow  (quick overview, 5-10 sources)
  2. Medium   (balanced research, 10-15 sources)
  3. Deep     (comprehensive, 15-25 sources)
> 2

Starting research...
```

## Advanced Features

### Concurrent Processing

Speed up research with parallel tool execution:

```bash
./research-agent research "AI in healthcare" \
  --concurrent 5 \
  --depth deep
```

**How it works:**
- Executes up to 5 tools simultaneously
- Faster completion for complex research
- Trade-off: Higher resource usage

### Document Patterns

Process multiple files with patterns:

```bash
# All PDFs in a directory
./research-agent research "Analysis of all papers" \
  --documents "research-papers/*.pdf"

# Specific pattern
./research-agent research "2024 reports summary" \
  --documents "reports/2024-*.pdf"
```

### Custom Configuration

Override defaults per command:

```bash
./research-agent research "Topic" \
  --depth deep \
  --max-sources 30 \
  --concurrent 5 \
  --cite-style Chicago \
  --format pdf \
  --output custom-report.pdf
```

### Document Indexing

Build a searchable document library:

```bash
# Add documents
./research-agent document add research-paper.pdf --index
./research-agent document add technical-report.pdf --index

# List indexed documents
./research-agent document list

# Search across all documents
./research-agent document search "neural networks"

# Analyze specific document
./research-agent document analyze important-paper.pdf
```

## Research Tools Explained

### 1. Web Search Tool

**Purpose:** Find current information on the internet

**How it works:**
- Uses DuckDuckGo (privacy-focused)
- Searches for query keywords
- Retrieves top N results
- Extracts text and metadata

**Example:**
```
Query: "Go concurrency patterns"
Results: 10 articles from Go blogs, Stack Overflow, documentation
```

**Best for:**
- Current information
- Multiple perspectives
- Tutorials and guides

### 2. Wikipedia Search Tool

**Purpose:** Get encyclopedic knowledge

**How it works:**
- Queries Wikipedia API
- Retrieves article content
- Extracts summaries and sections
- Provides authoritative overview

**Example:**
```
Query: "Docker"
Result: Docker article with history, architecture, usage
```

**Best for:**
- Background information
- Historical context
- Fundamental concepts

### 3. PDF Text Extraction Tool

**Purpose:** Read and analyze PDF documents

**How it works:**
- Opens PDF file
- Extracts text content
- Preserves structure (headings, paragraphs)
- Handles multi-page documents

**Example:**
```
Input: research-paper.pdf
Output: Full text content with metadata
```

**Best for:**
- Academic papers
- Technical reports
- Research publications

### 4. DOCX Text Extraction Tool

**Purpose:** Read Microsoft Word documents

**How it works:**
- Opens DOCX file
- Extracts text and formatting
- Preserves document structure
- Handles tables and lists

**Example:**
```
Input: meeting-notes.docx
Output: Text content with structure
```

**Best for:**
- Meeting notes
- Project documentation
- Requirements documents

### 5. Summarization Tool

**Purpose:** Create concise summaries

**How it works:**
- Analyzes text content
- Identifies key points
- Generates concise summary
- Maintains important details

**Example:**
```
Input: 5000-word article
Output: 200-word summary with key points
```

**Best for:**
- Long documents
- Executive summaries
- Quick overviews

### 6. Information Extraction Tool

**Purpose:** Extract specific facts and data

**How it works:**
- Identifies entities (people, places, dates)
- Extracts key facts
- Finds relationships
- Structures information

**Example:**
```
Input: "Docker was created by Solomon Hykes in 2013"
Output: {
  person: "Solomon Hykes",
  technology: "Docker",
  year: 2013,
  action: "created"
}
```

**Best for:**
- Finding specific data
- Building knowledge graphs
- Extracting structured info

### 7. Fact Checking Tool

**Purpose:** Verify information accuracy

**How it works:**
- Compares claims across sources
- Checks for consistency
- Assigns confidence scores
- Flags contradictions

**Example:**
```
Claim: "Go was released in 2009"
Sources: 3 sources confirm, 0 contradict
Confidence: 95%
Status: Verified
```

**Best for:**
- Validating claims
- Ensuring accuracy
- Building confidence

### 8. Synthesis Tool

**Purpose:** Combine information from multiple sources

**How it works:**
- Analyzes all gathered information
- Identifies themes and patterns
- Combines related findings
- Creates coherent narrative

**Example:**
```
Input: 15 sources on Go concurrency
Output: 4 key findings with supporting evidence
```

**Best for:**
- Creating reports
- Identifying patterns
- Generating insights

## Understanding Reports

### Report Structure

Generated reports follow a professional structure:

```markdown
# Research Report: [Query]

**Generated:** January 15, 2024
**Research Depth:** Medium
**Sources:** 15
**Citation Style:** APA

## Executive Summary
Brief overview of findings (2-3 paragraphs)

## Key Findings
1. **Finding Title**
   - Description
   - Sources: [1], [2]
   - Confidence: 85%

2. **Another Finding**
   - Description
   - Sources: [3], [4]
   - Confidence: 78%

## Detailed Analysis

### Section 1
Detailed discussion with inline citations [1][2]

### Section 2
More detailed analysis [3][4]

## Sources
[1] Smith, J. (2024). Title. URL
[2] Johnson, M. (2023). Title. URL

## Research Methodology
- Tools used: Web Search, Wikipedia, PDF Analysis
- Documents analyzed: paper1.pdf, paper2.pdf
- Web sources: 10
- Research steps: 8
- Duration: 3m 24s

## Appendix
Additional information
```

### Reading Reports

**Executive Summary**
- Start here for quick overview
- Understand main conclusions
- Decide if deeper reading needed

**Key Findings**
- Main discoveries
- Evidence and sources
- Confidence levels

**Detailed Analysis**
- In-depth discussion
- Supporting evidence
- Context and implications

**Sources**
- All references
- Formatted citations
- URLs for web sources

**Methodology**
- How research was conducted
- Tools and sources used
- Process transparency

### Confidence Scores

Reports include confidence levels:

- **High (80-100%)**: Multiple authoritative sources agree
- **Medium (50-79%)**: Credible sources but limited verification
- **Low (<50%)**: Preliminary findings, needs more verification

Example:
```
Finding: "Docker uses containerization"
Sources: 5 authoritative sources
Confidence: 95% (High)
```

## Best Practices

### Writing Good Queries

**Be Specific**
- Bad: "computers"
- Good: "modern computer architecture"
- Better: "x86 vs ARM architecture comparison"

**Use Keywords**
- Include important terms
- Avoid filler words
- Think like a search engine

**Consider Scope**
- Too broad: "programming"
- Too narrow: "golang mutex implementation in version 1.19.2"
- Just right: "Go concurrency patterns"

### Choosing the Right Depth

**Use Shallow When:**
- You need quick facts
- Time is limited
- Topic is simple

**Use Medium When:**
- You need understanding
- Balanced detail is needed
- Time is moderate (3-5 min)

**Use Deep When:**
- Comprehensive research required
- Academic rigor needed
- You have 5-10 minutes

### Combining Sources

**Web + Wikipedia**
```bash
./research-agent research "Topic" --use-web
```
Best for: Current topics with good background

**Documents Only**
```bash
./research-agent research "Summary" --pdf file.pdf --no-web
```
Best for: Analyzing specific documents

**Everything**
```bash
./research-agent research "Topic" \
  --pdf paper.pdf \
  --docx notes.docx \
  --use-web \
  --depth deep
```
Best for: Comprehensive research

### Managing Sessions

**Meaningful Names**
Use descriptive queries:
- Good: "Go concurrency best practices 2024"
- Bad: "test"

**Regular Cleanup**
```bash
# Review old sessions
./research-agent session list

# Delete completed ones you don't need
./research-agent session delete old-id
```

**Export Important Research**
```bash
# Save as multiple formats
./research-agent export abc123 --format markdown -o report.md
./research-agent export abc123 --format pdf -o report.pdf
```

## Common Patterns

### Pattern 1: Literature Review

```bash
# Gather papers
./research-agent research "Machine learning in healthcare: literature review" \
  --pdf paper1.pdf \
  --pdf paper2.pdf \
  --pdf paper3.pdf \
  --pdf paper4.pdf \
  --depth deep \
  --cite-style APA \
  --output literature-review.md
```

### Pattern 2: Topic Exploration

```bash
# Start with shallow research
./research-agent research "Quantum computing" --depth shallow

# Review results, then go deeper
./research-agent research "Quantum computing applications" --depth deep
```

### Pattern 3: Comparative Analysis

```bash
./research-agent research "Docker vs Kubernetes vs Docker Swarm comparison" \
  --depth medium \
  --max-sources 20 \
  --output container-comparison.md
```

### Pattern 4: Document Summarization

```bash
# Multiple related documents
./research-agent research "Summary of Q1 2024 reports" \
  --documents "reports/2024-Q1-*.pdf" \
  --no-web \
  --output q1-summary.md
```

### Pattern 5: Progressive Research

```bash
# Day 1: Initial research
./research-agent research "AI safety" --output ai-safety-v1.md

# Day 2: Add more sources
./research-agent session resume abc123
./research-agent document add new-paper.pdf

# Day 3: Final export
./research-agent export abc123 --format pdf --output final-report.pdf
```

## Troubleshooting

### Common Issues

**Issue: Command not found**
```bash
# Solution: Make executable
chmod +x research-agent

# Or use full path
./research-agent research "topic"
```

**Issue: PDF extraction fails**
```
Error: unable to read PDF

# Solution: Check file permissions
ls -l file.pdf
chmod 644 file.pdf

# Or: Verify PDF is not corrupted
file file.pdf
```

**Issue: Timeout errors**
```
Error: research timeout exceeded

# Solution: Increase timeout in config
agent:
  timeout: 600  # 10 minutes instead of 5
```

**Issue: No results found**
```
# Solution: Try different query
# Bad: "xyz123"
# Good: "introduction to topic xyz"

# Or increase sources
./research-agent research "topic" --max-sources 30
```

### Performance Issues

**Slow Research**
```bash
# Increase concurrency
./research-agent research "topic" --concurrent 5

# Or reduce depth
./research-agent research "topic" --depth shallow
```

**High Memory Usage**
```bash
# Reduce concurrent tools
./research-agent research "topic" --concurrent 2

# Or process fewer documents at once
```

## Learning Exercises

### Exercise 1: Basic Research

Conduct research on a topic you're interested in:

```bash
./research-agent research "YOUR TOPIC HERE" --output my-first-research.md
```

**Goals:**
- Understand the research process
- Read the generated report
- Examine sources

### Exercise 2: Depth Comparison

Research the same topic at different depths:

```bash
./research-agent research "Docker" --depth shallow --output docker-shallow.md
./research-agent research "Docker" --depth medium --output docker-medium.md
./research-agent research "Docker" --depth deep --output docker-deep.md
```

**Goals:**
- Compare report quality
- Understand depth trade-offs
- Note time differences

### Exercise 3: Document Analysis

Find a PDF online and analyze it:

```bash
# Download a paper (e.g., from arxiv.org)
./research-agent research "Summary of paper" --pdf downloaded-paper.pdf --output summary.md
```

**Goals:**
- Use document analysis
- See text extraction in action
- Compare with manual reading

### Exercise 4: Session Management

Practice managing research sessions:

```bash
# Conduct research
./research-agent research "Topic 1" --output topic1.md
./research-agent research "Topic 2" --output topic2.md

# List sessions
./research-agent session list

# Show details
./research-agent session show [ID]

# Export
./research-agent export [ID] --format pdf
```

**Goals:**
- Track multiple projects
- Resume previous work
- Export in formats

### Exercise 5: Advanced Research

Combine all features:

```bash
./research-agent research "Comprehensive topic analysis" \
  --pdf paper1.pdf \
  --pdf paper2.pdf \
  --use-web \
  --depth deep \
  --max-sources 25 \
  --cite-style APA \
  --concurrent 5 \
  --output comprehensive-report.md
```

**Goals:**
- Use all features together
- Create publication-quality report
- Master the tool

## Code Walkthrough

### Project Structure

```
deep-research-agent/
├── main.go                 # Entry point
├── cmd/                    # CLI commands
│   ├── root.go            # Root command setup
│   ├── research.go        # Research command
│   ├── session.go         # Session management
│   ├── document.go        # Document commands
│   ├── export.go          # Export functionality
│   └── interactive.go     # Interactive mode
├── config/                 # Configuration
│   └── config.go          # Config management
├── reporting/              # Report generation
│   ├── generator.go       # Main generator
│   ├── json.go            # JSON export
│   └── pdf.go             # PDF generation
├── ui/                     # User interface
│   ├── progress.go        # Progress tracking
│   └── interactive.go     # Interactive UI
├── tools/                  # Research tools
│   └── interface.go       # Tool interface
└── models/                 # Data models
    └── source.go          # Source model
```

### Key Components

**1. Main Entry Point (main.go)**
```go
func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```
- Starts the CLI
- Handles errors
- Exits appropriately

**2. Research Command (cmd/research.go)**
- Parses flags
- Validates input
- Orchestrates research
- Generates reports

**3. Progress Tracker (ui/progress.go)**
- Shows real-time progress
- Updates status
- Displays sources found
- Tracks time

**4. Report Generator (reporting/generator.go)**
- Structures report
- Formats citations
- Adds methodology
- Exports formats

### How Components Interact

```
User Input (CLI)
    ↓
Root Command (cmd/root.go)
    ↓
Research Command (cmd/research.go)
    ↓
Progress Tracker Started (ui/progress.go)
    ↓
Tools Executed (tools/)
    ↓
Results Collected
    ↓
Report Generated (reporting/generator.go)
    ↓
Output to File/Console
```

### Extending the Agent

**Add a New Command**
1. Create file in `cmd/`
2. Define cobra command
3. Add to root command
4. Implement logic

**Add a New Tool**
1. Implement Tool interface
2. Add to tools registry
3. Document in TOOLS_GUIDE.md

**Add a New Export Format**
1. Create exporter in `reporting/`
2. Add to export command
3. Update documentation

## Next Steps

Now that you understand the Deep Research Agent:

1. **Practice**: Run the learning exercises
2. **Explore**: Try different research topics
3. **Experiment**: Combine features in new ways
4. **Customize**: Adjust configuration to your needs
5. **Extend**: Add new tools or features

### Further Reading

- [TOOLS_GUIDE.md](TOOLS_GUIDE.md) - Detailed tool documentation
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [README.md](README.md) - Complete reference
- [examples/](examples/) - Example scripts

### Getting Help

- Use `--help` flag for any command
- Check documentation files
- Review example scripts
- Run with `--verbose` for debugging

Happy researching!
