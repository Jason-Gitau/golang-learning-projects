package tools

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// RandomTool generates random numbers and data
type RandomTool struct {
	rng *rand.Rand
}

// NewRandomTool creates a new random tool
func NewRandomTool() *RandomTool {
	return &RandomTool{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Name returns the tool name
func (t *RandomTool) Name() string {
	return "random"
}

// Description returns the tool description
func (t *RandomTool) Description() string {
	return "Generates random numbers, strings, and data"
}

// Execute executes the random tool
// Expected params:
// - type: string (int, float, bool, string, uuid, choice)
// - min: int (for int type)
// - max: int (for int type)
// - length: int (for string type)
// - choices: []interface{} (for choice type)
func (t *RandomTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	randType, ok := params["type"].(string)
	if !ok {
		randType = "int" // Default type
	}

	switch randType {
	case "int":
		return t.randomInt(params)
	case "float":
		return t.randomFloat(params)
	case "bool":
		return t.randomBool(params)
	case "string":
		return t.randomString(params)
	case "uuid":
		return t.randomUUID(params)
	case "choice":
		return t.randomChoice(params)
	case "dice":
		return t.rollDice(params)
	default:
		return nil, fmt.Errorf("unknown type: %s (supported: int, float, bool, string, uuid, choice, dice)", randType)
	}
}

func (t *RandomTool) randomInt(params map[string]interface{}) (interface{}, error) {
	min := 0
	max := 100

	if val, ok := params["min"]; ok {
		if minInt, ok := val.(float64); ok {
			min = int(minInt)
		}
	}

	if val, ok := params["max"]; ok {
		if maxInt, ok := val.(float64); ok {
			max = int(maxInt)
		}
	}

	if min >= max {
		return nil, fmt.Errorf("min must be less than max")
	}

	result := t.rng.Intn(max-min) + min

	return map[string]interface{}{
		"type":   "int",
		"min":    min,
		"max":    max,
		"result": result,
	}, nil
}

func (t *RandomTool) randomFloat(params map[string]interface{}) (interface{}, error) {
	min := 0.0
	max := 1.0

	if val, ok := params["min"]; ok {
		if minFloat, ok := val.(float64); ok {
			min = minFloat
		}
	}

	if val, ok := params["max"]; ok {
		if maxFloat, ok := val.(float64); ok {
			max = maxFloat
		}
	}

	if min >= max {
		return nil, fmt.Errorf("min must be less than max")
	}

	result := min + t.rng.Float64()*(max-min)

	return map[string]interface{}{
		"type":   "float",
		"min":    min,
		"max":    max,
		"result": result,
	}, nil
}

func (t *RandomTool) randomBool(params map[string]interface{}) (interface{}, error) {
	result := t.rng.Intn(2) == 1

	return map[string]interface{}{
		"type":   "bool",
		"result": result,
	}, nil
}

func (t *RandomTool) randomString(params map[string]interface{}) (interface{}, error) {
	length := 10

	if val, ok := params["length"]; ok {
		if lenInt, ok := val.(float64); ok {
			length = int(lenInt)
		}
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[t.rng.Intn(len(charset))]
	}

	return map[string]interface{}{
		"type":   "string",
		"length": length,
		"result": string(result),
	}, nil
}

func (t *RandomTool) randomUUID(params map[string]interface{}) (interface{}, error) {
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		t.rng.Uint32(),
		t.rng.Uint32()&0xffff,
		(t.rng.Uint32()&0x0fff)|0x4000,
		(t.rng.Uint32()&0x3fff)|0x8000,
		t.rng.Uint64()&0xffffffffffff,
	)

	return map[string]interface{}{
		"type":   "uuid",
		"result": uuid,
	}, nil
}

func (t *RandomTool) randomChoice(params map[string]interface{}) (interface{}, error) {
	choices, ok := params["choices"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("choices parameter is required (array)")
	}

	if len(choices) == 0 {
		return nil, fmt.Errorf("choices array cannot be empty")
	}

	idx := t.rng.Intn(len(choices))
	result := choices[idx]

	return map[string]interface{}{
		"type":    "choice",
		"choices": choices,
		"result":  result,
		"index":   idx,
	}, nil
}

func (t *RandomTool) rollDice(params map[string]interface{}) (interface{}, error) {
	sides := 6
	count := 1

	if val, ok := params["sides"]; ok {
		if sidesInt, ok := val.(float64); ok {
			sides = int(sidesInt)
		}
	}

	if val, ok := params["count"]; ok {
		if countInt, ok := val.(float64); ok {
			count = int(countInt)
		}
	}

	if sides < 2 {
		return nil, fmt.Errorf("dice must have at least 2 sides")
	}

	if count < 1 {
		return nil, fmt.Errorf("must roll at least 1 die")
	}

	rolls := make([]int, count)
	total := 0
	for i := 0; i < count; i++ {
		roll := t.rng.Intn(sides) + 1
		rolls[i] = roll
		total += roll
	}

	return map[string]interface{}{
		"type":  "dice",
		"sides": sides,
		"count": count,
		"rolls": rolls,
		"total": total,
	}, nil
}
