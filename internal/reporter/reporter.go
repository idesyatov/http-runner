package reporter

import (
	"fmt"
	"github.com/idesyatov/http-runner/pkg/color"
	"time"
)

// Report contains all data needed for generating a report.
type Report struct {
	URL             string            // The URL of the request
	Method          string            // The HTTP method used
	Count           int               // The number of requests made
	Concurrency     int               // The level of concurrency
	TotalDuration   time.Duration     // The total duration of the request execution
	ParsedHeaders   map[string]string // Headers passed to the request
	ParsedData      map[string]string // Data passed to the request
	AverageResponse float64           // The average response time
	MinResponse     float64           // The minimum response time
	MaxResponse     float64           // The maximum response time
	SuccessCount    int               // The count of successful requests
	SuccessRate     float64           // The success rate as a percentage
	StatusCodes     map[int]int       // A map to store status codes and their counts
}

// NewReport creates a new report from the given Report structure.
func NewReport(report Report) *Report {
	return &report // Return the report structure
}

// Generate outputs the report to the console.
func (r *Report) Generate() {
	fmt.Printf("Request URL: %s\n", color.Colorize(color.Green, r.URL))
	fmt.Printf("Request Method: %s\n", r.Method)

	// Output headers if they exist
	if len(r.ParsedHeaders) > 0 {
		fmt.Println("Request Headers:")
		for key, value := range r.ParsedHeaders {
			fmt.Printf("  - %s: %s\n", key, value)
		}
	}
	// Output data if they exist
	if len(r.ParsedData) > 0 {
		fmt.Println("Request Data:")
		for key, value := range r.ParsedData {
			fmt.Printf("  - %s: %s\n", key, value)
		}
	}
	fmt.Printf("Request Count: %d\n", r.Count)
	fmt.Printf("Request Concurrency: %d\n", r.Concurrency)
	fmt.Printf("Average Response Time: %.6f seconds\n", r.AverageResponse)
	fmt.Printf("Minimum Response Time: %.6f seconds\n", r.MinResponse)
	fmt.Printf("Maximum Response Time: %.6f seconds\n", r.MaxResponse)
	fmt.Printf("Success Count: %d\n", r.SuccessCount)
	fmt.Printf("Success Rate: %.2f%%\n", r.SuccessRate)

	// Output percentage of status codes
	for code, count := range r.StatusCodes {
		percentage := (float64(count) / float64(r.Count)) * 100
		fmt.Printf("Status Code %d: %.2f%%\n", code, percentage)
	}

	// Output total execution time
	fmt.Printf("Total Duration: %.6f seconds\n\n", r.TotalDuration.Seconds())
}
