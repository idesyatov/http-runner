package flags

import (
	"os"
	"testing"
	"time"
)

// Test for loading configuration from file
func TestLoadConfigFromFile_ValidFile(t *testing.T) {
	validYaml := `
endpoints:
  - url: "http://example.com"
    verbose: true
    method: "GET"
    headers:
      Authorization: "Bearer token"
    count: 10
    concurrency: 5
    data: {}
`
	tmpFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Delete the file after the test

	_, err = tmpFile.Write([]byte(validYaml))
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}
	tmpFile.Close()

	// Loading configuration from file
	config := loadConfigFromFile(tmpFile.Name())

	if config == nil {
		t.Fatal("Expected config to not be nil")
	}

	if len(config.Endpoints) != 1 {
		t.Fatalf("Expected 1 endpoint, got %d", len(config.Endpoints))
	}

	if config.Endpoints[0].URL != "http://example.com" {
		t.Fatalf("Expected URL to be 'http://example.com', got '%s'", config.Endpoints[0].URL)
	}
}

// Test that omitted fields get default values when loaded from file
func TestLoadConfigFromFile_AppliesDefaults(t *testing.T) {
	yamlWithoutDefaults := `
endpoints:
  - url: "http://example.com"
`
	tmpFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(yamlWithoutDefaults)); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}
	tmpFile.Close()

	config := loadConfigFromFile(tmpFile.Name())

	ep := config.Endpoints[0]
	if ep.Method != "GET" {
		t.Errorf("Expected default Method 'GET', got '%s'", ep.Method)
	}
	if ep.Count != 1 {
		t.Errorf("Expected default Count 1, got %d", ep.Count)
	}
	if ep.Concurrency != 10 {
		t.Errorf("Expected default Concurrency 10, got %d", ep.Concurrency)
	}
	if time.Duration(ep.Timeout) != 5*time.Second {
		t.Errorf("Expected default Timeout 5s, got %s", time.Duration(ep.Timeout))
	}
}

// Test that timeout/duration are parsed from YAML duration strings
func TestLoadConfigFromFile_ParsesDurations(t *testing.T) {
	yamlWithDurations := `
endpoints:
  - url: "http://example.com"
    timeout: "10s"
    duration: "30s"
    rate: 50
`
	tmpFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(yamlWithDurations)); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}
	tmpFile.Close()

	config := loadConfigFromFile(tmpFile.Name())
	ep := config.Endpoints[0]

	if time.Duration(ep.Timeout) != 10*time.Second {
		t.Errorf("Expected Timeout 10s, got %s", time.Duration(ep.Timeout))
	}
	if time.Duration(ep.Duration) != 30*time.Second {
		t.Errorf("Expected Duration 30s, got %s", time.Duration(ep.Duration))
	}
	if ep.Rate != 50 {
		t.Errorf("Expected Rate 50, got %d", ep.Rate)
	}
}

// Header Parsing Test
func TestParseHeadersFromCLI(t *testing.T) {
	rawHeaders := "Authorization: Bearer token, Content-Type: application/json"
	headers, err := parseHeadersFromCLI(rawHeaders)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(headers) != 2 {
		t.Fatalf("Expected 2 headers, got %d", len(headers))
	}

	if headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected Authorization header to be 'Bearer token', got '%s'", headers["Authorization"])
	}

	if headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', got '%s'", headers["Content-Type"])
	}
}

// Test that incorrect header formatting returns an error instead of being skipped silently
func TestParseHeadersFromCLI_InvalidFormat(t *testing.T) {
	rawHeaders := "Authorization: Bearer token, InvalidHeader"
	headers, err := parseHeadersFromCLI(rawHeaders)

	if err == nil {
		t.Fatal("Expected an error for invalid header format, got nil")
	}
	if headers != nil {
		t.Errorf("Expected nil map on error, got %v", headers)
	}
}

// Data parsing test
func TestParseDataFromCLI_ValidJSON(t *testing.T) {
	rawData := `{"key": "value"}`
	data, err := parseDataFromCLI(rawData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data["key"] != "value" {
		t.Errorf("Expected key to be 'value', got '%s'", data["key"])
	}
}

// Test that invalid data format returns an error instead of being ignored silently
func TestParseDataFromCLI_InvalidJSON(t *testing.T) {
	rawData := `{"key": "value",}`
	data, err := parseDataFromCLI(rawData) // Invalid JSON

	if err == nil {
		t.Fatal("Expected an error for invalid JSON, got nil")
	}
	if data != nil {
		t.Errorf("Expected nil map on error, got %v", data)
	}
}
