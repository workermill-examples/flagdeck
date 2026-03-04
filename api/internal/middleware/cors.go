package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func CORS(config ...CORSConfig) fiber.Handler {
	var cfg CORSConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
			"http://127.0.0.1:8080",
		}
	}

	if len(cfg.AllowedMethods) == 0 {
		cfg.AllowedMethods = []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		}
	}

	if len(cfg.AllowedHeaders) == 0 {
		cfg.AllowedHeaders = []string{
			"Content-Type",
			"Authorization",
			"X-API-Key",
			"X-Requested-With",
			"Origin",
			"Accept",
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.AllowedOrigins, ","),
		AllowMethods:     strings.Join(cfg.AllowedMethods, ","),
		AllowHeaders:     strings.Join(cfg.AllowedHeaders, ","),
		AllowCredentials: true,
		ExposeHeaders:    "X-RateLimit-Limit,X-RateLimit-Remaining,X-RateLimit-Reset",
	})
}

func DefaultCORS() fiber.Handler {
	return CORS()
}
