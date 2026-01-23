package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBPath     string
	JWTSecret  string
	JWTExpiry  string
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Port:      getEnv("PORT", "8000"),
		DBPath:    getEnv("DB_PATH", "./app.db"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry: getEnv("JWT_EXPIRY", "24h"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
