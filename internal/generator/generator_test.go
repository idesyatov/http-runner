package generator_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

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
	report := gen.GenerateRequests(context.Background(), cfg) // Generate requests

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
	if report.RequestsPerSec <= 0 {
		t.Errorf("expected positive requests/sec, got %f", report.RequestsPerSec)
	}
	if report.P95Response < report.P50Response {
		t.Errorf("expected p95 >= p50, got p95=%f p50=%f", report.P95Response, report.P50Response)
	}
}

// TestGenerateRequests_Non2xxNotSuccess verifies that a non-2xx response is
// recorded (status code + latency) but does not count as a success.
func TestGenerateRequests_Non2xxNotSuccess(t *testing.T) {
	mockResp := &http.Response{StatusCode: 500}
	mockClient := &MockClient{Response: mockResp}
	gen := generator.NewGenerator(&httpclient.Client{Client: http.Client{Transport: mockClient}})

	cfg := generator.RequestConfig{
		Method:      "GET",
		URL:         "https://example.com",
		Count:       5,
		Concurrency: 2,
	}
	report := gen.GenerateRequests(context.Background(), cfg)

	if report.SuccessCount != 0 {
		t.Errorf("expected 0 successful (2xx) requests, got %d", report.SuccessCount)
	}
	if report.SuccessRate != 0.0 {
		t.Errorf("expected success rate 0.0%%, got %f", report.SuccessRate)
	}
	if report.StatusCodes[500] != 5 {
		t.Errorf("expected 5 status code 500, got %d", report.StatusCodes[500])
	}
	// The responses completed, so their latency must still be measured.
	if report.AverageResponse <= 0 {
		t.Errorf("expected positive average response for completed non-2xx requests, got %f", report.AverageResponse)
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
	report := gen.GenerateRequests(context.Background(), cfg) // Generate requests

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
	// Failed requests must not contribute to latency metrics.
	if report.AverageResponse != 0 {
		t.Errorf("expected average response 0 when all requests fail, got %f", report.AverageResponse)
	}
	if report.ErrorCount != 5 {
		t.Errorf("expected 5 transport errors, got %d", report.ErrorCount)
	}
	if report.Errors["other"] != 5 {
		t.Errorf("expected 5 errors classified as 'other', got %d", report.Errors["other"])
	}
}

// TestGenerateRequests_Duration verifies that duration mode sends requests for
// the configured wall-clock time, bounded by the rate limit.
func TestGenerateRequests_Duration(t *testing.T) {
	mockClient := &MockClient{Response: &http.Response{StatusCode: 200}}
	gen := generator.NewGenerator(&httpclient.Client{Client: http.Client{Transport: mockClient}})

	cfg := generator.RequestConfig{
		Method:      "GET",
		URL:         "https://example.com",
		Concurrency: 5,
		Duration:    150 * time.Millisecond,
		Rate:        20, // keep the request count small and bounded
	}
	report := gen.GenerateRequests(context.Background(), cfg)

	if report.Count < 1 {
		t.Errorf("expected at least 1 request in duration mode, got %d", report.Count)
	}
	if report.SuccessCount != report.Count {
		t.Errorf("expected all %d requests to succeed, got %d", report.Count, report.SuccessCount)
	}
}

// TestGenerateRequests_ContextCancelled verifies that an already-cancelled
// context prevents any request from being launched.
func TestGenerateRequests_ContextCancelled(t *testing.T) {
	mockClient := &MockClient{Response: &http.Response{StatusCode: 200}}
	gen := generator.NewGenerator(&httpclient.Client{Client: http.Client{Transport: mockClient}})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before the run starts

	cfg := generator.RequestConfig{
		Method:      "GET",
		URL:         "https://example.com",
		Count:       100,
		Concurrency: 5,
	}
	report := gen.GenerateRequests(ctx, cfg)

	if report.Count != 0 {
		t.Errorf("expected 0 requests with a cancelled context, got %d", report.Count)
	}
}
