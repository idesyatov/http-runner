package main

import (
    "github.com/idesyatov/http-runner/internal/config"
    "github.com/idesyatov/http-runner/internal/generator"
    "github.com/idesyatov/http-runner/pkg/httpclient"
    "log"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) < 4 {
        log.Fatal("Usage: http-runner <method> <url> <count>")
    }

    method := os.Args[1]
    url := os.Args[2]
    count, err := strconv.Atoi(os.Args[3])
    if err != nil {
        log.Fatal("Count must be an integer:", err)
    }

    cfg := config.NewConfig() // Здесь можно добавить конфигурацию
    client := httpclient.NewClient(cfg.Timeout)
    gen := generator.NewGenerator(client)

    gen.GenerateRequests(method, url, count)
}
