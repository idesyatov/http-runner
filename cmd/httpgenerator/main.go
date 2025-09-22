package main

import (
	"flag"
	"github.com/idesyatov/http-runner/internal/config"
	"github.com/idesyatov/http-runner/internal/generator"
	"github.com/idesyatov/http-runner/pkg/httpclient"
	"log"
)

func main() {
	// Определение флагов
	method := flag.String("method", "GET", "HTTP method to use (e.g., GET, POST)")
	url := flag.String("url", "", "Target URL")
	count := flag.Int("count", 1, "Number of requests to send")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	// Парсинг флагов
	flag.Parse()

	if *url == "" {
		log.Fatal("URL must be provided")
	}

	cfg := config.NewConfig() // Здесь можно добавить конфигурацию
	client := httpclient.NewClient(cfg.Timeout)
	gen := generator.NewGenerator(client)

	gen.GenerateRequests(*method, *url, *count, *verbose) // Передаем флаги в генератор
}
