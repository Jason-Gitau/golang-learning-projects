package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// ProgressTracker tracks research progress
type ProgressTracker struct {
	query         string
	currentStep   string
	currentIndex  int
	totalSteps    int
	sourcesFound  int
	startTime     time.Time
	bar           *progressbar.ProgressBar
	mu            sync.Mutex
	stopped       bool
	completedSteps []string
	pendingSteps   []string
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(query string) *ProgressTracker {
	return &ProgressTracker{
		query:          query,
		totalSteps:     10,
		startTime:      time.Now(),
		completedSteps: []string{},
		pendingSteps: []string{
			"Web search",
			"Wikipedia lookup",
			"Document analysis",
			"Information extraction",
			"Fact checking",
			"Summarization",
		},
	}
}

// Start begins the progress display
func (pt *ProgressTracker) Start() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	fmt.Printf("\n%s Researching: %s\n", color.CyanString("🔍"), color.New(color.Bold).Sprint(pt.query))
	fmt.Println()

	pt.bar = progressbar.NewOptions(pt.totalSteps,
		progressbar.OptionSetDescription("Progress"),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "━",
			SaucerHead:    "━",
			SaucerPadding: "━",
			BarStart:      "",
			BarEnd:        "",
		}),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetElapsedTime(true),
	)
}

// UpdateStep updates the current step
func (pt *ProgressTracker) UpdateStep(step string, index, total int) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.stopped {
		return
	}

	pt.currentStep = step
	pt.currentIndex = index
	pt.totalSteps = total

	// Update progress bar
	if pt.bar != nil {
		pt.bar.Set(index)
	}

	// Print step status
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Printf("\n%s %s (%d/%d)\n",
		yellow.Sprint("⚡"),
		color.New(color.Bold).Sprint(step),
		index+1,
		total,
	)

	// Add to completed steps
	if index > 0 {
		prevStep := fmt.Sprintf("%s completed", step)
		if !contains(pt.completedSteps, prevStep) {
			pt.completedSteps = append(pt.completedSteps, prevStep)
			fmt.Printf("%s %s\n", green.Sprint("✓"), green.Sprint(prevStep))
		}
	}
}

// AddSource increments the sources found counter
func (pt *ProgressTracker) AddSource(sourceType string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.sourcesFound++
	green := color.New(color.FgGreen)
	fmt.Printf("%s Found %s source (total: %d)\n",
		green.Sprint("✓"),
		sourceType,
		pt.sourcesFound,
	)
}

// ShowStatus displays current status
func (pt *ProgressTracker) ShowStatus() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	cyan := color.New(color.FgCyan)
	fmt.Println()
	fmt.Println(cyan.Sprint("Current Status:"))
	fmt.Printf("  Step: %s (%d/%d)\n", pt.currentStep, pt.currentIndex+1, pt.totalSteps)
	fmt.Printf("  Sources: %d\n", pt.sourcesFound)
	fmt.Printf("  Elapsed: %s\n", time.Since(pt.startTime).Round(time.Second))
	fmt.Println()
}

// Stop completes the progress display
func (pt *ProgressTracker) Stop() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.stopped {
		return
	}

	pt.stopped = true

	if pt.bar != nil {
		pt.bar.Finish()
	}

	green := color.New(color.FgGreen, color.Bold)
	fmt.Println()
	fmt.Printf("%s Research completed!\n", green.Sprint("✓"))
	fmt.Printf("  Sources found: %d\n", pt.sourcesFound)
	fmt.Printf("  Time elapsed: %s\n", time.Since(pt.startTime).Round(time.Second))
	fmt.Println()
}

// DisplaySummary shows a final summary
func (pt *ProgressTracker) DisplaySummary(findings int, reportPath string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	bold.Println("Research Summary")
	fmt.Println(color.CyanString("================================================================================"))
	fmt.Println()

	cyan.Print("Query:         ")
	fmt.Println(pt.query)

	cyan.Print("Duration:      ")
	fmt.Println(time.Since(pt.startTime).Round(time.Second))

	cyan.Print("Sources:       ")
	fmt.Printf("%d\n", pt.sourcesFound)

	cyan.Print("Findings:      ")
	fmt.Printf("%d\n", findings)

	if reportPath != "" {
		cyan.Print("Report:        ")
		fmt.Println(reportPath)
	}

	fmt.Println()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SimpleProgress shows a simple progress indicator
func SimpleProgress(message string) *SimpleProgressIndicator {
	return &SimpleProgressIndicator{
		message:   message,
		startTime: time.Now(),
		stopped:   false,
	}
}

// SimpleProgressIndicator is a simple spinner
type SimpleProgressIndicator struct {
	message   string
	startTime time.Time
	stopped   bool
	mu        sync.Mutex
}

// Start begins the spinner
func (sp *SimpleProgressIndicator) Start() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	fmt.Printf("%s %s...\n", color.YellowString("⚡"), sp.message)
}

// Stop completes the spinner
func (sp *SimpleProgressIndicator) Stop() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.stopped {
		return
	}

	sp.stopped = true
	fmt.Printf("%s %s (took %s)\n",
		color.GreenString("✓"),
		sp.message,
		time.Since(sp.startTime).Round(time.Millisecond),
	)
}
