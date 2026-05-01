package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	ServerPort   string
	ProxyPort    string
	ProxyBaseURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	proxyPort := getEnv("PROXY_PORT", "9090")

	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/api_gateway?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		ProxyPort:    proxyPort,
		ProxyBaseURL: getEnv("PROXY_BASE_URL", "http://localhost:"+proxyPort),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
