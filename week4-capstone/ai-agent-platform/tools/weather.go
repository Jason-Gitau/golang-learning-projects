package tools

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// WeatherTool provides weather information (mock data)
type WeatherTool struct {
	rand *rand.Rand
}

// NewWeatherTool creates a new weather tool
func NewWeatherTool() *WeatherTool {
	return &WeatherTool{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Name returns the tool name
func (t *WeatherTool) Name() string {
	return "weather"
}

// Description returns the tool description
func (t *WeatherTool) Description() string {
	return "Retrieves current weather information for a specified location"
}

// Parameters returns the tool parameters schema
func (t *WeatherTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city or location to get weather for",
			},
			"units": map[string]interface{}{
				"type":        "string",
				"description": "Temperature units (celsius or fahrenheit)",
				"enum":        []string{"celsius", "fahrenheit"},
				"default":     "celsius",
			},
		},
		"required": []string{"location"},
	}
}

// Execute retrieves weather data
func (t *WeatherTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("location must be a non-empty string")
	}

	units := "celsius"
	if u, ok := params["units"].(string); ok {
		units = u
	}

	// Generate mock weather data
	weather := t.generateMockWeather(location, units)

	return weather, nil
}

// generateMockWeather creates realistic mock weather data
func (t *WeatherTool) generateMockWeather(location, units string) map[string]interface{} {
	conditions := []string{"Sunny", "Partly Cloudy", "Cloudy", "Rainy", "Stormy", "Snowy", "Foggy", "Windy"}
	condition := conditions[t.rand.Intn(len(conditions))]

	// Generate temperature based on units
	var temp float64
	if units == "fahrenheit" {
		temp = float64(t.rand.Intn(70) + 30) // 30-100°F
	} else {
		temp = float64(t.rand.Intn(35) - 5) // -5 to 30°C
	}

	humidity := t.rand.Intn(60) + 30  // 30-90%
	windSpeed := t.rand.Intn(30) + 5  // 5-35 km/h or mph
	precipitation := t.rand.Intn(100) // 0-100%

	return map[string]interface{}{
		"location":    strings.Title(location),
		"temperature": temp,
		"units":       units,
		"condition":   condition,
		"humidity":    humidity,
		"wind_speed":  windSpeed,
		"wind_unit":   map[string]string{"celsius": "km/h", "fahrenheit": "mph"}[units],
		"precipitation_chance": precipitation,
		"timestamp":            time.Now().Format(time.RFC3339),
		"description": fmt.Sprintf("%s with a temperature of %.1f°%s",
			condition,
			temp,
			map[string]string{"celsius": "C", "fahrenheit": "F"}[units],
		),
	}
}
