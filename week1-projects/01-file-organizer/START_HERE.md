# 📖 File Organizer - Code Reading Guide

Welcome! This guide walks you through the **File Organizer** codebase so you understand every part.

## 🎯 What This Project Does

Automatically organize files in any folder by their extension into categorized subdirectories.

**Example:**
```
Before:
├── document.txt
├── image.jpg
├── video.mp4
└── song.mp3

After:
├── Documents/
│   └── document.txt
├── Images/
│   └── image.jpg
├── Videos/
│   └── video.mp4
└── Music/
    └── song.mp3
```

## 📚 Reading Order (30-45 minutes)

Read the files in this exact order to progressively understand the codebase:

### Step 1: Understand the Entry Point (5 min)
**File:** `main.go`

**What to focus on:**
- Lines 1-15: Imports and global variables
- Lines 18-35: Flag definitions (--source, --output, --dry-run, etc.)
- Lines 37-60: Main logic flow - how commands are handled
- Lines 85-110: Helper functions (usage, stats display)

**Key Learning:** How CLI applications use flags and route commands.

**Questions to answer:**
- What flags can users pass?
- How does the program handle the -dry-run flag?
- Where does it create the FileOrganizer instance?

---

### Step 2: Learn Core Logic (10-15 min)
**File:** `organizer/organizer.go`

**What to focus on (in order):**

**Part A: Struct Definition (Lines 1-35)**
- `URLMapping struct` (lines 5-21) - What data gets stored about each file mapping
- `URLShortener struct` (lines 23-35) - How the organizer stores its state
- Understand: Why use maps? Why use pointers?

**Part B: Constructor (Lines 36-48)**
- `NewFileOrganizer()` function - Creates and initializes a new organizer
- Shows how to initialize a map with default values
- Demonstrates pointer returns

**Part C: Main Methods (Lines 50-130)**
- `ShortenURL()` (lines 50-100) - The core logic:
  - Validates input
  - Checks if already shortened
  - Generates new short code
  - Stores the mapping
- `GetURL()` (lines 102-120) - Retrieves original URL and tracks visits
- Demonstrates: Maps, error handling, loops, string operations

**Part D: Helper Methods (Lines 132-230)**
- `GetMapping()` - Retrieves full mapping details
- `ListAllMappings()` - Returns all mappings as a slice
- `generateCode()` - Creates random short codes (lines 195-205)
- `validateURL()` - Input validation (lines 207-240)

**Key Learning Concepts:**
1. **Structs** - How to define and use custom types
2. **Methods** - Functions that belong to types (note the `(us *URLShortener)` part)
3. **Pointers** - Why use `*` in method receivers
4. **Maps** - Fast lookups, comma-ok idiom
5. **Slices** - append() operations, iteration
6. **Error Handling** - Returning errors, error wrapping
7. **Random Generation** - Seeding, charset selection

**Questions to answer:**
- What does the organizer.Organize() method do?
- How are files stored in memory?
- What's the purpose of error wrapping with %w?

---

### Step 3: Understand Data Persistence (5-10 min)
**File:** `organizer/config.go`

**What to focus on:**

**Part A: Struct Definition (Lines 1-20)**
- `Config struct` - Represents JSON configuration structure
- Note the `json:"..."` tags - How Go maps JSON keys to struct fields

**Part B: JSON Operations (Lines 22-80)**
- `LoadConfig()` (lines 22-50) - Reads JSON from file, unmarshals to struct
  - json.Unmarshal (line 45) - JSON bytes → Go struct
- `SaveConfig()` (lines 52-70) - Marshals struct to JSON, writes to file
  - json.MarshalIndent (line 64) - Go struct → pretty JSON
- `CreateDefaultConfig()` (lines 72-90) - Generates starter config

**Key Learning Concepts:**
1. **JSON Marshaling** - Converting Go structs ↔ JSON
2. **File I/O** - os.ReadFile(), os.WriteFile()
3. **Struct Tags** - How `json:"field_name"` works
4. **Error Handling** - Proper error wrapping

**Questions to answer:**
- What happens if the config file doesn't exist?
- Why use MarshalIndent instead of Marshal?
- What does the json tags do?

---

### Step 4: Learn File Operations (5-10 min)
**File:** `utils/fileutils.go`

**What to focus on (skim through):**

**Important Functions:**

1. **DirectoryExists()** (lines 1-10)
   - Uses os.Stat() to check if path exists
   - Shows the blank identifier (_)

2. **ListFiles()** (lines 70-95)
   - Reads directory contents
   - Filters files (not directories)
   - Demonstrates slice append operations

3. **CountFilesByExtension()** (lines 97-135)
   - Shows map usage for statistics
   - Iterates over directory
   - Groups files by extension

4. **PrintDirectoryTree()** (lines 195-235)
   - Demonstrates recursion
   - Formatting output nicely

**Key Learning Concepts:**
1. **os.ReadDir()** - Efficient directory reading
2. **os.Stat()** - Get file information
3. **Slices** - Dynamic lists, append()
4. **Maps** - Collecting statistics
5. **Recursion** - Tree printing

---

## 🔑 Key Concepts This Project Teaches

### 1. Structs & Methods
```go
// Lines 1-35 of organizer/organizer.go
type FileOrganizer struct { /* fields */ }

func (fo *FileOrganizer) Organize() error { /* method */ }
```
✅ Learn: Custom data types and functions on types

### 2. Pointers
```go
// Line 43: Return pointer for efficiency
return &FileOrganizer{ ... }

// Line 53: Pointer receiver (modifies state)
func (fo *FileOrganizer) Organize() error
```
✅ Learn: Why pointers matter (efficiency + state modification)

### 3. Maps
```go
// Lines 25-30: Map initialization
mappings: map[string]*URLMapping{ /* defaults */ }

// Line 61: Map lookup with comma-ok
mapping, exists := us.mappings[shortCode]
```
✅ Learn: Fast lookups, comma-ok idiom

### 4. Slices
```go
// Lines 141-151: Slice operations
var mappings []*URLMapping
for _, mapping := range us.Mappings {
    mappings = append(mappings, mapping)
}
```
✅ Learn: Dynamic arrays, append, for...range

### 5. Error Handling
```go
// Line 55: Validate input
if err := us.validateURL(url); err != nil {
    return "", err  // Return error immediately
}

// Line 84: Error wrapping for context
return "", fmt.Errorf("failed to create: %w", err)
```
✅ Learn: Return errors as values, error wrapping

### 6. File I/O
```go
// config.go Line 30: Read file
data, err := os.ReadFile(filepath)

// config.go Line 63: Write file
os.WriteFile(filepath, data, 0644)
```
✅ Learn: Reading/writing files, permissions

### 7. JSON
```go
// config.go Line 45: JSON → struct
json.Unmarshal(data, &config)

// config.go Line 64: Struct → JSON
json.MarshalIndent(config, "", "  ")
```
✅ Learn: Serialization, struct tags

---

## 🚀 Quick Start (2 minutes)

```bash
# Build
go build -o file-organizer.exe

# Test with dry-run (preview)
./file-organizer.exe -source . -dry-run

# See what happens (list mappings)
./file-organizer.exe -list

# Get help
./file-organizer.exe -help
```

---

## 📝 Code Comments Guide

Every significant function has comments explaining:
- **What** the function does
- **Why** it uses certain patterns
- **How** to use the returned values

Look for comment lines starting with:
- `//` - Single line explanation
- `/* */` - Multi-line explanation
- Function comments above function definitions

---

## 💡 Learning Tips

### Tip 1: Understand the Flow
1. Open main.go
2. Find the handleOrganize() function
3. Trace what it calls step by step
4. Look up each function in other files

### Tip 2: Find Examples
- Look for `url, err :=` patterns for error handling
- Look for `for code, mapping := range` patterns for maps
- Look for `append(...)` for slices

### Tip 3: Don't Memorize Syntax
- Focus on understanding the pattern
- Why use pointers? (efficiency)
- Why use maps? (fast lookup)
- Why return errors? (Go's philosophy)

### Tip 4: Run It!
After reading main.go, build and run:
```bash
go build -o file-organizer.exe
./file-organizer.exe -help
./file-organizer.exe -list
```

---

## 🎯 What to Focus On

### Most Important Concepts (Must Learn)
- ✅ Structs definition and usage (organizer.go lines 1-35)
- ✅ Methods with pointer receivers (organizer.go line 53)
- ✅ Error handling pattern (all files, especially main.go)
- ✅ Maps (organizer.go line 25-30, 61)
- ✅ JSON marshaling (config.go line 45, 64)

### Important Concepts (Should Learn)
- ✅ Slices and append (utils/fileutils.go)
- ✅ File I/O operations (config.go, utils/fileutils.go)
- ✅ CLI flags (main.go line 18-35)

### Nice to Know (Can Learn Later)
- ✅ Recursion (utils/fileutils.go PrintDirectoryTree)
- ✅ Advanced map operations
- ✅ Optimization techniques

---

## 📊 Code Statistics

- **Total Lines:** 754
- **Packages:** 3 (main, organizer, utils)
- **Functions:** 15+
- **Concepts Demonstrated:** 10+
- **Error Cases:** 10+

---

## 🔗 Navigation

After reading this project, you might want to:

1. **See Code in Action** → Run the project with test files
2. **Read More** → Check `README.md` for full documentation
3. **Try Modifying** → Change code and see what happens
4. **Compare** → Read the URL Shortener project to compare approaches
5. **Move On** → Go to Week 2 and reference patterns from here

---

## ❓ Common Questions

**Q: Do I need to understand every line?**
A: No, focus on the big picture first, then dive into details.

**Q: What's a pointer?**
A: A variable that stores a memory address. `*Type` is a pointer, `&variable` gets the address.

**Q: Why `(fo *FileOrganizer)` in methods?**
A: The `(fo *FileOrganizer)` is the receiver, like `this` in other languages. The `*` means it can modify the struct.

**Q: What's the comma-ok idiom?**
A: `value, exists := map[key]` returns the value AND a boolean saying if the key exists.

**Q: Why error wrapping with %w?**
A: It preserves the original error so you can trace the full error chain.

---

## 🎓 After Reading This

You should understand:
- ✅ How Go structs work
- ✅ How methods attach behavior to types
- ✅ Why pointers matter
- ✅ How maps work and when to use them
- ✅ How Go handles errors
- ✅ How to read/write JSON
- ✅ How to work with files

---

**Status:** Ready to read! Start with main.go 📖

Go back to main repository → [../START_HERE.md](../START_HERE.md)
