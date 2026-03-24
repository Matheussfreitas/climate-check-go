package main

import (
	"fmt"
	"log"

	"github.com/Matheussfreitas/climate-check-go/config"
	"github.com/Matheussfreitas/climate-check-go/internal/controllers"
	"github.com/Matheussfreitas/climate-check-go/internal/repositories"
	"github.com/Matheussfreitas/climate-check-go/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	// Repository → Service → Controller (dependency injection)
	repo := repositories.NewWeatherRepository(cfg.WeatherBaseURL, cfg.GeocodingBaseURL)
	svc := services.NewWeatherService(repo)
	ctrl := controllers.NewWeatherController(svc)

	router := gin.Default()

	api := router.Group("/api/v1")
	ctrl.RegisterRoutes(api)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("climate-check-go server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
