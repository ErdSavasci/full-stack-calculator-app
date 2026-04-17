package main

import (
	"calculator/config"

	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type MathPayload struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

func setupRouter() *gin.Engine {
	// Create the Gin router
	r := gin.Default()

	// Computing the operation
	r.POST("/compute", func(c *gin.Context) {
		var payload MathPayload

		// Unpack the body
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid numbers provided"})
			return
		}

		// Do the math
		result, err := calculateAddition(payload.First, payload.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "addition calculation failed"})
		}

		// Return the result
		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})

	return r
}

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Create the Gin router
	r := setupRouter()

	// Running the server with the specified port inside the config.yaml
	err := r.Run(":" + cfg.Services.Addition.Port)
	if err != nil {
		log.Fatal(err)
	}
}

// Perform addition (result = first + second)
func calculateAddition(first string, second string) (string, error) {
	// To compute accurately, we are using string -> Decimal conversion
	f, err1 := decimal.NewFromString(first)
	s, err2 := decimal.NewFromString(second)

	// Validate that the strings were actually valid numbers
	if err1 != nil || err2 != nil {
		return "", errors.New("invalid number format provided")
	}

	// Perform the operation
	result := f.Add(s)

	// Convert back to string
	finalString := result.Round(8).String()

	// The screen size limit
	if len(finalString) > 50 {
		return "", errors.New("result is too large to display")
	}

	return finalString, nil
}
