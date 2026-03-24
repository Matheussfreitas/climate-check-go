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
			Name: "Rio de Janeiro",
			Sys:  struct{ Country string `json:"country"` }{Country: "BR"},
			Main: struct {
				Temp      float64 `json:"temp"`
				FeelsLike float64 `json:"feels_like"`
				TempMin   float64 `json:"temp_min"`
				TempMax   float64 `json:"temp_max"`
				Humidity  int     `json:"humidity"`
			}{Temp: 32.0, FeelsLike: 35.0, TempMin: 28.0, TempMax: 34.0, Humidity: 80},
			Weather: []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
			}{{Main: "Clear", Description: "céu limpo"}},
			Wind:       struct{ Speed float64 `json:"speed"` }{Speed: 5.0},
			Visibility: 10000,
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
			City: struct {
				Name    string `json:"name"`
				Country string `json:"country"`
			}{Name: "Fortaleza", Country: "BR"},
			List: []repositories.ForecastItem{
				{
					DtTxt: "2024-06-01 12:00:00",
					Main: struct {
						Temp     float64 `json:"temp"`
						TempMin  float64 `json:"temp_min"`
						TempMax  float64 `json:"temp_max"`
						Humidity int     `json:"humidity"`
					}{Temp: 30.0, TempMin: 27.0, TempMax: 33.0, Humidity: 75},
					Weather: []struct {
						Main        string `json:"main"`
						Description string `json:"description"`
					}{{Main: "Clear", Description: "céu limpo"}},
					Wind: struct{ Speed float64 `json:"speed"` }{Speed: 4.0},
				},
				{
					DtTxt: "2024-06-01 15:00:00", // same day — should be skipped
					Main: struct {
						Temp     float64 `json:"temp"`
						TempMin  float64 `json:"temp_min"`
						TempMax  float64 `json:"temp_max"`
						Humidity int     `json:"humidity"`
					}{Temp: 31.0, TempMin: 28.0, TempMax: 34.0, Humidity: 70},
					Weather: []struct {
						Main        string `json:"main"`
						Description string `json:"description"`
					}{{Main: "Clear", Description: "ensolarado"}},
					Wind: struct{ Speed float64 `json:"speed"` }{Speed: 3.0},
				},
				{
					DtTxt: "2024-06-02 12:00:00",
					Main: struct {
						Temp     float64 `json:"temp"`
						TempMin  float64 `json:"temp_min"`
						TempMax  float64 `json:"temp_max"`
						Humidity int     `json:"humidity"`
					}{Temp: 28.0, TempMin: 25.0, TempMax: 30.0, Humidity: 80},
					Weather: []struct {
						Main        string `json:"main"`
						Description string `json:"description"`
					}{{Main: "Rain", Description: "chuva"}},
					Wind: struct{ Speed float64 `json:"speed"` }{Speed: 6.0},
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
