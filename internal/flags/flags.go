package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const defaultTimeout = 5 * time.Second

// Config holds the configuration options for the HTTP client application.
type Config struct {
	ShowVersion bool       // Flag to indicate whether to display the application version.
	Output      string     // Output format: "text" or "json".
	Insecure    bool       // Skip TLS certificate verification.
	Redirects   bool       // Follow HTTP redirects.
	Endpoints   []Endpoint // List of endpoints to process.
}

// Duration wraps time.Duration so it can be unmarshalled from a YAML string
// such as "10s" or "500ms".
type Duration time.Duration

// UnmarshalYAML parses a duration written as a string (e.g. "10s").
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	*d = Duration(v)
	return nil
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
	Data        interface{}       `yaml:"data"`
	Timeout     Duration          `yaml:"timeout"`  // Per-request timeout (e.g. "10s").
	Duration    Duration          `yaml:"duration"` // Run for this wall-clock time instead of Count.
	Rate        int               `yaml:"rate"`     // Target requests per second (0 = unlimited).
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
	timeout := flag.String("timeout", defaultTimeout.String(), "Per-request timeout (e.g. 10s, 500ms).")
	loadDuration := flag.String("duration", "", "Run for this wall-clock duration instead of -count (e.g. 30s).")
	rate := flag.Int("rate", 0, "Target requests per second (0 = unlimited).")
	output := flag.String("output", "text", "Output format: text or json.")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification.")
	redirects := flag.Bool("redirects", true, "Follow HTTP redirects.")

	flag.Parse()

	if *output != "text" && *output != "json" {
		fmt.Fprintf(os.Stderr, "invalid -output %q (expected text or json)\n", *output)
		os.Exit(1)
	}

	var endpoints []Endpoint

	if *configFile != "" {
		config := loadConfigFromFile(*configFile)
		endpoints = config.Endpoints
	}

	// If the configuration file is not specified, we use flags.
	if len(endpoints) == 0 {
		parsedHeaders, err := parseHeadersFromCLI(*headers)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		parsedData, err := parseDataFromCLI(*data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		timeoutDur, err := parseDuration(*timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -timeout: %s\n", err)
			os.Exit(1)
		}
		if timeoutDur == 0 {
			timeoutDur = defaultTimeout
		}
		loadDur, err := parseDuration(*loadDuration)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -duration: %s\n", err)
			os.Exit(1)
		}
		endpoints = append(endpoints, Endpoint{
			URL:         *url,
			Verbose:     *verbose,
			Method:      *method,
			Headers:     parsedHeaders,
			Count:       *count,
			Concurrency: *concurrency,
			Data:        parsedData,
			Timeout:     Duration(timeoutDur),
			Duration:    Duration(loadDur),
			Rate:        *rate,
		})
	}

	return &Config{
		ShowVersion: *showVersion,
		Output:      *output,
		Insecure:    *insecure,
		Redirects:   *redirects,
		Endpoints:   endpoints,
	}
}

// parseDuration parses a duration string; an empty string yields 0 (disabled).
func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	return time.ParseDuration(s)
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

	// Apply default values for omitted fields.
	for i := range configFile.Endpoints {
		if configFile.Endpoints[i].Method == "" {
			configFile.Endpoints[i].Method = "GET"
		}
		if configFile.Endpoints[i].Count == 0 {
			configFile.Endpoints[i].Count = 1
		}
		if configFile.Endpoints[i].Concurrency == 0 {
			configFile.Endpoints[i].Concurrency = 10
		}
		if configFile.Endpoints[i].Timeout == 0 {
			configFile.Endpoints[i].Timeout = Duration(defaultTimeout)
		}
		configFile.Endpoints[i].Data = normalizeYAML(configFile.Endpoints[i].Data)
	}

	return &Config{
		ShowVersion: false, // No version flag in file
		Endpoints:   configFile.Endpoints,
	}
}

// parseHeadersFromCLI parses headers from a string and returns them as a map.
func parseHeadersFromCLI(headers string) (map[string]string, error) {
	parsedHeaders := make(map[string]string)
	if headers != "" {
		for _, header := range strings.Split(headers, ",") {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				parsedHeaders[key] = value
			} else {
				return nil, fmt.Errorf("invalid header format: %q (expected key:value)", header)
			}
		}
	}
	return parsedHeaders, nil
}

// parseDataFromCLI parses the -data value into an arbitrary JSON value. An empty
// string yields a nil body. A value starting with "@" is treated as a path to a
// file containing the JSON (curl style); otherwise the value itself is the JSON.
func parseDataFromCLI(data string) (interface{}, error) {
	if data == "" {
		return nil, nil
	}
	src := []byte(data)
	if strings.HasPrefix(data, "@") {
		b, err := os.ReadFile(strings.TrimPrefix(data, "@"))
		if err != nil {
			return nil, fmt.Errorf("reading data file: %w", err)
		}
		src = b
	}
	var v interface{}
	if err := json.Unmarshal(src, &v); err != nil {
		return nil, fmt.Errorf("invalid data format: %w", err)
	}
	return v, nil
}

// normalizeYAML converts a value decoded by gopkg.in/yaml.v2 into a JSON-encodable
// shape: yaml.v2 decodes nested objects as map[interface{}]interface{}, which
// json.Marshal cannot handle. Keys are stringified and the conversion recurses
// through nested maps and slices.
func normalizeYAML(v interface{}) interface{} {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(val))
		for k, item := range val {
			m[fmt.Sprintf("%v", k)] = normalizeYAML(item)
		}
		return m
	case []interface{}:
		for i, item := range val {
			val[i] = normalizeYAML(item)
		}
		return val
	default:
		return v
	}
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
