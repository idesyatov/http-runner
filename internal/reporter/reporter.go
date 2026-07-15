package reporter

import (
	"encoding/json"
	"fmt"
	"github.com/idesyatov/http-runner/pkg/color"
	"sort"
	"strings"
	"time"
)

// Report contains all data needed for generating a report.
type Report struct {
	URL             string            // The URL of the request
	Method          string            // The HTTP method used
	Count           int               // The number of requests made
	Concurrency     int               // The level of concurrency
	TotalDuration   time.Duration     // The total duration of the request execution
	RequestsPerSec  float64           // Throughput: requests per second over the whole run
	TotalBytes      int64             // Total response body bytes read across completed requests
	BytesPerSec     float64           // Throughput: response body bytes per second over the whole run
	ParsedHeaders   map[string]string // Headers passed to the request
	ParsedData      interface{}       // Data passed to the request (arbitrary JSON)
	AverageResponse float64           // The average response time
	P50Response     float64           // The 50th percentile (median) response time
	P90Response     float64           // The 90th percentile response time
	P95Response     float64           // The 95th percentile response time
	P99Response     float64           // The 99th percentile response time
	MinResponse     float64           // The minimum response time
	MaxResponse     float64           // The maximum response time
	AvgDNS          float64           // Average DNS resolution time over new connections
	AvgConnect      float64           // Average TCP connect time over new connections
	AvgTLS          float64           // Average TLS handshake time over new connections
	AvgTTFB         float64           // Average time to first response byte over completed requests
	ConnReuseRate   float64           // Percentage of completed requests served over a reused connection
	SuccessCount    int               // The count of successful (2xx) responses
	SuccessRate     float64           // The success rate as a percentage
	StatusCodes     map[int]int       // A map to store status codes and their counts
	ErrorCount      int               // The number of requests that failed with a transport error
	Errors          map[string]int    // Transport errors grouped by category
	Histogram       []Bucket          // Latency distribution over completed requests
}

// Bucket is one bar of the latency histogram: [Start, End] seconds and how many
// completed requests fell in that range.
type Bucket struct {
	Start float64 // Lower bound in seconds (inclusive)
	End   float64 // Upper bound in seconds
	Count int     // Number of completed requests in this range
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
	// Output data as JSON if it exists
	if r.ParsedData != nil {
		fmt.Println("Request Data:")
		if b, err := json.MarshalIndent(r.ParsedData, "  ", "  "); err == nil {
			fmt.Printf("  %s\n", b)
		}
	}
	fmt.Printf("Request Count: %d\n", r.Count)
	fmt.Printf("Request Concurrency: %d\n", r.Concurrency)
	fmt.Printf("Requests/sec: %.2f\n", r.RequestsPerSec)
	fmt.Printf("Bytes/sec: %.2f (%d total)\n", r.BytesPerSec, r.TotalBytes)
	fmt.Printf("Average Response Time: %.6f seconds\n", r.AverageResponse)
	fmt.Printf("p50 Response Time: %.6f seconds\n", r.P50Response)
	fmt.Printf("p90 Response Time: %.6f seconds\n", r.P90Response)
	fmt.Printf("p95 Response Time: %.6f seconds\n", r.P95Response)
	fmt.Printf("p99 Response Time: %.6f seconds\n", r.P99Response)
	fmt.Printf("Minimum Response Time: %.6f seconds\n", r.MinResponse)
	fmt.Printf("Maximum Response Time: %.6f seconds\n", r.MaxResponse)

	// Connection phase breakdown (averages). DNS/connect/TLS are zero when every
	// request reused a pooled connection.
	fmt.Println("Latency breakdown (avg):")
	fmt.Printf("  DNS:      %.6f seconds\n", r.AvgDNS)
	fmt.Printf("  Connect:  %.6f seconds\n", r.AvgConnect)
	fmt.Printf("  TLS:      %.6f seconds\n", r.AvgTLS)
	fmt.Printf("  TTFB:     %.6f seconds\n", r.AvgTTFB)
	fmt.Printf("  Conn reuse: %.2f%%\n", r.ConnReuseRate)

	fmt.Printf("Success Count: %d\n", r.SuccessCount)
	fmt.Printf("Success Rate: %.2f%%\n", r.SuccessRate)

	// Output percentage of status codes in ascending order for stable output
	codes := make([]int, 0, len(r.StatusCodes))
	for code := range r.StatusCodes {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	for _, code := range codes {
		percentage := (float64(r.StatusCodes[code]) / float64(r.Count)) * 100
		fmt.Printf("Status Code %d: %.2f%%\n", code, percentage)
	}

	// Output transport errors grouped by category, if any
	if r.ErrorCount > 0 {
		fmt.Printf("Errors: %d\n", r.ErrorCount)
		cats := make([]string, 0, len(r.Errors))
		for cat := range r.Errors {
			cats = append(cats, cat)
		}
		sort.Strings(cats)
		for _, cat := range cats {
			fmt.Printf("  - %s: %d\n", cat, r.Errors[cat])
		}
	}

	// Latency distribution as a text histogram (hey-style), scaled to the busiest
	// bucket. Skipped when there were no completed requests.
	if len(r.Histogram) > 0 {
		maxCount := 0
		for _, b := range r.Histogram {
			if b.Count > maxCount {
				maxCount = b.Count
			}
		}
		fmt.Println("Latency distribution:")
		const barWidth = 40
		for _, b := range r.Histogram {
			bar := 0
			if maxCount > 0 {
				bar = b.Count * barWidth / maxCount
			}
			fmt.Printf("  %.4f [%4d] |%s\n", b.Start, b.Count, strings.Repeat("■", bar))
		}
	}

	// Output total execution time
	fmt.Printf("Total Duration: %.6f seconds\n\n", r.TotalDuration.Seconds())
}

// jsonReport is the machine-readable shape of a report, with durations as
// seconds and stable field names.
type jsonReport struct {
	URL                string            `json:"url"`
	Method             string            `json:"method"`
	Count              int               `json:"count"`
	Concurrency        int               `json:"concurrency"`
	TotalDurationSec   float64           `json:"total_duration_sec"`
	RequestsPerSec     float64           `json:"requests_per_sec"`
	TotalBytes         int64             `json:"total_bytes"`
	BytesPerSec        float64           `json:"bytes_per_sec"`
	Headers            map[string]string `json:"headers,omitempty"`
	Data               interface{}       `json:"data,omitempty"`
	AverageResponseSec float64           `json:"average_response_sec"`
	P50Sec             float64           `json:"p50_sec"`
	P90Sec             float64           `json:"p90_sec"`
	P95Sec             float64           `json:"p95_sec"`
	P99Sec             float64           `json:"p99_sec"`
	MinSec             float64           `json:"min_sec"`
	MaxSec             float64           `json:"max_sec"`
	AvgDNSSec          float64           `json:"avg_dns_sec"`
	AvgConnectSec      float64           `json:"avg_connect_sec"`
	AvgTLSSec          float64           `json:"avg_tls_sec"`
	AvgTTFBSec         float64           `json:"avg_ttfb_sec"`
	ConnReuseRate      float64           `json:"conn_reuse_rate"`
	SuccessCount       int               `json:"success_count"`
	SuccessRate        float64           `json:"success_rate"`
	StatusCodes        map[int]int       `json:"status_codes,omitempty"`
	ErrorCount         int               `json:"error_count"`
	Errors             map[string]int    `json:"errors,omitempty"`
	Histogram          []jsonBucket      `json:"histogram,omitempty"`
}

// jsonBucket is the machine-readable shape of a histogram bucket.
type jsonBucket struct {
	StartSec float64 `json:"start_sec"`
	EndSec   float64 `json:"end_sec"`
	Count    int     `json:"count"`
}

// JSON returns the report marshalled as indented JSON.
func (r *Report) JSON() ([]byte, error) {
	var buckets []jsonBucket
	for _, b := range r.Histogram {
		buckets = append(buckets, jsonBucket{StartSec: b.Start, EndSec: b.End, Count: b.Count})
	}
	return json.MarshalIndent(jsonReport{
		URL:                r.URL,
		Method:             r.Method,
		Count:              r.Count,
		Concurrency:        r.Concurrency,
		TotalDurationSec:   r.TotalDuration.Seconds(),
		RequestsPerSec:     r.RequestsPerSec,
		TotalBytes:         r.TotalBytes,
		BytesPerSec:        r.BytesPerSec,
		Headers:            r.ParsedHeaders,
		Data:               r.ParsedData,
		AverageResponseSec: r.AverageResponse,
		P50Sec:             r.P50Response,
		P90Sec:             r.P90Response,
		P95Sec:             r.P95Response,
		P99Sec:             r.P99Response,
		MinSec:             r.MinResponse,
		MaxSec:             r.MaxResponse,
		AvgDNSSec:          r.AvgDNS,
		AvgConnectSec:      r.AvgConnect,
		AvgTLSSec:          r.AvgTLS,
		AvgTTFBSec:         r.AvgTTFB,
		ConnReuseRate:      r.ConnReuseRate,
		SuccessCount:       r.SuccessCount,
		SuccessRate:        r.SuccessRate,
		StatusCodes:        r.StatusCodes,
		ErrorCount:         r.ErrorCount,
		Errors:             r.Errors,
		Histogram:          buckets,
	}, "", "  ")
}

// GenerateJSON prints the report as JSON to the console.
func (r *Report) GenerateJSON() error {
	b, err := r.JSON()
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
