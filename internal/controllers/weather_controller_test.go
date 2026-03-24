package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Matheussfreitas/climate-check-go/internal/controllers"
	"github.com/Matheussfreitas/climate-check-go/internal/services"
	"github.com/gin-gonic/gin"
)

// mockWeatherService is a test double for WeatherService.
type mockWeatherService struct {
	currentWeather *services.CurrentWeather
	forecast       *services.ForecastSummary
	err            error
}

func (m *mockWeatherService) GetCurrentWeather(_ string) (*services.CurrentWeather, error) {
	return m.currentWeather, m.err
}

func (m *mockWeatherService) GetForecast(_ string) (*services.ForecastSummary, error) {
	return m.forecast, m.err
}

func setupRouter(svc services.WeatherService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	ctrl := controllers.NewWeatherController(svc)
	ctrl.RegisterRoutes(router.Group("/api/v1"))
	return router
}

func TestGetCurrentWeather_Success(t *testing.T) {
	svc := &mockWeatherService{
		currentWeather: &services.CurrentWeather{
			City:        "Brasília",
			Country:     "BR",
			Temperature: 24.0,
			Description: "parcialmente nublado",
			Suggestion:  "Temperatura agradável. Ótimo para atividades ao ar livre.",
		},
	}

	router := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather?city=Bras%C3%ADlia", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result services.CurrentWeather
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.City != "Brasília" {
		t.Errorf("expected city 'Brasília', got '%s'", result.City)
	}
}

func TestGetCurrentWeather_MissingCity(t *testing.T) {
	router := setupRouter(&mockWeatherService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetCurrentWeather_NotFound(t *testing.T) {
	svc := &mockWeatherService{err: fmt.Errorf("city 'XYZ' not found")}
	router := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather?city=XYZ", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestGetCurrentWeather_InternalError(t *testing.T) {
	svc := &mockWeatherService{err: fmt.Errorf("some internal error")}
	router := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather?city=Somewhere", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestGetForecast_Success(t *testing.T) {
	svc := &mockWeatherService{
		forecast: &services.ForecastSummary{
			City:    "Porto Alegre",
			Country: "BR",
			Days: []services.DaySummary{
				{Date: "2024-07-01", TempMin: 10.0, TempMax: 15.0, Description: "frio", Suggestion: "Faz frio."},
			},
		},
	}

	router := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather/forecast?city=Porto+Alegre", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result services.ForecastSummary
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.City != "Porto Alegre" {
		t.Errorf("expected city 'Porto Alegre', got '%s'", result.City)
	}
	if len(result.Days) != 1 {
		t.Errorf("expected 1 day, got %d", len(result.Days))
	}
}

func TestGetForecast_MissingCity(t *testing.T) {
	router := setupRouter(&mockWeatherService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather/forecast", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetForecast_NotFound(t *testing.T) {
	svc := &mockWeatherService{err: fmt.Errorf("city 'Unknown' not found")}
	router := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather/forecast?city=Unknown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
