package flags

import (
	"flag"
	"log"
	"strings"
)

type Config struct {
	ShowVersion  bool
	Method       string
	URL          string
	Count        int
	Verbose      bool
	Concurrency  int
	ParsedHeaders map[string]string
}

func ParseFlags() *Config {
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

	if *url == "" {
		log.Fatal("URL must be provided")
	}

	// Parsing headers
	parsedHeaders := make(map[string]string)
	if *headers != "" {
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

	return &Config{
		ShowVersion:  *showVersion,
		Method:       *method,
		URL:          *url,
		Count:        *count,
		Verbose:      *verbose,
		Concurrency:  *concurrency,
		ParsedHeaders: parsedHeaders,
	}
}
