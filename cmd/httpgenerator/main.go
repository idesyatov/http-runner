package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

const version = "1.0.1"
const gitUrl = "https://github.com/idesyatov/http-runner"

func main() {
	// Definition of flags
	showVersion := flag.Bool("version", false, "Show version")
	method := flag.String("method", "GET", "HTTP method to use (e.g., GET, POST)")
	url := flag.String("url", "", "Target URL")
	count := flag.Int("count", 1, "Number of requests to send")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent requests (default is 10)")

	// Parsing flags
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("URL: %s\n", gitUrl)
		return
	}

	if *url == "" {
		log.Fatal("URL must be provided")
	}

	cfg := config.NewConfig() // You can add a configuration here
	client := httpclient.NewClient(cfg.Timeout)
	gen := generator.NewGenerator(client)

	gen.GenerateRequests(*method, *url, *count, *verbose, *concurrency) 
}
