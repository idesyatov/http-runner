package generator

import (
	"fmt"
	"sync"
	"time"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

type Generator struct {
	Client *httpclient.Client // The HTTP client used for sending requests
}

// RequestConfig holds the configuration for generating requests.
type RequestConfig struct {
	Method        string            // The HTTP method to use
	URL           string            // The URL to send requests to
	Count         int               // The number of requests to generate
	Verbose       bool              // Flag to enable verbose output
	Concurrency   int               // The level of concurrency for requests
	ParsedHeaders map[string]string // Headers to include in the requests
	Data          map[string]string // Data to include in the request body
}

type GeneratorReport struct {
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

// NewGenerator creates a new Generator instance with the provided HTTP client.
func NewGenerator(client *httpclient.Client) *Generator {
	return &Generator{Client: client}
}

// GenerateRequests generates and sends HTTP requests based on the provided configuration.
func (g *Generator) GenerateRequests(cfg RequestConfig) GeneratorReport {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalResponseTime time.Duration // Total response time for all requests
	var minResponseTime time.Duration   // Minimum response time recorded
	var maxResponseTime time.Duration   // Maximum response time recorded
	var successCount int                // Count of successful requests
	var statusCodes = make(map[int]int) // Map for storing status codes

	startTime := time.Now() // Start of total execution time

	// Create a channel for the semaphore to limit concurrency
	semaphore := make(chan struct{}, cfg.Concurrency)

	for i := 0; i < cfg.Count; i++ {
		wg.Add(1)

		// Acquire semaphore
		semaphore <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			start := time.Now()
			// Send the request using the HTTP client
			resp, err := g.Client.SendRequest(cfg.Method, cfg.URL, cfg.ParsedHeaders, cfg.Data)
			responseTime := time.Since(start)

			mu.Lock()
			totalResponseTime += responseTime
			if err == nil {
				successCount++
				statusCodes[resp.StatusCode]++ // Increment the counter for the status code
				if minResponseTime == 0 || responseTime < minResponseTime {
					minResponseTime = responseTime
				}
				if responseTime > maxResponseTime {
					maxResponseTime = responseTime
				}
			}
			mu.Unlock()

			// Output response status only when verbose is enabled
			if cfg.Verbose {
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Response Status:", resp.Status)
				}
			}
		}()
	}
	wg.Wait()

	// Statistics output
	averageResponseTime := totalResponseTime.Seconds() / float64(cfg.Count)
	successRate := (float64(successCount) / float64(cfg.Count)) * 100
	totalDuration := time.Since(startTime) // Total execution time

	// Create a report using the unified Report structure
	return GeneratorReport{
		URL:             cfg.URL,
		Method:          cfg.Method,
		Count:           cfg.Count,
		Concurrency:     cfg.Concurrency,
		TotalDuration:   totalDuration,
		ParsedHeaders:   cfg.ParsedHeaders,
		ParsedData:      cfg.Data,
		AverageResponse: averageResponseTime,
		MinResponse:     minResponseTime.Seconds(),
		MaxResponse:     maxResponseTime.Seconds(),
		SuccessCount:    successCount,
		SuccessRate:     successRate,
		StatusCodes:     statusCodes,
	}
}
