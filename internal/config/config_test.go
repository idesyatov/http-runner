package config

import (
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config == nil {
		t.Fatal("Expected config to not be nil")
	}

	if config.Timeout != 5*time.Second {
		t.Fatalf("Expected timeout to be 5 seconds, got %v", config.Timeout)
	}
}
