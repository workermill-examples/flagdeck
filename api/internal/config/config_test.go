package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDefaultValues(t *testing.T) {
	// Clear any existing environment variables to test defaults
	os.Unsetenv("MONGODB_URI")
	os.Unsetenv("REDIS_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("CORS_ORIGINS")

	cfg, err := Load()
	require.NoError(t, err, "Loading config with defaults should not error")
	require.NotNil(t, cfg, "Config should not be nil")

	// Verify default values are applied correctly
	assert.Equal(t, "mongodb://localhost:27017/flagdeck", cfg.MongoDBURI, "MongoDBURI should have default value")
	assert.Equal(t, "redis://localhost:6379/0", cfg.RedisURL, "RedisURL should have default value")
	assert.Equal(t, "dev-secret-change-in-prod", cfg.JWTSecret, "JWTSecret should have default value")
	assert.Equal(t, "8080", cfg.Port, "Port should have default value")
	assert.Equal(t, "development", cfg.Environment, "Environment should have default value")
	assert.Equal(t, "http://localhost:3000", cfg.CORSOrigins, "CORSOrigins should have default value")
}

func TestConfigEnvVarOverrides(t *testing.T) {
	// Set test environment variables
	testMongoURI := "mongodb://test:27017/testdb"
	testRedisURL := "redis://test:6379/1"
	testJWTSecret := "test-secret"
	testPort := "9090"
	testEnv := "testing"
	testCORS := "http://test.example.com"

	os.Setenv("MONGODB_URI", testMongoURI)
	os.Setenv("REDIS_URL", testRedisURL)
	os.Setenv("JWT_SECRET", testJWTSecret)
	os.Setenv("PORT", testPort)
	os.Setenv("ENVIRONMENT", testEnv)
	os.Setenv("CORS_ORIGINS", testCORS)

	defer func() {
		// Clean up environment variables
		os.Unsetenv("MONGODB_URI")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("CORS_ORIGINS")
	}()

	cfg, err := Load()
	require.NoError(t, err, "Loading config with env vars should not error")
	require.NotNil(t, cfg, "Config should not be nil")

	// Verify environment variables override defaults
	assert.Equal(t, testMongoURI, cfg.MongoDBURI, "MongoDBURI should use env var value")
	assert.Equal(t, testRedisURL, cfg.RedisURL, "RedisURL should use env var value")
	assert.Equal(t, testJWTSecret, cfg.JWTSecret, "JWTSecret should use env var value")
	assert.Equal(t, testPort, cfg.Port, "Port should use env var value")
	assert.Equal(t, testEnv, cfg.Environment, "Environment should use env var value")
	assert.Equal(t, testCORS, cfg.CORSOrigins, "CORSOrigins should use env var value")
}

func TestConfigStructTags(t *testing.T) {
	// This is a trivial test that verifies the Config struct has the correct field types
	// It ensures the struct is properly defined for parsing
	var cfg Config

	// Verify fields are of expected types (strings in this case)
	assert.IsType(t, "", cfg.MongoDBURI, "MongoDBURI should be a string")
	assert.IsType(t, "", cfg.RedisURL, "RedisURL should be a string")
	assert.IsType(t, "", cfg.JWTSecret, "JWTSecret should be a string")
	assert.IsType(t, "", cfg.Port, "Port should be a string")
	assert.IsType(t, "", cfg.Environment, "Environment should be a string")
	assert.IsType(t, "", cfg.CORSOrigins, "CORSOrigins should be a string")
}
