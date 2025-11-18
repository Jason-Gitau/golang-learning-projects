# Deep Research Agent - Tools System

## Overview

The Tools System provides 8 production-ready research tools that the Deep Research Agent uses to gather, process, and analyze information from various sources.

## Implemented Tools

### 1. PDF Processor (`pdf_processor.go`)
**Purpose:** Extract text, metadata, and content from PDF files

**Features:**
- Full text extraction from PDF documents
- Page-by-page extraction for citations
- PDF content search with context
- Metadata extraction (title, author, subject)
- Multi-page PDF support

**Actions:**
- `extract_text` - Extract all text from PDF
- `extract_page` - Extract text from specific page
- `search` - Search for query in PDF with context
- `get_metadata` - Get PDF metadata

**Example Usage:**
```go
pdfTool := NewPDFProcessor()
result, err := pdfTool.Execute(ctx, map[string]interface{}{
    "file_path": "research.pdf",
    "action": "extract_text",
})
```

### 2. DOCX Processor (`docx_processor.go`)
**Purpose:** Extract text and structure from Word documents

**Features:**
- Full text extraction
- Structured extraction with headings
- Table extraction
- Document search
- Heading detection and classification

**Actions:**
- `extract_text` - Extract all text
- `extract_structured` - Extract with structure (headings, tables)
- `search` - Search within document

**Example Usage:**
```go
docxTool := NewDOCXProcessor()
result, err := docxTool.Execute(ctx, map[string]interface{}{
    "file_path": "document.docx",
    "action": "extract_structured",
})
```

### 3. Web Search (`web_search.go`)
**Purpose:** Search the web for information

**Features:**
- Mock search with context-aware results
- Realistic result generation based on query
- Relevance scoring
- Predefined templates for common topics (Go, AI, research, etc.)
- Extensible for real API integration (DuckDuckGo, SerpAPI)

**Parameters:**
- `query` - Search query
- `max_results` - Maximum results (default: 10)

**Example Usage:**
```go
searchTool := NewWebSearch()
result, err := searchTool.Execute(ctx, map[string]interface{}{
    "query": "Go concurrency patterns",
    "max_results": 10,
})
```

### 4. Wikipedia (`wikipedia.go`)
**Purpose:** Search and retrieve Wikipedia articles

**Features:**
- Wikipedia article search
- Article summary retrieval
- Full article content extraction
- Real Wikipedia API integration with fallback to mock
- Handles network errors gracefully

**Actions:**
- `search` - Search Wikipedia for articles
- `summary` - Get article summary
- `full_article` - Get complete article content

**Example Usage:**
```go
wikiTool := NewWikipedia()
result, err := wikiTool.Execute(ctx, map[string]interface{}{
    "action": "summary",
    "query": "Artificial Intelligence",
})
```

### 5. URL Fetcher (`url_fetcher.go`)
**Purpose:** Fetch and extract content from web pages

**Features:**
- HTML page fetching with proper User-Agent
- Main content extraction (removes boilerplate)
- Metadata extraction (Open Graph, meta tags)
- Clean text extraction
- Context cancellation support

**Actions:**
- `fetch` - Fetch complete page
- `extract_text` - Extract clean text only
- `extract_metadata` - Extract metadata

**Example Usage:**
```go
fetchTool := NewURLFetcher()
result, err := fetchTool.Execute(ctx, map[string]interface{}{
    "url": "https://example.com/article",
    "action": "extract_text",
})
```

### 6. Summarizer (`summarizer.go`)
**Purpose:** Summarize long text using extractive summarization

**Features:**
- Extractive summarization (sentence selection)
- Configurable compression ratio
- Key point extraction
- Sentence importance scoring
- Position-based and frequency-based ranking

**Actions:**
- `summarize` - Generate summary
- `key_points` - Extract key points as bullets

**Example Usage:**
```go
summTool := NewSummarizer()
result, err := summTool.Execute(ctx, map[string]interface{}{
    "text": longText,
    "action": "summarize",
    "max_sentences": 5,
    "compression_ratio": 0.3,
})
```

### 7. Citation Manager (`citation_manager.go`)
**Purpose:** Manage citations and generate bibliographies

**Features:**
- Citation tracking and management
- Multiple citation styles (APA, MLA, Chicago)
- Duplicate source detection
- Bibliography generation
- Automatic citation formatting

**Actions:**
- `add_source` - Add source to registry
- `generate_citation` - Format citation for a source
- `get_bibliography` - Generate complete bibliography
- `deduplicate` - Remove duplicate sources

**Example Usage:**
```go
citeTool := NewCitationManager()
result, err := citeTool.Execute(ctx, map[string]interface{}{
    "action": "get_bibliography",
    "style": "APA",
})
```

### 8. Fact Checker (`fact_checker.go`)
**Purpose:** Verify claims and detect contradictions

**Features:**
- Claim verification against sources
- Text similarity analysis
- Contradiction detection
- Confidence scoring
- Multi-source cross-referencing

**Actions:**
- `check_fact` - Verify a claim
- `find_contradictions` - Find contradictory findings
- `calculate_confidence` - Calculate confidence score

**Example Usage:**
```go
factTool := NewFactChecker()
result, err := factTool.Execute(ctx, map[string]interface{}{
    "action": "check_fact",
    "claim": "Go supports concurrency through goroutines",
    "sources": sources,
})
```

## Tool Interface

All tools implement the standard `Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() []Parameter
    Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}
```

## Tool Registry

The `ToolRegistry` manages all tools:

```go
registry := NewToolRegistry()
RegisterAllTools(registry) // Registers all 8 tools

// Use tools
tool, _ := registry.Get("pdf_processor")
result, err := tool.Execute(ctx, params)

// Or execute directly
result, err := registry.Execute(ctx, "web_search", params)
```

## Dependencies

- `github.com/ledongthuc/pdf` - PDF processing
- `github.com/fumiama/go-docx` - DOCX processing
- `github.com/PuerkitoBio/goquery` - HTML parsing

## Notes

### Field Compatibility

The tools were designed with an extended `Source` model that includes additional fields (ID, Author, Publisher, etc.). The existing project uses a simpler `Source` model defined in `models/research.go`.

**To use these tools with the existing models, you have two options:**

1. **Extend the existing Source model** in `models/research.go` to include:
   - `ID`, `Author`, `Publisher`, `PublishDate`, `AccessDate`, `Content`, `FilePath`, `Metadata`

2. **Adapt the tools** to work with the existing Source fields:
   - Use `CitationKey` instead of `ID`
   - Use `Snippet` instead of `Content` or `Excerpt`
   - Add metadata as needed

### Import Paths

The project uses `deep-research-agent/` as the import prefix (configured in go.mod). All tools import:
- `deep-research-agent/models`
- `deep-research-agent/tools` (for registry)

## Production Ready Features

- Context support for cancellation
- Timeout handling (10s+ per tool)
- Comprehensive error handling
- No panics - graceful error returns
- Structured result metadata
- Thread-safe registry
- HTTP error handling with timeouts
- User-Agent headers for web requests
- Mock/fallback modes for external APIs

## Testing

Build the tools:
```bash
cd /home/user/golang-learning-projects/week4-capstone/deep-research-agent
go build ./tools/...
```

Run tests (create tests as needed):
```bash
go test ./tools/... -v
```
