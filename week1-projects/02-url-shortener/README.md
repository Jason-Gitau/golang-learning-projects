# URL Shortener - Week 1 Project 2

A command-line URL shortener that demonstrates core Go concepts from Week 1 of your learning path.

## Overview

Convert long URLs into short, memorable codes. Track visit counts and manage your shortened URLs with ease.

```
https://www.example.com/very/long/path/to/resource → http://short.url/abc123
```

## Features

✅ **Shorten URLs** - Create 6-character short codes for any URL
✅ **Retrieve URLs** - Look up the original URL from a short code
✅ **Track Visits** - Automatic visit counter for each shortened URL
✅ **Persistent Storage** - All data saved to JSON file
✅ **List All** - View all shortened URLs in a table format
✅ **Statistics** - Get detailed stats including compression ratio
✅ **Delete URLs** - Remove shortened URLs when no longer needed

## Quick Start

### Build the Project

```bash
cd week1-projects/02-url-shortener
go build -o url-shortener.exe
```

### Create a Short URL

```bash
./url-shortener.exe shorten https://www.google.com
```

Output:
```
✅ URL Shortened Successfully!

Original URL: https://www.google.com
Short Code:  abc123
Short URL:   http://short.url/abc123
```

### Retrieve Original URL

```bash
./url-shortener.exe get abc123
```

### List All Shortened URLs

```bash
./url-shortener.exe list
```

### Get Statistics

```bash
./url-shortener.exe stats abc123
```

### Delete a URL

```bash
./url-shortener.exe delete abc123
```

## Project Structure

```
02-url-shortener/
├── main.go              # CLI entry point (170 lines)
├── shortener/
│   └── shortener.go     # Core logic (280 lines)
├── storage/
│   └── storage.go       # JSON persistence (180 lines)
├── go.mod              # Go module definition
└── urls.json           # Data storage (created on first use)
```

## Code Statistics

- **Total Lines:** 630 lines of code
- **Packages:** 2 (shortener, storage)
- **Functions:** 20+
- **Error Cases:** 15+

## Learning Concepts

This project demonstrates all Week 1 concepts:

### 1. **Structs & Methods** (`shortener/shortener.go`)
```go
type URLShortener struct {
    mappings map[string]*URLMapping
    codeLength int
    charset string
}

func (us *URLShortener) ShortenURL(url string) (string, error)
```

### 2. **Maps** (`shortener/shortener.go`)
```go
mappings := make(map[string]*URLMapping)
mapping, exists := mappings[shortCode]
for code, mapping := range mappings { ... }
```

### 3. **Slices** (`storage/storage.go`)
```go
var mappings []*shortener.URLMapping
mappings = append(mappings, mapping)
for _, mapping := range mappings { ... }
```

### 4. **Pointers** (throughout)
```go
func NewURLShortener() *URLShortener
us *URLShortener
mapping *shortener.URLMapping
```

### 5. **Error Handling** (all files)
```go
if err := us.validateURL(url); err != nil {
    return "", err
}
return fmt.Errorf("failed to shorten: %w", err)
```

### 6. **JSON Marshaling/Unmarshaling** (`storage/storage.go`)
```go
data, _ := json.MarshalIndent(mappings, "", "  ")
json.Unmarshal(data, &mappings)
```

### 7. **File I/O** (`storage/storage.go`)
```go
os.WriteFile(filepath, data, 0644)
data, _ := os.ReadFile(filepath)
os.Stat(filepath)
```

## Commands Reference

### Shorten a URL
```bash
./url-shortener.exe shorten <URL>

# Example
./url-shortener.exe shorten https://www.github.com/golang/go
```

Creates a new short code and saves it to `urls.json`.

### Get Original URL
```bash
./url-shortener.exe get <SHORT_CODE>

# Example
./url-shortener.exe get abc123
```

Retrieves the original URL and increments visit counter.

### List All URLs
```bash
./url-shortener.exe list
```

Displays all shortened URLs in a table:
```
CODE     ORIGINAL URL                          VISITS  CREATED
abc123   https://www.github.com/golang/go      5       2025-10-30 15:30:45
def456   https://www.google.com                2       2025-10-30 15:35:20
```

### Get Statistics
```bash
./url-shortener.exe stats <SHORT_CODE>

# Example
./url-shortener.exe stats abc123
```

Shows:
- Original URL
- Creation date
- Visit count
- Original URL length
- Code length
- Compression ratio

### Delete a URL
```bash
./url-shortener.exe delete <SHORT_CODE>

# Example
./url-shortener.exe delete abc123
```

### Help
```bash
./url-shortener.exe help
```

## How It Works

### Shortening Flow

1. **Validate** - Check that URL is valid (starts with http://, not empty, etc.)
2. **Check Existing** - Look for the URL in existing mappings
3. **Generate Code** - Create a random 6-character alphanumeric code
4. **Create Mapping** - Store the short code → original URL mapping
5. **Save** - Write to `urls.json` for persistence
6. **Return** - Show the user the short code

### Retrieval Flow

1. **Lookup** - Find the mapping in the in-memory map
2. **Increment Visits** - Add 1 to the visit counter
3. **Update Storage** - Save updated visits to file
4. **Return** - Show the original URL and visit count

### Data Persistence

All data is stored in JSON format:

```json
[
  {
    "short_code": "abc123",
    "original_url": "https://www.example.com",
    "created_at": "2025-10-30 15:30:45",
    "visits": 5
  },
  {
    "short_code": "def456",
    "original_url": "https://www.google.com",
    "created_at": "2025-10-30 15:35:20",
    "visits": 2
  }
]
```

When the app starts, it loads all mappings from `urls.json` into memory.

## Key Functions

### URLShortener Package

```go
NewURLShortener(codeLength int) *URLShortener
    // Create a new shortener instance

ShortenURL(originalURL string) (string, error)
    // Shorten a URL and return the short code

GetURL(shortCode string) (string, error)
    // Get original URL and increment visits

GetMapping(shortCode string) (*URLMapping, error)
    // Get full mapping details

ListAllMappings() []*URLMapping
    // Get all mappings as a slice

GetStats(shortCode string) (map[string]interface{}, error)
    // Get detailed statistics

DeleteURL(shortCode string) error
    // Remove a shortened URL
```

### Storage Package

```go
NewStorage(filePath string) *Storage
    // Create a new storage instance

SaveMappings(mappings []*URLMapping) error
    // Write all mappings to JSON

LoadMappings() ([]*URLMapping, error)
    // Read all mappings from JSON

AppendMapping(mapping *URLMapping) error
    // Add a new mapping to the file

RemoveMapping(shortCode string) error
    // Delete a mapping from the file

UpdateMapping(mapping *URLMapping) error
    // Update an existing mapping
```

## Code Examples

### Reading the Code

Start with this order:

1. **main.go** - Understand the CLI flow and command routing
2. **shortener/shortener.go** - Core logic with detailed comments
3. **storage/storage.go** - JSON persistence patterns

### Key Patterns

**Using Pointers:**
```go
us := shortener.NewURLShortener(6)  // Returns *URLShortener
mapping, err := us.GetMapping(code) // Returns *URLMapping
```

**Error Handling:**
```go
if err := us.ShortenURL(url); err != nil {
    log.Fatalf("Failed: %v", err)
}
```

**Map Operations:**
```go
// Check if key exists
mapping, exists := us.mappings[code]
if !exists {
    // Handle not found
}

// Iterate
for code, mapping := range us.mappings {
    // Process each mapping
}
```

**Slice Operations:**
```go
// Create empty slice
var mappings []*shortener.URLMapping

// Append
mappings = append(mappings, newMapping)

// Filter (create new slice without one element)
var filtered []*shortener.URLMapping
for _, m := range mappings {
    if m.ShortCode != deleteCode {
        filtered = append(filtered, m)
    }
}
```

## File Format

### urls.json

Auto-created on first run. Contains an array of URL mappings:

```json
[
  {
    "short_code": "abc123",
    "original_url": "https://example.com/long/path",
    "created_at": "2025-10-30 15:30:45",
    "visits": 3
  }
]
```

## Error Handling

The tool handles various error cases:

✅ Empty URL
✅ Invalid URL format
✅ URL too short
✅ URL with spaces
✅ Duplicate short codes
✅ Non-existent short codes
✅ File I/O errors
✅ JSON parsing errors

## Testing

### Manual Testing Workflow

```bash
# 1. Shorten several URLs
./url-shortener.exe shorten https://www.github.com
./url-shortener.exe shorten https://www.golang.org
./url-shortener.exe shorten https://www.example.com

# 2. List all URLs
./url-shortener.exe list

# 3. Access URLs (increases visit count)
./url-shortener.exe get <code1>
./url-shortener.exe get <code1>
./url-shortener.exe get <code2>

# 4. Check stats
./url-shortener.exe stats <code1>

# 5. Delete one
./url-shortener.exe delete <code1>

# 6. Verify deletion
./url-shortener.exe list

# 7. Check urls.json file
cat urls.json
```

## Improvements & Extensions

After completing Week 1, you could enhance this project with:

**Week 2 Ideas:**
- REST API instead of CLI
- Web interface to create/manage short URLs
- Database instead of JSON file

**Week 3 Ideas:**
- Concurrent URL fetching (check if URLs are valid)
- Rate limiting for URL creation
- Background analytics processing

**Week 4 Ideas:**
- User authentication
- URL expiration/TTL
- Custom short codes
- API key management

## Troubleshooting

**Issue: "short code not found"**
- Make sure you're using the correct code
- Check `urls.json` exists and contains data
- Rebuild the application with `go build`

**Issue: "URL already exists"**
- The exact URL is already shortened
- Use the existing short code instead
- Or try a different URL

**Issue: "failed to save mapping"**
- Check write permissions in the directory
- Make sure urls.json isn't locked by another program
- Try deleting urls.json and starting fresh

## Summary

This URL Shortener project demonstrates:

✅ Structs and methods for organizing code
✅ Maps for fast lookups (O(1) access)
✅ Slices for variable-length data
✅ Pointers for efficiency
✅ Error handling patterns
✅ JSON serialization
✅ File I/O operations
✅ Professional CLI design

**Total Learning Value:** High! This project combines all Week 1 concepts in a real, useful tool.

---

**Part of Week 1 Projects:** Completed alongside the File Organizer CLI

**Next:** Move to Week 2 for Web Development and REST APIs!
