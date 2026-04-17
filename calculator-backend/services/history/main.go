package main

import (
	"calculator/config"

	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CalculationLog For a single log line
type CalculationLog struct {
	ID        string    `json:"id"`
	First     string    `json:"first"`
	Second    string    `json:"second"`
	Operation string    `json:"operation"`
	Result    string    `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryStore For storage purposes
type HistoryStore struct {
	mu   sync.RWMutex
	logs []CalculationLog
}

var store = &HistoryStore{}

// Writing
func (s *HistoryStore) add(log CalculationLog) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prepend to keep newest at the top
	s.logs = append([]CalculationLog{log}, s.logs...)

	// storing top 5
	if len(s.logs) > 5 {
		s.logs = s.logs[:5]
	}
}

// Reading
func (s *HistoryStore) getAll() []CalculationLog {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.logs == nil {
		return []CalculationLog{}
	}

	return s.logs
}

func (s *HistoryStore) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = []CalculationLog{}
}

func setupRouter() *gin.Engine {
	// Create the Gin router
	r := gin.Default()

	// GET /get method gets the history logs
	r.GET("/get", func(c *gin.Context) {
		cLogs := store.getAll()
		c.JSON(http.StatusOK, cLogs)
	})

	// POST /save method saves the log into history storage
	r.POST("/save", func(c *gin.Context) {
		var cLog CalculationLog
		if err := c.ShouldBindJSON(&cLog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log format"})
			return
		}

		// Adding some extra information
		cLog.ID = uuid.New().String()
		cLog.Timestamp = time.Now()

		store.add(cLog)
		c.JSON(http.StatusCreated, cLog)
	})

	// DELETE /clear clears the history storage
	r.DELETE("/clear", func(c *gin.Context) {
		store.clear()
		c.JSON(http.StatusOK, nil)
	})

	return r
}

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Get router
	r := setupRouter()

	// Running the server with the specified port inside the config.yaml
	err := r.Run(":" + cfg.Services.History.Port)
	if err != nil {
		log.Fatal(err)
	}
}
