package repositories_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Matheussfreitas/climate-check-go/internal/repositories"
)

func newWeatherServer(statusCode int, body interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetCurrentWeather_Success(t *testing.T) {
	payload := map[string]interface{}{
		"name": "São Paulo",
		"sys":  map[string]string{"country": "BR"},
		"main": map[string]interface{}{
			"temp":       25.0,
			"feels_like": 26.0,
			"temp_min":   22.0,
			"temp_max":   28.0,
			"humidity":   70,
		},
		"weather":    []map[string]string{{"main": "Clouds", "description": "nublado"}},
		"wind":       map[string]float64{"speed": 3.5},
		"visibility": 10000,
	}

	ts := newWeatherServer(http.StatusOK, payload)
	defer ts.Close()

	repo := repositories.NewWeatherRepository("test-key", ts.URL)
	data, err := repo.GetCurrentWeather("São Paulo")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.Name != "São Paulo" {
		t.Errorf("expected city 'São Paulo', got '%s'", data.Name)
	}
	if data.Main.Temp != 25.0 {
		t.Errorf("expected temp 25.0, got %f", data.Main.Temp)
	}
	if data.Main.Humidity != 70 {
		t.Errorf("expected humidity 70, got %d", data.Main.Humidity)
	}
	if len(data.Weather) == 0 || data.Weather[0].Description != "nublado" {
		t.Errorf("unexpected weather description")
	}
}

func TestGetCurrentWeather_NotFound(t *testing.T) {
	ts := newWeatherServer(http.StatusNotFound, nil)
	defer ts.Close()

	repo := repositories.NewWeatherRepository("test-key", ts.URL)
	_, err := repo.GetCurrentWeather("InvalidCity")

	if err == nil {
		t.Fatal("expected error for not-found city, got nil")
	}
}

func TestGetCurrentWeather_InvalidAPIKey(t *testing.T) {
	ts := newWeatherServer(http.StatusUnauthorized, nil)
	defer ts.Close()

	repo := repositories.NewWeatherRepository("bad-key", ts.URL)
	_, err := repo.GetCurrentWeather("São Paulo")

	if err == nil {
		t.Fatal("expected error for invalid API key, got nil")
	}
}

func TestGetForecast_Success(t *testing.T) {
	payload := map[string]interface{}{
		"city": map[string]string{"name": "Curitiba", "country": "BR"},
		"list": []map[string]interface{}{
			{
				"dt_txt": "2024-01-01 12:00:00",
				"main": map[string]interface{}{
					"temp":     18.0,
					"temp_min": 15.0,
					"temp_max": 20.0,
					"humidity": 65,
				},
				"weather": []map[string]string{{"main": "Clear", "description": "céu limpo"}},
				"wind":    map[string]float64{"speed": 2.0},
			},
			{
				"dt_txt": "2024-01-02 12:00:00",
				"main": map[string]interface{}{
					"temp":     16.0,
					"temp_min": 13.0,
					"temp_max": 19.0,
					"humidity": 75,
				},
				"weather": []map[string]string{{"main": "Rain", "description": "chuva leve"}},
				"wind":    map[string]float64{"speed": 4.0},
			},
		},
	}

	ts := newWeatherServer(http.StatusOK, payload)
	defer ts.Close()

	repo := repositories.NewWeatherRepository("test-key", ts.URL)
	data, err := repo.GetForecast("Curitiba")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.City.Name != "Curitiba" {
		t.Errorf("expected city 'Curitiba', got '%s'", data.City.Name)
	}
	if len(data.List) != 2 {
		t.Errorf("expected 2 forecast items, got %d", len(data.List))
	}
}

func TestGetForecast_NotFound(t *testing.T) {
	ts := newWeatherServer(http.StatusNotFound, nil)
	defer ts.Close()

	repo := repositories.NewWeatherRepository("test-key", ts.URL)
	_, err := repo.GetForecast("NoSuchCity")

	if err == nil {
		t.Fatal("expected error for not-found city, got nil")
	}
}
