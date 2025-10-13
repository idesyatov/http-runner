package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config holds the configuration options for the HTTP client application.
type Config struct {
	ShowVersion bool       // Flag to indicate whether to display the application version.
	Endpoints   []Endpoint // List of endpoints to process.
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
	URL         string            `yaml:"url"`
	Verbose     bool              `yaml:"verbose"`
	Method      string            `yaml:"method"`
	Headers     map[string]string `yaml:"headers"`
	Count       int               `yaml:"count"`
	Concurrency int               `yaml:"concurrency"`
	Data        map[string]string `yaml:"data"`
}

// DefineFlags defines the flags and returns them as a Config structure.
func DefineFlags() *Config {
	showVersion := flag.Bool("version", false, "Show version")
	configFile := flag.String("config-file", "", "Path to the configuration file")

	// Defining flags for endpoints.
	url := flag.String("url", "", "Target URL for the requests.")
	count := flag.Int("count", 1, "Number of requests to send.")
	verbose := flag.Bool("verbose", false, "Enable verbose output.")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent requests to send.")
	method := flag.String("method", "GET", "HTTP method to use (e.g., GET, POST). Default is GET.")
	headers := flag.String("headers", "", "Comma-separated list of headers in the format key:value.")
	data := flag.String("data", "", "JSON string of data to send in the request body.")

	flag.Parse()

	var endpoints []Endpoint

	if *configFile != "" {
		config := loadConfigFromFile(*configFile)
		endpoints = config.Endpoints
	}

	// If the configuration file is not specified, we use flags.
	if len(endpoints) == 0 {
		endpoints = append(endpoints, Endpoint{
			URL:         *url,
			Verbose:     *verbose,
			Method:      *method,
			Headers:     parseHeadersFromCLI(*headers),
			Count:       *count,
			Concurrency: *concurrency,
			Data:        parseDataFromCLI(*data),
		})
	}

	return &Config{
		ShowVersion: *showVersion,
		Endpoints:   endpoints,
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

	// Set the default value for the request method.
	for i := range configFile.Endpoints {
		if configFile.Endpoints[i].Method == "" {
			configFile.Endpoints[i].Method = "GET"
		}
	}

	return &Config{
		ShowVersion: false, // No version flag in file
		Endpoints:   configFile.Endpoints,
	}
}

// parseHeadersFromCLI parses headers from a string and returns them as a map.
func parseHeadersFromCLI(headers string) map[string]string {
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

// parseDataFromCLI parses data from a JSON string and returns it as a map.
func parseDataFromCLI(data string) map[string]string {
	parsedData := make(map[string]string)
	if data != "" {
		if err := json.Unmarshal([]byte(data), &parsedData); err != nil {
			fmt.Printf("Invalid data format: %s\n", err)
		}
	}
	return parsedData
}

// ParseFlags combines flag definition and condition checking.
func ParseFlags(metadata Metadata) *Config {
	config := DefineFlags()

	if config.ShowVersion {
		fmt.Printf("Version: %s\n", metadata.Version)
		fmt.Printf("GitHub: %s\n", metadata.GitURL)
		os.Exit(0)
	}

	if len(config.Endpoints) == 0 || config.Endpoints[0].URL == "" {
		flag.Usage()
		os.Exit(1)
	}

	return config
}
