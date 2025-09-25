package main

import (
	"log"
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/flags"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
)

const version = "1.0.2"
const gitUrl = "https://github.com/idesyatov/http-runner"

func main() {
	cfg := flags.ParseFlags()

	if cfg.ShowVersion {
		log.Printf("Version: %s\n", version)
		log.Printf("GitHub: %s\n", gitUrl)
		return
	}

	// Creating a client and generator
	httpCfg := config.NewConfig() // You can add configuration here
	client := httpclient.NewClient(httpCfg.Timeout)
	gen := generator.NewGenerator(client)

	gen.GenerateRequests(cfg.Method, cfg.URL, cfg.Count, cfg.Verbose, cfg.Concurrency, cfg.ParsedHeaders)
}
