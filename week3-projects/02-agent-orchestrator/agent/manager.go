package agent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang-learning/agent-orchestrator/config"
	"github.com/golang-learning/agent-orchestrator/messaging"
	"github.com/golang-learning/agent-orchestrator/models"
	"github.com/golang-learning/agent-orchestrator/tools"
)

// Manager manages a pool of agents
type Manager struct {
	agents       []*Agent
	toolRegistry *tools.Registry
	stateManager *StateManager
	router       *messaging.Router
	config       *config.Config
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	startTime    time.Time

	// Statistics
	statsLock          sync.RWMutex
	totalRequests      int64
	successfulRequests int64
	failedRequests     int64
	totalProcessTime   time.Duration
}

// NewManager creates a new agent manager
func NewManager(cfg *config.Config, toolRegistry *tools.Registry) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		agents:       make([]*Agent, 0, cfg.NumAgents),
		toolRegistry: toolRegistry,
		stateManager: NewStateManager(),
		router:       messaging.NewRouter(ctx, cfg.RequestQueueSize),
		config:       cfg,
		ctx:          ctx,
		cancel:       cancel,
		startTime:    time.Now(),
	}
}

// Start starts the agent manager and all agents
func (m *Manager) Start() error {
	log.Printf("Starting Agent Manager with %d agents", m.config.NumAgents)

	// Start the router
	m.router.Start()

	// Create and start agents
	for i := 0; i < m.config.NumAgents; i++ {
		agentID := fmt.Sprintf("agent-%d", i+1)
		agent := NewAgent(agentID, m.toolRegistry, m.stateManager, m.router, m.ctx)
		m.agents = append(m.agents, agent)

		// Start agent in goroutine
		m.wg.Add(1)
		go func(a *Agent) {
			defer m.wg.Done()
			a.Start()
		}(agent)
	}

	log.Printf("Agent Manager started successfully with %d agents", len(m.agents))
	return nil
}

// Stop stops the agent manager and all agents
func (m *Manager) Stop() {
	log.Println("Stopping Agent Manager...")

	// Cancel context to signal all agents to stop
	m.cancel()

	// Wait for all agents to finish
	m.wg.Wait()

	// Stop the router
	m.router.Stop()

	log.Println("Agent Manager stopped")
}

// SubmitRequest submits a request to be processed by an agent
func (m *Manager) SubmitRequest(req *models.Request) (*models.Response, error) {
	// Update statistics
	m.statsLock.Lock()
	m.totalRequests++
	m.statsLock.Unlock()

	// Submit request through router
	resp, err := m.router.SubmitRequest(req)

	// Update statistics based on response
	if resp != nil {
		m.statsLock.Lock()
		if resp.Success {
			m.successfulRequests++
		} else {
			m.failedRequests++
		}
		m.totalProcessTime += resp.Duration
		m.statsLock.Unlock()
	}

	return resp, err
}

// GetAgentInfo returns information about a specific agent
func (m *Manager) GetAgentInfo(agentID string) (*models.AgentInfo, error) {
	info, exists := m.stateManager.GetInfo(agentID)
	if !exists {
		return nil, fmt.Errorf("agent %s not found", agentID)
	}
	return info, nil
}

// GetAllAgentInfo returns information about all agents
func (m *Manager) GetAllAgentInfo() []*models.AgentInfo {
	return m.stateManager.GetAllInfo()
}

// GetStatistics returns overall system statistics
func (m *Manager) GetStatistics() *models.Statistics {
	m.statsLock.RLock()
	defer m.statsLock.RUnlock()

	stateCounts := m.stateManager.GetStateCounts()

	avgProcessTime := time.Duration(0)
	if m.successfulRequests > 0 {
		avgProcessTime = m.totalProcessTime / time.Duration(m.successfulRequests)
	}

	return &models.Statistics{
		TotalRequests:      m.totalRequests,
		SuccessfulRequests: m.successfulRequests,
		FailedRequests:     m.failedRequests,
		AverageProcessTime: avgProcessTime,
		ActiveAgents:       m.stateManager.Count(),
		IdleAgents:         stateCounts[models.StateIdle],
		BusyAgents:         stateCounts[models.StateBusy],
		ErrorAgents:        stateCounts[models.StateError],
		Uptime:             time.Since(m.startTime),
	}
}

// HealthCheck performs a health check on the system
func (m *Manager) HealthCheck() bool {
	// Check if context is still active
	select {
	case <-m.ctx.Done():
		return false
	default:
		// Check if we have any active agents
		return m.stateManager.Count() > 0
	}
}

// GetToolRegistry returns the tool registry
func (m *Manager) GetToolRegistry() *tools.Registry {
	return m.toolRegistry
}
