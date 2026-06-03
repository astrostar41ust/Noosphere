package config

import "os"

type Config struct {
	Port             string
	InferenceURL     string
	DatabaseURL      string 
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
}

func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("SERVER_PORT", "8080"),
		InferenceURL:     getEnv("INFERENCE_URL", "http://localhost:8000"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5435/noosphere?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}