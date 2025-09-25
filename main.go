package main

import (
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

func main() {
	cfg := flags.ParseFlags()

	// Creating a client and generator
	httpCfg := config.NewConfig() // You can add configuration here
	client := httpclient.NewClient(httpCfg.Timeout)
	gen := generator.NewGenerator(client)

	gen.GenerateRequests(cfg.Method, cfg.URL, cfg.Count, cfg.Verbose, cfg.Concurrency, cfg.ParsedHeaders)
}
