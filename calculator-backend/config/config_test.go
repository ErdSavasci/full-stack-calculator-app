package config

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	// Create a fake YAML string
	mockYAML := []byte(`
server:
  port: "8080"
services:
  addition:
    port: "8081"
`)

	// Temporarily reset Viper to a clean slate
	viper.Reset()
	viper.SetConfigType("yaml")

	// Trick Viper into reading from memory instead of the hard drive!
	err := viper.ReadConfig(bytes.NewBuffer(mockYAML))
	if err != nil {
		t.Fatalf("Failed to read mock config from buffer: %v", err)
	}

	cfg := LoadConfig()

	// Assertions
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected server port to be '8080', got '%s'", cfg.Server.Port)
	}
	if cfg.Services.Addition.Port != "8081" {
		t.Errorf("Expected addition URL to be '8081', got '%s'", cfg.Services.Addition.Port)
	}
}
