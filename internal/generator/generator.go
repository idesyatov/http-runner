package generator

import (
	"fmt"
	"sync"
	"time"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

type Generator struct {
	Client *httpclient.Client
}

func NewGenerator(client *httpclient.Client) *Generator {
	return &Generator{Client: client}
}

func (g *Generator) GenerateRequests(method, url string, count int, verbose bool, maxConcurrent int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalResponseTime time.Duration
	var minResponseTime time.Duration
	var maxResponseTime time.Duration
	var successCount int
	var statusCodes = make(map[int]int) // For storing status codes

	startTime := time.Now() // Start of total execution time

	// Create a channel for the semaphore
	semaphore := make(chan struct{}, maxConcurrent)

	for i := 0; i < count; i++ {
		wg.Add(1)

		// Acquire semaphore
		semaphore <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			start := time.Now()
			resp, err := g.Client.SendRequest(method, url)
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
			if verbose {
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
	averageResponseTime := totalResponseTime.Seconds() / float64(count)
	successRate := (float64(successCount) / float64(count)) * 100
	totalDuration := time.Since(startTime) // Total execution time

	// Output statistics
	fmt.Printf("Request URL: \033[32m%s\033[0m\n", url)
	fmt.Printf("Request Method: %s\n", method)
	fmt.Printf("Request Count: %d\n", count)
	fmt.Printf("Request Concurrency: %d\n", maxConcurrent)
	fmt.Printf("Average Response Time: %.6f seconds\n", averageResponseTime)
	fmt.Printf("Minimum Response Time: %.6f seconds\n", minResponseTime.Seconds())
	fmt.Printf("Maximum Response Time: %.6f seconds\n", maxResponseTime.Seconds())
	fmt.Printf("Success Rate: %.2f%%\n", successRate)

	// Output percentage of status codes
	for code, count := range statusCodes {
		percentage := (float64(count) / float64(count)) * 100 // Corrected
		fmt.Printf("Status Code %d: %.2f%%\n", code, percentage)
	}

	// Display total execution time
	fmt.Printf("Total Duration: %.6f seconds\n", totalDuration.Seconds())
}
