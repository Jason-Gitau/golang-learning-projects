package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"deep-research-agent/config"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "research-agent",
	Short: "AI-powered research agent with document analysis and web search",
	Long: `Deep Research Agent is a production-ready AI research assistant that helps you:

- Conduct comprehensive research on any topic
- Analyze PDF and DOCX documents
- Search the web and Wikipedia
- Generate professional research reports with citations
- Manage research sessions and sources

The agent uses multiple specialized tools to gather, analyze, and synthesize
information from various sources, providing you with well-researched, cited reports.`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.research-agent.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if err := config.InitConfig(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
