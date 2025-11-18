package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Session represents a research session
type Session struct {
	ID          string
	Query       string
	Depth       string
	Sources     int
	Created     time.Time
	Updated     time.Time
	Status      string
	ReportPath  string
}

var (
	sessionFormat string
)

// sessionCmd represents the session management commands
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage research sessions",
	Long:  `List, show, resume, and delete research sessions.`,
}

var listSessionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all research sessions",
	Long:  `Display a list of all saved research sessions.`,
	RunE:  runListSessions,
}

var showSessionCmd = &cobra.Command{
	Use:   "show [session-id]",
	Short: "Show details of a research session",
	Long:  `Display detailed information about a specific research session.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runShowSession,
}

var resumeSessionCmd = &cobra.Command{
	Use:   "resume [session-id]",
	Short: "Resume a research session",
	Long:  `Continue research from a previous session, optionally adding new sources.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runResumeSession,
}

var deleteSessionCmd = &cobra.Command{
	Use:   "delete [session-id]",
	Short: "Delete a research session",
	Long:  `Remove a research session and its associated data.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteSession,
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(listSessionsCmd)
	sessionCmd.AddCommand(showSessionCmd)
	sessionCmd.AddCommand(resumeSessionCmd)
	sessionCmd.AddCommand(deleteSessionCmd)

	// Flags
	listSessionsCmd.Flags().StringVar(&sessionFormat, "format", "table", "output format: table or json")
}

func runListSessions(cmd *cobra.Command, args []string) error {
	// Mock sessions - in real implementation, this would query the database
	sessions := getMockSessions()

	if len(sessions) == 0 {
		fmt.Println("No research sessions found.")
		return nil
	}

	if sessionFormat == "json" {
		// Print as JSON
		fmt.Println("[")
		for i, s := range sessions {
			fmt.Printf("  {\"id\": \"%s\", \"query\": \"%s\", \"status\": \"%s\"}", s.ID, s.Query, s.Status)
			if i < len(sessions)-1 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Println("]")
		return nil
	}

	// Print as table
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	bold.Println("\nResearch Sessions:")
	fmt.Println(color.CyanString(fmt.Sprintf("%-12s %-40s %-10s %-8s %-20s", "ID", "Query", "Status", "Sources", "Created")))
	fmt.Println(color.CyanString("--------------------------------------------------------------------------------"))

	for _, s := range sessions {
		statusColor := green
		if s.Status == "in_progress" {
			statusColor = yellow
		}

		fmt.Printf("%-12s %-40s ", s.ID, truncateString(s.Query, 40))
		statusColor.Printf("%-10s ", s.Status)
		fmt.Printf("%-8d %-20s\n", s.Sources, s.Created.Format("2006-01-02 15:04"))
	}

	fmt.Printf("\nTotal sessions: %d\n\n", len(sessions))

	return nil
}

func runShowSession(cmd *cobra.Command, args []string) error {
	sessionID := args[0]

	// Mock session retrieval
	session := findSessionByID(sessionID)
	if session == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	bold.Println("Research Session Details")
	fmt.Println(color.CyanString("================================================================================"))
	fmt.Println()

	cyan.Print("ID:          ")
	fmt.Println(session.ID)

	cyan.Print("Query:       ")
	fmt.Println(session.Query)

	cyan.Print("Depth:       ")
	fmt.Println(session.Depth)

	cyan.Print("Status:      ")
	if session.Status == "completed" {
		color.Green(session.Status)
	} else {
		color.Yellow(session.Status)
	}

	cyan.Print("Sources:     ")
	fmt.Printf("%d\n", session.Sources)

	cyan.Print("Created:     ")
	fmt.Println(session.Created.Format("2006-01-02 15:04:05"))

	cyan.Print("Updated:     ")
	fmt.Println(session.Updated.Format("2006-01-02 15:04:05"))

	if session.ReportPath != "" {
		cyan.Print("Report:      ")
		fmt.Println(session.ReportPath)
	}

	fmt.Println()

	return nil
}

func runResumeSession(cmd *cobra.Command, args []string) error {
	sessionID := args[0]

	session := findSessionByID(sessionID)
	if session == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	fmt.Printf("Resuming research session: %s\n", sessionID)
	fmt.Printf("Original query: %s\n", session.Query)
	fmt.Println("\nContinuing research...")

	// In real implementation, this would resume the research
	time.Sleep(1 * time.Second)

	color.Green("\nSession resumed successfully!")
	fmt.Printf("Use 'research-agent show-session %s' to view details\n", sessionID)

	return nil
}

func runDeleteSession(cmd *cobra.Command, args []string) error {
	sessionID := args[0]

	session := findSessionByID(sessionID)
	if session == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	fmt.Printf("Deleting research session: %s\n", sessionID)
	fmt.Printf("Query: %s\n", session.Query)

	// In real implementation, this would delete from database
	time.Sleep(300 * time.Millisecond)

	color.Green("\nSession deleted successfully!")

	return nil
}

func getMockSessions() []Session {
	return []Session{
		{
			ID:      "abc123",
			Query:   "AI ethics in healthcare",
			Depth:   "deep",
			Sources: 15,
			Created: time.Now().Add(-2 * time.Hour),
			Updated: time.Now().Add(-1 * time.Hour),
			Status:  "completed",
			ReportPath: "./reports/ai-ethics-healthcare.md",
		},
		{
			ID:      "def456",
			Query:   "Quantum computing applications",
			Depth:   "medium",
			Sources: 10,
			Created: time.Now().Add(-24 * time.Hour),
			Updated: time.Now().Add(-20 * time.Hour),
			Status:  "completed",
			ReportPath: "./reports/quantum-computing.md",
		},
		{
			ID:      "ghi789",
			Query:   "Climate change impact on agriculture",
			Depth:   "deep",
			Sources: 8,
			Created: time.Now().Add(-1 * time.Hour),
			Updated: time.Now().Add(-30 * time.Minute),
			Status:  "in_progress",
		},
	}
}

func findSessionByID(id string) *Session {
	sessions := getMockSessions()
	for _, s := range sessions {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
