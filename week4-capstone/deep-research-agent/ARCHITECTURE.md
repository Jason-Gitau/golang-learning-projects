# Architecture Documentation

Complete system architecture and design documentation for the Deep Research Agent.

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Principles](#architecture-principles)
3. [Component Architecture](#component-architecture)
4. [Data Flow](#data-flow)
5. [Concurrency Model](#concurrency-model)
6. [Storage Architecture](#storage-architecture)
7. [Tool System](#tool-system)
8. [Report Generation](#report-generation)
9. [Error Handling](#error-handling)
10. [Performance Considerations](#performance-considerations)
11. [Security](#security)
12. [Deployment](#deployment)

## System Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      User Interface Layer                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ CLI Commands │  │ Interactive  │  │ Progress Display │  │
│  │   (Cobra)    │  │     Mode     │  │   (Progress UI)  │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Orchestration Layer                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │          Research Agent (Core Logic)                 │   │
│  │  - Query Planning   - Tool Orchestration             │   │
│  │  - Result Synthesis - Report Coordination            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────┬───────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼────────┐ ┌─────▼──────┐ ┌────────▼───────┐
│  Tools Layer   │ │  Storage   │ │   Reporting    │
│  ┌──────────┐  │ │  ┌──────┐  │ │  ┌──────────┐  │
│  │Web Search│  │ │  │SQLite│  │ │  │Markdown  │  │
│  │Wikipedia │  │ │  │  DB  │  │ │  │  JSON    │  │
│  │PDF Reader│  │ │  │      │  │ │  │  PDF     │  │
│  │DOCX Read │  │ │  │Files │  │ │  └──────────┘  │
│  │etc...    │  │ │  └──────┘  │ └────────────────┘
│  └──────────┘  │ └────────────┘
└────────────────┘
```

### Component Overview

| Layer | Components | Responsibility |
|-------|------------|----------------|
| UI Layer | CLI, Interactive Mode, Progress UI | User interaction, command processing, visual feedback |
| Orchestration | Research Agent | Research planning, tool coordination, synthesis |
| Tools Layer | 8 Research Tools | Specialized capabilities (search, extract, analyze) |
| Storage Layer | SQLite DB, File Storage | Session persistence, document indexing |
| Reporting Layer | Generators | Report creation in multiple formats |

## Architecture Principles

### 1. Modularity

Each component has a single, well-defined responsibility:

```
CLI Commands → Handle user input
Agent → Orchestrate research
Tools → Perform specific tasks
Storage → Persist data
Reporting → Generate output
```

### 2. Interface-Based Design

Components interact through well-defined interfaces:

```go
// Tools implement a common interface
type Tool interface {
    Name() string
    Description() string
    Parameters() []Parameter
    Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}

// Reporters implement a common interface
type Reporter interface {
    Generate(report *Report) (string, error)
}
```

### 3. Concurrent by Default

- Tools execute in parallel when possible
- Worker pools manage concurrent operations
- Context-based cancellation throughout

### 4. Configuration-Driven

- YAML configuration for defaults
- Environment variables for overrides
- Command-line flags for runtime control

### 5. Error Resilience

- Graceful degradation when tools fail
- Retry logic for transient errors
- Comprehensive error context

## Component Architecture

### 1. CLI Layer (`cmd/`)

**Responsibility:** Command parsing and user interaction

**Components:**
- `root.go` - Root command, configuration initialization
- `research.go` - Research command and flag handling
- `session.go` - Session management commands
- `document.go` - Document management commands
- `export.go` - Report export commands
- `interactive.go` - Interactive mode
- `stats.go` - Statistics display

**Technology:**
- Cobra for CLI framework
- Viper for configuration
- Flags for parameter handling

**Key Patterns:**

```go
// Command structure
var researchCmd = &cobra.Command{
    Use:   "research [query]",
    Short: "Conduct research",
    RunE:  runResearch,
}

// Flag binding
researchCmd.Flags().StringVar(&depth, "depth", "", "research depth")
```

### 2. Configuration Layer (`config/`)

**Responsibility:** Configuration management

**Components:**
- `config.go` - Config loading, validation, defaults

**Configuration Sources (Priority):**
1. Command-line flags (highest)
2. Environment variables
3. Config file
4. Defaults (lowest)

**Structure:**

```go
type Config struct {
    Agent struct {
        DefaultDepth     string
        MaxSources       int
        ConcurrentTools  int
        Timeout          int
    }
    Citation struct {
        DefaultStyle string
    }
    Output struct {
        DefaultFormat string
        Template      string
    }
    Storage struct {
        DatabasePath  string
        DocumentsDir  string
    }
}
```

### 3. UI Layer (`ui/`)

**Responsibility:** User interface and feedback

**Components:**
- `progress.go` - Real-time progress tracking
- `interactive.go` - Interactive session management

**Progress Tracking:**

```go
type ProgressTracker struct {
    query         string
    currentStep   string
    currentIndex  int
    totalSteps    int
    sourcesFound  int
    startTime     time.Time
    bar           *progressbar.ProgressBar
}

// Usage
progress := ui.NewProgressTracker(query)
progress.Start()
progress.UpdateStep("Searching web", 1, 10)
progress.AddSource("web")
progress.Stop()
```

**Features:**
- Real-time progress bars
- Step-by-step status updates
- Source counting
- Time tracking
- Colored output

### 4. Tools Layer (`tools/`)

**Responsibility:** Research capabilities

**Interface:**

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() []Parameter
    Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}

type ToolResult struct {
    Success  bool
    Data     interface{}
    Sources  []models.Source
    Error    error
    Metadata map[string]interface{}
}
```

**Tool Registry:**

```go
var toolRegistry = map[string]Tool{
    "web_search":      &WebSearchTool{},
    "wikipedia":       &WikipediaTool{},
    "pdf_extract":     &PDFTool{},
    "docx_extract":    &DOCXTool{},
    "summarize":       &SummarizeTool{},
    "extract_info":    &ExtractTool{},
    "fact_check":      &FactCheckTool{},
    "synthesize":      &SynthesisTool{},
}
```

### 5. Reporting Layer (`reporting/`)

**Responsibility:** Report generation in multiple formats

**Components:**
- `generator.go` - Main report generator
- `markdown.go` - Markdown formatting (implicit in generator)
- `json.go` - JSON export
- `pdf.go` - PDF generation

**Report Structure:**

```go
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
```

**Formatters:**

```go
// Markdown report
func FormatMarkdownReport(report *Report) string

// JSON export
func GenerateJSONReport(session interface{}) string

// PDF generation
func GeneratePDF(content, outputPath string) error
```

### 6. Models Layer (`models/`)

**Responsibility:** Data structures

**Key Models:**

```go
type Source struct {
    ID          string
    Type        string  // pdf, docx, web, wikipedia
    Title       string
    URL         string
    FilePath    string
    Author      string
    Publisher   string
    PublishDate time.Time
    AccessDate  time.Time
    Content     string
    Excerpt     string
    Metadata    map[string]interface{}
}

type Finding struct {
    ID         string
    Content    string
    Source     Source
    Confidence float64
    Timestamp  time.Time
    Tags       []string
    Metadata   map[string]interface{}
}
```

## Data Flow

### Research Flow

```
1. User Input
   ↓
   research-agent research "query" --depth deep --pdf file.pdf
   ↓
2. Command Parsing (cmd/research.go)
   ↓
   - Parse flags
   - Validate parameters
   - Load configuration
   ↓
3. Progress Tracker Started (ui/progress.go)
   ↓
4. Research Agent Initialization
   ↓
   - Create context with timeout
   - Initialize tool executor
   - Plan research strategy
   ↓
5. Tool Execution (Parallel)
   ┌──────────────┬──────────────┬──────────────┐
   │ Web Search   │ Wikipedia    │ PDF Extract  │
   │ (goroutine)  │ (goroutine)  │ (goroutine)  │
   └──────┬───────┴──────┬───────┴──────┬───────┘
          │              │              │
6. Result Collection (channels)
   ↓
7. Information Processing
   ├─ Summarization
   ├─ Information Extraction
   └─ Fact Checking
   ↓
8. Synthesis
   ↓
   - Combine findings
   - Resolve contradictions
   - Assign confidence
   ↓
9. Report Generation (reporting/)
   ↓
   - Format as markdown/JSON/PDF
   - Add citations
   - Include methodology
   ↓
10. Output
    ├─ Display to console
    └─ Save to file
```

### Session Persistence Flow

```
Research Session
   ↓
Create Session Object
   ├─ ID (generated)
   ├─ Query
   ├─ Timestamp
   ├─ Sources
   └─ Status
   ↓
Save to SQLite
   ↓
   sessions table:
   - id, query, created_at, updated_at, status
   ↓
   sources table:
   - id, session_id, type, title, url, content
   ↓
   findings table:
   - id, session_id, content, confidence
   ↓
Generate Report Path
   ↓
Update Session with Report Path
```

## Concurrency Model

### Worker Pool Pattern

```go
type WorkerPool struct {
    workers   int
    tasks     chan Task
    results   chan Result
    wg        sync.WaitGroup
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    for task := range wp.tasks {
        result := task.Execute()
        wp.results <- result
    }
}
```

### Tool Execution Pattern

```go
// Concurrent tool execution
func executeTools(ctx context.Context, tools []Tool, params map[string]interface{}) []ToolResult {
    results := make(chan *ToolResult, len(tools))

    // Start workers
    for _, tool := range tools {
        go func(t Tool) {
            result, err := t.Execute(ctx, params)
            if err != nil {
                result = &ToolResult{Success: false, Error: err}
            }
            results <- result
        }(tool)
    }

    // Collect results
    var allResults []ToolResult
    for i := 0; i < len(tools); i++ {
        select {
        case result := <-results:
            allResults = append(allResults, *result)
        case <-ctx.Done():
            return allResults  // Timeout
        }
    }

    return allResults
}
```

### Context-Based Cancellation

```go
// Research with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

// Tools respect context
func (t *WebSearchTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // Proceed with execution
    }
}
```

### Synchronization

**Channels for Communication:**
```go
results := make(chan *ToolResult, numTools)
errors := make(chan error, numTools)
```

**WaitGroups for Coordination:**
```go
var wg sync.WaitGroup
wg.Add(numWorkers)
// ... spawn workers ...
wg.Wait()
```

**Mutexes for Shared State:**
```go
type ProgressTracker struct {
    mu sync.Mutex
    // ... fields ...
}

func (pt *ProgressTracker) UpdateStep(step string) {
    pt.mu.Lock()
    defer pt.mu.Unlock()
    pt.currentStep = step
}
```

## Storage Architecture

### Database Schema

**SQLite Database:**

```sql
-- Sessions table
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    query TEXT NOT NULL,
    depth TEXT,
    max_sources INTEGER,
    cite_style TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT,  -- 'in_progress', 'completed', 'failed'
    report_path TEXT
);

-- Sources table
CREATE TABLE sources (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    type TEXT,  -- 'web', 'wikipedia', 'pdf', 'docx'
    title TEXT,
    url TEXT,
    file_path TEXT,
    author TEXT,
    publisher TEXT,
    publish_date TEXT,
    access_date DATETIME,
    content TEXT,
    excerpt TEXT,
    metadata TEXT,  -- JSON
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Findings table
CREATE TABLE findings (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    content TEXT,
    source_id TEXT,
    confidence REAL,
    timestamp DATETIME,
    tags TEXT,  -- JSON array
    metadata TEXT,  -- JSON
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (source_id) REFERENCES sources(id)
);

-- Document index table
CREATE TABLE documents (
    id TEXT PRIMARY KEY,
    filename TEXT,
    path TEXT,
    type TEXT,
    size INTEGER,
    indexed_at DATETIME,
    page_count INTEGER,
    word_count INTEGER,
    content TEXT,
    metadata TEXT  -- JSON
);

-- Create indexes
CREATE INDEX idx_sessions_created ON sessions(created_at);
CREATE INDEX idx_sources_session ON sources(session_id);
CREATE INDEX idx_findings_session ON findings(session_id);
CREATE INDEX idx_documents_filename ON documents(filename);
```

### File Storage

```
~/.research-agent/
├── data/
│   └── research.db          # SQLite database
├── documents/               # Indexed documents
│   ├── paper1.pdf
│   ├── paper2.pdf
│   └── notes.docx
├── cache/                   # Tool result cache
│   ├── web_search/
│   └── wikipedia/
└── reports/                 # Generated reports
    ├── session-abc123.md
    ├── session-def456.pdf
    └── session-ghi789.json
```

### Caching Strategy

```go
type Cache struct {
    store map[string]CacheEntry
    mu    sync.RWMutex
    ttl   time.Duration
}

type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, ok := c.store[key]
    if !ok || time.Now().After(entry.ExpiresAt) {
        return nil, false
    }

    return entry.Data, true
}

func (c *Cache) Set(key string, data interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.store[key] = CacheEntry{
        Data:      data,
        ExpiresAt: time.Now().Add(c.ttl),
    }
}
```

## Tool System

### Tool Lifecycle

```
1. Tool Registration
   ↓
   tools.Register("tool_name", &ToolImpl{})
   ↓
2. Tool Selection (by agent)
   ↓
   selectedTools := agent.SelectTools(query)
   ↓
3. Parameter Preparation
   ↓
   params := prepareParams(query, options)
   ↓
4. Tool Execution
   ↓
   result := tool.Execute(ctx, params)
   ↓
5. Result Validation
   ↓
   if !result.Success { retry or fail }
   ↓
6. Result Caching
   ↓
   cache.Set(toolKey, result)
   ↓
7. Result Aggregation
```

### Tool Communication

Tools communicate through structured results:

```go
type ToolResult struct {
    Success  bool                   // Did tool succeed?
    Data     interface{}            // Tool output data
    Sources  []models.Source        // Sources found
    Error    error                  // Error if failed
    Metadata map[string]interface{} // Additional info
}
```

### Tool Chaining

Tools can be chained for complex workflows:

```
WebSearch → Summarize → ExtractInfo → FactCheck
    ↓           ↓            ↓            ↓
  Articles   Summaries     Facts      Verified
```

## Report Generation

### Report Pipeline

```
Research Results
   ↓
Create Report Object
   ├─ Query
   ├─ Findings
   ├─ Sources
   └─ Methodology
   ↓
Format Selection
   ├─ Markdown
   ├─ JSON
   └─ PDF
   ↓
Citation Formatting
   ├─ APA
   ├─ MLA
   └─ Chicago
   ↓
Template Application
   ↓
Output Generation
```

### Markdown Generation

```go
func FormatMarkdownReport(report *Report) string {
    var sb strings.Builder

    // Header
    sb.WriteString(fmt.Sprintf("# Research Report: %s\n\n", report.Query))

    // Metadata
    sb.WriteString(fmt.Sprintf("**Generated:** %s\n", report.GeneratedAt))
    sb.WriteString(fmt.Sprintf("**Sources:** %d\n\n", len(report.Sources)))

    // Executive Summary
    sb.WriteString("## Executive Summary\n\n")
    sb.WriteString(report.Summary + "\n\n")

    // Key Findings
    sb.WriteString("## Key Findings\n\n")
    for i, finding := range report.Findings {
        sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, finding.Title))
        sb.WriteString(fmt.Sprintf("   - %s\n", finding.Description))
        // ... citations ...
    }

    // Sources
    sb.WriteString("## Sources\n\n")
    for i, source := range report.Sources {
        sb.WriteString(fmt.Sprintf("[%d] %s\n\n", i+1, source.Citation))
    }

    return sb.String()
}
```

### Citation Formatting

```go
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
    // Author, A. (Year). Title. Source.
    return fmt.Sprintf("%s. (%s). %s. %s",
        source.Author,
        source.PublishDate.Year(),
        source.Title,
        source.URL,
    )
}
```

## Error Handling

### Error Hierarchy

```
ApplicationError
   ├─ ConfigError
   │  ├─ InvalidConfigError
   │  └─ MissingConfigError
   ├─ ToolError
   │  ├─ ToolExecutionError
   │  ├─ ToolTimeoutError
   │  └─ ToolValidationError
   ├─ StorageError
   │  ├─ DatabaseError
   │  └─ FileError
   └─ ReportError
      ├─ GenerationError
      └─ ExportError
```

### Error Handling Pattern

```go
// Wrap errors with context
func (t *WebSearchTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    results, err := t.search(query)
    if err != nil {
        return nil, fmt.Errorf("web search failed: %w", err)
    }
    // ...
}

// Graceful degradation
results := make([]ToolResult, 0)
for _, tool := range tools {
    result, err := tool.Execute(ctx, params)
    if err != nil {
        log.Printf("Tool %s failed: %v", tool.Name(), err)
        continue  // Skip failed tool
    }
    results = append(results, result)
}

// If we have some results, continue
if len(results) == 0 {
    return fmt.Errorf("all tools failed")
}
```

### Retry Logic

```go
func executeWithRetry(fn func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }

        // Exponential backoff
        time.Sleep(time.Duration(1<<i) * time.Second)
    }
    return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}
```

## Performance Considerations

### Optimization Strategies

**1. Concurrent Tool Execution**
- Execute independent tools in parallel
- Use worker pools to limit concurrency
- Respect system resources

**2. Caching**
- Cache tool results (30-minute TTL)
- Cache document extractions
- Cache web search results

**3. Lazy Loading**
- Load documents only when needed
- Stream large files
- Paginate database queries

**4. Resource Limits**
- Limit max sources per tool
- Timeout long-running operations
- Cap memory usage for large documents

### Benchmarks

**Typical Performance (Medium Depth):**

| Operation | Time | Notes |
|-----------|------|-------|
| Web Search | 2-3s | Network dependent |
| Wikipedia | 1-2s | Usually cached |
| PDF Extract | 1-5s | Depends on size |
| Synthesis | 2-3s | CPU intensive |
| Total Research | 3-5 min | With 15 sources |

**Scaling:**
- Shallow: 1-2 min (5-10 sources)
- Medium: 3-5 min (10-15 sources)
- Deep: 5-10 min (15-25+ sources)

## Security

### Input Validation

```go
// Validate query
func validateQuery(query string) error {
    if len(query) == 0 {
        return errors.New("query cannot be empty")
    }
    if len(query) > 500 {
        return errors.New("query too long (max 500 chars)")
    }
    return nil
}

// Validate file paths
func validateFilePath(path string) error {
    // Check for path traversal
    cleanPath := filepath.Clean(path)
    if strings.Contains(cleanPath, "..") {
        return errors.New("invalid file path")
    }
    return nil
}
```

### Safe File Operations

```go
// Safe file reading
func safeReadFile(path string) ([]byte, error) {
    // Validate path
    if err := validateFilePath(path); err != nil {
        return nil, err
    }

    // Check file exists and is readable
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    // Check file size (limit to 100MB)
    if info.Size() > 100*1024*1024 {
        return nil, errors.New("file too large")
    }

    return os.ReadFile(path)
}
```

### Data Sanitization

```go
// Sanitize user input for database
func sanitizeInput(input string) string {
    // Remove control characters
    // Escape special characters
    // Limit length
    return sanitized
}
```

## Deployment

### Binary Distribution

```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o research-agent-linux
GOOS=darwin GOARCH=amd64 go build -o research-agent-macos
GOOS=windows GOARCH=amd64 go build -o research-agent.exe
```

### Docker Deployment

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o research-agent

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/research-agent /usr/local/bin/
ENTRYPOINT ["research-agent"]
```

### Configuration

```yaml
# Production config
agent:
  default_depth: medium
  max_sources: 20
  concurrent_tools: 5
  timeout: 600

storage:
  database_path: /var/lib/research-agent/research.db
  documents_dir: /var/lib/research-agent/documents
```

## Future Enhancements

### Planned Features

1. **Distributed Processing**
   - Distribute tool execution across machines
   - Shared cache layer
   - Queue-based job processing

2. **Advanced AI Integration**
   - LLM-powered query planning
   - Intelligent source selection
   - Better synthesis

3. **Web Interface**
   - Browser-based UI
   - Real-time collaboration
   - Report sharing

4. **Plugin System**
   - Custom tool plugins
   - Format plugins
   - Storage plugins

### Scalability Improvements

1. **Horizontal Scaling**
   - Multiple agent instances
   - Shared database
   - Load balancing

2. **Caching Layer**
   - Redis for distributed cache
   - Longer TTLs for stable content
   - Cache invalidation strategies

3. **Database Optimization**
   - PostgreSQL for production
   - Read replicas
   - Partitioning by date

## Summary

The Deep Research Agent is built with:

- **Modular Architecture**: Clear separation of concerns
- **Concurrent Design**: Parallel tool execution
- **Interface-Based**: Extensible tool system
- **Storage Layer**: SQLite with file storage
- **Multiple Formats**: Markdown, JSON, PDF output
- **Production-Ready**: Error handling, logging, configuration

For more details, see:
- [README.md](README.md) - Complete documentation
- [TOOLS_GUIDE.md](TOOLS_GUIDE.md) - Tool details
- [START_HERE.md](START_HERE.md) - Learning guide

---

**Architecture Documentation** - Understanding the Deep Research Agent's design
