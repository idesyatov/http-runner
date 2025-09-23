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

func (g *Generator) GenerateRequests(method, url string, count int, verbose bool) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalResponseTime time.Duration
	var minResponseTime time.Duration
	var maxResponseTime time.Duration
	var successCount int

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			resp, err := g.Client.SendRequest(method, url)
			responseTime := time.Since(start)

			mu.Lock()
			totalResponseTime += responseTime
			if err == nil {
				successCount++
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

	fmt.Printf("Request url: \033[32m%s\033[0m\n", url)
	fmt.Printf("Request method: %s\n", method)
	fmt.Printf("Request count: %d\n", count)
	fmt.Printf("Average Response Time: %.6f seconds\n", averageResponseTime)
	fmt.Printf("Minimum Response Time: %.6f seconds\n", minResponseTime.Seconds())
	fmt.Printf("Maximum Response Time: %.6f seconds\n", maxResponseTime.Seconds())
	fmt.Printf("Success Rate: %.2f%%\n", successRate)
}