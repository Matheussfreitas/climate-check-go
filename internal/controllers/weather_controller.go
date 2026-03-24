package controllers

import (
	"net/http"
	"strings"

	"github.com/Matheussfreitas/climate-check-go/internal/services"
	"github.com/gin-gonic/gin"
)

// WeatherController handles HTTP requests for weather endpoints.
type WeatherController struct {
	service services.WeatherService
}

// NewWeatherController creates a new WeatherController.
func NewWeatherController(service services.WeatherService) *WeatherController {
	return &WeatherController{service: service}
}

// RegisterRoutes registers all weather-related routes on the given router group.
func (c *WeatherController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/weather", c.GetCurrentWeather)
	rg.GET("/weather/forecast", c.GetForecast)
}

// GetCurrentWeather godoc
// @Summary     Get current weather
// @Description Returns current weather conditions and a routine suggestion for the given city.
// @Param       city query string true "City name (e.g. São Paulo)"
// @Success     200 {object} services.CurrentWeather
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /weather [get]
func (c *WeatherController) GetCurrentWeather(ctx *gin.Context) {
	city := ctx.Query("city")
	if city == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'city' is required"})
		return
	}

	weather, err := c.service.GetCurrentWeather(city)
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFound(err) {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, weather)
}

// GetForecast godoc
// @Summary     Get 5-day weather forecast
// @Description Returns a 5-day forecast summary with daily suggestions for routine planning.
// @Param       city query string true "City name (e.g. São Paulo)"
// @Success     200 {object} services.ForecastSummary
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /weather/forecast [get]
func (c *WeatherController) GetForecast(ctx *gin.Context) {
	city := ctx.Query("city")
	if city == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'city' is required"})
		return
	}

	forecast, err := c.service.GetForecast(city)
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFound(err) {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, forecast)
}

// isNotFound checks whether an error message indicates a "not found" situation.
func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}
