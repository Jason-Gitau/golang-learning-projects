package agent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"deep-research-agent/models"
	"deep-research-agent/tools"
)

// ResearchOrchestrator manages the execution of research steps
type ResearchOrchestrator struct {
	toolRegistry *tools.ToolRegistry
	memory       *ResearchMemory
	config       models.AgentConfig
	progressChan chan models.ResearchProgress
	cancelFunc   context.CancelFunc
}

// NewResearchOrchestrator creates a new research orchestrator
func NewResearchOrchestrator(
	registry *tools.ToolRegistry,
	memory *ResearchMemory,
	config models.AgentConfig,
	progressChan chan models.ResearchProgress,
) *ResearchOrchestrator {
	return &ResearchOrchestrator{
		toolRegistry: registry,
		memory:       memory,
		config:       config,
		progressChan: progressChan,
	}
}

// ExecutePlan executes a research plan
func (o *ResearchOrchestrator) ExecutePlan(ctx context.Context, plan *models.ResearchPlan) ([]models.StepResult, error) {
	// Create cancellable context
	execCtx, cancel := context.WithTimeout(ctx, o.config.Timeout)
	o.cancelFunc = cancel
	defer cancel()

	results := make([]models.StepResult, 0, len(plan.Steps))

	// Group steps by priority and dependencies
	stepGroups := o.groupStepsByPriority(plan.Steps)

	stepNum := 0
	totalSteps := len(plan.Steps)

	// Execute step groups
	for _, group := range stepGroups {
		// Check if context is cancelled
		select {
		case <-execCtx.Done():
			return results, fmt.Errorf("execution cancelled: %w", execCtx.Err())
		default:
		}

		// Execute steps in this group (potentially in parallel)
		groupResults, err := o.executeStepGroup(execCtx, group, totalSteps, &stepNum)
		results = append(results, groupResults...)

		if err != nil {
			log.Printf("Error executing step group: %v", err)
			// Continue with other groups even if one fails
		}
	}

	return results, nil
}

// groupStepsByPriority groups steps by priority for execution
func (o *ResearchOrchestrator) groupStepsByPriority(steps []models.ResearchStep) [][]models.ResearchStep {
	// Group steps by priority
	priorityMap := make(map[int][]models.ResearchStep)
	maxPriority := 0

	for _, step := range steps {
		priority := step.Priority
		if priority == 0 {
			priority = 1
		}
		priorityMap[priority] = append(priorityMap[priority], step)
		if priority > maxPriority {
			maxPriority = priority
		}
	}

	// Convert to ordered slice
	groups := make([][]models.ResearchStep, 0, maxPriority)
	for i := 1; i <= maxPriority; i++ {
		if steps, exists := priorityMap[i]; exists {
			groups = append(groups, steps)
		}
	}

	return groups
}

// executeStepGroup executes a group of steps (potentially in parallel)
func (o *ResearchOrchestrator) executeStepGroup(
	ctx context.Context,
	steps []models.ResearchStep,
	totalSteps int,
	currentStep *int,
) ([]models.StepResult, error) {
	results := make([]models.StepResult, 0, len(steps))
	resultsMu := sync.Mutex{}

	// Determine parallelism
	// Steps with no dependencies on each other can run in parallel
	canParallelize := o.canParallelizeGroup(steps)

	if canParallelize && len(steps) > 1 {
		// Execute in parallel with worker pool
		results = o.executeParallel(ctx, steps, totalSteps, currentStep)
	} else {
		// Execute sequentially
		for _, step := range steps {
			*currentStep++
			result := o.executeStep(ctx, step, *currentStep, totalSteps)
			results = append(results, result)

			// Add to memory
			o.memory.AddStepResult(result)
		}
	}

	resultsMu.Lock()
	defer resultsMu.Unlock()

	return results, nil
}

// canParallelizeGroup checks if steps in a group can be parallelized
func (o *ResearchOrchestrator) canParallelizeGroup(steps []models.ResearchStep) bool {
	// Check if any step depends on another step in the group
	stepNumbers := make(map[int]bool)
	for _, step := range steps {
		stepNumbers[step.StepNumber] = true
	}

	for _, step := range steps {
		for _, dep := range step.DependsOn {
			if stepNumbers[dep] {
				return false // Has dependency within group
			}
		}
	}

	return true
}

// executeParallel executes steps in parallel using a worker pool
func (o *ResearchOrchestrator) executeParallel(
	ctx context.Context,
	steps []models.ResearchStep,
	totalSteps int,
	currentStep *int,
) []models.StepResult {
	results := make([]models.StepResult, len(steps))
	var wg sync.WaitGroup

	// Create worker pool
	maxWorkers := o.config.ConcurrentTools
	if maxWorkers <= 0 {
		maxWorkers = 3
	}

	semaphore := make(chan struct{}, maxWorkers)

	for i, step := range steps {
		wg.Add(1)
		go func(idx int, s models.ResearchStep) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			*currentStep++
			result := o.executeStep(ctx, s, *currentStep, totalSteps)
			results[idx] = result

			// Add to memory
			o.memory.AddStepResult(result)
		}(i, step)
	}

	wg.Wait()
	return results
}

// executeStep executes a single research step
func (o *ResearchOrchestrator) executeStep(
	ctx context.Context,
	step models.ResearchStep,
	currentStepNum int,
	totalSteps int,
) models.StepResult {
	startTime := time.Now()

	// Send progress update
	o.sendProgress(currentStepNum, totalSteps, step.Action, "in_progress")

	// Check if tool exists
	if !o.toolRegistry.HasTool(step.Tool) {
		return models.StepResult{
			Step:      step,
			Success:   false,
			Error:     fmt.Errorf("tool not found: %s", step.Tool),
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime),
		}
	}

	// Execute with retries
	var result models.StepResult
	var lastErr error

	maxRetries := step.MaxRetries
	if maxRetries == 0 {
		maxRetries = o.config.RetryAttempts
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			time.Sleep(o.config.RetryDelay)
			log.Printf("Retrying step %d (attempt %d/%d)", step.StepNumber, attempt+1, maxRetries+1)
		}

		// Create step-specific context with timeout
		stepCtx, cancel := context.WithTimeout(ctx, o.config.StepTimeout)

		// Execute tool
		toolResult, err := o.toolRegistry.Execute(stepCtx, step.Tool, step.Parameters)
		cancel()

		if err == nil && toolResult != nil && toolResult.Success {
			// Success!
			result = models.StepResult{
				Step:      step,
				Success:   true,
				Data:      toolResult.Data,
				Sources:   toolResult.Sources,
				Error:     nil,
				StartTime: startTime,
				EndTime:   time.Now(),
				Duration:  time.Since(startTime),
				Retries:   attempt,
			}

			// Extract and add sources to memory
			for _, source := range toolResult.Sources {
				o.memory.AddSource(source)
			}

			break
		} else {
			lastErr = err
			if err == nil && toolResult != nil {
				lastErr = toolResult.Error
			}
		}
	}

	// If all retries failed
	if result.Success == false {
		result = models.StepResult{
			Step:      step,
			Success:   false,
			Error:     lastErr,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime),
			Retries:   maxRetries,
		}
	}

	// Send completion progress
	status := "completed"
	if !result.Success {
		status = "failed"
	}
	o.sendProgress(currentStepNum, totalSteps, step.Action, status)

	return result
}

// sendProgress sends a progress update
func (o *ResearchOrchestrator) sendProgress(currentStep, totalSteps int, description, status string) {
	if o.progressChan == nil {
		return
	}

	progress := float64(currentStep) / float64(totalSteps) * 100

	update := models.ResearchProgress{
		SessionID:       o.memory.GetSessionID(),
		CurrentStep:     currentStep,
		TotalSteps:      totalSteps,
		StepDescription: description,
		Progress:        progress,
		Status:          status,
		Message:         fmt.Sprintf("Step %d/%d: %s", currentStep, totalSteps, description),
		Timestamp:       time.Now(),
	}

	// Non-blocking send
	select {
	case o.progressChan <- update:
	default:
		// Channel full, skip this update
	}
}

// Cancel stops the execution
func (o *ResearchOrchestrator) Cancel() {
	if o.cancelFunc != nil {
		o.cancelFunc()
	}
}

// AggregateResults combines step results into a final research result
func (o *ResearchOrchestrator) AggregateResults(
	query models.ResearchQuery,
	results []models.StepResult,
	duration time.Duration,
) *models.ResearchResult {
	// Count successes and failures
	successCount := 0
	failedCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failedCount++
		}
	}

	// Get all findings and sources from memory
	findings := o.memory.GetFindings()
	sources := o.memory.GetSources()

	// Generate summary
	summary := o.generateSummary(query, findings, sources)

	// Calculate overall confidence
	confidence := o.memory.GetAverageConfidence()
	if confidence == 0 && len(sources) > 0 {
		// Default confidence based on number of sources
		confidence = 0.5 + (float64(len(sources)) * 0.05)
		if confidence > 0.9 {
			confidence = 0.9
		}
	}

	return &models.ResearchResult{
		Query:           query.Query,
		Summary:         summary,
		KeyFindings:     findings,
		Sources:         sources,
		Steps:           results,
		Confidence:      confidence,
		Duration:        duration,
		TotalSteps:      len(results),
		SuccessfulSteps: successCount,
		FailedSteps:     failedCount,
		ResearchType:    models.ResearchGeneral, // This should come from the plan
		CompletionTime:  time.Now(),
	}
}

// generateSummary creates a summary from findings and sources
func (o *ResearchOrchestrator) generateSummary(
	query models.ResearchQuery,
	findings []models.Finding,
	sources []models.Source,
) string {
	summary := fmt.Sprintf("Research Summary for: %s\n\n", query.Query)

	if len(findings) > 0 {
		summary += fmt.Sprintf("Found %d key findings from %d sources.\n\n", len(findings), len(sources))

		// Add top findings
		maxFindings := 5
		if len(findings) < maxFindings {
			maxFindings = len(findings)
		}

		summary += "Key Findings:\n"
		for i := 0; i < maxFindings; i++ {
			summary += fmt.Sprintf("- %s\n", findings[i].Content)
		}
	} else {
		summary += "No significant findings were discovered.\n"
	}

	if len(sources) > 0 {
		summary += fmt.Sprintf("\n\nInformation gathered from %d sources:\n", len(sources))

		// Group sources by type
		sourceTypes := make(map[string]int)
		for _, source := range sources {
			sourceTypes[source.Type]++
		}

		for sourceType, count := range sourceTypes {
			summary += fmt.Sprintf("- %s: %d sources\n", sourceType, count)
		}
	}

	return summary
}
