package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/internal/reporter"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

// version is the application version. It is overridden at build time by
// GoReleaser via -ldflags "-X main.version=...".
var version = "1.7.1"

func main() {
	metadata := flags.Metadata{
		GitURL:  "https://github.com/idesyatov/http-runner",
		Version: version,
	}

	cfg := flags.ParseFlags(metadata)

	// Cancel the run on Ctrl-C: stop launching new requests, let in-flight ones
	// finish, and still print the report for what was sent.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Iterate over all endpoints
	for _, endpoint := range cfg.Endpoints {
		client := httpclient.NewClient(time.Duration(endpoint.Timeout), cfg.Insecure, cfg.Redirects, endpoint.Concurrency)
		gen := generator.NewGenerator(client)

		// Create RequestConfig for each endpoint
		requestConfig := generator.RequestConfig{
			Method:        endpoint.Method,
			URL:           endpoint.URL,
			Count:         endpoint.Count,
			Verbose:       endpoint.Verbose,
			Concurrency:   endpoint.Concurrency,
			ParsedHeaders: endpoint.Headers,
			Data:          endpoint.Data,
			Duration:      time.Duration(endpoint.Duration),
			Rate:          endpoint.Rate,
		}

		// Generate requests based on the configuration
		generatorReport := gen.GenerateRequests(ctx, requestConfig)

		// Create a new report using the generated data
		report := &reporter.Report{
			URL:             generatorReport.URL,
			Method:          generatorReport.Method,
			Count:           generatorReport.Count,
			Concurrency:     generatorReport.Concurrency,
			TotalDuration:   generatorReport.TotalDuration,
			RequestsPerSec:  generatorReport.RequestsPerSec,
			ParsedHeaders:   generatorReport.ParsedHeaders,
			ParsedData:      generatorReport.ParsedData,
			AverageResponse: generatorReport.AverageResponse,
			P50Response:     generatorReport.P50Response,
			P90Response:     generatorReport.P90Response,
			P95Response:     generatorReport.P95Response,
			P99Response:     generatorReport.P99Response,
			MinResponse:     generatorReport.MinResponse,
			MaxResponse:     generatorReport.MaxResponse,
			SuccessCount:    generatorReport.SuccessCount,
			SuccessRate:     generatorReport.SuccessRate,
			StatusCodes:     generatorReport.StatusCodes,
			ErrorCount:      generatorReport.ErrorCount,
			Errors:          generatorReport.Errors,
		}

		if cfg.Output == "json" {
			if err := report.GenerateJSON(); err != nil {
				fmt.Fprintln(os.Stderr, "error writing JSON report:", err)
				os.Exit(1)
			}
		} else {
			report.Generate()
		}

		// Stop processing further endpoints if the run was interrupted.
		if ctx.Err() != nil {
			break
		}
	}
}
