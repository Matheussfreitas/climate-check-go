package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	WeatherAPIKey  string
	WeatherBaseURL string
	ServerPort     string
}

// Load reads environment variables (from a .env file if present) and returns a Config.
func Load() (*Config, error) {
	// Ignore error – .env is optional (environment variables may already be set).
	_ = godotenv.Load()

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("WEATHER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openweathermap.org/data/2.5"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		WeatherAPIKey:  apiKey,
		WeatherBaseURL: baseURL,
		ServerPort:     port,
	}, nil
}
