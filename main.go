package main

import (
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/internal/reporter"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

func main() {
	metadata := flags.Metadata{
		Version: "1.2.0",
		GitURL:  "https://github.com/idesyatov/http-runner",
	}

	cfg := flags.ParseFlags(metadata)

	httpCfg := config.NewConfig() // You can add configuration here.
	client := httpclient.NewClient(httpCfg.Timeout)
	gen := generator.NewGenerator(client)

	// Create a RequestConfig instance with the necessary parameters.
	requestConfig := generator.RequestConfig{
		Method:        cfg.Method,
		URL:           cfg.URL,
		Count:         cfg.Count,
		Verbose:       cfg.Verbose,
		Concurrency:   cfg.Concurrency,
		ParsedHeaders: cfg.ParsedHeaders,
	}

	// Generate requests based on the provided request configuration.
	generatorReport := gen.GenerateRequests(requestConfig)

	// Create a new report using the generated data and output it to the console.
	reporter.NewReport(reporter.Report{
		URL:             generatorReport.URL,
		Method:          generatorReport.Method,
		Count:           generatorReport.Count,
		Concurrency:     generatorReport.Concurrency,
		TotalDuration:   generatorReport.TotalDuration,
		ParsedHeaders:   generatorReport.ParsedHeaders,
		AverageResponse: generatorReport.AverageResponse,
		MinResponse:     generatorReport.MinResponse,
		MaxResponse:     generatorReport.MaxResponse,
		SuccessCount:    generatorReport.SuccessCount,
		SuccessRate:     generatorReport.SuccessRate,
		StatusCodes:     generatorReport.StatusCodes,
	}).Generate()
}
