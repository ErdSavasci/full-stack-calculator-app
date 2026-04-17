package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetSetEndpoint(t *testing.T) {
	// For test purpose
	gin.SetMode(gin.TestMode)

	// Create the Gin router
	router := setupRouter()

	t.Run("Saving log", func(t *testing.T) {
		// Reset global state
		store.logs = nil

		// Create the JSON payload
		payload := CalculationLog{
			First:     "100",
			Second:    "2",
			Operation: "multiply",
			Result:    "200",
		}
		jsonValue, _ := json.Marshal(payload)

		// Create a mock HTTP Request
		req, _ := http.NewRequest("POST", "/save", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		// Assert the HTTP Status Code
		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		// Assert Server-Generated Fields (UUID and Timestamp)
		var savedResponse CalculationLog
		err := json.Unmarshal(w.Body.Bytes(), &savedResponse)
		if err != nil {
			t.Fatalf("failed to parse response JSON: %v", err)
		}
		if savedResponse.ID == "" {
			t.Errorf("expected server to generate a UUID, but it was empty")
		}
		if savedResponse.Timestamp.IsZero() {
			t.Errorf("expected server to generate a Timestamp, but it was zero")
		}
	})

	t.Run("Getting logs", func(t *testing.T) {
		// Seeding: Instead of relying on the previous test, we are providing sample data
		store.logs = []CalculationLog{
			{ID: "123", Operation: "add", Result: "30"},
		}

		// Create a mock HTTP Request
		req, _ := http.NewRequest("GET", "/get", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create a mock HTTP Response Recorder
		w := httptest.NewRecorder()

		// Send the request to the router
		router.ServeHTTP(w, req)

		// Assert the HTTP Status Code
		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		// Assert GET Body
		var fetchedLogs []CalculationLog
		err := json.Unmarshal(w.Body.Bytes(), &fetchedLogs)
		if err != nil {
			t.Fatalf("failed to parse response JSON: %v", err)
		}
		if len(fetchedLogs) != 1 {
			t.Errorf("expected exactly 1 log from GET, got %d", len(fetchedLogs))
		}
		if fetchedLogs[0].Result != store.logs[0].Result {
			t.Errorf("expected fetched log result to be '200', got '%s'", fetchedLogs[0].Result)
		}
	})
}

func TestHistoryStore(t *testing.T) {
	// Reset global state for a clean test environment
	store.logs = nil

	// Create mock logs
	log1 := CalculationLog{First: "10", Second: "5", Operation: "add", Result: "15"}
	log2 := CalculationLog{First: "20", Second: "4", Operation: "divide", Result: "5"}

	// 3. Add them to the store
	store.add(log1)
	store.add(log2)

	// Retrieve them
	logs := store.getAll()

	// Size check
	if len(logs) != 2 {
		t.Errorf("expected 2 logs in store, got %d", len(logs))
	}

	// Verify "Prepend" Logic (The last added one will be in the first place)
	if logs[0].Operation != "divide" {
		t.Errorf("expected newest log to be 'divide' at index 0, got '%s'", logs[0].Operation)
	}
}
