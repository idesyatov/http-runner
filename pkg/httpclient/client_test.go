package httpclient

import (
	"time"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewClient checks that the NewClient function initializes the client with the correct timeout.
func TestNewClient(t *testing.T) {
	timeout := 5 * time.Second
	client := NewClient(timeout)

	if client.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.Timeout)
	}
}

// TestSendRequest checks that SendRequest sends a request with the correct method, URL, headers, and data.
func TestSendRequest(t *testing.T) {
	// Create a test server to mock requests
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method %s, got %s", http.MethodPost, r.Method)
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		// Check URL
		if r.URL.String() != "/test" {
			t.Errorf("Expected URL /test, got %s", r.URL.String())
			http.Error(w, "Invalid URL", http.StatusNotFound)
			return
		}

		// Check headers
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("Expected header Authorization: Bearer token, got %s", r.Header.Get("Authorization"))
			http.Error(w, "Invalid header", http.StatusUnauthorized)
			return
		}

		// Check body
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode body: %v", err)
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}

		if body["key"] != "value" {
			t.Errorf("Expected body key: value, got %v", body)
			http.Error(w, "Invalid body content", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer testServer.Close()

	client := NewClient(5 * time.Second)

	headers := map[string]string{
		"Authorization": "Bearer token",
	}

	data := map[string]string{
		"key": "value",
	}

	// Send request
	resp, err := client.SendRequest(http.MethodPost, testServer.URL+"/test", headers, data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

