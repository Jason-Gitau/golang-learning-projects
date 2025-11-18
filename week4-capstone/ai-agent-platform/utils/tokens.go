package utils

import (
	"strings"
)

// ModelPricing represents pricing for a model (per 1K tokens)
type ModelPricing struct {
	InputPricePerK  float64
	OutputPricePerK float64
}

// Pricing information for different models (in dollars per 1K tokens)
var modelPricing = map[string]ModelPricing{
	"gpt-4": {
		InputPricePerK:  0.03,
		OutputPricePerK: 0.06,
	},
	"gpt-4-turbo": {
		InputPricePerK:  0.01,
		OutputPricePerK: 0.03,
	},
	"gpt-3.5-turbo": {
		InputPricePerK:  0.0005,
		OutputPricePerK: 0.0015,
	},
	"claude-3-opus": {
		InputPricePerK:  0.015,
		OutputPricePerK: 0.075,
	},
	"claude-3-sonnet": {
		InputPricePerK:  0.003,
		OutputPricePerK: 0.015,
	},
	"claude-3-haiku": {
		InputPricePerK:  0.00025,
		OutputPricePerK: 0.00125,
	},
}

// CountTokens provides a simple token count estimation
// In production, you would use a proper tokenizer like tiktoken
func CountTokens(text string) int {
	// Simple estimation: ~4 characters per token
	// This is a rough approximation for English text
	words := strings.Fields(text)
	totalChars := len(text)

	// Estimate tokens based on both words and characters
	// Average of word count and character count / 4
	tokenEstimate := (len(words) + totalChars/4) / 2

	if tokenEstimate < 1 {
		return 1
	}

	return tokenEstimate
}

// CalculateCost calculates the cost for token usage
// Returns cost in cents
func CalculateCost(model string, inputTokens, outputTokens int) float64 {
	pricing, exists := modelPricing[model]
	if !exists {
		// Default pricing if model not found
		pricing = ModelPricing{
			InputPricePerK:  0.01,
			OutputPricePerK: 0.03,
		}
	}

	// Calculate cost in dollars
	inputCost := (float64(inputTokens) / 1000.0) * pricing.InputPricePerK
	outputCost := (float64(outputTokens) / 1000.0) * pricing.OutputPricePerK
	totalCost := inputCost + outputCost

	// Convert to cents
	return totalCost * 100.0
}

// GetModelPricing returns the pricing for a model
func GetModelPricing(model string) (ModelPricing, bool) {
	pricing, exists := modelPricing[model]
	return pricing, exists
}
