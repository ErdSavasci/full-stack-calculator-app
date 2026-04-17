package main

import (
	"calculator/config"
	"calculator/utils"

	"log"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter(cfg config.Config) *gin.Engine {
	// Create the Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Grouping all API routes together
	api := r.Group("/api")
	{
		// For calculations (POST only)
		api.POST("/calculate/:operation", func(c *gin.Context) {
			utils.HandleCalculation(c, cfg)
		})
		// Saving the history (as optional, again)
		api.POST("/history", func(c *gin.Context) {
			utils.HandleHistory(c, cfg)
		})
		// Getting the history (as optional)
		api.GET("/history", func(c *gin.Context) {
			utils.HandleHistory(c, cfg)
		})
		// Getting the history (as optional)
		api.DELETE("/history", func(c *gin.Context) {
			utils.HandleHistory(c, cfg)
		})
	}

	return r
}

func main() {
	slog.Info("Full-Stack Calculator App is initializing...")

	// Load config
	cfg := config.LoadConfig()

	// Get router
	r := setupRouter(cfg)

	// Running the server with the specified port inside the config.yaml
	err := r.Run(":" + cfg.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Full-Stack Calculator App is initialized.")
}
