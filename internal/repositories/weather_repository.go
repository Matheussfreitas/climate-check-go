package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// WeatherData holds the raw current-weather response from OpenWeatherMap.
type WeatherData struct {
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
	} `json:"sys"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Visibility int `json:"visibility"`
}

// ForecastData holds the raw 5-day/3-hour forecast response from OpenWeatherMap.
type ForecastData struct {
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`
	List []ForecastItem `json:"list"`
}

// ForecastItem represents a single forecast entry.
type ForecastItem struct {
	DtTxt string `json:"dt_txt"`
	Main  struct {
		Temp     float64 `json:"temp"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

// WeatherRepository defines the contract for weather data access.
type WeatherRepository interface {
	GetCurrentWeather(city string) (*WeatherData, error)
	GetForecast(city string) (*ForecastData, error)
}

type weatherRepository struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewWeatherRepository creates a new WeatherRepository.
func NewWeatherRepository(apiKey, baseURL string) WeatherRepository {
	return &weatherRepository{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// GetCurrentWeather fetches current weather for the given city.
func (r *weatherRepository) GetCurrentWeather(city string) (*WeatherData, error) {
	endpoint := fmt.Sprintf(
		"%s/weather?q=%s&appid=%s&units=metric&lang=pt_br",
		r.baseURL,
		url.QueryEscape(city),
		r.apiKey,
	)

	resp, err := r.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to call weather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("city '%s' not found", city)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid API key")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return &data, nil
}

// GetForecast fetches the 5-day forecast for the given city.
func (r *weatherRepository) GetForecast(city string) (*ForecastData, error) {
	endpoint := fmt.Sprintf(
		"%s/forecast?q=%s&appid=%s&units=metric&lang=pt_br",
		r.baseURL,
		url.QueryEscape(city),
		r.apiKey,
	)

	resp, err := r.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to call forecast API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("city '%s' not found", city)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid API key")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status %d", resp.StatusCode)
	}

	var data ForecastData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	return &data, nil
}
