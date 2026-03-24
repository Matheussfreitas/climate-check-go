package services_test

import (
	"fmt"
	"testing"

	"github.com/Matheussfreitas/climate-check-go/internal/repositories"
	"github.com/Matheussfreitas/climate-check-go/internal/services"
)

// mockWeatherRepository is a test double for WeatherRepository.
type mockWeatherRepository struct {
	currentWeather *repositories.WeatherData
	forecastData   *repositories.ForecastData
	err            error
}

func (m *mockWeatherRepository) GetCurrentWeather(_ string) (*repositories.WeatherData, error) {
	return m.currentWeather, m.err
}

func (m *mockWeatherRepository) GetForecast(_ string) (*repositories.ForecastData, error) {
	return m.forecastData, m.err
}

func TestWeatherService_GetCurrentWeather_Success(t *testing.T) {
	mock := &mockWeatherRepository{
		currentWeather: &repositories.WeatherData{
			City:        "Rio de Janeiro",
			Country:     "BR",
			Temperature: 32.0,
			FeelsLike:   35.0,
			TempMin:     28.0,
			TempMax:     34.0,
			Humidity:    80,
			Description: "céu limpo",
			WindSpeed:   5.0,
			Visibility:  10000,
		},
	}

	svc := services.NewWeatherService(mock)
	result, err := svc.GetCurrentWeather("Rio de Janeiro")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.City != "Rio de Janeiro" {
		t.Errorf("expected city 'Rio de Janeiro', got '%s'", result.City)
	}
	if result.Temperature != 32.0 {
		t.Errorf("expected temperature 32.0, got %f", result.Temperature)
	}
	if result.FeelsLike != 35.0 {
		t.Errorf("expected feels like 35.0, got %f", result.FeelsLike)
	}
	if result.TempMin != 28.0 || result.TempMax != 34.0 {
		t.Errorf("unexpected min/max temps: min=%f max=%f", result.TempMin, result.TempMax)
	}
	if result.Description != "céu limpo" {
		t.Errorf("unexpected description: %s", result.Description)
	}
	if result.Suggestion == "" {
		t.Error("expected a non-empty suggestion")
	}
}

func TestWeatherService_GetCurrentWeather_EmptyCity(t *testing.T) {
	svc := services.NewWeatherService(&mockWeatherRepository{})
	_, err := svc.GetCurrentWeather("")

	if err == nil {
		t.Fatal("expected error for empty city, got nil")
	}
}

func TestWeatherService_GetCurrentWeather_RepoError(t *testing.T) {
	mock := &mockWeatherRepository{err: fmt.Errorf("city 'X' not found")}
	svc := services.NewWeatherService(mock)
	_, err := svc.GetCurrentWeather("X")

	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}
}

func TestWeatherService_GetForecast_Success(t *testing.T) {
	mock := &mockWeatherRepository{
		forecastData: &repositories.ForecastData{
			City:    "Fortaleza",
			Country: "BR",
			List: []repositories.ForecastItem{
				{
					Date:        "2024-06-01",
					TempMin:     27.0,
					TempMax:     33.0,
					Humidity:    75,
					Description: "céu limpo",
					WindSpeed:   4.0,
				},
				{
					Date:        "2024-06-01", // same day — should be skipped
					TempMin:     28.0,
					TempMax:     34.0,
					Humidity:    70,
					Description: "ensolarado",
					WindSpeed:   3.0,
				},
				{
					Date:        "2024-06-02",
					TempMin:     25.0,
					TempMax:     30.0,
					Humidity:    80,
					Description: "chuva",
					WindSpeed:   6.0,
				},
			},
		},
	}

	svc := services.NewWeatherService(mock)
	result, err := svc.GetForecast("Fortaleza")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.City != "Fortaleza" {
		t.Errorf("expected city 'Fortaleza', got '%s'", result.City)
	}
	// The two entries for 2024-06-01 should be deduplicated to one day.
	if len(result.Days) != 2 {
		t.Errorf("expected 2 unique days, got %d", len(result.Days))
	}
	if result.Days[0].Date != "2024-06-01" {
		t.Errorf("expected first day '2024-06-01', got '%s'", result.Days[0].Date)
	}
}

func TestWeatherService_GetForecast_EmptyCity(t *testing.T) {
	svc := services.NewWeatherService(&mockWeatherRepository{})
	_, err := svc.GetForecast("")

	if err == nil {
		t.Fatal("expected error for empty city, got nil")
	}
}

func TestWeatherService_GetForecast_RepoError(t *testing.T) {
	mock := &mockWeatherRepository{err: fmt.Errorf("city 'Y' not found")}
	svc := services.NewWeatherService(mock)
	_, err := svc.GetForecast("Y")

	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}
}
