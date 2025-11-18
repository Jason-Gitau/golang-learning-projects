package tools

import (
	"context"
	"fmt"
	"math"
	"strings"

	"deep-research-agent/models"
)

// FactChecker verifies claims against sources
type FactChecker struct{}

// NewFactChecker creates a new fact checker tool
func NewFactChecker() *FactChecker {
	return &FactChecker{}
}

// Name returns the tool name
func (f *FactChecker) Name() string {
	return "fact_checker"
}

// Description returns the tool description
func (f *FactChecker) Description() string {
	return "Cross-reference claims against sources, check consistency, and identify contradictions. Assigns confidence scores."
}

// Parameters returns the tool parameters
func (f *FactChecker) Parameters() []Parameter {
	return []Parameter{
		{
			Name:        "action",
			Type:        "string",
			Required:    true,
			Description: "Action: check_fact, find_contradictions, calculate_confidence",
		},
		{
			Name:        "claim",
			Type:        "string",
			Required:    false,
			Description: "Claim to verify",
		},
		{
			Name:        "sources",
			Type:        "array",
			Required:    false,
			Description: "Array of sources to check against",
		},
		{
			Name:        "findings",
			Type:        "array",
			Required:    false,
			Description: "Array of findings to check for contradictions",
		},
		{
			Name:        "evidence",
			Type:        "array",
			Required:    false,
			Description: "Array of evidence strings",
		},
	}
}

// Execute runs the fact checker
func (f *FactChecker) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	// Extract parameters
	action, ok := params["action"].(string)
	if !ok || action == "" {
		return nil, fmt.Errorf("action parameter is required")
	}

	// Execute action
	switch action {
	case "check_fact":
		claim, ok := params["claim"].(string)
		if !ok || claim == "" {
			return nil, fmt.Errorf("claim parameter is required for check_fact action")
		}

		sourcesData, ok := params["sources"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("sources parameter is required for check_fact action")
		}

		sources := f.convertToSources(sourcesData)
		return f.checkFact(claim, sources)

	case "find_contradictions":
		findingsData, ok := params["findings"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("findings parameter is required for find_contradictions action")
		}

		findings := f.convertToFindings(findingsData)
		return f.findContradictions(findings)

	case "calculate_confidence":
		claim, ok := params["claim"].(string)
		if !ok || claim == "" {
			return nil, fmt.Errorf("claim parameter is required for calculate_confidence action")
		}

		evidenceData, ok := params["evidence"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("evidence parameter is required for calculate_confidence action")
		}

		evidence := f.convertToStrings(evidenceData)
		return f.calculateConfidence(claim, evidence)

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// checkFact verifies a claim against sources
func (f *FactChecker) checkFact(claim string, sources []models.Source) (*ToolResult, error) {
	claimLower := strings.ToLower(claim)
	claimWords := f.tokenize(claimLower)

	supportingSources := make([]models.Source, 0)
	contradictions := make([]string, 0)
	evidenceScores := make([]float64, 0)

	// Check each source
	for _, source := range sources {
		contentLower := strings.ToLower(source.Content)

		// Calculate similarity
		similarity := f.calculateTextSimilarity(claimLower, contentLower)

		if similarity > 0.3 {
			// High similarity - likely supporting
			supportingSources = append(supportingSources, source)
			evidenceScores = append(evidenceScores, similarity)
		}

		// Check for contradictions (negation words near claim words)
		if f.hasContradiction(claimLower, contentLower) {
			contradictions = append(contradictions, source.Title)
		}
	}

	// Calculate overall confidence
	confidence := f.calculateConfidenceScore(evidenceScores, len(sources))

	// Determine if verified
	verified := len(supportingSources) > 0 && confidence > 0.5 && len(contradictions) == 0

	result := &FactCheckResult{
		Claim:             claim,
		Verified:          verified,
		Confidence:        confidence,
		SupportingSources: supportingSources,
		Contradictions:    contradictions,
		SourceCount:       len(sources),
		SupportCount:      len(supportingSources),
	}

	return &ToolResult{
		Success: true,
		Data:    result,
		Sources: supportingSources,
		Metadata: map[string]interface{}{
			"verified":           verified,
			"confidence":         confidence,
			"supporting_count":   len(supportingSources),
			"contradiction_count": len(contradictions),
		},
	}, nil
}

// findContradictions identifies contradictions in findings
func (f *FactChecker) findContradictions(findings []models.Finding) (*ToolResult, error) {
	contradictions := make([]Contradiction, 0)

	// Compare each pair of findings
	for i := 0; i < len(findings); i++ {
		for j := i + 1; j < len(findings); j++ {
			finding1 := findings[i]
			finding2 := findings[j]

			// Check if they contradict
			if f.areContradictory(finding1.Content, finding2.Content) {
				contradictions = append(contradictions, Contradiction{
					Finding1:    finding1.Content,
					Finding2:    finding2.Content,
					Source1:     finding1.Source.Title,
					Source2:     finding2.Source.Title,
					Severity:    f.calculateContradictionSeverity(finding1, finding2),
					Description: "Potential contradiction detected between sources",
				})
			}
		}
	}

	result := &ContradictionResult{
		Contradictions:     contradictions,
		ContradictionCount: len(contradictions),
		TotalFindings:      len(findings),
		HasContradictions:  len(contradictions) > 0,
	}

	return &ToolResult{
		Success: true,
		Data:    result,
		Metadata: map[string]interface{}{
			"contradiction_count": len(contradictions),
			"total_findings":      len(findings),
		},
	}, nil
}

// calculateConfidence calculates confidence score for a claim
func (f *FactChecker) calculateConfidence(claim string, evidence []string) (*ToolResult, error) {
	if len(evidence) == 0 {
		return &ToolResult{
			Success: true,
			Data: map[string]interface{}{
				"confidence": 0.0,
				"factors": map[string]float64{
					"evidence_count": 0.0,
					"similarity":     0.0,
				},
			},
		}, nil
	}

	claimLower := strings.ToLower(claim)
	similarities := make([]float64, 0)

	// Calculate similarity with each evidence
	for _, ev := range evidence {
		similarity := f.calculateTextSimilarity(claimLower, strings.ToLower(ev))
		similarities = append(similarities, similarity)
	}

	// Calculate confidence factors
	evidenceCount := len(evidence)
	avgSimilarity := f.average(similarities)
	maxSimilarity := f.max(similarities)

	// Weighted confidence score
	confidence := 0.0
	confidence += math.Min(float64(evidenceCount)/10.0, 0.4) // Max 0.4 for evidence count
	confidence += avgSimilarity * 0.3                         // Max 0.3 for avg similarity
	confidence += maxSimilarity * 0.3                         // Max 0.3 for max similarity

	result := &ConfidenceResult{
		Claim:         claim,
		Confidence:    math.Min(confidence, 1.0),
		EvidenceCount: evidenceCount,
		Factors: map[string]float64{
			"evidence_count": math.Min(float64(evidenceCount)/10.0, 0.4),
			"avg_similarity": avgSimilarity * 0.3,
			"max_similarity": maxSimilarity * 0.3,
		},
	}

	return &ToolResult{
		Success: true,
		Data:    result,
		Metadata: map[string]interface{}{
			"confidence":     result.Confidence,
			"evidence_count": evidenceCount,
		},
	}, nil
}

// calculateTextSimilarity calculates similarity between two texts
func (f *FactChecker) calculateTextSimilarity(text1, text2 string) float64 {
	words1 := f.tokenize(text1)
	words2 := f.tokenize(text2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Create word sets
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, word := range words1 {
		set1[word] = true
	}
	for _, word := range words2 {
		set2[word] = true
	}

	// Calculate Jaccard similarity
	intersection := 0
	for word := range set1 {
		if set2[word] {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// hasContradiction checks if content contradicts claim
func (f *FactChecker) hasContradiction(claim, content string) bool {
	negationWords := []string{"not", "no", "never", "neither", "none", "nobody", "nothing", "false", "incorrect", "wrong"}

	claimWords := f.tokenize(claim)
	contentWords := f.tokenize(content)

	// Look for negation words near claim words in content
	for i, word := range contentWords {
		// Check if this is a negation word
		isNegation := false
		for _, neg := range negationWords {
			if word == neg {
				isNegation = true
				break
			}
		}

		if isNegation {
			// Check nearby words (window of 5)
			start := i - 5
			if start < 0 {
				start = 0
			}
			end := i + 5
			if end > len(contentWords) {
				end = len(contentWords)
			}

			// Check if any claim words are nearby
			for j := start; j < end; j++ {
				for _, claimWord := range claimWords {
					if contentWords[j] == claimWord {
						return true
					}
				}
			}
		}
	}

	return false
}

// areContradictory checks if two findings contradict each other
func (f *FactChecker) areContradictory(text1, text2 string) bool {
	text1Lower := strings.ToLower(text1)
	text2Lower := strings.ToLower(text2)

	// Check for explicit contradictions
	if f.hasContradiction(text1Lower, text2Lower) || f.hasContradiction(text2Lower, text1Lower) {
		return true
	}

	// Check for opposite statements (simple heuristic)
	opposites := map[string]string{
		"increase": "decrease",
		"rise":     "fall",
		"growth":   "decline",
		"more":     "less",
		"higher":   "lower",
		"better":   "worse",
		"positive": "negative",
		"true":     "false",
		"correct":  "incorrect",
	}

	words1 := f.tokenize(text1Lower)
	words2 := f.tokenize(text2Lower)

	for _, word1 := range words1 {
		if opposite, exists := opposites[word1]; exists {
			for _, word2 := range words2 {
				if word2 == opposite {
					// Check if they're talking about the same thing
					similarity := f.calculateTextSimilarity(text1Lower, text2Lower)
					if similarity > 0.3 {
						return true
					}
				}
			}
		}
	}

	return false
}

// calculateContradictionSeverity calculates how severe a contradiction is
func (f *FactChecker) calculateContradictionSeverity(finding1, finding2 models.Finding) string {
	// Higher confidence findings make more severe contradictions
	avgConfidence := (finding1.Confidence + finding2.Confidence) / 2

	if avgConfidence > 0.8 {
		return "high"
	} else if avgConfidence > 0.5 {
		return "medium"
	}
	return "low"
}

// calculateConfidenceScore calculates overall confidence from evidence scores
func (f *FactChecker) calculateConfidenceScore(scores []float64, totalSources int) float64 {
	if len(scores) == 0 {
		return 0.0
	}

	// Average similarity
	avgScore := f.average(scores)

	// Boost for multiple sources
	sourceBoost := math.Min(float64(len(scores))/float64(totalSources), 1.0)

	confidence := (avgScore + sourceBoost) / 2.0
	return math.Min(confidence, 1.0)
}

// tokenize splits text into words
func (f *FactChecker) tokenize(text string) []string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	})

	// Filter out stop words and short words
	stopWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
		"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
		"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
		"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
	}

	filtered := make([]string, 0)
	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// Helper functions
func (f *FactChecker) average(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (f *FactChecker) max(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// Conversion helpers
func (f *FactChecker) convertToSources(data []interface{}) []models.Source {
	sources := make([]models.Source, 0)
	for _, item := range data {
		if sourceMap, ok := item.(map[string]interface{}); ok {
			source := models.Source{}
			if id, ok := sourceMap["id"].(string); ok {
				source.ID = id
			}
			if title, ok := sourceMap["title"].(string); ok {
				source.Title = title
			}
			if content, ok := sourceMap["content"].(string); ok {
				source.Content = content
			}
			if url, ok := sourceMap["url"].(string); ok {
				source.URL = url
			}
			sources = append(sources, source)
		}
	}
	return sources
}

func (f *FactChecker) convertToFindings(data []interface{}) []models.Finding {
	findings := make([]models.Finding, 0)
	for _, item := range data {
		if findingMap, ok := item.(map[string]interface{}); ok {
			finding := models.Finding{}
			if id, ok := findingMap["id"].(string); ok {
				finding.ID = id
			}
			if content, ok := findingMap["content"].(string); ok {
				finding.Content = content
			}
			if confidence, ok := findingMap["confidence"].(float64); ok {
				finding.Confidence = confidence
			}
			// Extract source if present
			if sourceMap, ok := findingMap["source"].(map[string]interface{}); ok {
				if title, ok := sourceMap["title"].(string); ok {
					finding.Source.Title = title
				}
				if content, ok := sourceMap["content"].(string); ok {
					finding.Source.Content = content
				}
			}
			findings = append(findings, finding)
		}
	}
	return findings
}

func (f *FactChecker) convertToStrings(data []interface{}) []string {
	strings := make([]string, 0)
	for _, item := range data {
		if str, ok := item.(string); ok {
			strings = append(strings, str)
		}
	}
	return strings
}

// FactCheckResult represents the result of fact checking
type FactCheckResult struct {
	Claim             string          `json:"claim"`
	Verified          bool            `json:"verified"`
	Confidence        float64         `json:"confidence"`
	SupportingSources []models.Source `json:"supporting_sources"`
	Contradictions    []string        `json:"contradictions"`
	SourceCount       int             `json:"source_count"`
	SupportCount      int             `json:"support_count"`
}

// ContradictionResult represents contradictions found in findings
type ContradictionResult struct {
	Contradictions     []Contradiction `json:"contradictions"`
	ContradictionCount int             `json:"contradiction_count"`
	TotalFindings      int             `json:"total_findings"`
	HasContradictions  bool            `json:"has_contradictions"`
}

// Contradiction represents a contradiction between two findings
type Contradiction struct {
	Finding1    string `json:"finding1"`
	Finding2    string `json:"finding2"`
	Source1     string `json:"source1"`
	Source2     string `json:"source2"`
	Severity    string `json:"severity"` // low, medium, high
	Description string `json:"description"`
}

// ConfidenceResult represents confidence calculation result
type ConfidenceResult struct {
	Claim         string             `json:"claim"`
	Confidence    float64            `json:"confidence"`
	EvidenceCount int                `json:"evidence_count"`
	Factors       map[string]float64 `json:"factors"`
}
