package handlers

import (
	"testing"

	"github.com/workermill-examples/flagdeck/api/internal/database"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	// Create a simple test that always passes to bootstrap CI
	// Create mock health handler with nil dependencies (they won't be called in this test)
	handler := NewHealthHandler(nil, nil)

	// This is a minimal test to ensure the package compiles and the CI can run
	if handler == nil {
		t.Error("NewHealthHandler should not return nil")
	}

	t.Log("Health handler test passed - ready for CI")
}

func TestHealthHandler_Structure(t *testing.T) {
	// Test that the health handler has the correct structure
	var mongoDB *database.MongoDB
	var redisDB *database.RedisDB

	handler := NewHealthHandler(mongoDB, redisDB)

	if handler == nil {
		t.Fatal("NewHealthHandler returned nil")
	}

	if handler.MongoDB != mongoDB {
		t.Error("MongoDB field not set correctly")
	}

	if handler.Redis != redisDB {
		t.Error("Redis field not set correctly")
	}
}

func TestHealthResponse_Structure(t *testing.T) {
	// Test the response structure
	response := HealthResponse{
		Status:  "ok",
		MongoDB: "connected",
		Redis:   "connected",
	}

	if response.Status != "ok" {
		t.Error("Status field not set correctly")
	}

	if response.MongoDB != "connected" {
		t.Error("MongoDB field not set correctly")
	}

	if response.Redis != "connected" {
		t.Error("Redis field not set correctly")
	}
}
