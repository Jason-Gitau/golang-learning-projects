# Research Tools Guide

Complete documentation for all 8 research tools used by the Deep Research Agent.

## Table of Contents

1. [Tool Architecture](#tool-architecture)
2. [Web Search Tool](#1-web-search-tool)
3. [Wikipedia Search Tool](#2-wikipedia-search-tool)
4. [PDF Text Extraction Tool](#3-pdf-text-extraction-tool)
5. [DOCX Text Extraction Tool](#4-docx-text-extraction-tool)
6. [Text Summarization Tool](#5-text-summarization-tool)
7. [Information Extraction Tool](#6-information-extraction-tool)
8. [Fact Checking Tool](#7-fact-checking-tool)
9. [Research Synthesis Tool](#8-research-synthesis-tool)
10. [Tool Integration](#tool-integration)
11. [Tool Configuration](#tool-configuration)
12. [Custom Tools](#custom-tools)

## Tool Architecture

### Tool Interface

All tools implement a common interface:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() []Parameter
    Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}
```

### Tool Result

Tools return structured results:

```go
type ToolResult struct {
    Success  bool                   `json:"success"`
    Data     interface{}            `json:"data"`
    Sources  []models.Source        `json:"sources"`
    Error    error                  `json:"error,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}
```

### Tool Lifecycle

```
1. Tool Selection (by agent)
   ↓
2. Parameter Preparation
   ↓
3. Tool Execution
   ↓
4. Result Processing
   ↓
5. Source Aggregation
```

## 1. Web Search Tool

### Overview

Searches the internet for relevant information using web search APIs.

**Name:** `web_search`

**Purpose:** Find current, real-world information from across the web

**Primary Use Cases:**
- Researching current events
- Finding blog posts and articles
- Discovering multiple perspectives
- Locating technical documentation

### How It Works

```
Input: Query string
  ↓
Search Engine API (DuckDuckGo)
  ↓
Retrieve top N results
  ↓
Extract title, URL, snippet
  ↓
Fetch full content (optional)
  ↓
Return structured results
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Search query |
| `max_results` | int | No | Max results to return (default: 10) |
| `language` | string | No | Language code (default: "en") |
| `fetch_content` | bool | No | Fetch full page content (default: false) |

### Example Usage

**Command:**
```bash
research-agent research "Go concurrency patterns"
```

**Tool Execution:**
```json
{
  "tool": "web_search",
  "params": {
    "query": "Go concurrency patterns",
    "max_results": 10
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "results": [
      {
        "title": "Go Concurrency Patterns",
        "url": "https://go.dev/blog/pipelines",
        "snippet": "Go's concurrency primitives make it easy...",
        "relevance": 0.95
      }
    ],
    "total_results": 10
  },
  "sources": [...]
}
```

### Integration Points

- **Used By:** Research command (default)
- **Triggers:** Web search flag enabled (default)
- **Combines With:** Wikipedia, PDF extraction
- **Output:** Web sources with URLs and content

### Best Practices

**When to Use:**
- Need current information (news, updates)
- Researching technical topics
- Finding diverse perspectives
- Locating online resources

**When to Avoid:**
- Historical research (use Wikipedia)
- Academic papers (use PDFs)
- Internal documents (use document tools)

**Optimization Tips:**
1. Use specific keywords
2. Limit results for faster processing
3. Enable content fetching only when needed
4. Combine with fact-checking for verification

### Configuration

```yaml
tools:
  enable_web_search: true
  web_search_api: "duckduckgo"
  max_results_per_tool: 10
  tool_timeout: 30s
```

## 2. Wikipedia Search Tool

### Overview

Queries Wikipedia for encyclopedic knowledge and background information.

**Name:** `wikipedia_search`

**Purpose:** Retrieve authoritative, encyclopedic information

**Primary Use Cases:**
- Understanding fundamental concepts
- Getting historical context
- Finding authoritative overviews
- Researching well-established topics

### How It Works

```
Input: Query/Topic
  ↓
Wikipedia API Query
  ↓
Search for matching articles
  ↓
Retrieve article content
  ↓
Extract sections and summaries
  ↓
Return structured data
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Topic to search |
| `language` | string | No | Wikipedia language (default: "en") |
| `sections` | []string | No | Specific sections to extract |
| `summary_only` | bool | No | Return only summary (default: false) |

### Example Usage

**Command:**
```bash
research-agent research "Docker containers"
```

**Tool Execution:**
```json
{
  "tool": "wikipedia_search",
  "params": {
    "query": "Docker (software)",
    "summary_only": false
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "title": "Docker (software)",
    "summary": "Docker is a set of platform as a service...",
    "sections": {
      "History": "...",
      "Architecture": "...",
      "Usage": "..."
    },
    "url": "https://en.wikipedia.org/wiki/Docker_(software)"
  }
}
```

### Integration Points

- **Used By:** Research command (default)
- **Triggers:** Topic-based queries
- **Combines With:** Web search for current info
- **Output:** Structured article content

### Best Practices

**When to Use:**
- Need background information
- Researching established concepts
- Understanding fundamentals
- Historical context

**When to Avoid:**
- Very recent events (may not be covered)
- Highly specialized topics (may lack detail)
- Controversial topics (need multiple sources)

**Optimization Tips:**
1. Use proper article titles
2. Request summary for quick overviews
3. Specify sections for targeted extraction
4. Combine with web search for completeness

## 3. PDF Text Extraction Tool

### Overview

Extracts text content from PDF documents for analysis.

**Name:** `pdf_extract`

**Purpose:** Read and process PDF documents

**Primary Use Cases:**
- Analyzing research papers
- Extracting from technical reports
- Processing academic literature
- Reading formal documentation

### How It Works

```
Input: PDF file path
  ↓
Open PDF file
  ↓
Extract text from pages
  ↓
Preserve structure (headings, paragraphs)
  ↓
Extract metadata (title, author)
  ↓
Return structured text
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_path` | string | Yes | Path to PDF file |
| `pages` | []int | No | Specific pages to extract |
| `extract_metadata` | bool | No | Extract PDF metadata (default: true) |
| `preserve_formatting` | bool | No | Keep text formatting (default: true) |

### Example Usage

**Command:**
```bash
research-agent research "Summary of paper" --pdf research-paper.pdf
```

**Tool Execution:**
```json
{
  "tool": "pdf_extract",
  "params": {
    "file_path": "research-paper.pdf",
    "extract_metadata": true
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "text": "Full extracted text...",
    "metadata": {
      "title": "Research Paper Title",
      "author": "Author Name",
      "pages": 15,
      "created": "2024-01-15"
    },
    "sections": [
      {
        "heading": "Introduction",
        "content": "..."
      }
    ]
  }
}
```

### Integration Points

- **Used By:** Research command with `--pdf` flag
- **Triggers:** PDF files specified
- **Combines With:** Summarization, extraction tools
- **Output:** Structured text and metadata

### Best Practices

**When to Use:**
- Research papers and publications
- Technical documentation
- Formal reports
- Academic literature

**When to Avoid:**
- Scanned PDFs without OCR
- Image-heavy PDFs (limited text)
- Password-protected files

**Optimization Tips:**
1. Ensure PDF has selectable text
2. Use page ranges for large documents
3. Enable metadata extraction
4. Combine with summarization for long docs

### Supported Formats

- Standard PDF documents
- Multi-page PDFs
- PDFs with embedded fonts
- Text-based PDFs

**Not Supported:**
- Scanned images (needs OCR)
- Encrypted PDFs
- Forms with form fields

## 4. DOCX Text Extraction Tool

### Overview

Extracts text content from Microsoft Word documents.

**Name:** `docx_extract`

**Purpose:** Read and process Word documents

**Primary Use Cases:**
- Meeting notes and minutes
- Project documentation
- Requirements documents
- Business reports

### How It Works

```
Input: DOCX file path
  ↓
Open DOCX (ZIP archive)
  ↓
Parse document.xml
  ↓
Extract text and formatting
  ↓
Process tables and lists
  ↓
Return structured content
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_path` | string | Yes | Path to DOCX file |
| `extract_tables` | bool | No | Include tables (default: true) |
| `extract_comments` | bool | No | Include comments (default: false) |
| `preserve_formatting` | bool | No | Keep formatting (default: true) |

### Example Usage

**Command:**
```bash
research-agent research "Meeting summary" --docx meeting-notes.docx
```

**Tool Execution:**
```json
{
  "tool": "docx_extract",
  "params": {
    "file_path": "meeting-notes.docx",
    "extract_tables": true
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "text": "Full document text...",
    "metadata": {
      "title": "Meeting Notes",
      "author": "Author Name",
      "created": "2024-01-15",
      "modified": "2024-01-16"
    },
    "tables": [...],
    "lists": [...]
  }
}
```

### Integration Points

- **Used By:** Research command with `--docx` flag
- **Triggers:** DOCX files specified
- **Combines With:** Summarization, extraction
- **Output:** Structured document content

### Best Practices

**When to Use:**
- Meeting notes
- Project documentation
- Requirements
- Business reports

**When to Avoid:**
- Complex formatting (may lose layout)
- Heavy graphics (text only)
- Macros and embedded objects

**Optimization Tips:**
1. Simple documents work best
2. Extract tables for structured data
3. Include comments if needed
4. Combine with summarization

## 5. Text Summarization Tool

### Overview

Generates concise summaries of long text content.

**Name:** `summarize`

**Purpose:** Create brief, informative summaries

**Primary Use Cases:**
- Executive summaries
- Quick overviews
- Document condensation
- Key point extraction

### How It Works

```
Input: Long text
  ↓
Analyze content structure
  ↓
Identify key sentences
  ↓
Extract main points
  ↓
Generate coherent summary
  ↓
Return condensed text
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `text` | string | Yes | Text to summarize |
| `max_length` | int | No | Max summary length (default: 200 words) |
| `style` | string | No | Summary style: brief, detailed (default: brief) |
| `include_bullets` | bool | No | Use bullet points (default: false) |

### Example Usage

**Internal Use:**
```json
{
  "tool": "summarize",
  "params": {
    "text": "Long article text...",
    "max_length": 200,
    "include_bullets": true
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "summary": "Brief summary of content...",
    "key_points": [
      "Point 1",
      "Point 2",
      "Point 3"
    ],
    "word_count": 180
  }
}
```

### Integration Points

- **Used By:** All research workflows
- **Triggers:** Long content processing
- **Combines With:** All extraction tools
- **Output:** Condensed text and key points

### Best Practices

**When to Use:**
- Long documents (>1000 words)
- Multiple sources to combine
- Executive summaries needed
- Quick overviews required

**When to Avoid:**
- Very short text (<100 words)
- Technical details critical
- All nuance must be preserved

**Optimization Tips:**
1. Adjust max_length to needs
2. Use bullets for clarity
3. Choose appropriate style
4. Combine with extraction for structure

## 6. Information Extraction Tool

### Overview

Identifies and extracts specific facts, entities, and relationships.

**Name:** `extract_info`

**Purpose:** Extract structured information from text

**Primary Use Cases:**
- Finding specific facts
- Extracting entities (people, places, dates)
- Building knowledge graphs
- Structured data extraction

### How It Works

```
Input: Text + extraction targets
  ↓
Natural Language Processing
  ↓
Entity Recognition
  ↓
Relationship Extraction
  ↓
Fact Identification
  ↓
Return structured data
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `text` | string | Yes | Text to analyze |
| `targets` | []string | No | Specific entities to find |
| `extract_dates` | bool | No | Extract dates (default: true) |
| `extract_numbers` | bool | No | Extract numbers (default: true) |
| `extract_entities` | bool | No | Extract named entities (default: true) |

### Example Usage

**Internal Use:**
```json
{
  "tool": "extract_info",
  "params": {
    "text": "Docker was created by Solomon Hykes in 2013...",
    "extract_entities": true,
    "extract_dates": true
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "entities": {
      "people": ["Solomon Hykes"],
      "organizations": ["Docker"],
      "technologies": ["Docker", "containers"]
    },
    "dates": ["2013"],
    "facts": [
      {
        "subject": "Docker",
        "predicate": "created by",
        "object": "Solomon Hykes",
        "time": "2013"
      }
    ]
  }
}
```

### Integration Points

- **Used By:** Research synthesis
- **Triggers:** Structured data needs
- **Combines With:** All content sources
- **Output:** Structured facts and entities

### Best Practices

**When to Use:**
- Need specific facts
- Building databases
- Entity identification
- Relationship mapping

**When to Avoid:**
- Unstructured analysis preferred
- Context more important than facts
- Qualitative research

**Optimization Tips:**
1. Specify target entities
2. Enable relevant extractors
3. Combine with fact checking
4. Use for structured output

## 7. Fact Checking Tool

### Overview

Verifies information accuracy across multiple sources.

**Name:** `fact_check`

**Purpose:** Validate claims and ensure accuracy

**Primary Use Cases:**
- Verifying claims
- Cross-referencing sources
- Building confidence
- Ensuring accuracy

### How It Works

```
Input: Claim + sources
  ↓
Extract claim components
  ↓
Search for supporting evidence
  ↓
Compare across sources
  ↓
Check for contradictions
  ↓
Assign confidence score
  ↓
Return verification result
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `claim` | string | Yes | Statement to verify |
| `sources` | []Source | Yes | Sources to check against |
| `threshold` | float | No | Confidence threshold (default: 0.7) |
| `strict_mode` | bool | No | Require unanimous agreement (default: false) |

### Example Usage

**Internal Use:**
```json
{
  "tool": "fact_check",
  "params": {
    "claim": "Go was released in 2009",
    "sources": [...]
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "claim": "Go was released in 2009",
    "status": "verified",
    "confidence": 0.95,
    "supporting_sources": 5,
    "contradicting_sources": 0,
    "details": {
      "evidence": [
        {
          "source": "Wikipedia",
          "supports": true,
          "excerpt": "Go was announced in November 2009"
        }
      ]
    }
  }
}
```

### Integration Points

- **Used By:** Research synthesis
- **Triggers:** Conflicting information
- **Combines With:** All source tools
- **Output:** Verification status and confidence

### Best Practices

**When to Use:**
- Critical facts needed
- Contradictory information found
- High accuracy required
- Academic research

**When to Avoid:**
- Opinions and perspectives
- Very recent events (limited sources)
- Subjective claims

**Optimization Tips:**
1. Provide multiple sources
2. Set appropriate thresholds
3. Review contradictions
4. Combine with synthesis

### Confidence Levels

| Score | Status | Meaning |
|-------|--------|---------|
| 90-100% | Verified | Strong agreement across sources |
| 70-89% | Likely | Most sources agree |
| 50-69% | Uncertain | Mixed evidence |
| 0-49% | Disputed | Contradictory or insufficient evidence |

## 8. Research Synthesis Tool

### Overview

Combines information from multiple sources into coherent findings.

**Name:** `synthesize`

**Purpose:** Create unified insights from multiple sources

**Primary Use Cases:**
- Report generation
- Pattern identification
- Insight creation
- Knowledge consolidation

### How It Works

```
Input: Multiple source contents
  ↓
Analyze all information
  ↓
Identify common themes
  ↓
Group related findings
  ↓
Resolve contradictions
  ↓
Create coherent narrative
  ↓
Return structured findings
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sources` | []Source | Yes | Sources to synthesize |
| `query` | string | Yes | Research query/focus |
| `max_findings` | int | No | Max key findings (default: 5) |
| `include_contradictions` | bool | No | Note contradictions (default: true) |
| `cite_sources` | bool | No | Add source citations (default: true) |

### Example Usage

**Internal Use:**
```json
{
  "tool": "synthesize",
  "params": {
    "sources": [...],
    "query": "Go concurrency patterns",
    "max_findings": 5
  }
}
```

**Result:**
```json
{
  "success": true,
  "data": {
    "findings": [
      {
        "title": "Goroutines for Concurrency",
        "description": "Lightweight threads managed by Go runtime",
        "sources": [1, 2, 5],
        "confidence": 0.92
      }
    ],
    "themes": ["concurrency", "patterns", "best practices"],
    "contradictions": [],
    "summary": "Overall synthesis..."
  }
}
```

### Integration Points

- **Used By:** Report generation
- **Triggers:** Multiple sources gathered
- **Combines With:** All other tools
- **Output:** Structured findings and insights

### Best Practices

**When to Use:**
- Multiple sources gathered
- Report generation needed
- Insight creation required
- Knowledge consolidation

**Always Used:**
- This is the final step in research

**Optimization Tips:**
1. Provide diverse sources
2. Clear research focus
3. Enable contradiction detection
4. Cite sources for credibility

## Tool Integration

### Tool Execution Flow

```
Research Query
    ↓
Query Analysis
    ↓
Tool Selection
    ├─ Web Search (parallel)
    ├─ Wikipedia (parallel)
    ├─ PDF Extract (parallel)
    └─ DOCX Extract (parallel)
    ↓
Content Collection
    ↓
Summarization (per source)
    ↓
Information Extraction
    ↓
Fact Checking
    ↓
Synthesis
    ↓
Report Generation
```

### Parallel Execution

Tools run concurrently for performance:

```go
// Conceptual example
results := make(chan *ToolResult)

go webSearch(query, results)
go wikipediaSearch(query, results)
go pdfExtract(files, results)

// Collect results
allResults := collectResults(results)
```

### Error Handling

Tools handle errors gracefully:

```go
result := &ToolResult{
    Success: false,
    Error: err,
    Metadata: map[string]interface{}{
        "tool": "web_search",
        "retry_count": 2,
    },
}
```

## Tool Configuration

### Global Configuration

```yaml
tools:
  # Enable/disable tools
  enable_web_search: true
  enable_wikipedia: true
  enable_pdf: true
  enable_docx: true

  # Tool settings
  max_results_per_tool: 10
  tool_timeout: 30s

  # Web search
  web_search_api: "duckduckgo"
  web_search_key: ""  # API key if needed
```

### Per-Tool Configuration

```yaml
tools:
  web_search:
    api: "duckduckgo"
    max_results: 10
    fetch_content: true
    timeout: 30s

  wikipedia:
    language: "en"
    summary_only: false
    timeout: 20s

  pdf:
    extract_metadata: true
    max_pages: 100
    timeout: 60s
```

## Custom Tools

### Creating a Custom Tool

Implement the Tool interface:

```go
type CustomTool struct {
    name string
}

func (t *CustomTool) Name() string {
    return t.name
}

func (t *CustomTool) Description() string {
    return "Custom tool description"
}

func (t *CustomTool) Parameters() []Parameter {
    return []Parameter{
        {
            Name: "input",
            Type: "string",
            Required: true,
            Description: "Input parameter",
        },
    }
}

func (t *CustomTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    // Implementation
    return &ToolResult{
        Success: true,
        Data: result,
    }, nil
}
```

### Registering Custom Tools

```go
// Add to tools registry
tools.Register("custom_tool", &CustomTool{name: "custom_tool"})
```

### Best Practices for Custom Tools

1. **Follow Interface:** Implement all required methods
2. **Error Handling:** Return meaningful errors
3. **Timeout Support:** Respect context deadlines
4. **Documentation:** Describe parameters clearly
5. **Testing:** Add comprehensive tests

## Tool Performance

### Benchmarks

Typical execution times (medium depth):

| Tool | Average Time | Notes |
|------|--------------|-------|
| Web Search | 2-3s | Network dependent |
| Wikipedia | 1-2s | Usually fast |
| PDF Extract | 1-5s | Depends on size |
| DOCX Extract | 0.5-2s | Faster than PDF |
| Summarize | 1-2s | Text length dependent |
| Extract Info | 0.5-1s | Fast processing |
| Fact Check | 2-4s | Multiple source checks |
| Synthesize | 2-3s | Final processing |

### Optimization

**Concurrent Execution:**
```bash
research-agent research "topic" --concurrent 5
```

**Reduce Sources:**
```bash
research-agent research "topic" --max-sources 10
```

**Targeted Tools:**
```bash
research-agent research "topic" --pdf file.pdf --no-web
```

## Summary

### Tool Comparison

| Tool | Input | Output | Speed | Use Case |
|------|-------|--------|-------|----------|
| Web Search | Query | Articles | Medium | Current info |
| Wikipedia | Topic | Article | Fast | Background |
| PDF Extract | File | Text | Medium | Papers |
| DOCX Extract | File | Text | Fast | Documents |
| Summarize | Text | Summary | Fast | Condensation |
| Extract Info | Text | Facts | Fast | Data extraction |
| Fact Check | Claim | Verification | Slow | Accuracy |
| Synthesize | Sources | Findings | Medium | Integration |

### Quick Reference

**For Current Information:** Web Search
**For Background:** Wikipedia
**For Papers:** PDF Extract
**For Documents:** DOCX Extract
**For Overviews:** Summarize
**For Facts:** Extract Info
**For Accuracy:** Fact Check
**For Integration:** Synthesize

## Further Reading

- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [START_HERE.md](START_HERE.md) - Learning guide
- [README.md](README.md) - Complete documentation
- Code: `tools/` directory

---

**Tools Guide** - Understanding the Deep Research Agent's capabilities
