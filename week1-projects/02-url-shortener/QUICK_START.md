# Quick Start - URL Shortener

Get running in 3 minutes!

## 1️⃣ Build

```bash
cd week1-projects/02-url-shortener
go build -o url-shortener.exe
```

## 2️⃣ Create Short URLs

```bash
# Shorten a URL
./url-shortener.exe shorten https://www.github.com/golang/go

# You'll get a short code like: abc123
```

## 3️⃣ Try These Commands

```bash
# Get the original URL back
./url-shortener.exe get abc123

# List all your shortened URLs
./url-shortener.exe list

# See how many times it was accessed
./url-shortener.exe stats abc123

# Delete a URL you don't need
./url-shortener.exe delete abc123

# Get help
./url-shortener.exe help
```

## What You're Learning

| Concept | See It Here |
|---------|-------------|
| **Structs & Methods** | shortener/shortener.go (lines 1-30) |
| **Maps** | shortener/shortener.go (lines 40-60) |
| **Slices** | storage/storage.go (lines 50-80) |
| **Pointers** | All files (function returns) |
| **Error Handling** | All files (err != nil checks) |
| **JSON** | storage/storage.go (Marshal/Unmarshal) |
| **File I/O** | storage/storage.go (ReadFile/WriteFile) |

## Code Reading Order

1. **main.go** - See the CLI commands
2. **shortener/shortener.go** - Understand the core logic
3. **storage/storage.go** - See how data is saved

## Data

All URLs are saved in `urls.json` (auto-created on first run).

Check it out:
```bash
cat urls.json
```

## Try These Scenarios

### Scenario 1: Create & Retrieve

```bash
# Create
./url-shortener.exe shorten https://google.com

# Retrieve (note visits increases)
./url-shortener.exe get <code>
./url-shortener.exe get <code>  # visits = 2 now

# Check stats
./url-shortener.exe stats <code>
```

### Scenario 2: Multiple URLs

```bash
./url-shortener.exe shorten https://github.com
./url-shortener.exe shorten https://golang.org
./url-shortener.exe shorten https://stackoverflow.com

# List all
./url-shortener.exe list
```

### Scenario 3: Management

```bash
# Create a URL
./url-shortener.exe shorten https://example.com

# Use it a few times
./url-shortener.exe get <code>

# Check stats
./url-shortener.exe stats <code>

# Delete it
./url-shortener.exe delete <code>

# Verify deletion
./url-shortener.exe list
```

## How It Works

1. **Shorten** - Creates random 6-letter code
2. **Store** - Saves in urls.json
3. **Retrieve** - Looks up code, increments visits
4. **Manage** - List, delete, get stats

## Key Go Concepts

**Structs:**
```go
type URLMapping struct {
    ShortCode   string
    OriginalURL string
    CreatedAt   string
    Visits      int
}
```

**Methods:**
```go
func (us *URLShortener) ShortenURL(url string) (string, error)
```

**Maps:**
```go
mappings := make(map[string]*URLMapping)
mapping, exists := mappings[code]
```

**JSON:**
```go
json.MarshalIndent(data, "", "  ")
json.Unmarshal(bytes, &data)
```

## Files

- `main.go` - CLI interface
- `shortener/shortener.go` - Core logic
- `storage/storage.go` - JSON persistence
- `urls.json` - Your data (created automatically)

## Common Commands Cheat Sheet

```bash
# View all shortened URLs
./url-shortener.exe list

# Get original URL (and count visit)
./url-shortener.exe get <code>

# See detailed stats
./url-shortener.exe stats <code>

# Delete a shortened URL
./url-shortener.exe delete <code>

# Get help
./url-shortener.exe help
```

## Next Steps

1. **Understand the Code** - Read shortener.go with comments
2. **Experiment** - Try different URLs and commands
3. **Modify** - Change short code length from 6 to 8 characters
4. **Extend** - Add features like URL validation or expiration

## Troubleshooting

**"Command not found"**
- Make sure you're in the right directory
- Run `go build -o url-shortener.exe`

**"Short code not found"**
- You might have the wrong code
- Try `list` to see all codes

**"URL doesn't exist"**
- The URL you're shortening is invalid
- Must start with http:// or https://

---

Ready to dive into the code? Start with `main.go`! 🚀
