package main

import (
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/internal/reporter"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

// version is the application version. It is overridden at build time by
// GoReleaser via -ldflags "-X main.version=...".
var version = "1.5.0"

func main() {
	metadata := flags.Metadata{
		GitURL:  "https://github.com/idesyatov/http-runner",
		Version: version,
	}

	cfg := flags.ParseFlags(metadata)

	httpCfg := config.NewConfig() // You can add configuration here.
	client := httpclient.NewClient(httpCfg.Timeout)
	gen := generator.NewGenerator(client)

	// Iterate over all endpoints
	for _, endpoint := range cfg.Endpoints {
		// Create RequestConfig for each endpoint
		requestConfig := generator.RequestConfig{
			Method:        endpoint.Method,
			URL:           endpoint.URL,
			Count:         endpoint.Count,
			Verbose:       endpoint.Verbose,
			Concurrency:   endpoint.Concurrency,
			ParsedHeaders: endpoint.Headers,
			Data:          endpoint.Data,
		}

		// Generate requests based on the configuration
		generatorReport := gen.GenerateRequests(requestConfig)

		// Create a new report using the generated data
		reporter.NewReport(reporter.Report{
			URL:             generatorReport.URL,
			Method:          generatorReport.Method,
			Count:           generatorReport.Count,
			Concurrency:     generatorReport.Concurrency,
			TotalDuration:   generatorReport.TotalDuration,
			ParsedHeaders:   generatorReport.ParsedHeaders,
			ParsedData:      generatorReport.ParsedData,
			AverageResponse: generatorReport.AverageResponse,
			MinResponse:     generatorReport.MinResponse,
			MaxResponse:     generatorReport.MaxResponse,
			SuccessCount:    generatorReport.SuccessCount,
			SuccessRate:     generatorReport.SuccessRate,
			StatusCodes:     generatorReport.StatusCodes,
		}).Generate()
	}
}
