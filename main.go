package main

import (
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)



func main() {
	metadata := flags.Metadata{
		Version: "1.0.3",
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

	// Pass the RequestConfig instance to GenerateRequests.
	gen.GenerateRequests(requestConfig)
}