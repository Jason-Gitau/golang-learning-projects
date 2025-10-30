# 📖 URL Shortener - Code Reading Guide

Welcome! This guide walks you through the **URL Shortener** codebase so you understand every part.

## 🎯 What This Project Does

Create short, memorable codes for long URLs. Track how many times each shortened URL is accessed.

**Example:**
```
Input: https://www.github.com/golang/go
Output: http://short.url/abc123

Access abc123 → Shows original URL
Track visits → Shows "Accessed 5 times"
```

## 📚 Reading Order (45-60 minutes)

Read the files in this exact order to progressively understand the codebase:

### Step 1: Understand the CLI Interface (5 min)
**File:** `main.go`

**What to focus on:**

**Part A: Imports & Setup (Lines 1-20)**
- Global variables: `us` (shortener) and `st` (storage)
- Understand how they work together

**Part B: Main Function (Lines 22-50)**
- Lines 26-35: Initialize shortener and storage
- Lines 36-45: Load existing mappings from file
- Lines 47-57: Command routing (shorten, get, list, etc.)

**Part C: Command Handlers (Lines 59-200)**
- `handleShorten()` (lines 59-85) - Creates new short URL
- `handleGet()` (lines 87-110) - Retrieves original URL + increments visits
- `handleList()` (lines 112-135) - Shows all URLs in a table
- `handleStats()` (lines 137-160) - Shows detailed statistics
- `handleDelete()` (lines 162-180) - Removes a URL

**Key Learning:**
- How CLI routing works
- How commands pass data between functions
- How the shortener and storage work together

**Questions to answer:**
- How many commands does the tool support?
- What happens when you call `handleGet()`?
- Where are the data changes saved?

---

### Step 2: Learn Core Shortening Logic (15-20 min)
**File:** `shortener/shortener.go`

This is the MOST IMPORTANT file. It demonstrates all core Go concepts.

**What to focus on (in order):**

**Part A: Data Structures (Lines 1-35)**
- `URLMapping struct` (lines 5-21)
  - Stores: short code, original URL, creation date, visit count
  - Note the struct field tags: `json:"field_name"` for JSON serialization
- `URLShortener struct` (lines 23-35)
  - `Mappings map[string]*URLMapping` - In-memory cache of all shortened URLs
  - `CodeLength int` - How long should generated codes be
  - `charset string` - Characters used for code generation

**Key Learning:** Structs group related data, maps store it for fast lookup

---

**Part B: Constructor (Lines 37-48)**
- `NewURLShortener()` function
- Creates new struct instance with defaults
- Initializes map with `make(map[string]*URLMapping)`
- Returns pointer `*URLShortener` for efficiency

**Key Learning:** Constructors initialize state, pointers pass by reference

---

**Part C: Main ShortenURL Method (Lines 50-100)**
This is the CORE algorithm. Read this carefully!

1. **Validate input** (lines 55-57)
   - Check URL format
   - Return error if invalid

2. **Check if already shortened** (lines 59-66)
   - Iterate over all mappings with for...range
   - If URL exists, return existing code (avoid duplicates)

3. **Generate unique code** (lines 68-83)
   - Loop until we find unused code
   - Uses `generateCode()` (line 75)
   - Uses comma-ok idiom to check if code exists (line 77)

4. **Store mapping** (lines 85-96)
   - Create new URLMapping struct
   - Add to Mappings map with `us.Mappings[shortCode] = mapping`
   - Return the new short code

**Key Learning Concepts:**
- ✅ Maps for O(1) lookups
- ✅ Comma-ok idiom: `value, exists := map[key]`
- ✅ For...range loops over maps
- ✅ Error handling and validation
- ✅ String operations and time formatting

---

**Part D: Helper Methods (Lines 102-193)**

1. **GetURL()** (lines 102-120)
   - Looks up code in map
   - Increments visit counter
   - Returns original URL
   - Shows how to modify map values

2. **GetMapping()** (lines 122-130)
   - Returns full mapping details
   - Demonstrates pointer return

3. **ListAllMappings()** (lines 132-146)
   - Converts map to slice
   - Demonstrates: creating slice, for...range over map, append()

4. **GetStats()** (lines 148-168)
   - Creates stats map (map[string]interface{})
   - Shows flexible map types (values can be any type)
   - Calculates compression ratio

5. **DeleteURL()** (lines 170-178)
   - Uses `delete()` built-in function
   - Removes key from map

6. **generateCode()** (lines 195-208)
   - Creates random short code
   - Uses `rand.Intn()` to select random character
   - Demonstrates byte slices and string conversion
   - See init() at end (lines 272-276) for random seeding

7. **validateURL()** (lines 210-237)
   - Checks URL validity
   - Uses `strings.HasPrefix()` and `strings.Contains()`
   - Shows defensive programming

**Key Learning Concepts:**
- Maps, slices, for...range
- Error handling
- String operations
- Random generation
- Type conversions

---

### Step 3: Understand Data Persistence (10-15 min)
**File:** `storage/storage.go`

This file handles saving/loading to/from JSON file.

**What to focus on:**

**Part A: Struct Definition (Lines 1-10)**
- `Storage struct` - Just stores the file path
- Minimal design - single responsibility principle

---

**Part B: File Operations (Lines 12-60)**

1. **SaveMappings()** (lines 12-30)
   - Takes slice of mappings
   - Uses `json.MarshalIndent()` - Convert Go struct → JSON with pretty printing
   - Uses `os.WriteFile()` - Write to disk
   - `0644` is file permissions (rw-r--r--)

2. **LoadMappings()** (lines 32-60)
   - Checks if file exists with `os.Stat()`
   - Uses `os.ReadFile()` - Read entire file
   - Uses `json.Unmarshal()` - Convert JSON → Go struct
   - Handles empty files gracefully

**Key Learning Concepts:**
- File I/O: ReadFile, WriteFile, Stat
- JSON: Marshal/Unmarshal
- Error handling for file operations

---

**Part C: Data Mutation (Lines 62-130)**

1. **AppendMapping()** (lines 62-76)
   - Loads current mappings
   - Appends new one to slice
   - Saves all back to file

2. **RemoveMapping()** (lines 78-105)
   - Loads current mappings
   - Creates NEW slice without the removed mapping (filtering)
   - Saves to file
   - Demonstrates slice filtering pattern

3. **UpdateMapping()** (lines 107-130)
   - Loads current mappings
   - Finds and updates one mapping
   - Saves all back to file
   - Uses slice indexing to modify: `mappings[i] = mapping`

**Key Learning Concepts:**
- Slice filtering: create new slice without one element
- Slice indexing: modify element at index
- Stateless operations: load, modify, save

---

**Part D: Utility Methods (Lines 132-160)**
- `BackupFile()` - Creates backup
- `FileExists()` - Checks if file exists
- `GetFileSize()` - Gets file size
- `ClearAll()` - Deletes all mappings

**Key Learning:** How to work with file metadata

---

## 🔑 Key Concepts This Project Teaches

### 1. Maps for Fast Lookups
```go
// shortener.go Lines 29
Mappings: map[string]*URLMapping

// Line 77: Comma-ok idiom
if _, exists := us.Mappings[shortCode]; !exists {
    break  // Code doesn't exist, use it
}
```
✅ Learn: O(1) lookups, comma-ok idiom

### 2. Slices for Collections
```go
// shortener.go Lines 141-151
var mappings []*URLMapping
for _, mapping := range us.Mappings {
    mappings = append(mappings, mapping)
}
```
✅ Learn: Dynamic arrays, append, for...range

### 3. Pointers for Efficiency
```go
// shortener.go Line 43: Return pointer
return &URLShortener{ ... }

// Line 53: Method with pointer receiver
func (us *URLShortener) ShortenURL(url string) (string, error)
```
✅ Learn: Pointer receivers, memory efficiency

### 4. Error Handling
```go
// shortener.go Line 55-57
if err := us.validateURL(url); err != nil {
    return "", err
}
```
✅ Learn: Error returns, early returns, error propagation

### 5. JSON Marshaling
```go
// storage.go Line 20: Struct → JSON
data, _ := json.MarshalIndent(mappings, "", "  ")

// Line 45: JSON → Struct
json.Unmarshal(data, &mappings)
```
✅ Learn: Serialization with struct tags

### 6. File I/O
```go
// storage.go Line 22: Write file
os.WriteFile(s.filePath, data, 0644)

// Line 40: Read file
data, _ := os.ReadFile(s.filePath)
```
✅ Learn: File operations, permissions

### 7. Random Generation
```go
// shortener.go Lines 203-205
for i := range code {
    code[i] = us.charset[rand.Intn(len(us.charset))]
}
```
✅ Learn: Random numbers, seeding (line 273), byte manipulation

---

## 🚀 Quick Start (2 minutes)

```bash
# Build
go build -o url-shortener.exe

# Create URLs
./url-shortener.exe shorten https://github.com
./url-shortener.exe shorten https://golang.org

# List all
./url-shortener.exe list

# Get original
./url-shortener.exe get abc123

# See stats
./url-shortener.exe stats abc123

# Delete
./url-shortener.exe delete abc123
```

---

## 📝 Code Comments Guide

Every significant section has comments explaining:
- **What** the code does
- **Why** it uses certain patterns (e.g., "Demonstrates MAPS in Go")
- **How** to read/use it

Look for these comment patterns:
```go
// Demonstrates STRUCTS
type URLShortener struct { ... }

// Demonstrates MAPS
Mappings map[string]*URLMapping

// Demonstrates ERROR HANDLING
if err := someFunction(); err != nil {
    return "", err
}
```

---

## 💡 Learning Tips

### Tip 1: Understand Data Flow
1. `main.go` - User enters command
2. `handleShorten()` - Routes to handler
3. `shortener.ShortenURL()` - Creates mapping
4. `storage.AppendMapping()` - Saves to disk
5. `urls.json` - Data persists

Trace this flow by opening files and following function calls.

### Tip 2: Compare Methods
- `AppendMapping()` - Loads entire file, modifies, saves
- `UpdateMapping()` - Same pattern, different modification
- `RemoveMapping()` - Same pattern, filters slice

See the pattern? Load → Modify → Save

### Tip 3: Study the Idioms
- Comma-ok: `value, exists := map[key]`
- Error handling: `if err != nil { return err }`
- Slice filtering: Create new slice without item
- Struct tags: `json:"field_name"`

### Tip 4: Run and Experiment
```bash
# Build
go build -o url-shortener.exe

# Try commands
./url-shortener.exe shorten https://example.com
./url-shortener.exe list

# Look at urls.json
cat urls.json  # View saved data

# Try deleting, adding, etc.
```

---

## 🎯 What to Focus On

### Most Important (Must Understand)
- ✅ Maps for fast lookups (shortener.go line 29, 77)
- ✅ Slices and append (shortener.go line 141, 141-151)
- ✅ Error handling pattern (all files)
- ✅ Pointer receivers in methods (shortener.go line 53)
- ✅ JSON marshaling/unmarshaling (storage.go line 20, 45)

### Important (Should Understand)
- ✅ For...range over maps (shortener.go line 61)
- ✅ Comma-ok idiom (shortener.go line 77)
- ✅ File I/O operations (storage.go)
- ✅ Struct initialization (shortener.go line 43)

### Nice to Know (Can Learn Later)
- ✅ Random number generation (shortener.go line 203)
- ✅ Byte slice manipulation (shortener.go line 200)
- ✅ String operations (shortener.go line 226-237)
- ✅ Table formatting (main.go line 130)

---

## 📊 Code Statistics

- **Total Lines:** 630
- **Packages:** 2 (shortener, storage)
- **Functions:** 20+
- **Concepts Demonstrated:** 10+
- **Error Cases:** 15+

---

## 🔗 Navigation

After reading this project, you might want to:

1. **See Code in Action** → Run the project with various commands
2. **Read More** → Check `README.md` for full documentation
3. **Try Modifying** → Change code and rebuild to see effects
4. **Compare** → Read File Organizer project to compare approaches
5. **Move On** → Go to Week 2 and reference patterns from here

---

## ❓ Common Questions

**Q: Why use *URLShortener in methods?**
A: It's a pointer receiver, allowing methods to modify the struct's state. Without it, modifications wouldn't persist.

**Q: What's the comma-ok idiom?**
A: `value, exists := map[key]` returns both the value AND a boolean. It's the safe way to check if a key exists.

**Q: Why json.MarshalIndent instead of json.Marshal?**
A: MarshalIndent adds pretty printing (indentation), making the JSON readable. Marshal produces compact JSON.

**Q: How does persistence work?**
A: Every change loads the file, modifies it in memory, then saves back to disk. Simple but effective.

**Q: Why use maps instead of slices?**
A: Maps give O(1) lookup time by code, slices would require O(n) searching.

---

## 🎓 After Reading This

You should understand:
- ✅ How maps enable fast lookups
- ✅ How slices handle dynamic collections
- ✅ How pointers allow state modification
- ✅ How to handle errors properly
- ✅ How to persist data with JSON
- ✅ How to coordinate multiple packages (shortener + storage)
- ✅ Random code generation
- ✅ File I/O and persistence patterns

---

**Status:** Ready to read! Start with main.go 📖

Then read files in this order:
1. main.go (understand CLI flow)
2. shortener/shortener.go (core logic)
3. storage/storage.go (persistence)

Go back to main repository → [../START_HERE.md](../START_HERE.md)
