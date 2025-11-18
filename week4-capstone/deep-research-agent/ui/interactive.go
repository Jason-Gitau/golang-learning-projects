package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// InteractiveSession manages an interactive research session
type InteractiveSession struct {
	reader *bufio.Reader
	query  string
	depth  string
	useWeb bool
}

// NewInteractiveSession creates a new interactive session
func NewInteractiveSession() *InteractiveSession {
	return &InteractiveSession{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Start begins the interactive session
func (is *InteractiveSession) Start() error {
	// Welcome message
	is.printWelcome()

	// Get research query
	query, err := is.askQuestion("What would you like to research?")
	if err != nil {
		return err
	}
	is.query = query

	// Confirm research
	fmt.Println()
	fmt.Printf("Analyzing query: %s\n", color.CyanString(query))
	fmt.Println()

	// Ask about web sources
	useWeb, err := is.askYesNo("Should I include web sources?")
	if err != nil {
		return err
	}
	is.useWeb = useWeb

	// Ask about depth
	depth, err := is.askDepth()
	if err != nil {
		return err
	}
	is.depth = depth

	// Confirm and start
	fmt.Println()
	is.printResearchPlan()
	fmt.Println()

	confirm, err := is.askYesNo("Start research?")
	if err != nil {
		return err
	}

	if !confirm {
		fmt.Println("Research cancelled.")
		return nil
	}

	// Start research (placeholder - in real implementation, this would call the agent)
	fmt.Println()
	fmt.Println(color.GreenString("Starting research..."))
	fmt.Println()

	return nil
}

// RunMenu displays and handles the interactive menu
func (is *InteractiveSession) RunMenu() error {
	for {
		fmt.Println()
		is.printMenu()

		choice, err := is.askQuestion("> ")
		if err != nil {
			return err
		}

		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			is.viewSummary()
		case "2":
			is.exportReport()
		case "3":
			is.addSources()
		case "4":
			// New research - start over
			return is.Start()
		case "5", "quit", "exit", "q":
			fmt.Println(color.GreenString("Goodbye!"))
			return nil
		default:
			fmt.Println(color.RedString("Invalid choice. Please try again."))
		}
	}
}

func (is *InteractiveSession) printWelcome() {
	bold := color.New(color.Bold, color.FgCyan)

	fmt.Println()
	fmt.Println(color.CyanString("================================================================================"))
	bold.Println("                    Welcome to Research Agent 🔍")
	fmt.Println(color.CyanString("================================================================================"))
	fmt.Println()
	fmt.Println("An AI-powered research assistant that helps you:")
	fmt.Println("  • Conduct comprehensive research on any topic")
	fmt.Println("  • Analyze documents and web sources")
	fmt.Println("  • Generate professional reports with citations")
	fmt.Println()
}

func (is *InteractiveSession) printResearchPlan() {
	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)

	bold.Println("Research Plan:")

	cyan.Print("  Query:       ")
	fmt.Println(is.query)

	cyan.Print("  Depth:       ")
	fmt.Println(is.depth)

	cyan.Print("  Web sources: ")
	if is.useWeb {
		fmt.Println(color.GreenString("Enabled"))
	} else {
		fmt.Println(color.YellowString("Disabled"))
	}
}

func (is *InteractiveSession) printMenu() {
	bold := color.New(color.Bold)

	bold.Println("What would you like to do?")
	fmt.Println("  1. View summary")
	fmt.Println("  2. Export report")
	fmt.Println("  3. Add more sources")
	fmt.Println("  4. New research")
	fmt.Println("  5. Quit")
}

func (is *InteractiveSession) viewSummary() {
	fmt.Println()
	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)

	bold.Println("Research Summary")
	fmt.Println(color.CyanString("================================================================================"))
	fmt.Println()

	cyan.Println("Key Findings:")
	fmt.Println()
	fmt.Println("1. Significant progress has been made in recent years")
	fmt.Println("   • Multiple breakthroughs reported")
	fmt.Println("   • Sources: [1], [2], [3]")
	fmt.Println()
	fmt.Println("2. Challenges remain in implementation")
	fmt.Println("   • Technical limitations identified")
	fmt.Println("   • Sources: [4], [5]")
	fmt.Println()

	cyan.Println("Sources:")
	fmt.Println("  [1] Smith, J. (2024). Recent Advances")
	fmt.Println("  [2] Johnson, M. (2023). Technical Review")
	fmt.Println("  [3] Williams, K. (2024). Comprehensive Analysis")
	fmt.Println()
}

func (is *InteractiveSession) exportReport() {
	fmt.Println()

	format, err := is.askFormat()
	if err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		return
	}

	filename, err := is.askQuestion("Enter filename (or press Enter for default): ")
	if err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		return
	}

	if filename == "" {
		ext := format
		if format == "markdown" {
			ext = "md"
		}
		filename = fmt.Sprintf("research-report.%s", ext)
	}

	fmt.Println()
	fmt.Printf("Exporting report as %s to %s...\n", format, filename)

	// Simulate export
	fmt.Println(color.GreenString("✓ Report exported successfully!"))
}

func (is *InteractiveSession) addSources() {
	fmt.Println()
	fmt.Println("Add more sources:")
	fmt.Println("  1. Add PDF document")
	fmt.Println("  2. Add DOCX document")
	fmt.Println("  3. Search web")
	fmt.Println("  4. Back to menu")

	choice, err := is.askQuestion("> ")
	if err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		return
	}

	switch strings.TrimSpace(choice) {
	case "1", "2":
		path, err := is.askQuestion("Enter file path: ")
		if err != nil {
			fmt.Println(color.RedString("Error: %v", err))
			return
		}
		fmt.Printf("Adding document: %s\n", path)
		fmt.Println(color.GreenString("✓ Document added successfully!"))
	case "3":
		fmt.Println(color.YellowString("Web search integration coming soon!"))
	case "4":
		return
	default:
		fmt.Println(color.RedString("Invalid choice"))
	}
}

func (is *InteractiveSession) askQuestion(prompt string) (string, error) {
	fmt.Print(prompt + " ")

	response, err := is.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

func (is *InteractiveSession) askYesNo(question string) (bool, error) {
	for {
		response, err := is.askQuestion(question + " (y/n)")
		if err != nil {
			return false, err
		}

		response = strings.ToLower(response)

		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}

		fmt.Println(color.YellowString("Please enter 'y' or 'n'"))
	}
}

func (is *InteractiveSession) askDepth() (string, error) {
	fmt.Println()
	fmt.Println("How deep should I research?")
	fmt.Println("  1. Shallow  (quick overview, 5-10 sources)")
	fmt.Println("  2. Medium   (balanced research, 10-15 sources)")
	fmt.Println("  3. Deep     (comprehensive, 15-25 sources)")

	for {
		choice, err := is.askQuestion("> ")
		if err != nil {
			return "", err
		}

		switch strings.TrimSpace(choice) {
		case "1", "shallow":
			return "shallow", nil
		case "2", "medium":
			return "medium", nil
		case "3", "deep":
			return "deep", nil
		default:
			fmt.Println(color.YellowString("Please enter 1, 2, or 3"))
		}
	}
}

func (is *InteractiveSession) askFormat() (string, error) {
	fmt.Println("Select export format:")
	fmt.Println("  1. Markdown")
	fmt.Println("  2. PDF")
	fmt.Println("  3. JSON")

	for {
		choice, err := is.askQuestion("> ")
		if err != nil {
			return "", err
		}

		switch strings.TrimSpace(choice) {
		case "1", "markdown", "md":
			return "markdown", nil
		case "2", "pdf":
			return "pdf", nil
		case "3", "json":
			return "json", nil
		default:
			fmt.Println(color.YellowString("Please enter 1, 2, or 3"))
		}
	}
}
