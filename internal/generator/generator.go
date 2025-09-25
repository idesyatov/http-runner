package generator

import (
	"fmt"
	"sync"
	"time"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

// RequestConfig holds the configuration for generating requests.
type RequestConfig struct {
	Method        string
	URL           string
	Count         int
	Verbose       bool
	Concurrency   int
	ParsedHeaders map[string]string
}

type Generator struct {
	Client *httpclient.Client
}

// NewGenerator creates a new Generator instance with the provided HTTP client.
func NewGenerator(client *httpclient.Client) *Generator {
	return &Generator{Client: client}
}

// GenerateRequests generates and sends HTTP requests based on the provided configuration.
func (g *Generator) GenerateRequests(cfg RequestConfig) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalResponseTime time.Duration
	var minResponseTime time.Duration
	var maxResponseTime time.Duration
	var successCount int
	var statusCodes = make(map[int]int) // For storing status codes

	startTime := time.Now() // Start of total execution time

	// Create a channel for the semaphore
	semaphore := make(chan struct{}, cfg.Concurrency)

	for i := 0; i < cfg.Count; i++ {
		wg.Add(1)

		// Acquire semaphore
		semaphore <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			start := time.Now()
			// Pass headers to the SendRequest method
			resp, err := g.Client.SendRequest(cfg.Method, cfg.URL, cfg.ParsedHeaders)
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

	// Output statistics
	fmt.Printf("Request URL: \033[32m%s\033[0m\n", cfg.URL)
	fmt.Printf("Request Method: %s\n", cfg.Method)

	// Output headers if they exist
	if len(cfg.ParsedHeaders) > 0 {
		fmt.Println("Request Headers:")
		for key, value := range cfg.ParsedHeaders {
			fmt.Printf("  - %s: %s\n", key, value)
		}
	}

	fmt.Printf("Request Count: %d\n", cfg.Count)
	fmt.Printf("Request Concurrency: %d\n", cfg.Concurrency)
	fmt.Printf("Average Response Time: %.6f seconds\n", averageResponseTime)
	fmt.Printf("Minimum Response Time: %.6f seconds\n", minResponseTime.Seconds())
	fmt.Printf("Maximum Response Time: %.6f seconds\n", maxResponseTime.Seconds())
	fmt.Printf("Success Rate: %.2f%%\n", successRate)

	// Output percentage of status codes
	for code, count := range statusCodes {
		percentage := (float64(count) / float64(cfg.Count)) * 100 // Corrected
		fmt.Printf("Status Code %d: %.2f%%\n", code, percentage)
	}

	// Display total execution time
	fmt.Printf("Total Duration: %.6f seconds\n", totalDuration.Seconds())
}
