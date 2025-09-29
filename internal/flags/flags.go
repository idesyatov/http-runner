package flags

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
)

// Config holds the configuration options for the HTTP client application.
type Config struct {
	ShowVersion   bool              // Flag to indicate whether to display the application version.
	Method        string            // The HTTP method to use for requests (e.g., GET, POST).
	URL           string            // The target URL for the HTTP requests.
	Count         int               // The number of requests to send.
	Verbose       bool              // Flag to enable verbose output for debugging purposes.
	Concurrency   int               // The number of concurrent requests to be made.
	ParsedHeaders map[string]string // A map of HTTP headers to include in the requests.
}

// Metadata contains information about the application version and its source repository.
type Metadata struct {
	Version string // The version of the application.
	GitURL  string // The URL of the application's Git repository.
}

// ConfigFile holds the structure of the configuration file.
type ConfigFile struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Endpoint represents a single endpoint configuration.
type Endpoint struct {
	URL        string `yaml:"url"`
	Verbose    bool   `yaml:"verbose"`
	Method     string `yaml:"method"`
	Headers    string `yaml:"headers"`
	Count      int    `yaml:"count"`
	Concurrency int   `yaml:"concurrency"`
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
	configFile := flag.String("config-file", "", "Path to the configuration file")

	flag.Parse()

	if *configFile != "" {
		return loadConfigFromFile(*configFile)
	}

	return &Config{
		ShowVersion:   *showVersion,
		Method:        *method,
		URL:           *url,
		Count:         *count,
		Verbose:       *verbose,
		Concurrency:   *concurrency,
		ParsedHeaders: parseHeaders(*headers),
	}
}

// loadConfigFromFile loads configuration from a YAML file.
func loadConfigFromFile(filePath string) *Config {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	var configFile ConfigFile
	if err := yaml.Unmarshal(data, &configFile); err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		os.Exit(1)
	}

	if len(configFile.Endpoints) == 0 {
		fmt.Println("No endpoints found in the configuration file.")
		os.Exit(1)
	}

	// Assuming we take the first endpoint for simplicity
	endpoint := configFile.Endpoints[0]

	return &Config{
		ShowVersion:   false, // No version flag in file
		Method:        endpoint.Method,
		URL:           endpoint.URL,
		Count:         endpoint.Count,
		Verbose:       endpoint.Verbose,
		Concurrency:   endpoint.Concurrency,
		ParsedHeaders: parseHeaders(endpoint.Headers),
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

	if config.URL == "" && config.ParsedHeaders == nil {
		fmt.Println("The URL must be provided or a configuration file must be specified. Please use the --help flag for usage instructions.")
		os.Exit(1)
	}

	return config
}
