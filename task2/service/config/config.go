package config

import (
	"os"
	"time"
)

type GatewayConfig struct {
	HTTPAddr string
	GRPCAddr string
}

type ServiceConfig struct {
	GRPCAddr         string
	GitHubAPIBaseURL string
	GitHubTimeout    time.Duration
}

func LoadGatewayConfig() GatewayConfig {
	return GatewayConfig{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		GRPCAddr: getEnv("GRPC_ADDR", "service:9090"),
	}
}

func LoadServiceConfig() ServiceConfig {
	return ServiceConfig{
		GRPCAddr:         getEnv("GRPC_ADDR", ":9090"),
		GitHubAPIBaseURL: getEnv("GITHUB_API_BASE_URL", "https://api.github.com"),
		GitHubTimeout:    getEnvDuration("GITHUB_TIMEOUT", 10*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}

	return fallback
}
