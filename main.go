package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

const version = "1.0.2"
const gitUrl = "https://github.com/idesyatov/http-runner"

func main() {
	// Definition of flags
	showVersion := flag.Bool("version", false, "Show version")
	method := flag.String("method", "GET", "HTTP method to use (e.g., GET, POST)")
	url := flag.String("url", "", "Target URL")
	count := flag.Int("count", 1, "Number of requests to send")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent requests (default is 10)")
	headers := flag.String("headers", "", "Comma-separated list of headers in the format key:value")

	// Parsing flags
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("GitHub: %s\n", gitUrl)
		return
	}

	if *url == "" {
		log.Fatal("URL must be provided")
	}

	// Parse headers
	var parsedHeaders map[string]string
	if *headers != "" {
		parsedHeaders = make(map[string]string)
		for _, header := range strings.Split(*headers, ",") {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				parsedHeaders[key] = value
			} else {
				log.Printf("Invalid header format: %s", header)
			}
		}
	}

	cfg := config.NewConfig() // You can add a configuration here
	client := httpclient.NewClient(cfg.Timeout)
	gen := generator.NewGenerator(client)

	gen.GenerateRequests(*method, *url, *count, *verbose, *concurrency, parsedHeaders) 
}
