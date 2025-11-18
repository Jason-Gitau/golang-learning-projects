package tools

import (
	"context"
	"fmt"
	"time"
)

// DateTimeTool provides date and time utilities
type DateTimeTool struct{}

// NewDateTimeTool creates a new datetime tool
func NewDateTimeTool() *DateTimeTool {
	return &DateTimeTool{}
}

// Name returns the tool name
func (t *DateTimeTool) Name() string {
	return "datetime"
}

// Description returns the tool description
func (t *DateTimeTool) Description() string {
	return "Provides current time, date information, and timezone conversions"
}

// Parameters returns the tool parameters schema
func (t *DateTimeTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"description": "The operation to perform",
				"enum":        []string{"current", "timezone", "format", "parse", "add", "diff"},
				"default":     "current",
			},
			"timezone": map[string]interface{}{
				"type":        "string",
				"description": "Timezone for the operation (e.g., 'UTC', 'America/New_York')",
				"default":     "UTC",
			},
			"format": map[string]interface{}{
				"type":        "string",
				"description": "Format string for time formatting",
			},
			"time_string": map[string]interface{}{
				"type":        "string",
				"description": "Time string to parse",
			},
			"duration": map[string]interface{}{
				"type":        "string",
				"description": "Duration to add (e.g., '2h', '30m', '1h30m')",
			},
		},
		"required": []string{"operation"},
	}
}

// Execute performs datetime operations
func (t *DateTimeTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		operation = "current"
	}

	timezone := "UTC"
	if tz, ok := params["timezone"].(string); ok {
		timezone = tz
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %s", timezone)
	}

	switch operation {
	case "current":
		return t.getCurrentTime(location)

	case "timezone":
		return t.convertTimezone(params, location)

	case "format":
		return t.formatTime(params, location)

	case "parse":
		return t.parseTime(params, location)

	case "add":
		return t.addDuration(params, location)

	case "diff":
		return t.timeDifference(params, location)

	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}
}

// getCurrentTime returns the current time in the specified timezone
func (t *DateTimeTool) getCurrentTime(location *time.Location) (interface{}, error) {
	now := time.Now().In(location)

	return map[string]interface{}{
		"timestamp":    now.Unix(),
		"iso8601":      now.Format(time.RFC3339),
		"readable":     now.Format("Monday, January 2, 2006 at 3:04:05 PM MST"),
		"timezone":     location.String(),
		"date":         now.Format("2006-01-02"),
		"time":         now.Format("15:04:05"),
		"day_of_week":  now.Weekday().String(),
		"day_of_year":  now.YearDay(),
		"week_of_year": getWeekOfYear(now),
		"unix":         now.Unix(),
		"unix_nano":    now.UnixNano(),
	}, nil
}

// convertTimezone converts current time to a different timezone
func (t *DateTimeTool) convertTimezone(params map[string]interface{}, location *time.Location) (interface{}, error) {
	now := time.Now().In(location)

	return map[string]interface{}{
		"from_timezone": "Local",
		"to_timezone":   location.String(),
		"converted":     now.Format(time.RFC3339),
		"readable":      now.Format("Monday, January 2, 2006 at 3:04:05 PM MST"),
	}, nil
}

// formatTime formats the current time using a custom format
func (t *DateTimeTool) formatTime(params map[string]interface{}, location *time.Location) (interface{}, error) {
	format, ok := params["format"].(string)
	if !ok {
		format = time.RFC3339
	}

	now := time.Now().In(location)
	formatted := now.Format(format)

	return map[string]interface{}{
		"format":    format,
		"formatted": formatted,
		"timezone":  location.String(),
	}, nil
}

// parseTime parses a time string
func (t *DateTimeTool) parseTime(params map[string]interface{}, location *time.Location) (interface{}, error) {
	timeStr, ok := params["time_string"].(string)
	if !ok {
		return nil, fmt.Errorf("time_string is required for parse operation")
	}

	// Try common formats
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"01/02/2006",
		"01/02/2006 3:04 PM",
	}

	var parsed time.Time
	var parseErr error

	for _, format := range formats {
		parsed, parseErr = time.Parse(format, timeStr)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return nil, fmt.Errorf("unable to parse time string: %s", timeStr)
	}

	parsed = parsed.In(location)

	return map[string]interface{}{
		"input":     timeStr,
		"parsed":    parsed.Format(time.RFC3339),
		"timestamp": parsed.Unix(),
		"timezone":  location.String(),
	}, nil
}

// addDuration adds a duration to the current time
func (t *DateTimeTool) addDuration(params map[string]interface{}, location *time.Location) (interface{}, error) {
	durationStr, ok := params["duration"].(string)
	if !ok {
		return nil, fmt.Errorf("duration is required for add operation")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	now := time.Now().In(location)
	future := now.Add(duration)

	return map[string]interface{}{
		"current":  now.Format(time.RFC3339),
		"duration": durationStr,
		"result":   future.Format(time.RFC3339),
		"timezone": location.String(),
	}, nil
}

// timeDifference calculates the difference between current time and a future time
func (t *DateTimeTool) timeDifference(params map[string]interface{}, location *time.Location) (interface{}, error) {
	timeStr, ok := params["time_string"].(string)
	if !ok {
		return nil, fmt.Errorf("time_string is required for diff operation")
	}

	targetTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format, use RFC3339: %s", timeStr)
	}

	now := time.Now().In(location)
	diff := targetTime.Sub(now)

	return map[string]interface{}{
		"from":              now.Format(time.RFC3339),
		"to":                targetTime.Format(time.RFC3339),
		"difference":        diff.String(),
		"difference_hours":  diff.Hours(),
		"difference_minutes": diff.Minutes(),
		"difference_seconds": diff.Seconds(),
	}, nil
}

// getWeekOfYear calculates the ISO week number
func getWeekOfYear(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}
