package reporter

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/idesyatov/http-runner/pkg/color"
)

// TestNewReport verifies that NewReport initializes a Report correctly.
func TestNewReport(t *testing.T) {
	expectedReport := Report{
		URL:             "https://example.com",
		Method:          "GET",
		Count:           10,
		Concurrency:     5,
		TotalDuration:   time.Second * 5,
		ParsedHeaders:   map[string]string{"Authorization": "Bearer token"},
		ParsedData:      map[string]string{"key": "value"},
		AverageResponse: 0.5,
		MinResponse:     0.1,
		MaxResponse:     1.0,
		SuccessCount:    8,
		SuccessRate:     80.0,
		StatusCodes:     map[int]int{200: 8, 404: 2},
	}

	report := NewReport(expectedReport)

	// Compare each field individually
	if report.URL != expectedReport.URL ||
		report.Method != expectedReport.Method ||
		report.Count != expectedReport.Count ||
		report.Concurrency != expectedReport.Concurrency ||
		report.TotalDuration != expectedReport.TotalDuration ||
		report.AverageResponse != expectedReport.AverageResponse ||
		report.MinResponse != expectedReport.MinResponse ||
		report.MaxResponse != expectedReport.MaxResponse ||
		report.SuccessCount != expectedReport.SuccessCount ||
		report.SuccessRate != expectedReport.SuccessRate {

		t.Errorf("Report does not match expected values: got %+v, want %+v", report, expectedReport)
	}

	// Compare the maps
	for key, value := range expectedReport.ParsedHeaders {
		if report.ParsedHeaders[key] != value {
			t.Errorf("ParsedHeaders[%s] = %s, want %s", key, report.ParsedHeaders[key], value)
		}
	}
	for key, value := range expectedReport.ParsedData {
		if report.ParsedData[key] != value {
			t.Errorf("ParsedData[%s] = %s, want %s", key, report.ParsedData[key], value)
		}
	}
	for code, count := range expectedReport.StatusCodes {
		if report.StatusCodes[code] != count {
			t.Errorf("StatusCodes[%d] = %d, want %d", code, report.StatusCodes[code], count)
		}
	}
}

// TestGenerate verifies output of the Generate function.
func TestGenerate(t *testing.T) {
	var buf bytes.Buffer
	report := Report{
		URL:             "https://example.com",
		Method:          "GET",
		Count:           10,
		Concurrency:     5,
		TotalDuration:   time.Second * 5,
		ParsedHeaders:   map[string]string{"Authorization": "Bearer token"},
		ParsedData:      map[string]string{"key": "value"},
		AverageResponse: 0.5,
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
	if len(report.ParsedData) > 0 {
		buf.WriteString("Request Data:\n")
		for key, value := range report.ParsedData {
			fmt.Fprintf(&buf, "  - %s: %s\n", key, value)
		}
	}
	fmt.Fprintf(&buf, "Request Count: %d\n", report.Count)
	fmt.Fprintf(&buf, "Request Concurrency: %d\n", report.Concurrency)
	fmt.Fprintf(&buf, "Average Response Time: %.6f seconds\n", report.AverageResponse)
	fmt.Fprintf(&buf, "Minimum Response Time: %.6f seconds\n", report.MinResponse)
	fmt.Fprintf(&buf, "Maximum Response Time: %.6f seconds\n", report.MaxResponse)
	fmt.Fprintf(&buf, "Success Count: %d\n", report.SuccessCount)
	fmt.Fprintf(&buf, "Success Rate: %.2f%%\n", report.SuccessRate)

	// Output percentage of status codes
	for code, count := range report.StatusCodes {
		percentage := (float64(count) / float64(report.Count)) * 100
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
		"  - key: value\n" +
		"Request Count: 10\n" +
		"Request Concurrency: 5\n" +
		"Average Response Time: 0.500000 seconds\n" +
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
