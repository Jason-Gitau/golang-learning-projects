package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-learning/agent-orchestrator/agent"
	"github.com/golang-learning/agent-orchestrator/api"
	"github.com/golang-learning/agent-orchestrator/config"
	"github.com/golang-learning/agent-orchestrator/tools"
)

func main() {
	log.Println("=== Agent Orchestrator ===")
	log.Println("Starting up...")

	// Load configuration
	cfg := config.Default()

	// Override with environment variables if present
	if port := os.Getenv("PORT"); port != "" {
		cfg.WithPort(port)
	}
	if weatherKey := os.Getenv("WEATHER_API_KEY"); weatherKey != "" {
		cfg.WithWeatherAPIKey(weatherKey)
	}

	// Create tool registry
	registry := tools.NewRegistry()

	// Register tools
	if err := registerTools(registry, cfg); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	log.Printf("Registered %d tools: %v", registry.Count(), registry.ListNames())

	// Create agent manager
	manager := agent.NewManager(cfg, registry)

	// Start agent manager
	if err := manager.Start(); err != nil {
		log.Fatalf("Failed to start agent manager: %v", err)
	}

	// Create and start API server
	server := api.NewServer(cfg, manager)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop API server
	if err := server.Stop(shutdownCtx); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	// Stop agent manager
	manager.Stop()

	log.Println("Agent Orchestrator shut down successfully")
}

// registerTools registers all available tools
func registerTools(registry *tools.Registry, cfg *config.Config) error {
	// Calculator tool
	if err := registry.Register(tools.NewCalculatorTool()); err != nil {
		return err
	}

	// Time tool
	if err := registry.Register(tools.NewTimeTool()); err != nil {
		return err
	}

	// Random tool
	if err := registry.Register(tools.NewRandomTool()); err != nil {
		return err
	}

	// Weather tool
	if err := registry.Register(tools.NewWeatherTool(cfg.WeatherAPIKey)); err != nil {
		return err
	}

	// Text tool
	if err := registry.Register(tools.NewTextTool()); err != nil {
		return err
	}

	return nil
}
