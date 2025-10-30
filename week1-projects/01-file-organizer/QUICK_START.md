# Quick Start Guide - File Organizer

Get started in 2 minutes!

## 1️⃣ Run the Application

```bash
# Navigate to project
cd C:\Users\Jason\Golang_lessons\week1-file-organizer

# See what files would be organized (no changes made)
./file-organizer.exe -source test_files -dry-run

# Actually organize files
./file-organizer.exe -source test_files

# Check the results
ls test_files/
```

## 2️⃣ Try Different Commands

```bash
# Show all extension mappings
./file-organizer.exe -list

# Create a custom config file
./file-organizer.exe -create-config -config my-config.json

# Organize using your custom config
./file-organizer.exe -source test_files2 -config my-config.json

# Get help
./file-organizer.exe -help
```

## 3️⃣ Understand the Code

### Start with the simplest file first:

**`organizer/config.go`** - How to work with JSON
```go
// Load a config file
config, err := LoadConfig("config.json")

// Save a config file
SaveConfig("config.json", config)

// The key concepts:
// - json.Unmarshal (JSON → Go struct)
// - json.Marshal (Go struct → JSON)
```

**`utils/fileutils.go`** - File operations
```go
// Check if directory exists
if DirectoryExists(path) { ... }

// Get all files in a directory
files, _ := ListFiles(dirPath)

// Count files by extension
counts, _ := CountFilesByExtension(dirPath)
```

**`organizer/organizer.go`** - Core logic with STRUCTS and METHODS
```go
// Create a new organizer
organizer := NewFileOrganizer(sourceDir, outputDir)

// Do the organizing
err := organizer.Organize()
```

**`main.go`** - How it all comes together
```go
// Parse command line flags
sourceDir := flag.String("source", "", "...")
flag.Parse()

// Use the organizer
fo := organizer.NewFileOrganizer(sourceDir, outputDir)
fo.Organize()
```

## 4️⃣ Key Go Concepts You're Using

| Concept | Where | Example |
|---------|-------|---------|
| **Structs** | `organizer/organizer.go` | `type FileOrganizer struct { ... }` |
| **Methods** | `organizer/organizer.go` | `func (fo *FileOrganizer) Organize() { ... }` |
| **Pointers** | All files | `*FileOrganizer`, `*Config` |
| **Maps** | `organizer/organizer.go` | `map[string]string` for extensions |
| **Slices** | `utils/fileutils.go` | `[]string`, `[]FileInfo` |
| **Error Handling** | All files | `if err != nil { return err }` |
| **File I/O** | `utils/fileutils.go` | `os.ReadDir`, `os.Stat` |
| **JSON** | `organizer/config.go` | `json.Unmarshal`, `json.Marshal` |
| **CLI Flags** | `main.go` | `flag.String`, `flag.Bool` |

## 5️⃣ Try Modifying the Code

### Add a new extension mapping

**In `organizer/organizer.go`:**
```go
// In NewFileOrganizer(), add:
".iso": "Media",
".mkv": "Videos",
```

Then rebuild:
```bash
go build -o file-organizer.exe
```

### Change the default folder name

**In `organizer/organizer.go`:**
```go
// Change this:
folderName = "Other"

// To this:
folderName = "Miscellaneous"
```

### Add file size information

**In `utils/fileutils.go`:**
```go
// Use this in CountFilesByExtension:
info, _ := entry.Info()
size := info.Size()
```

## 6️⃣ Common Commands Cheat Sheet

```bash
# Build the project
go build -o file-organizer.exe

# Run tests (when you add them)
go test ./...

# Show dependencies
go list -m all

# Check for unused imports
go fmt ./...

# Get help
./file-organizer.exe -help

# List all mappings
./file-organizer.exe -list

# Preview what would happen
./file-organizer.exe -source ./Downloads -dry-run

# Actually organize
./file-organizer.exe -source ./Downloads

# Use custom config
./file-organizer.exe -source ./Downloads -config config.json

# Organize into different directory
./file-organizer.exe -source ./Downloads -output ./Organized
```

## 7️⃣ Next: Learn More

When you're ready to go deeper:

1. **Read the Comments** - Every significant code section has explanations
2. **Try the Examples** - Modify code and rebuild to see what changes
3. **Build on It** - Add features like undo, redo, or logging
4. **Test It** - Add test files and test different scenarios
5. **Deploy It** - Build for different platforms

## 8️⃣ Troubleshooting

**Problem:** "source directory does not exist"
```bash
# Make sure the directory exists
mkdir -p your-directory
./file-organizer.exe -source your-directory
```

**Problem:** "failed to move file"
```bash
# Files with same names might already exist there
./file-organizer.exe -source your-directory -dry-run
```

**Problem:** Can't rebuild
```bash
# Make sure you're in the right directory
cd week1-file-organizer
go build -o file-organizer.exe
```

## 9️⃣ Project Layout

```
week1-file-organizer/
├── main.go              ← Entry point (start here to understand flow)
├── organizer/           ← Core business logic
│   ├── organizer.go     ← Main struct and Organize() method
│   └── config.go        ← JSON configuration handling
├── utils/               ← Reusable utilities
│   └── fileutils.go     ← File system operations
└── README.md            ← Full documentation
```

## 🔟 Learning Path

**You've completed:**
- ✅ Week 1: Foundations (Structs, Error Handling, File I/O, JSON)

**Next:**
- ⬜ Week 2: Web Development (REST APIs, HTTP, Databases)
- ⬜ Week 3: Concurrency (Goroutines, Channels)
- ⬜ Week 4: Real-World Apps (Testing, Deployment, Production practices)

---

## 🚀 You're Ready!

You now understand:
- How to write Go code
- How to structure projects
- How to handle errors properly
- How to work with files and JSON
- How to build CLI tools

**Time to build something awesome!** 💪

For more details, see `README.md` in this directory.
