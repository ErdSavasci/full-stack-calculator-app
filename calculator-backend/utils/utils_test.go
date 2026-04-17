package utils

import (
	"calculator/config"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleCalculation(t *testing.T) {
	// For test purpose
	gin.SetMode(gin.TestMode)

	// Passing a dummy config for testing
	dummyCfg := config.Config{}
	dummyCfg.Services.Addition.Port = "8081"

	// Creating a test router
	router := gin.Default()
	router.POST("/api/calculate/:operation", func(c *gin.Context) {
		HandleCalculation(c, dummyCfg)
	})

	// Create the JSON payload
	payload := `{"first": "5", "second": "6"}`

	// Create a mock HTTP Request
	req, _ := http.NewRequest("POST", "/api/calculate/add", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a mock HTTP Response Recorder
	w := httptest.NewRecorder()

	// Send the request to the router
	router.ServeHTTP(w, req)

	// Assertion
	if w.Code == http.StatusNotFound {
		t.Errorf("Expected proxy to attempt connection, but got 404")
	}
}

func TestHandleHistory(t *testing.T) {
	// For test purpose
	gin.SetMode(gin.TestMode)

	// Passing a dummy config for testing
	dummyCfg := config.Config{}
	dummyCfg.Services.History.Port = "8088"

	// Creating a test router
	router := gin.Default()
	router.POST("/api/history", func(c *gin.Context) {
		HandleHistory(c, dummyCfg)
	})

	// Create the JSON payload
	payload := `{"first": "5", "second": "6", "operation": "add", "result": "11"}`

	// Create a mock HTTP Request
	req, _ := http.NewRequest("POST", "/api/history", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a mock HTTP Response Recorder
	w := httptest.NewRecorder()

	// Send the request to the router
	router.ServeHTTP(w, req)

	// Assertion
	if w.Code == http.StatusNotFound {
		t.Errorf("Expected proxy to attempt connection, but got 404")
	}
}
