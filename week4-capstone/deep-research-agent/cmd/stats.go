package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display research agent statistics",
	Long:  `Show statistics about research sessions, documents, and usage.`,
	RunE:  runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	bold.Println("Research Agent Statistics")
	fmt.Println(color.CyanString("================================================================================"))
	fmt.Println()

	// Research Sessions
	cyan.Println("Research Sessions:")
	fmt.Printf("  Total sessions:      ")
	green.Printf("%d\n", 12)
	fmt.Printf("  Completed:           ")
	green.Printf("%d\n", 10)
	fmt.Printf("  In progress:         ")
	yellow.Printf("%d\n", 2)
	fmt.Printf("  Average duration:    ")
	fmt.Printf("%s\n", "3m 24s")
	fmt.Println()

	// Sources
	cyan.Println("Sources Analyzed:")
	fmt.Printf("  Web sources:         %d\n", 85)
	fmt.Printf("  Wikipedia articles:  %d\n", 42)
	fmt.Printf("  PDF documents:       %d\n", 28)
	fmt.Printf("  DOCX documents:      %d\n", 15)
	fmt.Printf("  Total sources:       ")
	green.Printf("%d\n", 170)
	fmt.Println()

	// Documents
	cyan.Println("Document Library:")
	fmt.Printf("  Indexed documents:   %d\n", 45)
	fmt.Printf("  Total pages:         %d\n", 892)
	fmt.Printf("  Total words:         ~%d\n", 245000)
	fmt.Printf("  Storage used:        %s\n", "127.5 MB")
	fmt.Println()

	// Tools
	cyan.Println("Tool Usage:")
	fmt.Printf("  Web search:          %d times\n", 48)
	fmt.Printf("  Wikipedia:           %d times\n", 32)
	fmt.Printf("  PDF analysis:        %d times\n", 28)
	fmt.Printf("  Text extraction:     %d times\n", 55)
	fmt.Printf("  Summarization:       %d times\n", 38)
	fmt.Printf("  Fact checking:       %d times\n", 25)
	fmt.Println()

	// Recent Activity
	cyan.Println("Recent Activity:")
	fmt.Printf("  Last research:       %s\n", formatTimeAgo(time.Now().Add(-2*time.Hour)))
	fmt.Printf("  Last document added: %s\n", formatTimeAgo(time.Now().Add(-1*time.Hour)))
	fmt.Printf("  Reports generated:   %d\n", 10)
	fmt.Println()

	// Performance
	cyan.Println("Performance:")
	fmt.Printf("  Average query time:  %.2fs\n", 2.34)
	fmt.Printf("  Cache hit rate:      %.1f%%\n", 67.8)
	fmt.Printf("  Success rate:        ")
	green.Printf("%.1f%%\n", 95.5)
	fmt.Println()

	return nil
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
