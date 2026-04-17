package main

import (
	"calculator/config"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	// For test purpose
	gin.SetMode(gin.TestMode)

	// Passing a dummy config for testing
	dummyCfg := config.Config{}
	router := setupRouter(dummyCfg)

	t.Run("Origin allowance", func(t *testing.T) {
		// Create a mock HTTP Request
		req, _ := http.NewRequest("OPTIONS", "/api/history", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		// Assert CORS responds with 204 No Content for preflight checks
		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204 for OPTIONS, got %d", w.Code)
		}

		// Assert the Origin was explicitly allowed
		allowedOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowedOrigin != "http://localhost:3000" {
			t.Errorf("expected CORS header to allow localhost:3000, got '%s'", allowedOrigin)
		}
	})

	t.Run("Origin blockage", func(t *testing.T) {
		// Create a mock HTTP Request
		req, _ := http.NewRequest("OPTIONS", "/api/history", nil)
		req.Header.Set("Origin", "http://malicious-site.com")
		req.Header.Set("Access-Control-Request-Method", "GET")

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		// Assert CORS header is NOT set for malicious sites
		allowedOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowedOrigin == "http://malicious-site.com" {
			t.Errorf("expected CORS to block malicious site, but it was allowed")
		}
	})
}

func TestRoute(t *testing.T) {
	// For test purpose
	gin.SetMode(gin.TestMode)

	// Passing a dummy config for testing
	dummyCfg := config.Config{}
	router := setupRouter(dummyCfg)

	// Sending a PUT to /api/history
	// Expectation: It should return 404 or 405 HTTP code
	t.Run("Wrong method type part1", func(t *testing.T) {
		// Create a mock HTTP Request
		req, _ := http.NewRequest("PUT", "/api/history", nil)

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound && w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 404 or 405 for wrong method, got %d", w.Code)
		}
	})

	// Sending a GET to /api/histories
	// Expectation: It should return 404 HTTP code
	t.Run("Wrong method name part2", func(t *testing.T) {
		// Create a mock HTTP Request
		req, _ := http.NewRequest("GET", "/api/histories", nil)

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound && w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 404 or 405 for wrong method, got %d", w.Code)
		}
	})

	// Sending a POST to /api/calculate/add.
	// Expectation: It should return non 404 HTTP code
	t.Run("Correct method name and type part1", func(t *testing.T) {
		// Create the JSON payload
		payload := `{"first": "5", "second": "6"}`

		// Create a mock HTTP Request
		req, _ := http.NewRequest("POST", "/api/calculate/add", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		if w.Code == http.StatusNotFound {
			t.Errorf("expected route to be registered, but got 404 Not Found")
		}
	})

	// Sending a GET to /api/history
	// Expectation: It should return non 404 HTTP code
	t.Run("Correct method name and type part2", func(t *testing.T) {
		// Create a mock HTTP Request
		req, _ := http.NewRequest("GET", "/api/history", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		if w.Code == http.StatusNotFound {
			t.Errorf("expected route to be registered, but got 404 Not Found")
		}
	})

	// Sending a POST to /api/history
	// Expectation: It should return non 404 HTTP code
	t.Run("Correct method name and type part3", func(t *testing.T) {
		// Create the JSON payload
		payload := `{"first": "5", "second": "6", "operation": "add", "result": "11"}`

		// Create a mock HTTP Request
		req, _ := http.NewRequest("POST", "/api/history", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		if w.Code == http.StatusNotFound {
			t.Errorf("expected route to be registered, but got 404 Not Found")
		}
	})
}
