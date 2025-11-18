package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-learning/agent-orchestrator/messaging"
	"github.com/golang-learning/agent-orchestrator/models"
	"github.com/golang-learning/agent-orchestrator/tools"
)

// Agent represents a worker agent that processes requests
type Agent struct {
	id           string
	toolRegistry *tools.Registry
	stateManager *StateManager
	router       *messaging.Router
	ctx          context.Context
}

// NewAgent creates a new agent
func NewAgent(id string, toolRegistry *tools.Registry, stateManager *StateManager, router *messaging.Router, ctx context.Context) *Agent {
	return &Agent{
		id:           id,
		toolRegistry: toolRegistry,
		stateManager: stateManager,
		router:       router,
		ctx:          ctx,
	}
}

// Start starts the agent's processing loop
func (a *Agent) Start() {
	// Register agent in state manager
	a.stateManager.Register(a.id)
	log.Printf("Agent %s: Started and ready to process requests", a.id)

	// Main processing loop
	for {
		select {
		case req := <-a.router.GetRequestChannel():
			a.processRequest(req)
		case <-a.ctx.Done():
			log.Printf("Agent %s: Shutting down", a.id)
			a.stateManager.Unregister(a.id)
			return
		}
	}
}

// processRequest processes a single request
func (a *Agent) processRequest(req *models.Request) {
	startTime := time.Now()

	// Update state to busy
	a.stateManager.SetState(a.id, models.StateBusy)
	a.stateManager.SetCurrentRequest(a.id, req.ID)

	log.Printf("Agent %s: Processing request %s (tool: %s)", a.id, req.ID, req.ToolName)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(a.ctx, req.Timeout)
	defer cancel()

	// Get the tool
	tool, err := a.toolRegistry.Get(req.ToolName)
	if err != nil {
		a.sendErrorResponse(req.ID, err, time.Since(startTime))
		return
	}

	// Execute the tool
	result, err := tool.Execute(ctx, req.Params)

	duration := time.Since(startTime)

	if err != nil {
		a.sendErrorResponse(req.ID, err, duration)
	} else {
		a.sendSuccessResponse(req.ID, result, duration)
	}

	// Update state back to idle
	a.stateManager.SetState(a.id, models.StateIdle)
	a.stateManager.SetCurrentRequest(a.id, "")
	a.stateManager.IncrementRequestsHandled(a.id)

	log.Printf("Agent %s: Completed request %s in %v", a.id, req.ID, duration)
}

// sendSuccessResponse sends a successful response
func (a *Agent) sendSuccessResponse(requestID string, result interface{}, duration time.Duration) {
	resp := models.NewResponse(requestID, a.id, true, result, nil, duration)
	a.router.SendResponse(resp)
}

// sendErrorResponse sends an error response
func (a *Agent) sendErrorResponse(requestID string, err error, duration time.Duration) {
	log.Printf("Agent %s: Error processing request %s: %v", a.id, requestID, err)
	resp := models.NewResponse(requestID, a.id, false, nil, err, duration)
	a.router.SendResponse(resp)
	a.stateManager.SetState(a.id, models.StateIdle) // Reset to idle on error
}

// ID returns the agent ID
func (a *Agent) ID() string {
	return a.id
}

// GetInfo returns the agent's current information
func (a *Agent) GetInfo() (*models.AgentInfo, bool) {
	return a.stateManager.GetInfo(a.id)
}
