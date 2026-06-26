package generator

import (
	"context"
	"errors"
	"fmt"
	"github.com/idesyatov/http-runner/pkg/httpclient"
	"io"
	"math"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
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
	Duration      time.Duration     // If >0, run for this wall-clock time instead of Count
	Rate          int               // Target requests per second (0 = unlimited)
}

type GeneratorReport struct {
	URL             string            // The URL of the request
	Method          string            // The HTTP method used
	Count           int               // The number of requests made
	Concurrency     int               // The level of concurrency
	TotalDuration   time.Duration     // The total duration of the request execution
	RequestsPerSec  float64           // Throughput: requests per second over the whole run
	ParsedHeaders   map[string]string // Headers passed to the request
	ParsedData      map[string]string // Data passed to the request
	AverageResponse float64           // The average response time
	P50Response     float64           // The 50th percentile (median) response time
	P90Response     float64           // The 90th percentile response time
	P95Response     float64           // The 95th percentile response time
	P99Response     float64           // The 99th percentile response time
	MinResponse     float64           // The minimum response time
	MaxResponse     float64           // The maximum response time
	SuccessCount    int               // The count of successful (2xx) responses
	SuccessRate     float64           // The success rate as a percentage
	StatusCodes     map[int]int       // A map to store status codes and their counts
	ErrorCount      int               // The number of requests that failed with a transport error
	Errors          map[string]int    // Transport errors grouped by category
}

// NewGenerator creates a new Generator instance with the provided HTTP client.
func NewGenerator(client *httpclient.Client) *Generator {
	return &Generator{Client: client}
}

// GenerateRequests generates and sends HTTP requests based on the provided
// configuration. It stops early (after in-flight requests finish) when ctx is
// cancelled, returning a report for whatever was sent.
func (g *Generator) GenerateRequests(ctx context.Context, cfg RequestConfig) GeneratorReport {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalResponseTime time.Duration // Total response time across completed requests
	var minResponseTime time.Duration   // Minimum response time recorded
	var maxResponseTime time.Duration   // Maximum response time recorded
	var completedCount int              // Requests that got an HTTP response (no transport error)
	var successCount int                // Responses with a 2xx status code
	var errorCount int                  // Requests that failed with a transport error
	var statusCodes = make(map[int]int) // Map for storing status codes
	var errorTypes = make(map[string]int)
	var responseTimes []time.Duration // Per-request response times (completed only) for percentiles
	var sentCount int                 // Requests actually launched

	startTime := time.Now() // Start of total execution time

	// Create a channel for the semaphore to limit concurrency
	semaphore := make(chan struct{}, cfg.Concurrency)

	// Optional rate limiter: at most cfg.Rate requests started per second.
	var rateCh <-chan time.Time
	if cfg.Rate > 0 {
		interval := time.Second / time.Duration(cfg.Rate)
		if interval <= 0 {
			interval = time.Nanosecond
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		rateCh = ticker.C
	}

	worker := func() {
		defer wg.Done()
		defer func() { <-semaphore }() // Release semaphore

		start := time.Now()
		// Send the request using the HTTP client
		resp, err := g.Client.SendRequest(cfg.Method, cfg.URL, cfg.ParsedHeaders, cfg.Data)
		responseTime := time.Since(start)

		// Drain and close the body so the connection can be reused (keep-alive).
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		mu.Lock()
		// Latency metrics cover every completed request (one that returned an
		// HTTP response); transport errors carry no meaningful response time.
		// "Success" is narrower: only 2xx responses count toward SuccessCount.
		if err == nil {
			completedCount++
			totalResponseTime += responseTime
			responseTimes = append(responseTimes, responseTime)
			statusCodes[resp.StatusCode]++ // Increment the counter for the status code
			if minResponseTime == 0 || responseTime < minResponseTime {
				minResponseTime = responseTime
			}
			if responseTime > maxResponseTime {
				maxResponseTime = responseTime
			}
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				successCount++
			}
		} else {
			errorCount++
			errorTypes[classifyError(err)]++
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
	}

	// launch starts a single request, honouring the rate limit and concurrency
	// semaphore. It returns false if the run should stop (context cancelled).
	launch := func() bool {
		// Prioritise cancellation: a plain select below could otherwise pick the
		// semaphore branch even when ctx is already done (Go chooses randomly
		// among ready cases).
		select {
		case <-ctx.Done():
			return false
		default:
		}
		if rateCh != nil {
			select {
			case <-ctx.Done():
				return false
			case <-rateCh:
			}
		}
		select {
		case <-ctx.Done():
			return false
		case semaphore <- struct{}{}:
		}
		wg.Add(1)
		sentCount++
		go worker()
		return true
	}

	if cfg.Duration > 0 {
		deadline := startTime.Add(cfg.Duration)
		for time.Now().Before(deadline) {
			if !launch() {
				break
			}
		}
	} else {
		for i := 0; i < cfg.Count; i++ {
			if !launch() {
				break
			}
		}
	}
	wg.Wait()

	totalDuration := time.Since(startTime) // Total execution time

	// Statistics output
	var averageResponseTime, successRate, requestsPerSec float64
	if completedCount > 0 {
		averageResponseTime = totalResponseTime.Seconds() / float64(completedCount)
	}
	if sentCount > 0 {
		successRate = (float64(successCount) / float64(sentCount)) * 100
	}
	if totalDuration.Seconds() > 0 {
		requestsPerSec = float64(sentCount) / totalDuration.Seconds()
	}

	sort.Slice(responseTimes, func(i, j int) bool { return responseTimes[i] < responseTimes[j] })

	// Create a report using the unified Report structure
	return GeneratorReport{
		URL:             cfg.URL,
		Method:          cfg.Method,
		Count:           sentCount,
		Concurrency:     cfg.Concurrency,
		TotalDuration:   totalDuration,
		RequestsPerSec:  requestsPerSec,
		ParsedHeaders:   cfg.ParsedHeaders,
		ParsedData:      cfg.Data,
		AverageResponse: averageResponseTime,
		P50Response:     percentile(responseTimes, 50),
		P90Response:     percentile(responseTimes, 90),
		P95Response:     percentile(responseTimes, 95),
		P99Response:     percentile(responseTimes, 99),
		MinResponse:     minResponseTime.Seconds(),
		MaxResponse:     maxResponseTime.Seconds(),
		SuccessCount:    successCount,
		SuccessRate:     successRate,
		StatusCodes:     statusCodes,
		ErrorCount:      errorCount,
		Errors:          errorTypes,
	}
}

// classifyError groups a transport error into a short, human-readable category
// for the report (e.g. "timeout", "connection refused", "dns", "other").
func classifyError(err error) string {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout"
	}
	switch s := err.Error(); {
	case strings.Contains(s, "connection refused"):
		return "connection refused"
	case strings.Contains(s, "no such host"):
		return "dns"
	case strings.Contains(s, "context deadline exceeded"):
		return "timeout"
	default:
		return "other"
	}
}

// percentile returns the p-th percentile (0-100) of the ascending-sorted
// durations, in seconds, using the nearest-rank method. Returns 0 if empty.
func percentile(sorted []time.Duration, p float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	rank := int(math.Ceil(p / 100 * float64(n)))
	if rank < 1 {
		rank = 1
	}
	if rank > n {
		rank = n
	}
	return sorted[rank-1].Seconds()
}
