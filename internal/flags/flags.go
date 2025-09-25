package flags

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ShowVersion   bool
	Method        string
	URL           string
	Count         int
	Verbose       bool
	Concurrency   int
	ParsedHeaders map[string]string
}

type Metadata struct {
	Version string
	GitURL  string
}

// DefineFlags defines the flags and returns them as a Config structure.
func DefineFlags() *Config {
	showVersion := flag.Bool("version", false, "Show version")
	method := flag.String("method", "GET", "HTTP method to use (e.g., GET, POST)")
	url := flag.String("url", "", "Target URL")
	count := flag.Int("count", 1, "Number of requests to send")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent requests (default is 10)")
	headers := flag.String("headers", "", "Comma-separated list of headers in the format key:value")

	flag.Parse()

	return &Config{
		ShowVersion: *showVersion,
		Method:      *method,
		URL:         *url,
		Count:       *count,
		Verbose:     *verbose,
		Concurrency: *concurrency,
		ParsedHeaders: parseHeaders(*headers),
	}
}

// parseHeaders parses headers from a string and returns them as a map.
func parseHeaders(headers string) map[string]string {
	parsedHeaders := make(map[string]string)
	if headers != "" {
		for _, header := range strings.Split(headers, ",") {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				parsedHeaders[key] = value
			} else {
				fmt.Printf("Invalid header format: %s\n", header)
			}
		}
	}
	return parsedHeaders
}

// ParseFlags combines flag definition and condition checking.
func ParseFlags(metadata Metadata) *Config {
	config := DefineFlags()

	if config.ShowVersion {
		fmt.Printf("Version: %s\n", metadata.Version)
		fmt.Printf("GitHub: %s\n", metadata.GitURL)
		os.Exit(0)
	}

	if config.URL == "" {
		fmt.Println("The URL must be provided. Please use the --help flag for usage instructions.")
		os.Exit(1)
	}

	return config
}
