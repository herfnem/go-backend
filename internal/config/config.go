package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	JWTSecret      string
	JWTExpiry      time.Duration
	RequestTimeout time.Duration
	AllowedOrigins []string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtExpiry := getEnv("JWT_EXPIRY", "24h")
	parsedExpiry, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		parsedExpiry = 24 * time.Hour
	}

	return Config{
		Port:           getEnv("PORT", "8000"),
		DBPath:         getEnv("DB_PATH", "./app.db"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:      parsedExpiry,
		RequestTimeout: getDuration("REQUEST_TIMEOUT", 10*time.Second),
		AllowedOrigins: parseCSV(getEnv("ALLOWED_ORIGINS", "*")),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
