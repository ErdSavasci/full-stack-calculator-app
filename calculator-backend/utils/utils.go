package utils

import (
	"calculator/config"

	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HandleCalculation(c *gin.Context, cfg config.Config) {
	// Get the operation from the URL
	operation := c.Param("operation")

	// Determining which microservices to proxy the request
	var targetHost string
	switch operation {
	case "add":
		targetHost = "http://addition-service:" + cfg.Services.Addition.Port
	case "subtract":
		targetHost = "http://subtraction-service:" + cfg.Services.Subtraction.Port
	case "multiply":
		targetHost = "http://multiplication-service:" + cfg.Services.Multiplication.Port
	case "divide":
		targetHost = "http://division-service:" + cfg.Services.Division.Port
	case "exponential":
		targetHost = "http://exponentiation-service:" + cfg.Services.Exponentiation.Port
	case "squareroot":
		targetHost = "http://squareroot-service:" + cfg.Services.SquareRoot.Port
	case "percentage":
		targetHost = "http://percentage-service:" + cfg.Services.Percentage.Port
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Operation not found"})
		return
	}

	// string to URL
	targetURL, err := url.Parse(targetHost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal gateway error"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Path Rewrite: (from /api/calculate/<operation> to /compute)
	c.Request.URL.Path = "/compute"

	proxy.ServeHTTP(c.Writer, c.Request)

	return
}

func HandleHistory(c *gin.Context, cfg config.Config) {
	// Determining the microservice to proxy the request
	targetHost := "http://history-service:" + cfg.Services.History.Port

	// string to URL
	targetURL, err := url.Parse(targetHost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal gateway error"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Check the request method, and proxy accordingly
	if c.Request.Method == "GET" {
		// Path Rewrite: (from /api/history to /get)
		c.Request.URL.Path = "/get"
	} else if c.Request.Method == "POST" {
		// Path Rewrite: (from /api/history to /save)
		c.Request.URL.Path = "/save"
	} else if c.Request.Method == "DELETE" {
		// Path Rewrite: (from /api/history to /clear)
		c.Request.URL.Path = "/clear"
	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
