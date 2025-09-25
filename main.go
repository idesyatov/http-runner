package main

import (
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

type Metadata struct {
	Version string
	GitURL  string
}

func main() {
	metadata := flags.Metadata{
		Version: "1.0.2",
		GitURL:  "https://github.com/idesyatov/http-runner",
	}
	
	cfg := flags.ParseFlags(metadata)

	// Creating a client and generator
	httpCfg := config.NewConfig() // You can add configuration here
	client := httpclient.NewClient(httpCfg.Timeout)
	gen := generator.NewGenerator(client)

	gen.GenerateRequests(cfg.Method, cfg.URL, cfg.Count, cfg.Verbose, cfg.Concurrency, cfg.ParsedHeaders)
}
