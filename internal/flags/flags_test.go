package flags

import (
	"os"
	"testing"
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

// Header Parsing Test
func TestParseHeadersFromCLI(t *testing.T) {
	rawHeaders := "Authorization: Bearer token, Content-Type: application/json"
	headers := parseHeadersFromCLI(rawHeaders)

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

// Test for incorrect header formatting
func TestParseHeadersFromCLI_InvalidFormat(t *testing.T) {
	rawHeaders := "Authorization: Bearer token, InvalidHeader"
	headers := parseHeadersFromCLI(rawHeaders)

	if len(headers) != 1 {
		t.Fatalf("Expected 1 header, got %d", len(headers))
	}

	if headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected Authorization header to be 'Bearer token', got '%s'", headers["Authorization"])
	}
}

// Data parsing test
func TestParseDataFromCLI_ValidJSON(t *testing.T) {
	rawData := `{"key": "value"}`
	data := parseDataFromCLI(rawData)

	if data["key"] != "value" {
		t.Errorf("Expected key to be 'value', got '%s'", data["key"])
	}
}

// Test for invalid data format
func TestParseDataFromCLI_InvalidJSON(t *testing.T) {
	rawData := `{"key": "value",}`
	data := parseDataFromCLI(rawData) // Invalid JSON

	if len(data) != 0 {
		t.Error("Expected empty map for invalid JSON")
	}
}
