package repositories_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Matheussfreitas/climate-check-go/internal/repositories"
)

func TestGetCurrentWeather_Success(t *testing.T) {
	weatherMux := http.NewServeMux()
	weatherMux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"temperature_2m":       25.0,
				"relative_humidity_2m": 70,
				"weather_code":         3,
				"wind_speed_10m":       3.5,
			},
		})
	})
	weatherTS := httptest.NewServer(weatherMux)
	defer weatherTS.Close()

	geoMux := http.NewServeMux()
	geoMux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"name":         "São Paulo",
					"country_code": "br",
					"latitude":     -23.55,
					"longitude":    -46.63,
				},
			},
		})
	})
	geoTS := httptest.NewServer(geoMux)
	defer geoTS.Close()

	repo := repositories.NewWeatherRepository(weatherTS.URL, geoTS.URL)
	data, err := repo.GetCurrentWeather("São Paulo")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.City != "São Paulo" {
		t.Errorf("expected city 'São Paulo', got '%s'", data.City)
	}
	if data.Temperature != 25.0 {
		t.Errorf("expected temp 25.0, got %f", data.Temperature)
	}
	if data.Humidity != 70 {
		t.Errorf("expected humidity 70, got %d", data.Humidity)
	}
	if data.Description != "nublado" {
		t.Errorf("expected description 'nublado', got '%s'", data.Description)
	}
}

func TestGetCurrentWeather_CityNotFound(t *testing.T) {
	weatherTS := httptest.NewServer(http.NewServeMux())
	defer weatherTS.Close()

	geoMux := http.NewServeMux()
	geoMux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}})
	})
	geoTS := httptest.NewServer(geoMux)
	defer geoTS.Close()

	repo := repositories.NewWeatherRepository(weatherTS.URL, geoTS.URL)
	_, err := repo.GetCurrentWeather("InvalidCity")

	if err == nil {
		t.Fatal("expected error for not-found city, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestGetForecast_Success(t *testing.T) {
	weatherMux := http.NewServeMux()
	weatherMux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"daily": map[string]interface{}{
				"time":                     []string{"2024-01-01", "2024-01-02"},
				"temperature_2m_min":       []float64{15.0, 13.0},
				"temperature_2m_max":       []float64{20.0, 19.0},
				"relative_humidity_2m_mean": []int{65, 75},
				"wind_speed_10m_max":       []float64{2.0, 4.0},
				"weather_code":             []int{0, 63},
				"precipitation_sum":        []float64{0.0, 3.1},
			},
		})
	})
	weatherTS := httptest.NewServer(weatherMux)
	defer weatherTS.Close()

	geoMux := http.NewServeMux()
	geoMux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"name":         "Curitiba",
					"country_code": "BR",
					"latitude":     -25.42,
					"longitude":    -49.27,
				},
			},
		})
	})
	geoTS := httptest.NewServer(geoMux)
	defer geoTS.Close()

	repo := repositories.NewWeatherRepository(weatherTS.URL, geoTS.URL)
	data, err := repo.GetForecast("Curitiba")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.City != "Curitiba" {
		t.Errorf("expected city 'Curitiba', got '%s'", data.City)
	}
	if len(data.List) != 2 {
		t.Errorf("expected 2 forecast items, got %d", len(data.List))
	}
	if data.List[1].Description != "chuva" {
		t.Errorf("expected rain description for 2nd day, got '%s'", data.List[1].Description)
	}
}

func TestGetForecast_CityNotFound(t *testing.T) {
	weatherTS := httptest.NewServer(http.NewServeMux())
	defer weatherTS.Close()

	geoMux := http.NewServeMux()
	geoMux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}})
	})
	geoTS := httptest.NewServer(geoMux)
	defer geoTS.Close()

	repo := repositories.NewWeatherRepository(weatherTS.URL, geoTS.URL)
	_, err := repo.GetForecast("NoSuchCity")

	if err == nil {
		t.Fatal("expected error for not-found city, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}
