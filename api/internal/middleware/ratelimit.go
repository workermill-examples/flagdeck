package middleware

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RateLimitConfig struct {
	RedisClient   *redis.Client
	RequestsLimit int           // Number of requests allowed
	WindowSize    time.Duration // Time window
	KeyPrefix     string        // Redis key prefix for rate limiting
}

type RateLimitInfo struct {
	Remaining int64
	ResetTime time.Time
}

func RateLimit(config RateLimitConfig) fiber.Handler {
	if config.KeyPrefix == "" {
		config.KeyPrefix = "rate_limit"
	}

	return func(c *fiber.Ctx) error {
		var key string

		ip := c.IP()
		if ip != "" {
			key = fmt.Sprintf("%s:ip:%s", config.KeyPrefix, ip)
		} else {
			xForwardedFor := c.Get("X-Forwarded-For")
			if xForwardedFor != "" {
				if host, _, err := net.SplitHostPort(xForwardedFor); err == nil {
					key = fmt.Sprintf("%s:ip:%s", config.KeyPrefix, host)
				} else {
					key = fmt.Sprintf("%s:ip:%s", config.KeyPrefix, xForwardedFor)
				}
			} else {
				key = fmt.Sprintf("%s:ip:unknown", config.KeyPrefix)
			}
		}

		ctx := context.Background()
		now := time.Now()
		windowStart := now.Truncate(config.WindowSize)

		pipe := config.RedisClient.Pipeline()

		incCmd := pipe.Incr(ctx, key)
		pipe.ExpireAt(ctx, key, windowStart.Add(config.WindowSize))
		ttlCmd := pipe.TTL(ctx, key)

		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("Redis error in rate limiter: %v", err)
			return c.Next()
		}

		currentCount := incCmd.Val()
		ttl := ttlCmd.Val()

		if currentCount > int64(config.RequestsLimit) {
			resetTime := now.Add(ttl)

			c.Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsLimit))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			c.Set("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))

			return NewRateLimitError("Rate limit exceeded. Please try again later")
		}

		remaining := int64(config.RequestsLimit) - currentCount
		if remaining < 0 {
			remaining = 0
		}

		resetTime := now.Add(ttl)

		c.Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsLimit))
		c.Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		info := RateLimitInfo{
			Remaining: remaining,
			ResetTime: resetTime,
		}
		c.Locals("rateLimit", info)

		return c.Next()
	}
}

func AuthRateLimit(redisClient *redis.Client) fiber.Handler {
	return RateLimit(RateLimitConfig{
		RedisClient:   redisClient,
		RequestsLimit: 300,
		WindowSize:    time.Minute,
		KeyPrefix:     "auth_rate_limit",
	})
}

func APIRateLimit(redisClient *redis.Client) fiber.Handler {
	return RateLimit(RateLimitConfig{
		RedisClient:   redisClient,
		RequestsLimit: 1000,
		WindowSize:    time.Minute,
		KeyPrefix:     "api_rate_limit",
	})
}

func EvaluateRateLimit(redisClient *redis.Client) fiber.Handler {
	return RateLimit(RateLimitConfig{
		RedisClient:   redisClient,
		RequestsLimit: 1000,
		WindowSize:    time.Minute,
		KeyPrefix:     "eval_rate_limit",
	})
}
