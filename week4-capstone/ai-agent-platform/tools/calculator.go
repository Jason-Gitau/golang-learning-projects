package tools

import (
	"context"
	"fmt"
	"math"
)

// CalculatorTool performs mathematical operations
type CalculatorTool struct{}

// NewCalculatorTool creates a new calculator tool
func NewCalculatorTool() *CalculatorTool {
	return &CalculatorTool{}
}

// Name returns the tool name
func (t *CalculatorTool) Name() string {
	return "calculator"
}

// Description returns the tool description
func (t *CalculatorTool) Description() string {
	return "Performs mathematical operations including basic arithmetic, powers, and square roots"
}

// Parameters returns the tool parameters schema
func (t *CalculatorTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"description": "The mathematical operation to perform",
				"enum":        []string{"add", "subtract", "multiply", "divide", "power", "sqrt", "mod"},
			},
			"a": map[string]interface{}{
				"type":        "number",
				"description": "The first operand",
			},
			"b": map[string]interface{}{
				"type":        "number",
				"description": "The second operand (not required for sqrt)",
			},
		},
		"required": []string{"operation", "a"},
	}
}

// Execute runs the calculator operation
func (t *CalculatorTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation must be a string")
	}

	a, err := getNumber(params, "a")
	if err != nil {
		return nil, err
	}

	var result float64

	switch operation {
	case "sqrt":
		if a < 0 {
			return nil, fmt.Errorf("cannot take square root of negative number")
		}
		result = math.Sqrt(a)

	case "add", "subtract", "multiply", "divide", "power", "mod":
		b, err := getNumber(params, "b")
		if err != nil {
			return nil, err
		}

		switch operation {
		case "add":
			result = a + b
		case "subtract":
			result = a - b
		case "multiply":
			result = a * b
		case "divide":
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			result = a / b
		case "power":
			result = math.Pow(a, b)
		case "mod":
			if b == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			result = math.Mod(a, b)
		}

	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}

	return map[string]interface{}{
		"operation": operation,
		"result":    result,
		"operands": map[string]interface{}{
			"a": a,
			"b": params["b"],
		},
	}, nil
}

// getNumber extracts a number from params
func getNumber(params map[string]interface{}, key string) (float64, error) {
	val, exists := params[key]
	if !exists {
		return 0, fmt.Errorf("%s is required", key)
	}

	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("%s must be a number", key)
	}
}
