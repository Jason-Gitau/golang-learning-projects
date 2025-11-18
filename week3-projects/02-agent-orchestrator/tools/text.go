package tools

import (
	"context"
	"fmt"
	"strings"
	"unicode"
)

// TextTool provides text manipulation operations
type TextTool struct{}

// NewTextTool creates a new text tool
func NewTextTool() *TextTool {
	return &TextTool{}
}

// Name returns the tool name
func (t *TextTool) Name() string {
	return "text"
}

// Description returns the tool description
func (t *TextTool) Description() string {
	return "Provides text manipulation operations (uppercase, lowercase, reverse, count, etc.)"
}

// Execute executes the text tool
// Expected params:
// - operation: string (uppercase, lowercase, reverse, count, trim, replace, split)
// - text: string
// - find: string (for replace operation)
// - replace_with: string (for replace operation)
// - delimiter: string (for split operation)
func (t *TextTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required (string)")
	}

	text, ok := params["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text parameter is required (string)")
	}

	switch operation {
	case "uppercase":
		return t.uppercase(text)
	case "lowercase":
		return t.lowercase(text)
	case "title":
		return t.titleCase(text)
	case "reverse":
		return t.reverse(text)
	case "count":
		return t.count(text)
	case "trim":
		return t.trim(text)
	case "replace":
		return t.replace(text, params)
	case "split":
		return t.split(text, params)
	case "contains":
		return t.contains(text, params)
	case "word_count":
		return t.wordCount(text)
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

func (t *TextTool) uppercase(text string) (interface{}, error) {
	return map[string]interface{}{
		"operation": "uppercase",
		"original":  text,
		"result":    strings.ToUpper(text),
	}, nil
}

func (t *TextTool) lowercase(text string) (interface{}, error) {
	return map[string]interface{}{
		"operation": "lowercase",
		"original":  text,
		"result":    strings.ToLower(text),
	}, nil
}

func (t *TextTool) titleCase(text string) (interface{}, error) {
	return map[string]interface{}{
		"operation": "title",
		"original":  text,
		"result":    strings.Title(text),
	}, nil
}

func (t *TextTool) reverse(text string) (interface{}, error) {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return map[string]interface{}{
		"operation": "reverse",
		"original":  text,
		"result":    string(runes),
	}, nil
}

func (t *TextTool) count(text string) (interface{}, error) {
	charCount := len([]rune(text))
	letterCount := 0
	digitCount := 0
	spaceCount := 0
	specialCount := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			letterCount++
		} else if unicode.IsDigit(r) {
			digitCount++
		} else if unicode.IsSpace(r) {
			spaceCount++
		} else {
			specialCount++
		}
	}

	return map[string]interface{}{
		"operation":     "count",
		"text":          text,
		"total_chars":   charCount,
		"letters":       letterCount,
		"digits":        digitCount,
		"spaces":        spaceCount,
		"special_chars": specialCount,
	}, nil
}

func (t *TextTool) trim(text string) (interface{}, error) {
	return map[string]interface{}{
		"operation": "trim",
		"original":  text,
		"result":    strings.TrimSpace(text),
	}, nil
}

func (t *TextTool) replace(text string, params map[string]interface{}) (interface{}, error) {
	find, ok := params["find"].(string)
	if !ok {
		return nil, fmt.Errorf("find parameter is required for replace operation")
	}

	replaceWith, ok := params["replace_with"].(string)
	if !ok {
		return nil, fmt.Errorf("replace_with parameter is required for replace operation")
	}

	result := strings.ReplaceAll(text, find, replaceWith)
	count := strings.Count(text, find)

	return map[string]interface{}{
		"operation":    "replace",
		"original":     text,
		"find":         find,
		"replace_with": replaceWith,
		"result":       result,
		"occurrences":  count,
	}, nil
}

func (t *TextTool) split(text string, params map[string]interface{}) (interface{}, error) {
	delimiter, ok := params["delimiter"].(string)
	if !ok {
		delimiter = " " // Default to space
	}

	parts := strings.Split(text, delimiter)

	return map[string]interface{}{
		"operation": "split",
		"original":  text,
		"delimiter": delimiter,
		"result":    parts,
		"count":     len(parts),
	}, nil
}

func (t *TextTool) contains(text string, params map[string]interface{}) (interface{}, error) {
	search, ok := params["search"].(string)
	if !ok {
		return nil, fmt.Errorf("search parameter is required for contains operation")
	}

	contains := strings.Contains(text, search)
	count := strings.Count(text, search)

	return map[string]interface{}{
		"operation":   "contains",
		"text":        text,
		"search":      search,
		"contains":    contains,
		"occurrences": count,
	}, nil
}

func (t *TextTool) wordCount(text string) (interface{}, error) {
	words := strings.Fields(text)

	return map[string]interface{}{
		"operation":  "word_count",
		"text":       text,
		"word_count": len(words),
		"words":      words,
	}, nil
}
