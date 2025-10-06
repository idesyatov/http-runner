package main

import (
	"testing"

	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

// TestMain checks that the main function runs without errors.
func TestMain(t *testing.T) {
	// Setup mock metadata
	metadata := flags.Metadata{
		GitURL:  "https://github.com/idesyatov/http-runner",
		Version: "0.0.1",
	}

	// Mock configuration
	cfg := flags.ParseFlags(metadata)
	cfg.Endpoints = []flags.Endpoint{
		{
			Method:      "GET",
			URL:         "https://example.com",
			Count:       10,
			Concurrency: 2,
			Headers:     map[string]string{},
			Data:        map[string]string{},
		},
	}

	httpCfg := config.NewConfig() // Mock configuration can be placed here.
	client := httpclient.NewClient(httpCfg.Timeout)
	gen := generator.NewGenerator(client)

	// Generate request for the first endpoint
	requestConfig := generator.RequestConfig{
		Method:        cfg.Endpoints[0].Method,
		URL:           cfg.Endpoints[0].URL,
		Count:         cfg.Endpoints[0].Count,
		Verbose:       false,
		Concurrency:   cfg.Endpoints[0].Concurrency,
		ParsedHeaders: cfg.Endpoints[0].Headers,
		Data:          cfg.Endpoints[0].Data,
	}

	// Generate requests based on the configuration
	generatorReport := gen.GenerateRequests(requestConfig)

	if generatorReport.URL != cfg.Endpoints[0].URL {
		t.Errorf("Expected URL %s, got %s", cfg.Endpoints[0].URL, generatorReport.URL)
	}

	if generatorReport.Method != cfg.Endpoints[0].Method {
		t.Errorf("Expected Method %s, got %s", cfg.Endpoints[0].Method, generatorReport.Method)
	}
}
