package tools

import (
	"context"
	"fmt"
)

// CalculatorTool performs basic mathematical operations
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
	return "Performs basic mathematical operations (add, subtract, multiply, divide)"
}

// Execute executes the calculator tool
// Expected params:
// - operation: string (add, subtract, multiply, divide)
// - a: float64
// - b: float64
func (t *CalculatorTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required (string)")
	}

	a, ok := getFloat64(params["a"])
	if !ok {
		return nil, fmt.Errorf("parameter 'a' is required (number)")
	}

	b, ok := getFloat64(params["b"])
	if !ok {
		return nil, fmt.Errorf("parameter 'b' is required (number)")
	}

	var result float64
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
	default:
		return nil, fmt.Errorf("unknown operation: %s (supported: add, subtract, multiply, divide)", operation)
	}

	return map[string]interface{}{
		"operation": operation,
		"a":         a,
		"b":         b,
		"result":    result,
	}, nil
}

// getFloat64 converts interface{} to float64
func getFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	default:
		return 0, false
	}
}
