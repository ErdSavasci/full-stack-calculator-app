package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestComputeEndpoint(t *testing.T) {
	// Expected result
	expectedResult := "1.0000000e+20"

	// For test purpose
	gin.SetMode(gin.TestMode)

	// Create the Gin router
	router := setupRouter()

	// Create the JSON payload
	payload := MathPayload{
		First:  "10",
		Second: "20",
	}
	jsonValue, _ := json.Marshal(payload)

	// Create a mock HTTP Request
	req, _ := http.NewRequest("POST", "/compute", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Create a mock HTTP Response Recorder
	w := httptest.NewRecorder()

	// Send the request to the router
	router.ServeHTTP(w, req)

	// Assert the HTTP Status Code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Assert the JSON Response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if response["result"] != expectedResult {
		t.Errorf("expected result '%s', got '%s'", expectedResult, response["result"])
	}
}

func TestCalculateExponentiation(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		first         string
		second        string
		expected      string
		expectedError bool
	}{
		{
			name:          "Positive numbers part1",
			first:         "5",
			second:        "10",
			expected:      "9765625",
			expectedError: false,
		},
		{
			name:          "Positive numbers part2",
			first:         "500",
			second:        "125",
			expected:      "2.3509887e+337",
			expectedError: false,
		},
		{
			name:          "Negative numbers part1",
			first:         "-50",
			second:        "25",
			expected:      "-2.9802322e+42",
			expectedError: false,
		},
		{
			name:          "Negative numbers part2",
			first:         "-50",
			second:        "-25",
			expected:      "-3.3554432e-43",
			expectedError: false,
		},
		{
			name:          "Decimal numbers part1",
			first:         "0.1",
			second:        "0.2",
			expected:      "0.63095734",
			expectedError: false,
		},
		{
			name:          "Decimal numbers part2",
			first:         "0.1",
			second:        "2",
			expected:      "0.01",
			expectedError: false,
		},
		{
			name:          "Invalid inputs part1",
			first:         "test1",
			second:        "test2",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "Invalid inputs part2",
			first:         "test1",
			second:        "3",
			expected:      "",
			expectedError: true,
		},
	}

	// Loop through the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculateExponentiation(tc.first, tc.second)

			// Check if the error is given correctly
			if (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, got: %v", tc.expectedError, err)
			}

			// Check if the result matches
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}
