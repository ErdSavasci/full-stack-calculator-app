package main

import (
	"calculator/config"

	"errors"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type MathPayload struct {
	First string `json:"first"`
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
		result, err := calculateSquareRoot(payload.First)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "squareroot calculation failed"})
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
	err := r.Run(":" + cfg.Services.SquareRoot.Port)
	if err != nil {
		log.Fatal(err)
	}
}

// Perform square root (result = √ first)
func calculateSquareRoot(first string) (string, error) {
	// To compute accurately, we are using string -> Decimal conversion
	f, err1 := strconv.ParseFloat(first, 64)

	// Validate that the strings were actually valid numbers
	if err1 != nil {
		return "", errors.New("invalid number format provided")
	}

	if f < 0 {
		return "", errors.New("cannot calculate square root of a negative number")
	}

	// This check is done before performing the operation
	// If the parameter is too massive, fallback to decimal and format it in scientific notation.
	if math.IsInf(f, 0) {
		fDec, _ := decimal.NewFromString(first)

		massiveStr := fDec.Pow(decimal.NewFromFloat(0.5)).String()
		bf, _ := new(big.Float).SetString(massiveStr)

		// Convert back to string
		// .Text('e', 7) formats to scientific notation with 7 decimal places
		return bf.Text('e', 7), nil
	}

	// Perform the operation
	result := math.Sqrt(f)

	// If the result is a tiny fraction, format it in scientific notation.
	absResult := math.Abs(result)
	if (absResult > 0 && absResult < 0.000001) || absResult > 1e15 {
		bf := new(big.Float).SetFloat64(result)

		// Convert back to string
		// .Text('e', 7) formats to scientific notation with 7 decimal places
		return bf.Text('e', 7), nil
	}

	result2 := decimal.NewFromFloat(result)

	// Convert back to string
	finalString := result2.Round(8).String()

	return finalString, nil
}
