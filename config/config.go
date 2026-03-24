package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	WeatherBaseURL   string
	GeocodingBaseURL string
	ServerPort       string
}

// Load reads environment variables (from a .env file if present) and returns a Config.
func Load() (*Config, error) {
	// Ignore error – .env is optional (environment variables may already be set).
	_ = godotenv.Load()

	baseURL := os.Getenv("WEATHER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.open-meteo.com/v1"
	}

	geocodingBaseURL := os.Getenv("GEOCODING_BASE_URL")
	if geocodingBaseURL == "" {
		geocodingBaseURL = "https://geocoding-api.open-meteo.com/v1"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		WeatherBaseURL:   baseURL,
		GeocodingBaseURL: geocodingBaseURL,
		ServerPort:       port,
	}, nil
}
