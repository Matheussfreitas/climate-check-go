package services

import (
	"fmt"

	"github.com/Matheussfreitas/climate-check-go/internal/repositories"
)

// CurrentWeather holds processed current weather data for the user.
type CurrentWeather struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
	WindSpeed   float64 `json:"wind_speed"`
	Visibility  int     `json:"visibility"`
	Suggestion  string  `json:"suggestion"`
}

// ForecastSummary holds a daily forecast summary for routine planning.
type ForecastSummary struct {
	City    string        `json:"city"`
	Country string        `json:"country"`
	Days    []DaySummary  `json:"days"`
}

// DaySummary represents weather conditions for a single day.
type DaySummary struct {
	Date        string  `json:"date"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
	WindSpeed   float64 `json:"wind_speed"`
	Suggestion  string  `json:"suggestion"`
}

// WeatherService defines business-logic operations on weather data.
type WeatherService interface {
	GetCurrentWeather(city string) (*CurrentWeather, error)
	GetForecast(city string) (*ForecastSummary, error)
}

type weatherService struct {
	repo repositories.WeatherRepository
}

// NewWeatherService creates a new WeatherService backed by the given repository.
func NewWeatherService(repo repositories.WeatherRepository) WeatherService {
	return &weatherService{repo: repo}
}

// GetCurrentWeather returns enriched current weather with a routine suggestion.
func (s *weatherService) GetCurrentWeather(city string) (*CurrentWeather, error) {
	if city == "" {
		return nil, fmt.Errorf("city name is required")
	}

	data, err := s.repo.GetCurrentWeather(city)
	if err != nil {
		return nil, err
	}

	return &CurrentWeather{
		City:        data.City,
		Country:     data.Country,
		Temperature: data.Temperature,
		FeelsLike:   data.Temperature,
		TempMin:     data.Temperature,
		TempMax:     data.Temperature,
		Humidity:    data.Humidity,
		Description: data.Description,
		WindSpeed:   data.WindSpeed,
		Visibility:  data.Visibility,
		Suggestion:  buildSuggestion(data.Temperature, data.Humidity, data.WindSpeed, data.Description),
	}, nil
}

// GetForecast returns a 5-day forecast summary with daily suggestions.
func (s *weatherService) GetForecast(city string) (*ForecastSummary, error) {
	if city == "" {
		return nil, fmt.Errorf("city name is required")
	}

	data, err := s.repo.GetForecast(city)
	if err != nil {
		return nil, err
	}

	// Group forecast items by date (YYYY-MM-DD) and pick one entry per day.
	seen := map[string]bool{}
	var days []DaySummary

	for _, item := range data.List {
		date := item.Date
		if seen[date] {
			continue
		}
		seen[date] = true

		days = append(days, DaySummary{
			Date:        date,
			TempMin:     item.TempMin,
			TempMax:     item.TempMax,
			Humidity:    item.Humidity,
			Description: item.Description,
			WindSpeed:   item.WindSpeed,
			Suggestion:  buildSuggestion(item.TempMax, item.Humidity, item.WindSpeed, item.Description),
		})
	}

	return &ForecastSummary{
		City:    data.City,
		Country: data.Country,
		Days:    days,
	}, nil
}

// buildSuggestion generates a routine recommendation based on weather conditions.
func buildSuggestion(temp float64, humidity int, windSpeed float64, description string) string {
	switch {
	case temp >= 30 && humidity > 70:
		return "Muito quente e úmido. Hidrate-se bastante e evite atividades físicas intensas ao ar livre."
	case temp >= 30:
		return "Dia quente. Use protetor solar, roupas leves e mantenha-se hidratado."
	case temp >= 20:
		return "Temperatura agradável. Ótimo para atividades ao ar livre."
	case temp >= 10:
		return "Temperatura amena. Use uma jaqueta leve ao sair."
	case temp < 10:
		return "Faz frio. Vista-se em camadas e mantenha-se aquecido."
	default:
		return "Verifique as condições do tempo antes de sair."
	}
}
