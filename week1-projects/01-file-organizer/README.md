# File Organizer - Week 1 Go Learning Project

A production-ready CLI tool that organizes files into directories by extension. Built to demonstrate core Go concepts from Week 1 of your learning path.

## What You'll Learn

This project teaches all Week 1 concepts from your Go learning roadmap:

### Core Go Concepts Covered

1. **Structs & Methods** (`organizer/organizer.go`)
   - Define custom data types with `struct`
   - Attach methods to structs with receivers
   - Difference between struct and pointer receivers

2. **Functions & Error Handling** (all files)
   - Error handling with return values
   - Error wrapping with `fmt.Errorf`
   - Multiple return values

3. **Slices & Maps** (`organizer/organizer.go`, `utils/fileutils.go`)
   - Slice operations with `append`
   - Map initialization and lookups
   - Comma-ok idiom for map access
   - Iterating with `for...range`

4. **Pointers** (all files)
   - Creating pointers with `&`
   - Dereferencing pointers
   - Method receivers as pointers vs values
   - Why pointers matter for efficiency

5. **File I/O** (`utils/fileutils.go`, `organizer/config.go`)
   - Reading directories with `os.ReadDir`
   - Getting file information with `os.Stat`
   - Working with file paths using `filepath`
   - Creating directories with `os.MkdirAll`

6. **JSON Marshaling/Unmarshaling** (`organizer/config.go`)
   - Converting Go structs to JSON
   - Converting JSON to Go structs
   - JSON struct tags for field mapping
   - Pretty-printing JSON with `MarshalIndent`

7. **Command-Line Arguments** (`main.go`)
   - Parsing flags with the `flag` package
   - Required vs optional arguments
   - Providing user-friendly help text

## Project Structure

```
week1-file-organizer/
├── main.go                 # Entry point with CLI
├── organizer/
│   ├── organizer.go       # Core file organization logic
│   └── config.go          # JSON configuration support
├── utils/
│   └── fileutils.go       # File utility functions
├── config.json            # Example configuration file
├── test_files/            # Test directory with samples
└── README.md              # This file
```

## How to Use

### Basic Usage

```bash
# Organize files in a directory
./file-organizer -source /path/to/files

# Dry run (preview changes without moving files)
./file-organizer -source /path/to/files -dry-run

# Organize into a different output directory
./file-organizer -source ./Downloads -output ./Organized

# List all file extension mappings
./file-organizer -list
```

### Advanced Usage

```bash
# Create a default config file
./file-organizer -create-config -config my-config.json

# Use a custom configuration
./file-organizer -source ./Downloads -config my-config.json

# Show help
./file-organizer -help
```

## Building from Source

```bash
# Build the application
go build -o file-organizer ./main.go

# Run tests (when you add them)
go test ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o file-organizer-linux
GOOS=darwin GOARCH=amd64 go build -o file-organizer-mac
GOOS=windows GOARCH=amd64 go build -o file-organizer.exe
```

## Configuration

Create a `config.json` file to customize file extension mappings:

```json
{
  "extensions": {
    ".txt": "Documents",
    ".pdf": "Documents",
    ".jpg": "Images",
    ".png": "Images",
    ".mp4": "Videos",
    ".mp3": "Music",
    ".zip": "Archives"
  },
  "defaultFolder": "Other"
}
```

Run `./file-organizer -create-config -config config.json` to generate a default config.

## Code Examples

### Working with Structs and Methods

```go
// Define a struct
type FileOrganizer struct {
    SourceDir string
    OutputDir string
}

// Create an instance with a function
organizer := NewFileOrganizer(sourceDir, outputDir)

// Call a method on the struct
err := organizer.Organize()
```

### Error Handling

```go
// Go's idiomatic error handling
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### Working with Maps

```go
// Create a map
mappings := map[string]string{
    ".txt": "Documents",
    ".jpg": "Images",
}

// Look up a value with comma-ok
folder, exists := mappings[".txt"]
if exists {
    fmt.Println("Found:", folder)
}

// Iterate over a map
for ext, folder := range mappings {
    fmt.Printf("%s -> %s\n", ext, folder)
}
```

### Working with Slices

```go
// Create an empty slice
var files []string

// Append to a slice
files = append(files, "document.txt")
files = append(files, "image.jpg")

// Iterate over a slice
for _, filename := range files {
    fmt.Println(filename)
}
```

### JSON Marshaling

```go
// Unmarshal JSON to struct
var config Config
json.Unmarshal([]byte(jsonStr), &config)

// Marshal struct to JSON
jsonBytes, _ := json.Marshal(config)

// Pretty print JSON
jsonBytes, _ := json.MarshalIndent(config, "", "  ")
```

## Features

✅ **Organize files by extension** - Automatically sort files into folders
✅ **Dry run mode** - Preview changes before applying them
✅ **Custom configuration** - Define your own extension mappings via JSON
✅ **Error handling** - Graceful error messages and validation
✅ **Directory statistics** - Shows file count and types before organizing
✅ **Directory tree view** - Visual representation of organized structure
✅ **Helpful CLI interface** - Flags, help text, and user-friendly messages

## Learning Path Alignment

This project covers:

- **Days 1-2**: Syntax & basics ✓
  - Variables, functions, structs (all implemented)
  - Running and building Go programs ✓

- **Days 3-5**: Deep dive ✓
  - Error handling (demonstrated throughout)
  - File operations (utils/fileutils.go)
  - JSON marshaling (organizer/config.go)

- **Days 6-7**: Practice ✓
  - Write production-ready code (not just exercises)
  - Real CLI application with multiple files
  - Test different scenarios (dry-run, config, etc.)

## Next Steps (Week 2+)

To extend this project:

1. **Add tests** - Implement `*_test.go` files with Go testing
2. **Add logging** - Use structured logging instead of `fmt.Println`
3. **Add more features** - Undo functionality, scheduled runs, watch mode
4. **Package it** - Create a `go.sum`, publish to GitHub, build releases
5. **Deploy** - Build for multiple platforms and create installers

## Key Takeaways

By building this project, you've learned:

1. ✅ How to structure a Go project with packages
2. ✅ The importance of error handling in Go
3. ✅ How to work with files and directories
4. ✅ JSON serialization and deserialization
5. ✅ Command-line argument parsing
6. ✅ How to write clear, commented code
7. ✅ Testing your CLI application manually

## Files to Study

For learning purposes, read the files in this order:

1. **organizer/organizer.go** - Core logic with detailed comments
2. **organizer/config.go** - JSON handling
3. **utils/fileutils.go** - File operations and utility functions
4. **main.go** - CLI interface and orchestration

Each file contains comments explaining Go concepts as they're used.

## Common Errors & Solutions

**Error: "source directory does not exist"**
```bash
# Make sure the directory exists
mkdir -p ./my-files
./file-organizer -source ./my-files
```

**Error: "failed to move file: file exists"**
```bash
# Files with the same name already exist in the destination
# Run with -dry-run first to see what would happen
./file-organizer -source ./Downloads -dry-run
```

**Error: "failed to create directory"**
```bash
# Check if you have write permissions in the output directory
ls -la ./path/to/output
```

## Testing

To test the application:

```bash
# Create test files
mkdir test_files
touch test_files/{document.txt,image.jpg,video.mp4,song.mp3}

# Test dry-run
./file-organizer -source test_files -dry-run

# Actually organize
./file-organizer -source test_files

# Verify the structure
ls -R test_files
```

## Resources

- [A Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Standard Library Documentation](https://pkg.go.dev/std)

## License

This is a learning project from the 1-Month Go Learning Path.

---

**Questions?** Review the code comments for detailed explanations of each Go concept!
