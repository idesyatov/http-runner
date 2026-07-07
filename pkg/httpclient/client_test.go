package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient checks that the NewClient function initializes the client with the correct timeout.
func TestNewClient(t *testing.T) {
	timeout := 5 * time.Second
	client := NewClient(timeout, false, true, 10)

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
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer testServer.Close()

	client := NewClient(5*time.Second, false, true, 10)

	headers := map[string]string{
		"Authorization": "Bearer token",
	}

	data := map[string]string{
		"key": "value",
	}

	// Send request
	resp, _, err := client.SendRequest(http.MethodPost, testServer.URL+"/test", headers, data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestNewClient_HTTP2Enabled checks that HTTP/2 is force-enabled even though a
// custom TLSClientConfig is set (which would otherwise disable it).
func TestNewClient_HTTP2Enabled(t *testing.T) {
	client := NewClient(5*time.Second, false, true, 10)
	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.Transport)
	}
	if !tr.ForceAttemptHTTP2 {
		t.Error("expected ForceAttemptHTTP2 to be true")
	}
}

// TestSendRequest_Trace checks that connection phase timings are captured: a
// fresh connection records a positive TTFB and connect time and is not reused.
func TestSendRequest_Trace(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	client := NewClient(5*time.Second, false, true, 10)

	resp, trace, err := client.SendRequest(http.MethodGet, testServer.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if trace == nil {
		t.Fatal("expected non-nil trace")
	}
	if trace.TTFB <= 0 {
		t.Errorf("expected positive TTFB, got %v", trace.TTFB)
	}
	if trace.Connect <= 0 {
		t.Errorf("expected positive connect time on a new connection, got %v", trace.Connect)
	}
	if trace.Reused {
		t.Error("expected a fresh connection not to be reused")
	}
}

// TestSendRequest_UserContentTypeHonoured checks that a user-supplied
// Content-Type in any case is not overwritten by the default application/json.
func TestSendRequest_UserContentTypeHonoured(t *testing.T) {
	var got string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	client := NewClient(5*time.Second, false, true, 10)

	headers := map[string]string{"content-type": "text/plain"}
	data := map[string]string{"key": "value"}

	resp, _, err := client.SendRequest(http.MethodPost, testServer.URL, headers, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if got != "text/plain" {
		t.Errorf("Expected Content-Type text/plain to be honoured, got %q", got)
	}
}

// TestSendRequest_NestedBody checks that an arbitrary nested JSON body is sent
// verbatim (objects, arrays, numbers), not just flat string:string pairs.
func TestSendRequest_NestedBody(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode body: %v", err)
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}

		user, ok := body["user"].(map[string]interface{})
		if !ok || user["id"] != float64(1) {
			t.Errorf("Expected body.user.id 1, got %v", body["user"])
			http.Error(w, "Invalid body content", http.StatusBadRequest)
			return
		}
		roles, ok := user["roles"].([]interface{})
		if !ok || len(roles) != 2 || roles[0] != "admin" {
			t.Errorf("Expected body.user.roles [admin ops], got %v", user["roles"])
			http.Error(w, "Invalid body content", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	client := NewClient(5*time.Second, false, true, 10)

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    1,
			"roles": []interface{}{"admin", "ops"},
		},
	}

	resp, _, err := client.SendRequest(http.MethodPost, testServer.URL, nil, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
