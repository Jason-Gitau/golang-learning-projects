package agent

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"deep-research-agent/models"
)

// ResearchMemory manages context and findings during research
type ResearchMemory struct {
	sessionID       string
	query           string
	findings        []models.Finding
	sources         []models.Source
	visitedSources  map[string]bool // Track visited URLs/documents
	stepHistory     []models.StepResult
	contextData     map[string]interface{}
	createdAt       time.Time
	mu              sync.RWMutex
}

// NewResearchMemory creates a new research memory instance
func NewResearchMemory(sessionID, query string) *ResearchMemory {
	return &ResearchMemory{
		sessionID:      sessionID,
		query:          query,
		findings:       make([]models.Finding, 0),
		sources:        make([]models.Source, 0),
		visitedSources: make(map[string]bool),
		stepHistory:    make([]models.StepResult, 0),
		contextData:    make(map[string]interface{}),
		createdAt:      time.Now(),
	}
}

// AddFinding adds a new finding to memory
func (m *ResearchMemory) AddFinding(finding models.Finding) {
	m.mu.Lock()
	defer m.mu.Unlock()

	finding.Timestamp = time.Now()
	m.findings = append(m.findings, finding)

	// Also add sources from this finding
	for _, source := range finding.Sources {
		m.addSourceUnsafe(source)
	}
}

// AddFindings adds multiple findings to memory
func (m *ResearchMemory) AddFindings(findings []models.Finding) {
	for _, finding := range findings {
		m.AddFinding(finding)
	}
}

// AddSource adds a source to memory
func (m *ResearchMemory) AddSource(source models.Source) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.addSourceUnsafe(source)
}

// addSourceUnsafe adds a source without locking (internal use)
func (m *ResearchMemory) addSourceUnsafe(source models.Source) {
	source.Timestamp = time.Now()

	// Check if we've already visited this source
	key := m.getSourceKey(source)
	if !m.visitedSources[key] {
		m.sources = append(m.sources, source)
		m.visitedSources[key] = true
	}
}

// getSourceKey generates a unique key for a source
func (m *ResearchMemory) getSourceKey(source models.Source) string {
	if source.URL != "" {
		return source.URL
	}
	// For documents, use filename + page number
	return fmt.Sprintf("%s:page%d", source.Title, source.PageNumber)
}

// HasVisitedSource checks if a source has been visited
func (m *ResearchMemory) HasVisitedSource(sourceKey string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.visitedSources[sourceKey]
}

// AddStepResult records the result of a research step
func (m *ResearchMemory) AddStepResult(result models.StepResult) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stepHistory = append(m.stepHistory, result)

	// Extract sources from step result
	for _, source := range result.Sources {
		m.addSourceUnsafe(source)
	}
}

// GetFindings returns all findings
func (m *ResearchMemory) GetFindings() []models.Finding {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy
	findings := make([]models.Finding, len(m.findings))
	copy(findings, m.findings)
	return findings
}

// GetSources returns all sources
func (m *ResearchMemory) GetSources() []models.Source {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy
	sources := make([]models.Source, len(m.sources))
	copy(sources, m.sources)
	return sources
}

// GetStepHistory returns the history of research steps
func (m *ResearchMemory) GetStepHistory() []models.StepResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy
	history := make([]models.StepResult, len(m.stepHistory))
	copy(history, m.stepHistory)
	return history
}

// GetSuccessfulSteps returns count of successful steps
func (m *ResearchMemory) GetSuccessfulSteps() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, step := range m.stepHistory {
		if step.Success {
			count++
		}
	}
	return count
}

// GetFailedSteps returns count of failed steps
func (m *ResearchMemory) GetFailedSteps() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, step := range m.stepHistory {
		if !step.Success {
			count++
		}
	}
	return count
}

// SetContextData stores arbitrary context data
func (m *ResearchMemory) SetContextData(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.contextData[key] = value
}

// GetContextData retrieves context data
func (m *ResearchMemory) GetContextData(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.contextData[key]
	return value, exists
}

// QueryRelatedFindings searches for findings related to a query
func (m *ResearchMemory) QueryRelatedFindings(queryTerms []string) []models.Finding {
	m.mu.RLock()
	defer m.mu.RUnlock()

	related := make([]models.Finding, 0)

	for _, finding := range m.findings {
		// Simple keyword matching - check if finding contains any query terms
		contentLower := strings.ToLower(finding.Content)
		for _, term := range queryTerms {
			if strings.Contains(contentLower, strings.ToLower(term)) {
				related = append(related, finding)
				break
			}
		}
	}

	return related
}

// GetSourcesByType returns sources filtered by type
func (m *ResearchMemory) GetSourcesByType(sourceType string) []models.Source {
	m.mu.RLock()
	defer m.mu.RUnlock()

	filtered := make([]models.Source, 0)
	for _, source := range m.sources {
		if source.Type == sourceType {
			filtered = append(filtered, source)
		}
	}

	return filtered
}

// GetTopSources returns the N most relevant sources
func (m *ResearchMemory) GetTopSources(n int) []models.Source {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Sort by relevance and return top N
	sources := make([]models.Source, len(m.sources))
	copy(sources, m.sources)

	// Simple sort by relevance (bubble sort for small lists)
	for i := 0; i < len(sources)-1; i++ {
		for j := 0; j < len(sources)-i-1; j++ {
			if sources[j].Relevance < sources[j+1].Relevance {
				sources[j], sources[j+1] = sources[j+1], sources[j]
			}
		}
	}

	if n > len(sources) {
		n = len(sources)
	}

	return sources[:n]
}

// GetSummary generates a summary of the research memory
func (m *ResearchMemory) GetSummary() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	summary := fmt.Sprintf("Research Session: %s\n", m.sessionID)
	summary += fmt.Sprintf("Query: %s\n", m.query)
	summary += fmt.Sprintf("Total Findings: %d\n", len(m.findings))
	summary += fmt.Sprintf("Total Sources: %d\n", len(m.sources))
	summary += fmt.Sprintf("Steps Executed: %d\n", len(m.stepHistory))
	summary += fmt.Sprintf("Successful Steps: %d\n", m.GetSuccessfulSteps())
	summary += fmt.Sprintf("Failed Steps: %d\n", m.GetFailedSteps())

	return summary
}

// Clear clears all memory
func (m *ResearchMemory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.findings = make([]models.Finding, 0)
	m.sources = make([]models.Source, 0)
	m.visitedSources = make(map[string]bool)
	m.stepHistory = make([]models.StepResult, 0)
	m.contextData = make(map[string]interface{})
}

// GetSessionID returns the session ID
func (m *ResearchMemory) GetSessionID() string {
	return m.sessionID
}

// GetQuery returns the research query
func (m *ResearchMemory) GetQuery() string {
	return m.query
}

// GetTotalSteps returns total number of steps executed
func (m *ResearchMemory) GetTotalSteps() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.stepHistory)
}

// GetAverageConfidence calculates average confidence across all findings
func (m *ResearchMemory) GetAverageConfidence() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.findings) == 0 {
		return 0.0
	}

	totalConfidence := 0.0
	for _, finding := range m.findings {
		totalConfidence += finding.Confidence
	}

	return totalConfidence / float64(len(m.findings))
}

// GetSourceCount returns the total number of unique sources
func (m *ResearchMemory) GetSourceCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.sources)
}

// GetDuration returns how long the research session has been running
func (m *ResearchMemory) GetDuration() time.Duration {
	return time.Since(m.createdAt)
}
