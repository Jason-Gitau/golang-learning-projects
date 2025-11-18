package cmd

import (
	"github.com/spf13/cobra"
	"deep-research-agent/ui"
)

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start an interactive research session",
	Long: `Start an interactive session where you can research topics,
view results, and export reports through a guided interface.

This mode is ideal for:
- Exploratory research
- Iterative investigation
- Learning how to use the research agent

The interactive mode will guide you through:
1. Defining your research query
2. Selecting research depth and sources
3. Viewing results in real-time
4. Exporting reports in various formats`,
	RunE: runInteractive,
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

func runInteractive(cmd *cobra.Command, args []string) error {
	session := ui.NewInteractiveSession()

	// Start the interactive session
	if err := session.Start(); err != nil {
		return err
	}

	// Run the menu
	return session.RunMenu()
}
