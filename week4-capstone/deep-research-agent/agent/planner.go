package agent

import (
	"fmt"
	"strings"
	"time"

	"deep-research-agent/models"
	"deep-research-agent/tools"
)

// ResearchPlanner creates intelligent research plans
type ResearchPlanner struct {
	toolRegistry *tools.ToolRegistry
	config       models.AgentConfig
}

// NewResearchPlanner creates a new research planner
func NewResearchPlanner(registry *tools.ToolRegistry, config models.AgentConfig) *ResearchPlanner {
	return &ResearchPlanner{
		toolRegistry: registry,
		config:       config,
	}
}

// PlanResearch analyzes a query and creates a research plan
func (p *ResearchPlanner) PlanResearch(query models.ResearchQuery) (*models.ResearchPlan, error) {
	// Determine research type
	researchType := p.determineResearchType(query)

	// Create steps based on research type
	steps := p.createSteps(query, researchType)

	// Estimate execution time
	estimatedTime := p.estimateTime(steps)

	plan := &models.ResearchPlan{
		Query:         query,
		Steps:         steps,
		EstimatedTime: estimatedTime,
		ResearchType:  researchType,
		Strategy:      p.describeStrategy(query, researchType),
		CreatedAt:     time.Now(),
	}

	return plan, nil
}

// determineResearchType analyzes the query to determine the research type
func (p *ResearchPlanner) determineResearchType(query models.ResearchQuery) models.ResearchType {
	hasDocuments := len(query.Documents) > 0
	hasWeb := query.UseWeb
	hasWiki := query.UseWiki

	// Check for academic/scientific keywords
	academicKeywords := []string{"research", "study", "paper", "academic", "scientific", "analysis", "experiment"}
	queryLower := strings.ToLower(query.Query)
	isAcademic := false
	for _, keyword := range academicKeywords {
		if strings.Contains(queryLower, keyword) {
			isAcademic = true
			break
		}
	}

	// Determine type based on inputs
	if hasDocuments && (hasWeb || hasWiki) {
		return models.ResearchMulti
	} else if hasDocuments {
		return models.ResearchDocument
	} else if isAcademic {
		return models.ResearchAcademic
	} else {
		return models.ResearchGeneral
	}
}

// createSteps generates research steps based on query and type
func (p *ResearchPlanner) createSteps(query models.ResearchQuery, researchType models.ResearchType) []models.ResearchStep {
	steps := make([]models.ResearchStep, 0)
	stepNumber := 1

	switch researchType {
	case models.ResearchDocument:
		steps = p.createDocumentSteps(query, stepNumber)

	case models.ResearchAcademic:
		steps = p.createAcademicSteps(query, stepNumber)

	case models.ResearchMulti:
		steps = p.createMultiSourceSteps(query, stepNumber)

	case models.ResearchGeneral:
		steps = p.createGeneralSteps(query, stepNumber)
	}

	// Add depth-based additional steps
	if query.Depth == models.ResearchDeep {
		steps = append(steps, p.createDeepResearchSteps(query, len(steps)+1)...)
	}

	// Limit to max steps
	maxSteps := query.MaxSteps
	if maxSteps == 0 {
		maxSteps = p.config.MaxSteps
	}
	if len(steps) > maxSteps {
		steps = steps[:maxSteps]
	}

	return steps
}

// createDocumentSteps creates steps for document analysis
func (p *ResearchPlanner) createDocumentSteps(query models.ResearchQuery, startNum int) []models.ResearchStep {
	steps := make([]models.ResearchStep, 0)

	for i, docPath := range query.Documents {
		// Determine document type
		toolName := "pdf_processor"
		if strings.HasSuffix(strings.ToLower(docPath), ".docx") ||
		   strings.HasSuffix(strings.ToLower(docPath), ".doc") {
			toolName = "docx_processor"
		}

		steps = append(steps, models.ResearchStep{
			StepNumber: startNum + i,
			Tool:       toolName,
			Action:     "analyze_document",
			Parameters: map[string]interface{}{
				"file_path": docPath,
				"query":     query.Query,
			},
			Reasoning:  fmt.Sprintf("Analyze document %s for relevant information", docPath),
			DependsOn:  []int{},
			Priority:   1,
			MaxRetries: 2,
		})
	}

	return steps
}

// createAcademicSteps creates steps for academic research
func (p *ResearchPlanner) createAcademicSteps(query models.ResearchQuery, startNum int) []models.ResearchStep {
	steps := make([]models.ResearchStep, 0)

	// Step 1: Wikipedia for background
	if query.UseWiki && p.toolRegistry.HasTool("wikipedia") {
		steps = append(steps, models.ResearchStep{
			StepNumber: startNum,
			Tool:       "wikipedia",
			Action:     "search",
			Parameters: map[string]interface{}{
				"query":       query.Query,
				"max_results": 3,
			},
			Reasoning:  "Get comprehensive background from Wikipedia",
			DependsOn:  []int{},
			Priority:   1,
			MaxRetries: 2,
		})
		startNum++
	}

	// Step 2: Web search for academic sources
	if query.UseWeb && p.toolRegistry.HasTool("web_search") {
		steps = append(steps, models.ResearchStep{
			StepNumber: startNum,
			Tool:       "web_search",
			Action:     "search",
			Parameters: map[string]interface{}{
				"query":       query.Query,
				"max_results": 5,
			},
			Reasoning:  "Search web for academic papers and research",
			DependsOn:  []int{},
			Priority:   1,
			MaxRetries: 2,
		})
	}

	return steps
}

// createGeneralSteps creates steps for general research
func (p *ResearchPlanner) createGeneralSteps(query models.ResearchQuery, startNum int) []models.ResearchStep {
	steps := make([]models.ResearchStep, 0)

	// Step 1: Web search
	if query.UseWeb && p.toolRegistry.HasTool("web_search") {
		steps = append(steps, models.ResearchStep{
			StepNumber: startNum,
			Tool:       "web_search",
			Action:     "search",
			Parameters: map[string]interface{}{
				"query":       query.Query,
				"max_results": 5,
			},
			Reasoning:  "Search web for current information",
			DependsOn:  []int{},
			Priority:   1,
			MaxRetries: 2,
		})
		startNum++
	}

	// Step 2: Wikipedia for context
	if query.UseWiki && p.toolRegistry.HasTool("wikipedia") {
		steps = append(steps, models.ResearchStep{
			StepNumber: startNum,
			Tool:       "wikipedia",
			Action:     "search",
			Parameters: map[string]interface{}{
				"query":       query.Query,
				"max_results": 2,
			},
			Reasoning:  "Get background information from Wikipedia",
			DependsOn:  []int{},
			Priority:   2,
			MaxRetries: 2,
		})
	}

	return steps
}

// createMultiSourceSteps creates steps for multi-source research
func (p *ResearchPlanner) createMultiSourceSteps(query models.ResearchQuery, startNum int) []models.ResearchStep {
	steps := make([]models.ResearchStep, 0)

	// Phase 1: Analyze documents (sequential, high priority)
	docSteps := p.createDocumentSteps(query, startNum)
	steps = append(steps, docSteps...)
	startNum += len(docSteps)

	// Phase 2: Web research (parallel, can run concurrently)
	_ = startNum // webStep for future use
	if query.UseWeb && p.toolRegistry.HasTool("web_search") {
		steps = append(steps, models.ResearchStep{
			StepNumber: startNum,
			Tool:       "web_search",
			Action:     "search",
			Parameters: map[string]interface{}{
				"query":       query.Query,
				"max_results": 5,
			},
			Reasoning:  "Search web to supplement document findings",
			DependsOn:  []int{},
			Priority:   2,
			MaxRetries: 2,
		})
		startNum++
	}

	// Phase 3: Wikipedia (parallel with web search)
	if query.UseWiki && p.toolRegistry.HasTool("wikipedia") {
		steps = append(steps, models.ResearchStep{
			StepNumber: startNum,
			Tool:       "wikipedia",
			Action:     "search",
			Parameters: map[string]interface{}{
				"query":       query.Query,
				"max_results": 3,
			},
			Reasoning:  "Get encyclopedic knowledge to provide context",
			DependsOn:  []int{},
			Priority:   2,
			MaxRetries: 2,
		})
	}

	return steps
}

// createDeepResearchSteps creates additional steps for deep research
func (p *ResearchPlanner) createDeepResearchSteps(query models.ResearchQuery, startNum int) []models.ResearchStep {
	steps := make([]models.ResearchStep, 0)

	// For deep research, add cross-referencing and validation steps
	// These depend on earlier steps completing

	// Note: In a real implementation, we might add steps like:
	// - Follow up searches based on initial findings
	// - Cross-reference facts across sources
	// - Deeper analysis of specific topics discovered

	// For now, we'll add an additional web search with refined query
	if query.UseWeb && p.toolRegistry.HasTool("web_search") {
		steps = append(steps, models.ResearchStep{
			StepNumber: startNum,
			Tool:       "web_search",
			Action:     "search",
			Parameters: map[string]interface{}{
				"query":       query.Query + " detailed analysis",
				"max_results": 5,
			},
			Reasoning:  "Deep dive search for comprehensive coverage",
			DependsOn:  []int{1}, // Depends on first step
			Priority:   3,
			MaxRetries: 2,
		})
	}

	return steps
}

// describeStrategy generates a description of the research strategy
func (p *ResearchPlanner) describeStrategy(query models.ResearchQuery, researchType models.ResearchType) string {
	var strategy strings.Builder

	switch researchType {
	case models.ResearchDocument:
		strategy.WriteString("Document-focused research: Analyzing provided documents for relevant information. ")

	case models.ResearchAcademic:
		strategy.WriteString("Academic research: Gathering scholarly information from encyclopedic and web sources. ")

	case models.ResearchMulti:
		strategy.WriteString("Multi-source research: Combining document analysis with web and encyclopedic sources. ")

	case models.ResearchGeneral:
		strategy.WriteString("General research: Gathering information from various online sources. ")
	}

	switch query.Depth {
	case models.ResearchShallow:
		strategy.WriteString("Quick overview with essential information only.")
	case models.ResearchMedium:
		strategy.WriteString("Standard research depth with key sources.")
	case models.ResearchDeep:
		strategy.WriteString("Comprehensive deep-dive with extensive cross-referencing.")
	}

	return strategy.String()
}

// estimateTime estimates how long the research will take
func (p *ResearchPlanner) estimateTime(steps []models.ResearchStep) time.Duration {
	// Simple estimation: each step takes average 10 seconds
	// Adjust based on tool type
	totalSeconds := 0

	for _, step := range steps {
		switch step.Tool {
		case "pdf_processor", "docx_processor":
			totalSeconds += 15 // Document processing is slower
		case "web_search":
			totalSeconds += 8
		case "wikipedia":
			totalSeconds += 5
		default:
			totalSeconds += 10
		}
	}

	// Add some overhead for orchestration
	totalSeconds += 5

	return time.Duration(totalSeconds) * time.Second
}

// OptimizeSteps optimizes the execution order of steps
func (p *ResearchPlanner) OptimizeSteps(steps []models.ResearchStep) []models.ResearchStep {
	// Group by priority and dependencies
	// Steps with same priority and no dependencies can run in parallel

	// For now, just return steps as-is
	// In a more advanced implementation, we could:
	// 1. Sort by priority
	// 2. Identify parallelizable steps
	// 3. Reorder for optimal execution

	return steps
}
