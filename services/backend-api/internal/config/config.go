package config

import "os"

type Config struct {
	Port         string
	InferenceURL string
}

func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("SERVER_PORT", "8080"),
		InferenceURL: getEnv("INFERENCE_URL", "http://localhost:8000"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}