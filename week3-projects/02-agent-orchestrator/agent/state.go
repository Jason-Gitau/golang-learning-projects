package agent

import (
	"sync"
	"time"

	"github.com/golang-learning/agent-orchestrator/models"
)

// StateManager manages agent states thread-safely
type StateManager struct {
	states map[string]*models.AgentInfo
	mu     sync.RWMutex
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[string]*models.AgentInfo),
	}
}

// Register registers a new agent
func (sm *StateManager) Register(agentID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.states[agentID] = &models.AgentInfo{
		ID:              agentID,
		State:           models.StateIdle,
		RequestsHandled: 0,
		LastActive:      time.Now(),
	}
}

// Unregister removes an agent
func (sm *StateManager) Unregister(agentID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.states, agentID)
}

// SetState sets the state of an agent
func (sm *StateManager) SetState(agentID string, state models.AgentState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if info, exists := sm.states[agentID]; exists {
		info.State = state
		info.LastActive = time.Now()
	}
}

// SetCurrentRequest sets the current request being processed by an agent
func (sm *StateManager) SetCurrentRequest(agentID, requestID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if info, exists := sm.states[agentID]; exists {
		info.CurrentRequest = requestID
		info.LastActive = time.Now()
	}
}

// IncrementRequestsHandled increments the requests handled counter
func (sm *StateManager) IncrementRequestsHandled(agentID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if info, exists := sm.states[agentID]; exists {
		info.RequestsHandled++
		info.LastActive = time.Now()
	}
}

// GetState returns the state of an agent
func (sm *StateManager) GetState(agentID string) (models.AgentState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if info, exists := sm.states[agentID]; exists {
		return info.State, true
	}

	return "", false
}

// GetInfo returns information about an agent
func (sm *StateManager) GetInfo(agentID string) (*models.AgentInfo, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if info, exists := sm.states[agentID]; exists {
		// Return a copy to avoid race conditions
		infoCopy := *info
		return &infoCopy, true
	}

	return nil, false
}

// GetAllInfo returns information about all agents
func (sm *StateManager) GetAllInfo() []*models.AgentInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	infos := make([]*models.AgentInfo, 0, len(sm.states))
	for _, info := range sm.states {
		// Return copies to avoid race conditions
		infoCopy := *info
		infos = append(infos, &infoCopy)
	}

	return infos
}

// GetStateCounts returns counts of agents in each state
func (sm *StateManager) GetStateCounts() map[models.AgentState]int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	counts := map[models.AgentState]int{
		models.StateIdle:  0,
		models.StateBusy:  0,
		models.StateError: 0,
	}

	for _, info := range sm.states {
		counts[info.State]++
	}

	return counts
}

// Count returns the total number of registered agents
func (sm *StateManager) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.states)
}
