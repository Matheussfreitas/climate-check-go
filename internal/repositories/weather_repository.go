package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// WeatherData holds normalized current-weather data.
type WeatherData struct {
	City        string
	Country     string
	Temperature float64
	FeelsLike   float64
	TempMin     float64
	TempMax     float64
	Humidity    int
	WindSpeed   float64
	Visibility  int
	Description string
}

// ForecastItem represents a single forecast day.
type ForecastItem struct {
	Date        string
	TempMin     float64
	TempMax     float64
	Humidity    int
	WindSpeed   float64
	Description string
}

// ForecastData holds normalized forecast data.
type ForecastData struct {
	City    string
	Country string
	List    []ForecastItem
}

// WeatherRepository defines the contract for weather data access.
type WeatherRepository interface {
	GetCurrentWeather(city string) (*WeatherData, error)
	GetForecast(city string) (*ForecastData, error)
}

type weatherRepository struct {
	weatherBaseURL   string
	geocodingBaseURL string
	client           *http.Client
}

type geocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Country   string  `json:"country_code"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type openMeteoCurrentResponse struct {
	Current struct {
		Temperature2m      float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		RelativeHumidity2m int     `json:"relative_humidity_2m"`
		WindSpeed10m       float64 `json:"wind_speed_10m"`
		Visibility         float64 `json:"visibility"`
		WeatherCode        int     `json:"weather_code"`
	} `json:"current"`
	Daily struct {
		Temperature2mMin []float64 `json:"temperature_2m_min"`
		Temperature2mMax []float64 `json:"temperature_2m_max"`
	} `json:"daily"`
}

type openMeteoForecastResponse struct {
	Daily struct {
		Time                []string  `json:"time"`
		Temperature2mMin    []float64 `json:"temperature_2m_min"`
		Temperature2mMax    []float64 `json:"temperature_2m_max"`
		PrecipitationSum    []float64 `json:"precipitation_sum"`
		WindSpeed10mMax     []float64 `json:"wind_speed_10m_max"`
		RelativeHumidityAvg []int     `json:"relative_humidity_2m_mean"`
		WeatherCode         []int     `json:"weather_code"`
	} `json:"daily"`
}

type cityLocation struct {
	Name      string
	Country   string
	Latitude  float64
	Longitude float64
}

// NewWeatherRepository creates a new WeatherRepository.
func NewWeatherRepository(weatherBaseURL, geocodingBaseURL string) WeatherRepository {
	return &weatherRepository{
		weatherBaseURL:   weatherBaseURL,
		geocodingBaseURL: geocodingBaseURL,
		client:           &http.Client{},
	}
}

// GetCurrentWeather fetches current weather for the given city.
func (r *weatherRepository) GetCurrentWeather(city string) (*WeatherData, error) {
	location, err := r.resolveCity(city)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(
		"%s/forecast?latitude=%f&longitude=%f&current=temperature_2m,apparent_temperature,relative_humidity_2m,weather_code,wind_speed_10m,visibility&daily=temperature_2m_min,temperature_2m_max&timezone=auto&forecast_days=1",
		r.weatherBaseURL,
		location.Latitude,
		location.Longitude,
	)

	resp, err := r.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to call weather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var data openMeteoCurrentResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	description := weatherDescriptionFromCode(data.Current.WeatherCode)

	return &WeatherData{
		City:        location.Name,
		Country:     location.Country,
		Temperature: data.Current.Temperature2m,
		FeelsLike:   data.Current.ApparentTemperature,
		TempMin:     safeFloatAt(data.Daily.Temperature2mMin, 0),
		TempMax:     safeFloatAt(data.Daily.Temperature2mMax, 0),
		Humidity:    data.Current.RelativeHumidity2m,
		WindSpeed:   data.Current.WindSpeed10m,
		Visibility:  int(data.Current.Visibility),
		Description: description,
	}, nil
}

// GetForecast fetches the forecast for the given city.
func (r *weatherRepository) GetForecast(city string) (*ForecastData, error) {
	location, err := r.resolveCity(city)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(
		"%s/forecast?latitude=%f&longitude=%f&daily=temperature_2m_min,temperature_2m_max,relative_humidity_2m_mean,wind_speed_10m_max,weather_code,precipitation_sum&timezone=auto&forecast_days=5",
		r.weatherBaseURL,
		location.Latitude,
		location.Longitude,
	)

	resp, err := r.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to call forecast API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status %d", resp.StatusCode)
	}

	var data openMeteoForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	size := len(data.Daily.Time)
	items := make([]ForecastItem, 0, size)
	for i := 0; i < size; i++ {
		items = append(items, ForecastItem{
			Date:        safeStringAt(data.Daily.Time, i),
			TempMin:     safeFloatAt(data.Daily.Temperature2mMin, i),
			TempMax:     safeFloatAt(data.Daily.Temperature2mMax, i),
			Humidity:    safeIntAt(data.Daily.RelativeHumidityAvg, i),
			WindSpeed:   safeFloatAt(data.Daily.WindSpeed10mMax, i),
			Description: weatherDescriptionFromCode(safeIntAt(data.Daily.WeatherCode, i)),
		})
	}

	return &ForecastData{
		City:    location.Name,
		Country: location.Country,
		List:    items,
	}, nil
}

func (r *weatherRepository) resolveCity(city string) (*cityLocation, error) {
	endpoint := fmt.Sprintf("%s/search?name=%s&count=1&language=pt&format=json",
		r.geocodingBaseURL,
		url.QueryEscape(city),
	)

	resp, err := r.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to call geocoding API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var payload geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	if len(payload.Results) == 0 {
		return nil, fmt.Errorf("city '%s' not found", city)
	}

	result := payload.Results[0]
	return &cityLocation{
		Name:      result.Name,
		Country:   strings.ToUpper(result.Country),
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
	}, nil
}

func safeFloatAt(values []float64, idx int) float64 {
	if idx < 0 || idx >= len(values) {
		return 0
	}
	return values[idx]
}

func safeIntAt(values []int, idx int) int {
	if idx < 0 || idx >= len(values) {
		return 0
	}
	return values[idx]
}

func safeStringAt(values []string, idx int) string {
	if idx < 0 || idx >= len(values) {
		return ""
	}
	return values[idx]
}

func weatherDescriptionFromCode(code int) string {
	switch code {
	case 0:
		return "céu limpo"
	case 1:
		return "predominantemente limpo"
	case 2:
		return "parcialmente nublado"
	case 3:
		return "nublado"
	case 45, 48:
		return "neblina"
	case 51, 53, 55:
		return "garoa"
	case 56, 57:
		return "garoa congelante"
	case 61, 63, 65:
		return "chuva"
	case 66, 67:
		return "chuva congelante"
	case 71, 73, 75, 77:
		return "neve"
	case 80, 81, 82:
		return "pancadas de chuva"
	case 85, 86:
		return "pancadas de neve"
	case 95:
		return "trovoadas"
	case 96, 99:
		return "trovoadas com granizo"
	default:
		return "condições variáveis"
	}
}
