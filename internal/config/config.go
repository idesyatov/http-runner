package config

import "time"

type Config struct {
	Timeout time.Duration
}

func NewConfig() *Config {
	return &Config{
		Timeout: 5 * time.Second, // Установите таймаут по умолчанию
	}
}