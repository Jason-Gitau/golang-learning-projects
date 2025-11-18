# Quick Start Guide

Get up and running with Deep Research Agent in 5 minutes.

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/deep-research-agent.git
cd deep-research-agent

# Install dependencies
go mod download

# Build the agent
go build -o research-agent

# Verify installation
./research-agent --version
```

## Your First Research

### Example 1: Basic Research

```bash
./research-agent research "Go programming best practices"
```

This will:
1. Search the web for relevant sources
2. Query Wikipedia for background
3. Extract and synthesize information
4. Display a research report with citations

### Example 2: Save to File

```bash
./research-agent research "AI ethics" --output ai-ethics-report.md
```

The report will be saved as `ai-ethics-report.md`.

### Example 3: Analyze Documents

```bash
./research-agent research "Summary of paper" --pdf research_paper.pdf --no-web
```

This analyzes the PDF without using web sources.

### Example 4: Deep Research

```bash
./research-agent research "quantum computing" --depth deep --max-sources 25
```

This conducts comprehensive research with up to 25 sources.

## Common Commands

### Research Commands

```bash
# Basic research
./research-agent research "topic"

# With specific depth
./research-agent research "topic" --depth [shallow|medium|deep]

# Analyze PDFs
./research-agent research "summary" --pdf file.pdf

# Multiple documents
./research-agent research "analysis" --pdf file1.pdf --pdf file2.pdf --docx notes.docx

# Custom citation style
./research-agent research "topic" --cite-style [APA|MLA|Chicago]
```

### Session Management

```bash
# List all research sessions
./research-agent session list

# View session details
./research-agent session show abc123

# Resume previous session
./research-agent session resume abc123
```

### Export Reports

```bash
# Export as markdown
./research-agent export latest --format markdown --output report.md

# Export as PDF
./research-agent export latest --format pdf --output report.pdf

# Export as JSON
./research-agent export latest --format json --output data.json
```

### Interactive Mode

```bash
./research-agent interactive
```

Starts a guided interactive session where you can:
- Define research queries
- Choose sources and depth
- View results in real-time
- Export reports

## Configuration

Create `~/.research-agent.yaml`:

```yaml
agent:
  default_depth: medium
  max_sources: 15
  concurrent_tools: 3

citation:
  default_style: APA

output:
  default_format: markdown
```

## 5 Essential Examples

### 1. Quick Web Research

```bash
./research-agent research "Golang concurrency patterns" --output go-concurrency.md
```

**Use case:** Quick overview of a technical topic

### 2. Document Analysis

```bash
./research-agent research "Key findings from papers" \
  --pdf paper1.pdf \
  --pdf paper2.pdf \
  --no-web \
  --output findings.md
```

**Use case:** Analyze local documents without web search

### 3. Comprehensive Research

```bash
./research-agent research "Climate change impact on agriculture" \
  --depth deep \
  --max-sources 30 \
  --cite-style APA \
  --output climate-agriculture.md
```

**Use case:** In-depth research with many sources and academic citations

### 4. Mixed Sources

```bash
./research-agent research "AI in healthcare overview" \
  --pdf ai-healthcare-paper.pdf \
  --docx meeting-notes.docx \
  --use-web \
  --depth medium \
  --output ai-healthcare-report.md
```

**Use case:** Combine documents with web research

### 5. Interactive Session

```bash
./research-agent interactive
```

Then follow the prompts:
```
What would you like to research?
> Blockchain technology

Should I include web sources? (y/n) y

How deep should I research? (1-3)
> 2

Starting research...
```

**Use case:** Guided research for beginners

## Understanding Research Depth

| Depth   | Sources | Time    | Use Case                        |
|---------|---------|---------|--------------------------------|
| Shallow | 5-10    | 1-2 min | Quick overview                 |
| Medium  | 10-15   | 3-5 min | Balanced research              |
| Deep    | 15-25+  | 5-10 min| Comprehensive investigation    |

## Understanding Output

### Console Output

When you run a research command, you'll see:

```
Starting research: Go programming best practices
Depth: medium | Max Sources: 15 | Citation: APA
Web search: Enabled

🔍 Researching: Go programming best practices

Progress ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 45%

⚡ Analyzing sources (3/10)
✓ Web search completed (10 results)
✓ Wikipedia search completed
⚡ Processing information...

Sources found: 15
Time elapsed: 2m 34s

✓ Research completed!
  Sources found: 15
  Time elapsed: 2m 45s
```

### Markdown Report Structure

Generated reports include:

```markdown
# Research Report: [Your Query]

## Executive Summary
Brief overview of findings

## Key Findings
1. Finding with citations
2. Another finding

## Detailed Analysis
In-depth analysis with inline citations

## Sources
[1] Properly formatted citation
[2] Another citation

## Research Methodology
Details about the research process
```

## Next Steps

1. **Try Different Depths**: Experiment with `--depth shallow/medium/deep`
2. **Analyze Your Documents**: Use `--pdf` or `--docx` flags
3. **Manage Sessions**: Explore `session list` and `session show`
4. **Export Formats**: Try exporting as PDF or JSON
5. **Interactive Mode**: Run `./research-agent interactive`

## Common Issues

### Issue: "command not found"

Make sure the binary is executable:
```bash
chmod +x research-agent
./research-agent research "test"
```

Or move it to your PATH:
```bash
sudo mv research-agent /usr/local/bin/
research-agent research "test"
```

### Issue: PDF export fails

Install pandoc or pdflatex:
```bash
# Ubuntu/Debian
sudo apt-get install pandoc

# macOS
brew install pandoc
```

### Issue: Config file not found

Create the config directory:
```bash
mkdir -p ~/.research-agent
cp .research-agent.yaml ~/.research-agent.yaml
```

## Help and Documentation

### Get Help

```bash
# General help
./research-agent --help

# Command-specific help
./research-agent research --help
./research-agent session --help
./research-agent export --help
```

### Documentation Files

- **README.md** - Complete project documentation
- **START_HERE.md** - Comprehensive learning guide (start here!)
- **TOOLS_GUIDE.md** - Detailed tool documentation
- **ARCHITECTURE.md** - System design and architecture
- **examples/** - Example bash scripts

### View Statistics

```bash
./research-agent stats
```

Shows usage statistics, source counts, and performance metrics.

## Tips for Better Research

1. **Be Specific**: More specific queries yield better results
   - Good: "Go concurrency patterns for web servers"
   - Okay: "Go concurrency"

2. **Choose Appropriate Depth**:
   - Shallow: Quick facts or overview
   - Medium: Balanced investigation
   - Deep: Comprehensive academic research

3. **Use Citation Styles**:
   - APA: Scientific papers
   - MLA: Humanities papers
   - Chicago: Historical research

4. **Combine Sources**:
   - Use `--pdf` for academic papers
   - Use `--use-web` for current information
   - Use both for comprehensive research

5. **Save Sessions**:
   - Use `session list` to track your research
   - Use `session resume` to continue later
   - Use `export` to save in multiple formats

## Example Workflow

Here's a complete research workflow:

```bash
# 1. Conduct initial research
./research-agent research "Machine learning in finance" \
  --depth medium \
  --output ml-finance.md

# 2. View all sessions
./research-agent session list

# 3. Show session details (use ID from list)
./research-agent session show abc123

# 4. Add more analysis to documents
./research-agent document add finance-report.pdf --index

# 5. Search indexed documents
./research-agent document search "neural networks"

# 6. Export in multiple formats
./research-agent export abc123 --format pdf --output report.pdf
./research-agent export abc123 --format json --output data.json

# 7. View statistics
./research-agent stats
```

## You're Ready!

You now know how to:
- ✓ Conduct basic and advanced research
- ✓ Analyze documents
- ✓ Manage sessions
- ✓ Export reports
- ✓ Use interactive mode

For more details, see:
- [START_HERE.md](START_HERE.md) - Complete learning guide
- [TOOLS_GUIDE.md](TOOLS_GUIDE.md) - Tool documentation
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [README.md](README.md) - Full documentation

Happy researching!
