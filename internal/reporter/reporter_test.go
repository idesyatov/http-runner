package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/idesyatov/http-runner/pkg/color"
)

// TestGenerate verifies output of the Generate function.
func TestGenerate(t *testing.T) {
	var buf bytes.Buffer
	report := Report{
		URL:             "https://example.com",
		Method:          "GET",
		Count:           10,
		Concurrency:     5,
		TotalDuration:   time.Second * 5,
		RequestsPerSec:  20.0,
		ParsedHeaders:   map[string]string{"Authorization": "Bearer token"},
		ParsedData:      map[string]string{"key": "value"},
		AverageResponse: 0.5,
		P50Response:     0.4,
		P90Response:     0.8,
		P95Response:     0.9,
		P99Response:     0.99,
		MinResponse:     0.1,
		MaxResponse:     1.0,
		SuccessCount:    8,
		SuccessRate:     80.0,
		StatusCodes:     map[int]int{200: 8, 404: 2},
	}

	// Capture the output
	fmt.Fprintf(&buf, "Request URL: %s\n", color.Colorize(color.Green, report.URL))
	fmt.Fprintf(&buf, "Request Method: %s\n", report.Method)
	if len(report.ParsedHeaders) > 0 {
		buf.WriteString("Request Headers:\n")
		for key, value := range report.ParsedHeaders {
			fmt.Fprintf(&buf, "  - %s: %s\n", key, value)
		}
	}
	if report.ParsedData != nil {
		buf.WriteString("Request Data:\n")
		if b, err := json.MarshalIndent(report.ParsedData, "  ", "  "); err == nil {
			fmt.Fprintf(&buf, "  %s\n", b)
		}
	}
	fmt.Fprintf(&buf, "Request Count: %d\n", report.Count)
	fmt.Fprintf(&buf, "Request Concurrency: %d\n", report.Concurrency)
	fmt.Fprintf(&buf, "Requests/sec: %.2f\n", report.RequestsPerSec)
	fmt.Fprintf(&buf, "Average Response Time: %.6f seconds\n", report.AverageResponse)
	fmt.Fprintf(&buf, "p50 Response Time: %.6f seconds\n", report.P50Response)
	fmt.Fprintf(&buf, "p90 Response Time: %.6f seconds\n", report.P90Response)
	fmt.Fprintf(&buf, "p95 Response Time: %.6f seconds\n", report.P95Response)
	fmt.Fprintf(&buf, "p99 Response Time: %.6f seconds\n", report.P99Response)
	fmt.Fprintf(&buf, "Minimum Response Time: %.6f seconds\n", report.MinResponse)
	fmt.Fprintf(&buf, "Maximum Response Time: %.6f seconds\n", report.MaxResponse)
	fmt.Fprintf(&buf, "Success Count: %d\n", report.SuccessCount)
	fmt.Fprintf(&buf, "Success Rate: %.2f%%\n", report.SuccessRate)

	// Output percentage of status codes in ascending order for stable output
	codes := make([]int, 0, len(report.StatusCodes))
	for code := range report.StatusCodes {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	for _, code := range codes {
		percentage := (float64(report.StatusCodes[code]) / float64(report.Count)) * 100
		fmt.Fprintf(&buf, "Status Code %d: %.2f%%\n", code, percentage)
	}

	// Total execution time
	fmt.Fprintf(&buf, "Total Duration: %.6f seconds\n\n", report.TotalDuration.Seconds())

	// Check the expected output
	expectedOutput := "Request URL: " + color.Colorize(color.Green, report.URL) + "\n" +
		"Request Method: " + report.Method + "\n" +
		"Request Headers:\n" +
		"  - Authorization: Bearer token\n" +
		"Request Data:\n" +
		"  {\n" +
		"    \"key\": \"value\"\n" +
		"  }\n" +
		"Request Count: 10\n" +
		"Request Concurrency: 5\n" +
		"Requests/sec: 20.00\n" +
		"Average Response Time: 0.500000 seconds\n" +
		"p50 Response Time: 0.400000 seconds\n" +
		"p90 Response Time: 0.800000 seconds\n" +
		"p95 Response Time: 0.900000 seconds\n" +
		"p99 Response Time: 0.990000 seconds\n" +
		"Minimum Response Time: 0.100000 seconds\n" +
		"Maximum Response Time: 1.000000 seconds\n" +
		"Success Count: 8\n" +
		"Success Rate: 80.00%\n" +
		"Status Code 200: 80.00%\n" +
		"Status Code 404: 20.00%\n" +
		"Total Duration: 5.000000 seconds\n\n"

	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot output:\n%s", expectedOutput, buf.String())
	}
}

// TestReportJSON verifies the machine-readable JSON output.
func TestReportJSON(t *testing.T) {
	report := Report{
		URL:             "https://example.com",
		Method:          "GET",
		Count:           10,
		Concurrency:     5,
		TotalDuration:   time.Second * 5,
		RequestsPerSec:  20.0,
		AverageResponse: 0.5,
		P95Response:     0.9,
		SuccessCount:    8,
		SuccessRate:     80.0,
		StatusCodes:     map[int]int{200: 8, 404: 2},
		ErrorCount:      1,
		Errors:          map[string]int{"timeout": 1},
		ParsedData: map[string]interface{}{
			"user": map[string]interface{}{"id": 1},
		},
	}

	b, err := report.JSON()
	if err != nil {
		t.Fatalf("JSON() returned error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if out["success_rate"] != 80.0 {
		t.Errorf("expected success_rate 80, got %v", out["success_rate"])
	}
	if out["requests_per_sec"] != 20.0 {
		t.Errorf("expected requests_per_sec 20, got %v", out["requests_per_sec"])
	}
	if out["total_duration_sec"] != 5.0 {
		t.Errorf("expected total_duration_sec 5, got %v", out["total_duration_sec"])
	}
	sc, ok := out["status_codes"].(map[string]interface{})
	if !ok || sc["200"] != float64(8) {
		t.Errorf("expected status_codes[200]=8, got %v", out["status_codes"])
	}
	errs, ok := out["errors"].(map[string]interface{})
	if !ok || errs["timeout"] != float64(1) {
		t.Errorf("expected errors[timeout]=1, got %v", out["errors"])
	}
	data, ok := out["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got %v", out["data"])
	}
	user, ok := data["user"].(map[string]interface{})
	if !ok || user["id"] != float64(1) {
		t.Errorf("expected data.user.id=1, got %v", data["user"])
	}
}
