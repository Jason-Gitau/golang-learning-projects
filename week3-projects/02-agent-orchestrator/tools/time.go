package tools

import (
	"context"
	"fmt"
	"time"
)

// TimeTool provides time-related operations
type TimeTool struct{}

// NewTimeTool creates a new time tool
func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

// Name returns the tool name
func (t *TimeTool) Name() string {
	return "time"
}

// Description returns the tool description
func (t *TimeTool) Description() string {
	return "Provides current time, date, and timezone information"
}

// Execute executes the time tool
// Expected params:
// - action: string (current, timezone, format, add_duration)
// - timezone: string (optional, e.g., "America/New_York", "UTC")
// - format: string (optional, e.g., "2006-01-02 15:04:05")
// - duration: string (optional, for add_duration, e.g., "2h30m")
func (t *TimeTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	action, ok := params["action"].(string)
	if !ok {
		action = "current" // Default action
	}

	switch action {
	case "current":
		return t.getCurrentTime(params)
	case "timezone":
		return t.getTimeInTimezone(params)
	case "format":
		return t.formatTime(params)
	case "add_duration":
		return t.addDuration(params)
	case "unix":
		return t.getUnixTime(params)
	default:
		return nil, fmt.Errorf("unknown action: %s (supported: current, timezone, format, add_duration, unix)", action)
	}
}

func (t *TimeTool) getCurrentTime(params map[string]interface{}) (interface{}, error) {
	now := time.Now()
	return map[string]interface{}{
		"timestamp":     now.Format(time.RFC3339),
		"unix":          now.Unix(),
		"date":          now.Format("2006-01-02"),
		"time":          now.Format("15:04:05"),
		"timezone":      now.Location().String(),
		"day_of_week":   now.Weekday().String(),
		"day_of_year":   now.YearDay(),
		"week_of_year":  getWeekOfYear(now),
	}, nil
}

func (t *TimeTool) getTimeInTimezone(params map[string]interface{}) (interface{}, error) {
	timezone, ok := params["timezone"].(string)
	if !ok {
		return nil, fmt.Errorf("timezone parameter is required (string)")
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %v", err)
	}

	now := time.Now().In(loc)
	return map[string]interface{}{
		"timezone":   timezone,
		"timestamp":  now.Format(time.RFC3339),
		"date":       now.Format("2006-01-02"),
		"time":       now.Format("15:04:05"),
		"offset":     now.Format("-07:00"),
	}, nil
}

func (t *TimeTool) formatTime(params map[string]interface{}) (interface{}, error) {
	format, ok := params["format"].(string)
	if !ok {
		format = time.RFC3339 // Default format
	}

	now := time.Now()
	return map[string]interface{}{
		"format":    format,
		"formatted": now.Format(format),
	}, nil
}

func (t *TimeTool) addDuration(params map[string]interface{}) (interface{}, error) {
	durationStr, ok := params["duration"].(string)
	if !ok {
		return nil, fmt.Errorf("duration parameter is required (string, e.g., '2h30m')")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %v", err)
	}

	now := time.Now()
	future := now.Add(duration)

	return map[string]interface{}{
		"current":   now.Format(time.RFC3339),
		"duration":  durationStr,
		"future":    future.Format(time.RFC3339),
		"added_seconds": duration.Seconds(),
	}, nil
}

func (t *TimeTool) getUnixTime(params map[string]interface{}) (interface{}, error) {
	now := time.Now()
	return map[string]interface{}{
		"unix":        now.Unix(),
		"unix_milli":  now.UnixMilli(),
		"unix_micro":  now.UnixMicro(),
		"unix_nano":   now.UnixNano(),
	}, nil
}

// getWeekOfYear calculates the ISO 8601 week number
func getWeekOfYear(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}
