package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// WeatherTool provides weather information
type WeatherTool struct {
	apiKey string
	client *http.Client
}

// NewWeatherTool creates a new weather tool
func NewWeatherTool(apiKey string) *WeatherTool {
	return &WeatherTool{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Name returns the tool name
func (t *WeatherTool) Name() string {
	return "weather"
}

// Description returns the tool description
func (t *WeatherTool) Description() string {
	return "Provides weather information for a given location"
}

// Execute executes the weather tool
// Expected params:
// - location: string (city name or coordinates)
// - units: string (metric, imperial, standard) - optional, default: metric
func (t *WeatherTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	location, ok := params["location"].(string)
	if !ok {
		return nil, fmt.Errorf("location parameter is required (string)")
	}

	units := "metric"
	if val, ok := params["units"].(string); ok {
		units = val
	}

	// If no API key is set, return mock data
	if t.apiKey == "" {
		return t.getMockWeather(location, units), nil
	}

	return t.getRealWeather(ctx, location, units)
}

func (t *WeatherTool) getRealWeather(ctx context.Context, location, units string) (interface{}, error) {
	// OpenWeatherMap API
	baseURL := "https://api.openweathermap.org/data/2.5/weather"

	params := url.Values{}
	params.Add("q", location)
	params.Add("appid", t.apiKey)
	params.Add("units", units)

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("weather API error (status %d): %s", resp.StatusCode, string(body))
	}

	var weatherData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %v", err)
	}

	return weatherData, nil
}

func (t *WeatherTool) getMockWeather(location, units string) interface{} {
	// Return mock weather data when no API key is configured
	tempSymbol := "°C"
	windSpeedUnit := "m/s"
	if units == "imperial" {
		tempSymbol = "°F"
		windSpeedUnit = "mph"
	} else if units == "standard" {
		tempSymbol = "K"
	}

	// Generate somewhat realistic mock data
	temp := 20.0
	if units == "imperial" {
		temp = 68.0
	} else if units == "standard" {
		temp = 293.15
	}

	return map[string]interface{}{
		"location":    location,
		"temperature": fmt.Sprintf("%.1f%s", temp, tempSymbol),
		"feels_like":  fmt.Sprintf("%.1f%s", temp-2, tempSymbol),
		"condition":   "Partly Cloudy",
		"description": "Scattered clouds with mild temperatures",
		"humidity":    65,
		"wind_speed":  fmt.Sprintf("%.1f %s", 5.5, windSpeedUnit),
		"pressure":    1013,
		"visibility":  10000,
		"clouds":      40,
		"mock":        true,
		"note":        "Mock data - configure WEATHER_API_KEY for real weather data",
	}
}
