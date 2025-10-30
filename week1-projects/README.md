# Week 1 Projects - Go Learning Path

Two complete, production-ready CLI projects demonstrating all core Go concepts from Week 1.

## 📂 Project Structure

```
week1-projects/
├── 01-file-organizer/          ✅ Complete
│   ├── main.go                 (CLI entry point)
│   ├── organizer/              (Core logic)
│   ├── utils/                  (File utilities)
│   ├── README.md               (Full documentation)
│   └── QUICK_START.md          (Quick guide)
│
└── 02-url-shortener/           ✅ Complete
    ├── main.go                 (CLI entry point)
    ├── shortener/              (Core logic)
    ├── storage/                (JSON persistence)
    ├── README.md               (Full documentation)
    └── QUICK_START.md          (Quick guide)
```

## 🎯 Projects Overview

### Project 1: File Organizer CLI

**What it does:** Automatically organize files in any folder by their extension.

**Quick Example:**
```bash
cd 01-file-organizer
./file-organizer.exe -source ~/Downloads -dry-run
./file-organizer.exe -source ~/Downloads
```

**What You Learn:**
- ✅ Structs and methods
- ✅ File I/O operations
- ✅ Maps for organizing data
- ✅ Error handling
- ✅ Command-line flags
- ✅ JSON configuration
- ✅ Directory traversal

**Stats:** 754 lines, 3 packages, 15+ functions

---

### Project 2: URL Shortener

**What it does:** Create short codes for long URLs, track visits, manage your URLs.

**Quick Example:**
```bash
cd 02-url-shortener
./url-shortener.exe shorten https://www.github.com/golang/go
./url-shortener.exe list
./url-shortener.exe stats yz7lIe
```

**What You Learn:**
- ✅ Structs and methods
- ✅ Maps for fast lookups
- ✅ Slices for dynamic data
- ✅ Pointers for efficiency
- ✅ Error handling
- ✅ JSON persistence
- ✅ File I/O operations
- ✅ Random code generation
- ✅ Visit tracking

**Stats:** 630 lines, 2 packages, 20+ functions

---

## 🚀 Quick Start - Both Projects

### Build Both Projects

```bash
# Project 1
cd 01-file-organizer
go build -o file-organizer.exe

# Project 2
cd ../02-url-shortener
go build -o url-shortener.exe
```

### File Organizer - Examples

```bash
cd 01-file-organizer

# Preview (dry-run)
./file-organizer.exe -source C:\Users\Jason\Downloads -dry-run

# Actually organize
./file-organizer.exe -source C:\Users\Jason\Downloads

# Create custom config
./file-organizer.exe -create-config -config myconfig.json

# Use custom config
./file-organizer.exe -source C:\Users\Jason\Downloads -config myconfig.json

# List all mappings
./file-organizer.exe -list
```

### URL Shortener - Examples

```bash
cd 02-url-shortener

# Create short URLs
./url-shortener.exe shorten https://github.com/golang/go
./url-shortener.exe shorten https://golang.org

# List all
./url-shortener.exe list

# Get original URL (and increment visits)
./url-shortener.exe get abc123

# View statistics
./url-shortener.exe stats abc123

# Delete a URL
./url-shortener.exe delete abc123
```

## 📚 Learning Concepts

### Week 1 Goals - All Covered! ✅

| Concept | Project 1 | Project 2 |
|---------|-----------|-----------|
| **Structs** | ✅ | ✅ |
| **Methods** | ✅ | ✅ |
| **Pointers** | ✅ | ✅ |
| **Maps** | ✅ | ✅ |
| **Slices** | ✅ | ✅ |
| **Error Handling** | ✅ | ✅ |
| **File I/O** | ✅ | ✅ |
| **JSON** | ✅ | ✅ |
| **CLI Flags** | ✅ | - |
| **Visit Tracking** | - | ✅ |
| **Random Generation** | - | ✅ |

### Code Organization Patterns

Both projects demonstrate:
- ✅ Multiple packages (separation of concerns)
- ✅ Proper error handling throughout
- ✅ Function comments explaining concepts
- ✅ Type safety with structs
- ✅ Defensive programming (validation)
- ✅ Data persistence
- ✅ Professional CLI design

## 📖 Reading Guide

Start with this order:

**Project 1 (File Organizer):**
1. `01-file-organizer/QUICK_START.md` - Get it running
2. `01-file-organizer/main.go` - Understand CLI flow
3. `01-file-organizer/organizer/organizer.go` - Core logic
4. `01-file-organizer/organizer/config.go` - JSON handling
5. `01-file-organizer/utils/fileutils.go` - File operations

**Project 2 (URL Shortener):**
1. `02-url-shortener/QUICK_START.md` - Get it running
2. `02-url-shortener/main.go` - Understand CLI flow
3. `02-url-shortener/shortener/shortener.go` - Core logic
4. `02-url-shortener/storage/storage.go` - Persistence

## 🎓 Key Patterns Explained

### Using Pointers

```go
// Project 1
organizer := organizer.NewFileOrganizer(sourceDir, outputDir)

// Project 2
shortener := shortener.NewURLShortener(6)
```

### Error Handling

```go
// Validate input
if err := us.validateURL(url); err != nil {
    return "", err
}

// Wrap errors for context
return "", fmt.Errorf("failed to shorten: %w", err)
```

### Maps for Data Storage

```go
// Project 1: Extension → Folder mapping
ExtensionMap := map[string]string{
    ".txt": "Documents",
    ".jpg": "Images",
}

// Project 2: Short code → URL mapping
Mappings := map[string]*URLMapping{
    "abc123": &URLMapping{...},
}
```

### Slices for Collections

```go
// Project 1: List of files
var files []string
files = append(files, filename)

// Project 2: List of all mappings
var mappings []*shortener.URLMapping
for _, mapping := range us.Mappings {
    mappings = append(mappings, mapping)
}
```

### JSON Persistence

```go
// Save
data, _ := json.MarshalIndent(data, "", "  ")
os.WriteFile("data.json", data, 0644)

// Load
data, _ := os.ReadFile("data.json")
json.Unmarshal(data, &mappings)
```

## 💻 Code Statistics

| Metric | Project 1 | Project 2 | Total |
|--------|-----------|-----------|-------|
| Lines of Code | 754 | 630 | 1,384 |
| Packages | 3 | 2 | 5 |
| Functions | 15+ | 20+ | 35+ |
| Error Cases | 10+ | 15+ | 25+ |

## 🧪 Testing Both Projects

### File Organizer Testing

```bash
cd 01-file-organizer
mkdir test_files
touch test_files/{doc.txt,photo.jpg,video.mp4,song.mp3}

# Preview
./file-organizer.exe -source test_files -dry-run

# Run
./file-organizer.exe -source test_files

# Verify
ls -R test_files
```

### URL Shortener Testing

```bash
cd 02-url-shortener

# Create URLs
./url-shortener.exe shorten https://github.com
./url-shortener.exe shorten https://golang.org
./url-shortener.exe shorten https://google.com

# Use them
./url-shortener.exe get <code1>
./url-shortener.exe get <code1>  # visits increases
./url-shortener.exe stats <code1>

# Clean up
./url-shortener.exe delete <code1>
./url-shortener.exe list
```

## 🔗 Next Steps

After mastering Week 1:

1. **Extend the projects** - Add features, improve error handling
2. **Study other Go projects** - Review code on GitHub
3. **Move to Week 2** - REST APIs and web development
4. **Deploy these tools** - Build for different platforms, share them

## 📝 Files to Study

### Most Important (Read First)
- `01-file-organizer/organizer/organizer.go`
- `02-url-shortener/shortener/shortener.go`

### Supporting Files
- Both `main.go` files - CLI design patterns
- Both `QUICK_START.md` - Learning guides
- Both `README.md` - Complete references

### Data Files
- `01-file-organizer/config.json` - Configuration example
- `02-url-shortener/urls.json` - Data storage example

## 🎯 Success Criteria

By the end of Week 1, you can:

✅ Understand and write structs
✅ Create methods on types
✅ Use pointers effectively
✅ Work with maps and slices
✅ Handle errors properly
✅ Read and write files
✅ Parse JSON
✅ Build CLI applications
✅ Organize code into packages
✅ Write well-commented code

## 🌟 Summary

Two complete, working projects that:
- ✅ Solve real problems
- ✅ Demonstrate all Week 1 concepts
- ✅ Use professional coding practices
- ✅ Include comprehensive documentation
- ✅ Are ready for your portfolio

**Total value:** 1,384 lines of production-ready Go code with full documentation!

---

## 📚 Additional Resources

- [Project 1 Docs](./01-file-organizer/README.md)
- [Project 2 Docs](./02-url-shortener/README.md)
- [Project 1 Quick Start](./01-file-organizer/QUICK_START.md)
- [Project 2 Quick Start](./02-url-shortener/QUICK_START.md)

## 🚀 Ready?

Pick a project and start exploring the code!

```bash
# Jump in!
cd 01-file-organizer
# OR
cd 02-url-shortener
```

Happy learning! 🎓
