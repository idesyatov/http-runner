package flags

import (
	"encoding/json"
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

	m, ok := data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map body, got %T", data)
	}
	if m["key"] != "value" {
		t.Errorf("Expected key to be 'value', got '%v'", m["key"])
	}
}

// Test that nested JSON objects and arrays survive parsing.
func TestParseDataFromCLI_NestedJSON(t *testing.T) {
	rawData := `{"user": {"id": 1, "roles": ["admin", "ops"]}}`
	data, err := parseDataFromCLI(rawData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	m := data.(map[string]interface{})
	user := m["user"].(map[string]interface{})
	if user["id"] != float64(1) {
		t.Errorf("Expected user.id 1, got %v", user["id"])
	}
	roles := user["roles"].([]interface{})
	if len(roles) != 2 || roles[0] != "admin" {
		t.Errorf("Expected roles [admin ops], got %v", roles)
	}
}

// Test that an empty -data value yields a nil body (untyped nil).
func TestParseDataFromCLI_Empty(t *testing.T) {
	data, err := parseDataFromCLI("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if data != nil {
		t.Errorf("Expected nil body, got %v", data)
	}
}

// Test that -data @file.json reads and parses the file (curl style).
func TestParseDataFromCLI_File(t *testing.T) {
	path := t.TempDir() + "/payload.json"
	if err := os.WriteFile(path, []byte(`{"key": "value"}`), 0o600); err != nil {
		t.Fatalf("Error writing temp file: %v", err)
	}

	data, err := parseDataFromCLI("@" + path)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if data.(map[string]interface{})["key"] != "value" {
		t.Errorf("Expected key to be 'value', got %v", data)
	}
}

// Test that a missing data file returns an error.
func TestParseDataFromCLI_MissingFile(t *testing.T) {
	data, err := parseDataFromCLI("@/no/such/file.json")
	if err == nil {
		t.Fatal("Expected an error for missing file, got nil")
	}
	if data != nil {
		t.Errorf("Expected nil body on error, got %v", data)
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

// Test that normalizeYAML converts yaml.v2's map[interface{}]interface{} into a
// JSON-encodable shape, recursing through nested maps and slices.
func TestNormalizeYAML(t *testing.T) {
	in := map[interface{}]interface{}{
		"user": map[interface{}]interface{}{
			"id":    1,
			"roles": []interface{}{"admin", map[interface{}]interface{}{"scope": "all"}},
		},
	}

	out := normalizeYAML(in)
	if _, err := json.Marshal(out); err != nil {
		t.Fatalf("normalized value is not JSON-encodable: %v", err)
	}

	m, ok := out.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", out)
	}
	user := m["user"].(map[string]interface{})
	if user["id"] != 1 {
		t.Errorf("Expected user.id 1, got %v", user["id"])
	}
	roles := user["roles"].([]interface{})
	if roles[1].(map[string]interface{})["scope"] != "all" {
		t.Errorf("Expected nested roles[1].scope 'all', got %v", roles[1])
	}
}
