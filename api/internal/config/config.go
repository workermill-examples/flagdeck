package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoDBURI   string `envconfig:"MONGODB_URI" default:"mongodb://localhost:27017/flagdeck"`
	RedisURL     string `envconfig:"REDIS_URL" default:"redis://localhost:6379/0"`
	JWTSecret    string `envconfig:"JWT_SECRET" default:"dev-secret-change-in-prod"`
	Port         string `envconfig:"PORT" default:"8080"`
	Environment  string `envconfig:"ENVIRONMENT" default:"development"`
	CORSOrigins  string `envconfig:"CORS_ORIGINS" default:"http://localhost:3000"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}