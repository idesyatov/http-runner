package generator_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

// MockClient simulates the behavior of an HTTP client for testing
type MockClient struct {
	Response *http.Response
	Error    error
}

// RoundTrip simulates sending an HTTP request and returns a response or error
func (m *MockClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

// TestGenerateRequests_Success tests the successful request generation
func TestGenerateRequests_Success(t *testing.T) {
	// Mock successful response
	mockResp := &http.Response{StatusCode: 200}
	// Create a mock client
	mockClient := &MockClient{Response: mockResp}
	// Create generator with mocked client
	gen := generator.NewGenerator(&httpclient.Client{Client: http.Client{Transport: mockClient}})

	cfg := generator.RequestConfig{
		Method:      "GET",
		URL:         "https://example.com",
		Count:       10,
		Verbose:     false,
		Concurrency: 5,
	}
	report := gen.GenerateRequests(cfg) // Generate requests

	// Assertions to verify expected outcomes
	if report.SuccessCount != 10 {
		t.Errorf("expected 10 successful requests, got %d", report.SuccessCount)
	}
	if report.SuccessRate != 100.0 {
		t.Errorf("expected success rate 100.0%%, got %f", report.SuccessRate)
	}
	if report.StatusCodes[200] != 10 {
		t.Errorf("expected 10 status code 200, got %d", report.StatusCodes[200])
	}
}

// TestGenerateRequests_Failure tests the handling of request failures
func TestGenerateRequests_Failure(t *testing.T) {
	// Mock client to simulate an error
	mockClient := &MockClient{Error: errors.New("network error")}
	// Create generator with mocked client
	gen := generator.NewGenerator(&httpclient.Client{Client: http.Client{Transport: mockClient}})

	cfg := generator.RequestConfig{
		Method:      "GET",
		URL:         "https://example.com",
		Count:       5,
		Verbose:     false,
		Concurrency: 2,
	}
	report := gen.GenerateRequests(cfg) // Generate requests

	// Assertions to verify expected outcomes
	if report.SuccessCount != 0 {
		t.Errorf("expected 0 successful requests, got %d", report.SuccessCount)
	}
	if report.SuccessRate != 0.0 {
		t.Errorf("expected success rate 0.0%%, got %f", report.SuccessRate)
	}
	if len(report.StatusCodes) != 0 {
		t.Errorf("expected no status codes, got %d", len(report.StatusCodes))
	}
}
