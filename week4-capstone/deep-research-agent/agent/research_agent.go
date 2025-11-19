package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"deep-research-agent/models"
	"deep-research-agent/storage"
	"deep-research-agent/tools"
)

// ResearchAgent is the main research orchestration engine
type ResearchAgent struct {
	ID           string
	Config       models.AgentConfig
	ToolRegistry *tools.ToolRegistry
	Storage      *storage.ResearchDB
	Metrics      *models.AgentMetrics
	Status       models.AgentStatus
}

// NewResearchAgent creates a new research agent
func NewResearchAgent(
	config models.AgentConfig,
	registry *tools.ToolRegistry,
	db *storage.ResearchDB,
) *ResearchAgent {
	return &ResearchAgent{
		ID:           uuid.New().String(),
		Config:       config,
		ToolRegistry: registry,
		Storage:      db,
		Metrics:      models.NewAgentMetrics(),
		Status:       models.AgentIdle,
	}
}

// Research performs a complete research operation
func (a *ResearchAgent) Research(
	query models.ResearchQuery,
	options models.ResearchOptions,
) (*models.ResearchResult, error) {
	// Set defaults
	if query.Depth == "" {
		query.Depth = models.ResearchMedium
	}
	if query.MaxSteps == 0 {
		query.MaxSteps = a.Config.MaxSteps
	}

	// Create session ID
	sessionID := options.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// Update agent status
	a.Status = models.AgentPlanning
	startTime := time.Now()

	// Step 1: Create research plan
	log.Printf("Agent %s: Planning research for query: %s", a.ID, query.Query)
	planner := NewResearchPlanner(a.ToolRegistry, a.Config)
	plan, err := planner.PlanResearch(query)
	if err != nil {
		a.Status = models.AgentFailed
		return nil, fmt.Errorf("failed to create research plan: %w", err)
	}

	log.Printf("Agent %s: Created plan with %d steps (estimated time: %s)",
		a.ID, len(plan.Steps), plan.EstimatedTime)

	// Step 2: Initialize memory and orchestrator
	memory := NewResearchMemory(sessionID, query.Query)

	progressChan := options.ProgressChannel
	if progressChan == nil && options.StreamProgress {
		progressChan = make(chan models.ResearchProgress, 10)
		defer close(progressChan)
	}

	orchestrator := NewResearchOrchestrator(
		a.ToolRegistry,
		memory,
		a.Config,
		progressChan,
	)

	// Save initial session to database if requested
	if options.SaveSession && a.Storage != nil {
		session, err := a.Storage.SaveSession("", query, nil, "in_progress")
		if err != nil {
			log.Printf("Warning: Failed to save initial session: %v", err)
		} else {
			sessionID = session.ID
			memory = NewResearchMemory(sessionID, query.Query)
			orchestrator = NewResearchOrchestrator(
				a.ToolRegistry,
				memory,
				a.Config,
				progressChan,
			)
		}
	}

	// Step 3: Execute research plan
	a.Status = models.AgentResearching
	log.Printf("Agent %s: Executing research plan", a.ID)

	ctx := context.Background()
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	results, err := orchestrator.ExecutePlan(ctx, plan)
	if err != nil {
		a.Status = models.AgentFailed
		log.Printf("Agent %s: Execution failed: %v", a.ID, err)

		// Save failed session
		if options.SaveSession && a.Storage != nil {
			partialResult := orchestrator.AggregateResults(query, results, time.Since(startTime))
			_ = a.Storage.UpdateSession(sessionID, partialResult, "failed")
		}

		return nil, fmt.Errorf("failed to execute research plan: %w", err)
	}

	// Step 4: Aggregate results
	a.Status = models.AgentAggregating
	log.Printf("Agent %s: Aggregating results", a.ID)

	duration := time.Since(startTime)
	result := orchestrator.AggregateResults(query, results, duration)
	result.ResearchType = plan.ResearchType

	// Step 5: Save final session to database
	if options.SaveSession && a.Storage != nil {
		err := a.Storage.UpdateSession(sessionID, result, "completed")
		if err != nil {
			log.Printf("Warning: Failed to update session: %v", err)
		}
	}

	// Update agent status
	a.Status = models.AgentCompleted

	// Record metrics
	success := result.FailedSteps == 0
	a.Metrics.RecordResearch(result, success)

	log.Printf("Agent %s: Research completed in %s with %d/%d successful steps",
		a.ID, duration, result.SuccessfulSteps, result.TotalSteps)

	return result, nil
}

// PlanResearch creates a research plan without executing it
func (a *ResearchAgent) PlanResearch(query models.ResearchQuery) (*models.ResearchPlan, error) {
	planner := NewResearchPlanner(a.ToolRegistry, a.Config)
	return planner.PlanResearch(query)
}

// ExecuteStep executes a single research step
func (a *ResearchAgent) ExecuteStep(
	step models.ResearchStep,
	memory *ResearchMemory,
) (*models.StepResult, error) {
	orchestrator := NewResearchOrchestrator(
		a.ToolRegistry,
		memory,
		a.Config,
		nil,
	)

	ctx := context.Background()
	result := orchestrator.executeStep(ctx, step, 1, 1)

	return &result, nil
}

// GetMetrics returns current agent metrics
func (a *ResearchAgent) GetMetrics() models.AgentMetrics {
	return a.Metrics.GetMetrics()
}

// GetStatus returns the current agent status
func (a *ResearchAgent) GetStatus() models.AgentStatus {
	return a.Status
}

// ListSessions retrieves research sessions from the database
func (a *ResearchAgent) ListSessions(limit, offset int) ([]models.ResearchSession, error) {
	if a.Storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}
	return a.Storage.ListSessions(limit, offset)
}

// GetSession retrieves a specific research session
func (a *ResearchAgent) GetSession(sessionID string) (*models.ResearchSession, error) {
	if a.Storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}
	return a.Storage.GetSession(sessionID)
}

// DeleteSession deletes a research session
func (a *ResearchAgent) DeleteSession(sessionID string) error {
	if a.Storage == nil {
		return fmt.Errorf("storage not configured")
	}
	return a.Storage.DeleteSession(sessionID)
}

// IndexDocument adds a document to the index
func (a *ResearchAgent) IndexDocument(
	filename, filePath, fileType string,
	fileSize int64,
	pageCount int,
) (*models.Document, error) {
	if a.Storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}
	return a.Storage.IndexDocument("", filename, filePath, fileType, fileSize, pageCount)
}

// ListDocuments retrieves indexed documents
func (a *ResearchAgent) ListDocuments(limit, offset int) ([]models.Document, error) {
	if a.Storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}
	return a.Storage.ListDocuments(limit, offset)
}

// GetAvailableTools returns information about available tools
func (a *ResearchAgent) GetAvailableTools() map[string]tools.ToolInfo {
	return a.ToolRegistry.GetToolInfo()
}

// ValidateQuery checks if a research query is valid
func (a *ResearchAgent) ValidateQuery(query models.ResearchQuery) error {
	if query.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	// Check if any tools are available
	if !query.UseWeb && !query.UseWiki && len(query.Documents) == 0 {
		return fmt.Errorf("at least one research source must be enabled")
	}

	// Validate documents exist
	for _, docPath := range query.Documents {
		// In a real implementation, check if file exists
		if docPath == "" {
			return fmt.Errorf("document path cannot be empty")
		}
	}

	// Validate depth
	validDepths := map[models.ResearchDepth]bool{
		models.ResearchShallow: true,
		models.ResearchMedium:  true,
		models.ResearchDeep:    true,
	}
	if query.Depth != "" && !validDepths[query.Depth] {
		return fmt.Errorf("invalid research depth: %s", query.Depth)
	}

	return nil
}

// Close closes the agent and releases resources
func (a *ResearchAgent) Close() error {
	if a.Storage != nil {
		return a.Storage.Close()
	}
	return nil
}

// String returns a string representation of the agent
func (a *ResearchAgent) String() string {
	return fmt.Sprintf("ResearchAgent{ID: %s, Status: %s, Tools: %d}",
		a.ID, a.Status, len(a.ToolRegistry.ListToolNames()))
}
