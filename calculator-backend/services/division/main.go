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
		result, err := calculateDivision(payload.First, payload.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "division calculation failed"})
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

	// Get router
	r := setupRouter()

	// Running the server with the specified port inside the config.yaml
	err := r.Run(":" + cfg.Services.Division.Port)
	if err != nil {
		log.Fatal(err)
	}
}

// Perform division (result = first / second)
func calculateDivision(first string, second string) (string, error) {
	// To compute accurately, we are using string -> Decimal conversion
	f, err1 := decimal.NewFromString(first)
	s, err2 := decimal.NewFromString(second)

	// Validate that the strings were actually valid numbers
	if err1 != nil || err2 != nil {
		return "", errors.New("invalid number format provided")
	}

	if s.Sign() == 0 {
		return "", errors.New("division by zero is undefined")
	}

	// Perform the operation
	result := f.Div(s)

	// Convert back to string
	finalString := result.Round(8).String()

	// The screen size limit
	if len(finalString) > 50 {
		return "", errors.New("result is too large to display")
	}

	return finalString, nil
}
