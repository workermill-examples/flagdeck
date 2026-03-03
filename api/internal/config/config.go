package config

import (
	"os"
)

type Config struct {
	Port       string
	MongodbURI string
	RedisURL   string
	JWTSecret  string
}

func Load() *Config {
	cfg := &Config{
		Port:       getEnvOrDefault("PORT", "3000"),
		MongodbURI: getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017"),
		RedisURL:   getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:  getEnvOrDefault("JWT_SECRET", "your-secret-key"),
	}
	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
