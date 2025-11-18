package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"deep-research-agent/api"
)

var (
	port        int
	workers     int
	dbPath      string
	uploadDir   string
	exportDir   string
	debugMode   bool
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server with WebSocket support",
	Long: `Start the Deep Research Agent API server.

This command starts a REST API server with WebSocket support for real-time
research progress updates. The server provides:

- REST API endpoints for research operations
- WebSocket streaming for real-time progress
- Document upload and analysis
- Session management
- Export functionality (Markdown, JSON, PDF)
- Job queue with async processing

Example:
  research-agent serve --port 8080 --workers 10
  research-agent serve --port 3000 --debug

The server will start on the specified port (default: 8080) and spawn
worker processes to handle research jobs asynchronously.`,
	RunE: runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Server configuration flags
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	serveCmd.Flags().IntVarP(&workers, "workers", "w", 10, "Number of worker processes")
	serveCmd.Flags().StringVar(&dbPath, "db", "data/research.db", "Path to the database file")
	serveCmd.Flags().StringVar(&uploadDir, "upload-dir", "uploads/files", "Directory for uploaded files")
	serveCmd.Flags().StringVar(&exportDir, "export-dir", "exports", "Directory for exported files")
	serveCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode")

	// Bind flags to viper
	viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.workers", serveCmd.Flags().Lookup("workers"))
	viper.BindPFlag("server.db_path", serveCmd.Flags().Lookup("db"))
	viper.BindPFlag("server.upload_dir", serveCmd.Flags().Lookup("upload-dir"))
	viper.BindPFlag("server.export_dir", serveCmd.Flags().Lookup("export-dir"))
	viper.BindPFlag("server.debug", serveCmd.Flags().Lookup("debug"))
}

func runServe(cmd *cobra.Command, args []string) error {
	// Create server configuration
	config := api.ServerConfig{
		Port:        port,
		WorkerCount: workers,
		DBPath:      dbPath,
		UploadDir:   uploadDir,
		ExportDir:   exportDir,
		Debug:       debugMode,
	}

	// Print startup banner
	printBanner(config)

	// Create server
	server, err := api.NewServer(config)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Setup graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.Start()
	}()

	// Wait for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		// Give outstanding requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
			return err
		}

		log.Println("Server stopped gracefully")
	}

	return nil
}

func printBanner(config api.ServerConfig) {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║     Deep Research Agent API Server                          ║
║     AI-Powered Research with Real-time Streaming            ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
	fmt.Printf("Server Configuration:\n")
	fmt.Printf("  Port:        %d\n", config.Port)
	fmt.Printf("  Workers:     %d\n", config.WorkerCount)
	fmt.Printf("  Database:    %s\n", config.DBPath)
	fmt.Printf("  Upload Dir:  %s\n", config.UploadDir)
	fmt.Printf("  Export Dir:  %s\n", config.ExportDir)
	fmt.Printf("  Debug Mode:  %v\n", config.Debug)
	fmt.Println()
	fmt.Printf("API Endpoints:\n")
	fmt.Printf("  Health:      http://localhost:%d/health\n", config.Port)
	fmt.Printf("  Research:    http://localhost:%d/api/v1/research/start\n", config.Port)
	fmt.Printf("  WebSocket:   ws://localhost:%d/api/v1/research/:id/stream\n", config.Port)
	fmt.Printf("  Documents:   http://localhost:%d/api/v1/documents/upload\n", config.Port)
	fmt.Printf("  Sessions:    http://localhost:%d/api/v1/sessions\n", config.Port)
	fmt.Printf("  Statistics:  http://localhost:%d/api/v1/stats\n", config.Port)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the server")
	fmt.Println()
}
